package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"x-novel/internal/llm"
	"x-novel/internal/repository"
	"x-novel/pkg/logger"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// ========== 数据结构 ==========

// DetectionIssue 检测到的问题
type DetectionIssue struct {
	Type        string `json:"type"`        // typo, grammar, logic, repetition
	Severity    string `json:"severity"`    // error, warning, info
	Position    string `json:"position"`    // 问题位置描述
	Original    string `json:"original"`    // 原文
	Suggestion  string `json:"suggestion"`  // 修改建议
	Explanation string `json:"explanation"` // 解释
}

// DetectionResult 检测结果
type DetectionResult struct {
	Issues     []DetectionIssue `json:"issues"`
	Summary    string           `json:"summary"`
	TotalCount int              `json:"total_count"`
	TypeCounts map[string]int   `json:"type_counts"`
}

// ReviewScore 审阅评分项
type ReviewScore struct {
	Dimension string `json:"dimension"` // plot, character, writing, pacing, creativity, readability
	Score     int    `json:"score"`     // 1-100
	Comment   string `json:"comment"`
}

// ReviewResult AI 审阅结果
type ReviewResult struct {
	Scores      []ReviewScore `json:"scores"`
	OverallScore int          `json:"overall_score"`
	Highlights  []string      `json:"highlights"`
	Issues      []string      `json:"issues"`
	Suggestions []string      `json:"suggestions"`
	Summary     string        `json:"summary"`
}

// MarketPrediction 市场预测结果
type MarketPrediction struct {
	MarketScore     int                `json:"market_score"`      // 市场适配度 1-100
	TargetAudience  string             `json:"target_audience"`   // 目标读者群体
	CompetitiveEdge string             `json:"competitive_edge"`  // 竞争优势
	MarketTrends    []MarketTrendItem  `json:"market_trends"`     // 市场趋势分析
	ReaderAppeal    []ReaderAppealItem `json:"reader_appeal"`     // 读者吸引力维度
	Monetization    MonetizationAdvice `json:"monetization"`      // 变现建议
	Risks           []string           `json:"risks"`             // 风险提示
	Recommendations []string           `json:"recommendations"`   // 改进建议
	Summary         string             `json:"summary"`           // 总结
}

// MarketTrendItem 市场趋势条目
type MarketTrendItem struct {
	Trend    string `json:"trend"`    // 趋势名称
	Fit      int    `json:"fit"`      // 契合度 1-100
	Analysis string `json:"analysis"` // 分析
}

// ReaderAppealItem 读者吸引力维度
type ReaderAppealItem struct {
	Dimension string `json:"dimension"` // 维度
	Score     int    `json:"score"`     // 1-100
	Comment   string `json:"comment"`   // 说明
}

// MonetizationAdvice 变现建议
type MonetizationAdvice struct {
	Platforms     []string `json:"platforms"`     // 推荐平台
	PricingModel  string   `json:"pricing_model"` // 定价模式
	IPPotential   int      `json:"ip_potential"`   // IP 潜力 1-100
	Suggestion    string   `json:"suggestion"`    // 核心建议
}

// ========== 服务 ==========

type ReviewService struct {
	projectRepo *repository.ProjectRepository
	chapterRepo *repository.ChapterRepository
	modelRepo   *repository.ModelConfigRepository
	llmManager  *llm.Manager
}

func NewReviewService(
	projectRepo *repository.ProjectRepository,
	chapterRepo *repository.ChapterRepository,
	modelRepo *repository.ModelConfigRepository,
	llmManager *llm.Manager,
) *ReviewService {
	return &ReviewService{
		projectRepo: projectRepo,
		chapterRepo: chapterRepo,
		modelRepo:   modelRepo,
		llmManager:  llmManager,
	}
}

// DetectErrors 检测章节内容中的问题
func (s *ReviewService) DetectErrors(ctx context.Context, deviceID uuid.UUID, content string, types []string) (*DetectionResult, error) {
	if len(types) == 0 {
		types = []string{"typo", "grammar", "logic", "repetition"}
	}

	prompt := getDetectionPrompt(content, types)
	result, err := s.callLLM(ctx, deviceID, prompt, 0.2, "review")
	if err != nil {
		logger.Error("错误检测 LLM 调用失败", zap.Error(err))
		return s.getMockDetection(content), nil
	}

	detection, err := parseDetectionResult(result)
	if err != nil {
		logger.Error("解析检测结果失败", zap.Error(err))
		return s.getMockDetection(content), nil
	}

	return detection, nil
}

// ReviewChapter 审阅章节
func (s *ReviewService) ReviewChapter(ctx context.Context, deviceID uuid.UUID, projectID string, chapterNumber int) (*ReviewResult, error) {
	project, err := s.projectRepo.GetByID(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("项目不存在: %w", err)
	}

	chapter, err := s.chapterRepo.GetByProjectAndNumber(ctx, projectID, chapterNumber)
	if err != nil {
		return nil, fmt.Errorf("章节不存在: %w", err)
	}

	if chapter.Content == "" {
		return nil, fmt.Errorf("章节内容为空")
	}

	prompt := getReviewPrompt(project.Title, chapterNumber, chapter.Title, chapter.Content)
	result, err := s.callLLM(ctx, deviceID, prompt, 0.3, "review")
	if err != nil {
		logger.Error("AI 审阅 LLM 调用失败", zap.Error(err))
		return s.getMockReview(), nil
	}

	review, err := parseReviewResult(result)
	if err != nil {
		logger.Error("解析审阅结果失败", zap.Error(err))
		return s.getMockReview(), nil
	}

	return review, nil
}

// ReviewProject 审阅整个项目（基于前几章）
func (s *ReviewService) ReviewProject(ctx context.Context, deviceID uuid.UUID, projectID string) (*ReviewResult, error) {
	project, err := s.projectRepo.GetByID(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("项目不存在: %w", err)
	}

	chapters, err := s.chapterRepo.ListByProject(ctx, projectID)
	if err != nil {
		return nil, err
	}

	var contentParts []string
	wordCount := 0
	for _, ch := range chapters {
		if ch.Content == "" {
			continue
		}
		contentParts = append(contentParts, fmt.Sprintf("## 第%d章 %s\n%s", ch.ChapterNumber, ch.Title, ch.Content))
		wordCount += ch.WordCount
		if wordCount > 10000 {
			break
		}
	}

	if len(contentParts) == 0 {
		return nil, fmt.Errorf("没有可审阅的章节内容")
	}

	fullContent := strings.Join(contentParts, "\n\n---\n\n")
	prompt := getProjectReviewPrompt(project.Title, fullContent)
	result, err := s.callLLM(ctx, deviceID, prompt, 0.3, "review")
	if err != nil {
		logger.Error("项目审阅 LLM 调用失败", zap.Error(err))
		return s.getMockReview(), nil
	}

	review, err := parseReviewResult(result)
	if err != nil {
		return s.getMockReview(), nil
	}

	return review, nil
}

// MarketPredict 市场预测
func (s *ReviewService) MarketPredict(ctx context.Context, deviceID uuid.UUID, projectID string) (*MarketPrediction, error) {
	project, err := s.projectRepo.GetByID(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("项目不存在: %w", err)
	}

	chapters, err := s.chapterRepo.ListByProject(ctx, projectID)
	if err != nil {
		return nil, err
	}

	var contentParts []string
	wordCount := 0
	for _, ch := range chapters {
		if ch.Content == "" {
			continue
		}
		contentParts = append(contentParts, fmt.Sprintf("第%d章 %s（%d字）", ch.ChapterNumber, ch.Title, ch.WordCount))
		wordCount += ch.WordCount
		if len(contentParts) >= 5 {
			contentParts = append(contentParts, ch.Content)
		}
	}

	prompt := getMarketPredictPrompt(project.Title, project.Genre, project.CoreSeed, project.PlotArchitecture, strings.Join(contentParts, "\n"))
	result, err := s.callLLM(ctx, deviceID, prompt, 0.4, "review")
	if err != nil {
		logger.Error("市场预测 LLM 调用失败", zap.Error(err))
		return s.getMockMarketPrediction(), nil
	}

	prediction, err := parseMarketPrediction(result)
	if err != nil {
		logger.Error("解析市场预测结果失败", zap.Error(err))
		return s.getMockMarketPrediction(), nil
	}

	return prediction, nil
}

func (s *ReviewService) callLLM(ctx context.Context, deviceID uuid.UUID, prompt string, temperature float32, purpose string) (string, error) {
	modelConfig, err := s.modelRepo.GetByPurpose(ctx, deviceID.String(), purpose)
	if err != nil {
		modelConfig, err = s.modelRepo.GetByPurpose(ctx, deviceID.String(), "general")
		if err != nil {
			return "", fmt.Errorf("未配置模型: %w", err)
		}
	}

	messages := []llm.ChatMessage{
		{Role: "user", Content: prompt},
	}
	options := llm.ChatOptions{
		Temperature: temperature,
		MaxTokens:   4096,
		APIKey:      modelConfig.APIKey,
	}

	if modelConfig.BaseURL != "" {
		adapter := llm.NewOpenAIAdapter(modelConfig.BaseURL, modelConfig.ModelName)
		return adapter.ChatCompletion(ctx, messages, options)
	}

	provider := "openai"
	if modelConfig.Provider != nil {
		provider = modelConfig.Provider.Name
	}
	return s.llmManager.ChatCompletion(ctx, provider, messages, options)
}

// ========== 提示词 ==========

func getDetectionPrompt(content string, types []string) string {
	typeDesc := strings.Join(types, "、")
	return fmt.Sprintf(`你是一位专业的文学编辑和校对专家。请对以下小说文本进行错误检测。

## 检测范围
%s

## 检测类型说明
- typo：错别字、同音字误用
- grammar：病句、语法错误、标点符号问题
- logic：前后矛盾、逻辑不通
- repetition：重复表达、冗余语句

## 输出格式
请严格按照以下 JSON 格式输出，不要添加其他文字：

{
  "issues": [
    {
      "type": "typo|grammar|logic|repetition",
      "severity": "error|warning|info",
      "position": "问题所在位置的简短描述",
      "original": "原文片段（10-30字）",
      "suggestion": "修改建议",
      "explanation": "为什么这是一个问题"
    }
  ],
  "summary": "整体质量评述（20-50字）"
}

## 注意
1. severity 分级：error=严重错误必须修改，warning=建议修改，info=可以考虑
2. 如果文本质量很好没有问题，issues 返回空数组
3. 每个问题的 original 要精确引用原文

## 待检测文本
%s`, typeDesc, content)
}

func getReviewPrompt(title string, chapterNumber int, chapterTitle, content string) string {
	return fmt.Sprintf(`你是一位资深的小说评论家和文学编辑。请对以下章节进行全面审阅。

## 作品信息
- 标题：%s
- 章节：第%d章 %s

## 评分维度（每项 1-100 分）
1. plot（情节）：故事发展是否合理、吸引人
2. character（人物）：角色是否立体、有说服力
3. writing（文笔）：语言是否流畅、有表现力
4. pacing（节奏）：张弛是否有度、结构是否合理
5. creativity（创意）：是否有新意、能否引发思考
6. readability（可读性）：是否易读、让人想继续

## 输出格式
请严格按照以下 JSON 格式输出：

{
  "scores": [
    {"dimension": "plot", "score": 85, "comment": "简评"},
    {"dimension": "character", "score": 80, "comment": "简评"},
    {"dimension": "writing", "score": 75, "comment": "简评"},
    {"dimension": "pacing", "score": 70, "comment": "简评"},
    {"dimension": "creativity", "score": 82, "comment": "简评"},
    {"dimension": "readability", "score": 78, "comment": "简评"}
  ],
  "overall_score": 78,
  "highlights": ["亮点1", "亮点2", "亮点3"],
  "issues": ["问题1", "问题2"],
  "suggestions": ["建议1", "建议2", "建议3"],
  "summary": "总体评价（50-100字）"
}

## 章节内容
%s`, title, chapterNumber, chapterTitle, content)
}

func getProjectReviewPrompt(title string, content string) string {
	return fmt.Sprintf(`你是一位资深的小说评论家。请对以下小说作品进行整体审阅。

## 作品：%s

## 评分维度（每项 1-100 分）
1. plot（情节）：整体故事架构和发展
2. character（人物）：角色塑造的深度和一致性
3. writing（文笔）：语言风格和表现力
4. pacing（节奏）：整体节奏把控
5. creativity（创意）：题材和叙事的新颖程度
6. readability（可读性）：整体阅读体验

## 输出格式
请严格按照以下 JSON 格式输出：

{
  "scores": [
    {"dimension": "plot", "score": 85, "comment": "简评"},
    {"dimension": "character", "score": 80, "comment": "简评"},
    {"dimension": "writing", "score": 75, "comment": "简评"},
    {"dimension": "pacing", "score": 70, "comment": "简评"},
    {"dimension": "creativity", "score": 82, "comment": "简评"},
    {"dimension": "readability", "score": 78, "comment": "简评"}
  ],
  "overall_score": 78,
  "highlights": ["亮点1", "亮点2", "亮点3"],
  "issues": ["问题1", "问题2"],
  "suggestions": ["建议1", "建议2", "建议3"],
  "summary": "总体评价（50-100字）"
}

## 作品内容
%s`, title, content)
}

func getMarketPredictPrompt(title, genre, coreSeed, plotArchitecture, chaptersSummary string) string {
	return fmt.Sprintf(`你是一位网络文学市场分析专家，对中国网文市场（包括起点、番茄、晋江等平台）有深入了解。
请基于以下小说信息进行市场预测分析。

## 作品信息
- 标题：%s
- 类型：%s
- 核心设定：%s

## 情节架构
%s

## 章节概况
%s

## 输出格式
请严格按照以下 JSON 格式输出：

{
  "market_score": 75,
  "target_audience": "18-30岁男性读者，偏好XX类型",
  "competitive_edge": "与同类作品相比的独特优势",
  "market_trends": [
    {"trend": "趋势名称", "fit": 80, "analysis": "分析说明"},
    {"trend": "趋势名称", "fit": 60, "analysis": "分析说明"}
  ],
  "reader_appeal": [
    {"dimension": "情感共鸣", "score": 80, "comment": "说明"},
    {"dimension": "代入感", "score": 75, "comment": "说明"},
    {"dimension": "爽点密度", "score": 70, "comment": "说明"},
    {"dimension": "悬念设置", "score": 85, "comment": "说明"},
    {"dimension": "社交传播性", "score": 65, "comment": "说明"}
  ],
  "monetization": {
    "platforms": ["推荐平台1", "推荐平台2"],
    "pricing_model": "推荐的定价/付费模式",
    "ip_potential": 70,
    "suggestion": "核心变现建议"
  },
  "risks": ["风险1", "风险2"],
  "recommendations": ["建议1", "建议2", "建议3"],
  "summary": "50-100字的市场分析总结"
}

注意：
1. market_trends 至少给出 3 个趋势分析
2. reader_appeal 包含情感共鸣、代入感、爽点密度、悬念设置、社交传播性 5 个维度
3. 所有 score/fit 范围为 1-100
4. 分析要基于当前网文市场实际情况
`, title, genre, coreSeed, plotArchitecture, chaptersSummary)
}

// ========== 解析 ==========

func parseDetectionResult(raw string) (*DetectionResult, error) {
	cleaned := cleanJSONStr(raw)
	var result DetectionResult
	if err := json.Unmarshal([]byte(cleaned), &result); err != nil {
		return nil, err
	}
	if result.Issues == nil {
		result.Issues = []DetectionIssue{}
	}
	result.TotalCount = len(result.Issues)
	result.TypeCounts = make(map[string]int)
	for _, issue := range result.Issues {
		result.TypeCounts[issue.Type]++
	}
	return &result, nil
}

func parseReviewResult(raw string) (*ReviewResult, error) {
	cleaned := cleanJSONStr(raw)
	var result ReviewResult
	if err := json.Unmarshal([]byte(cleaned), &result); err != nil {
		return nil, err
	}
	if result.Scores == nil {
		result.Scores = []ReviewScore{}
	}
	if result.Highlights == nil {
		result.Highlights = []string{}
	}
	if result.Issues == nil {
		result.Issues = []string{}
	}
	if result.Suggestions == nil {
		result.Suggestions = []string{}
	}
	return &result, nil
}

func parseMarketPrediction(raw string) (*MarketPrediction, error) {
	cleaned := cleanJSONStr(raw)
	var result MarketPrediction
	if err := json.Unmarshal([]byte(cleaned), &result); err != nil {
		return nil, err
	}
	if result.MarketTrends == nil {
		result.MarketTrends = []MarketTrendItem{}
	}
	if result.ReaderAppeal == nil {
		result.ReaderAppeal = []ReaderAppealItem{}
	}
	if result.Risks == nil {
		result.Risks = []string{}
	}
	if result.Recommendations == nil {
		result.Recommendations = []string{}
	}
	return &result, nil
}

func cleanJSONStr(s string) string {
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "```json")
	s = strings.TrimPrefix(s, "```")
	s = strings.TrimSuffix(s, "```")
	return strings.TrimSpace(s)
}

// ========== 模拟数据 ==========

func (s *ReviewService) getMockDetection(content string) *DetectionResult {
	issues := []DetectionIssue{
		{Type: "grammar", Severity: "warning", Position: "段落开头", Original: "（从原文中提取）", Suggestion: "建议调整语序使表达更流畅", Explanation: "当前句式略显冗长，可以拆分"},
		{Type: "repetition", Severity: "info", Position: "中间段落", Original: "（从原文中提取）", Suggestion: "考虑使用近义词替换", Explanation: "连续使用了相同的描述词"},
	}
	typeCounts := map[string]int{"grammar": 1, "repetition": 1}
	return &DetectionResult{
		Issues:     issues,
		Summary:    "文本整体质量良好，有少量语法和重复问题可以优化（模拟数据）",
		TotalCount: len(issues),
		TypeCounts: typeCounts,
	}
}

func (s *ReviewService) getMockMarketPrediction() *MarketPrediction {
	return &MarketPrediction{
		MarketScore:    72,
		TargetAudience: "18-35岁男性读者，偏好都市异能、科幻类型",
		CompetitiveEdge: "将程序员文化与超能力设定结合，贴近年轻读者群体，代入感强",
		MarketTrends: []MarketTrendItem{
			{Trend: "都市异能热度回升", Fit: 85, Analysis: "近年都市异能类作品在多平台表现优秀，本作设定契合"},
			{Trend: "科技与超能力融合", Fit: 78, Analysis: "AI+超能力题材新颖，符合当下科技焦点"},
			{Trend: "短视频改编潜力", Fit: 65, Analysis: "节奏紧凑，有一定短剧改编潜力"},
		},
		ReaderAppeal: []ReaderAppealItem{
			{Dimension: "情感共鸣", Score: 78, Comment: "程序员主角设定易引发职场读者共鸣"},
			{Dimension: "代入感", Score: 82, Comment: "现代都市背景提供良好代入体验"},
			{Dimension: "爽点密度", Score: 70, Comment: "超能力获得节奏尚可，建议增加反转"},
			{Dimension: "悬念设置", Score: 75, Comment: "黑暗组织线索铺设合理"},
			{Dimension: "社交传播性", Score: 68, Comment: "需要增加更多社交媒体友好的金句和名场面"},
		},
		Monetization: MonetizationAdvice{
			Platforms:    []string{"番茄小说", "起点中文网", "七猫"},
			PricingModel: "免费+广告 或 VIP 章节付费",
			IPPotential:  65,
			Suggestion:   "建议先在番茄小说免费连载积累读者，后续考虑付费转化（模拟数据）",
		},
		Risks:           []string{"都市异能类竞品较多，需要强化差异化", "程序员题材可能限制受众范围", "前期节奏需要更紧凑"},
		Recommendations: []string{"开头前三章加入更强的钩子，提高留存率", "增加更多爽点和高潮情节", "考虑设计具有传播力的名场面"},
		Summary:         "本作在都市异能赛道具有一定竞争力，程序员+超能力设定有新意。建议优化开头节奏和爽点密度，可以先在免费平台积累口碑。（模拟数据）",
	}
}

func (s *ReviewService) getMockReview() *ReviewResult {
	return &ReviewResult{
		Scores: []ReviewScore{
			{Dimension: "plot", Score: 78, Comment: "情节发展合理，但转折可以更自然"},
			{Dimension: "character", Score: 82, Comment: "主角形象鲜明，配角可进一步丰富"},
			{Dimension: "writing", Score: 75, Comment: "文笔流畅，部分描写可以更细腻"},
			{Dimension: "pacing", Score: 70, Comment: "前半段节奏偏快，建议适当放缓"},
			{Dimension: "creativity", Score: 80, Comment: "设定有新意，可以进一步深挖"},
			{Dimension: "readability", Score: 85, Comment: "可读性强，适合目标受众"},
		},
		OverallScore: 78,
		Highlights:   []string{"开篇引人入胜，悬念设置巧妙", "角色对话自然生动", "世界观设定有独创性"},
		Issues:       []string{"部分过渡段落衔接不够自然", "次要角色形象略显单薄"},
		Suggestions:  []string{"建议在转折处增加角色心理描写", "可以通过细节描写增强场景沉浸感", "考虑加入更多伏笔增强故事深度"},
		Summary:      "这是一部有潜力的作品，情节设计和角色塑造都有亮点，文笔流畅可读性强。建议在节奏把控和细节描写上进一步打磨。（模拟数据）",
	}
}
