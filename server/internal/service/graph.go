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

// GraphNode 图谱节点
type GraphNode struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Type        string   `json:"type"` // protagonist, antagonist, supporting, minor
	Description string   `json:"description"`
	Traits      []string `json:"traits"`
	Group       string   `json:"group"`
}

// GraphEdge 图谱边
type GraphEdge struct {
	Source      string `json:"source"`
	Target      string `json:"target"`
	Relation    string `json:"relation"`
	Description string `json:"description"`
	Weight      int    `json:"weight"`
}

// GraphData 完整图谱数据
type GraphData struct {
	Nodes    []GraphNode    `json:"nodes"`
	Edges    []GraphEdge    `json:"edges"`
	Snapshots []GraphSnapshot `json:"snapshots,omitempty"`
}

// GraphSnapshot 章节快照
type GraphSnapshot struct {
	ChapterNumber int    `json:"chapter_number"`
	Summary       string `json:"summary"`
	NodesCount    int    `json:"nodes_count"`
	EdgesCount    int    `json:"edges_count"`
}

// ChapterDelta 章节增量数据
type ChapterDelta struct {
	NewNodes       []GraphNode `json:"new_nodes"`
	NewEdges       []GraphEdge `json:"new_edges"`
	UpdatedEdges   []GraphEdge `json:"updated_edges"`
	ChapterSummary string      `json:"chapter_summary"`
}

type GraphService struct {
	projectRepo *repository.ProjectRepository
	chapterRepo *repository.ChapterRepository
	modelRepo   *repository.ModelConfigRepository
	llmManager  *llm.Manager
}

func NewGraphService(
	projectRepo *repository.ProjectRepository,
	chapterRepo *repository.ChapterRepository,
	modelRepo *repository.ModelConfigRepository,
	llmManager *llm.Manager,
) *GraphService {
	return &GraphService{
		projectRepo: projectRepo,
		chapterRepo: chapterRepo,
		modelRepo:   modelRepo,
		llmManager:  llmManager,
	}
}

// GenerateGraph 从项目架构生成初始图谱
func (s *GraphService) GenerateGraph(ctx context.Context, deviceID uuid.UUID, projectID string) (*GraphData, error) {
	project, err := s.projectRepo.GetByID(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("项目不存在: %w", err)
	}

	if project.CoreSeed == "" {
		return nil, fmt.Errorf("请先生成小说架构")
	}

	prompt := GetExtractGraphPrompt(project.Title, project.CoreSeed, project.CharacterDynamics, project.WorldBuilding)

	result, err := s.callLLM(ctx, deviceID, prompt)
	if err != nil {
		logger.Error("LLM 提取图谱失败，使用模拟数据", zap.Error(err))
		result = s.getMockGraphJSON(project.Title)
	}

	graphData, err := s.parseGraphData(result)
	if err != nil {
		logger.Error("解析图谱数据失败，使用模拟数据", zap.Error(err))
		graphData = s.getMockGraph(project.Title)
	}

	// 保存到项目
	graphJSON, _ := json.Marshal(graphData)
	project.GraphData = string(graphJSON)
	if err := s.projectRepo.Update(ctx, project); err != nil {
		return nil, fmt.Errorf("保存图谱失败: %w", err)
	}

	return graphData, nil
}

// UpdateGraphFromChapter 从章节内容更新图谱（增量）
func (s *GraphService) UpdateGraphFromChapter(ctx context.Context, deviceID uuid.UUID, projectID string, chapterNumber int) (*GraphData, error) {
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

	// 加载已有图谱
	var graphData GraphData
	if project.GraphData != "" {
		json.Unmarshal([]byte(project.GraphData), &graphData)
	}

	existingJSON, _ := json.Marshal(graphData)
	prompt := GetExtractChapterGraphPrompt(project.Title, chapterNumber, chapter.Content, string(existingJSON))

	result, err := s.callLLM(ctx, deviceID, prompt)
	if err != nil {
		logger.Error("LLM 提取章节增量失败", zap.Error(err))
		return s.applyMockDelta(&graphData, chapterNumber), nil
	}

	delta, err := s.parseChapterDelta(result)
	if err != nil {
		logger.Error("解析章节增量失败", zap.Error(err))
		return s.applyMockDelta(&graphData, chapterNumber), nil
	}

	s.applyDelta(&graphData, delta, chapterNumber)

	// 保存
	graphJSON, _ := json.Marshal(graphData)
	project.GraphData = string(graphJSON)
	if err := s.projectRepo.Update(ctx, project); err != nil {
		return nil, fmt.Errorf("保存图谱失败: %w", err)
	}

	// 保存章节图谱快照
	chapterSnapshot, _ := json.Marshal(graphData)
	chapter.ChapterGraph = string(chapterSnapshot)
	s.chapterRepo.Update(ctx, chapter)

	return &graphData, nil
}

// GetGraph 获取项目图谱
func (s *GraphService) GetGraph(ctx context.Context, projectID string) (*GraphData, error) {
	project, err := s.projectRepo.GetByID(ctx, projectID)
	if err != nil {
		return nil, err
	}

	if project.GraphData == "" {
		return &GraphData{Nodes: []GraphNode{}, Edges: []GraphEdge{}}, nil
	}

	var graphData GraphData
	if err := json.Unmarshal([]byte(project.GraphData), &graphData); err != nil {
		return nil, fmt.Errorf("解析图谱数据失败: %w", err)
	}

	return &graphData, nil
}

// GetChapterSnapshot 获取某章节时的图谱快照
func (s *GraphService) GetChapterSnapshot(ctx context.Context, projectID string, chapterNumber int) (*GraphData, error) {
	chapter, err := s.chapterRepo.GetByProjectAndNumber(ctx, projectID, chapterNumber)
	if err != nil {
		return nil, err
	}

	if chapter.ChapterGraph == "" {
		return s.GetGraph(ctx, projectID)
	}

	var graphData GraphData
	if err := json.Unmarshal([]byte(chapter.ChapterGraph), &graphData); err != nil {
		return nil, err
	}

	return &graphData, nil
}

func (s *GraphService) applyDelta(graph *GraphData, delta *ChapterDelta, chapterNumber int) {
	existingNodeIDs := make(map[string]bool)
	for _, n := range graph.Nodes {
		existingNodeIDs[n.ID] = true
	}

	for _, node := range delta.NewNodes {
		if !existingNodeIDs[node.ID] {
			graph.Nodes = append(graph.Nodes, node)
			existingNodeIDs[node.ID] = true
		}
	}

	graph.Edges = append(graph.Edges, delta.NewEdges...)

	for _, updated := range delta.UpdatedEdges {
		found := false
		for i, edge := range graph.Edges {
			if (edge.Source == updated.Source && edge.Target == updated.Target) ||
				(edge.Source == updated.Target && edge.Target == updated.Source) {
				graph.Edges[i] = updated
				found = true
				break
			}
		}
		if !found {
			graph.Edges = append(graph.Edges, updated)
		}
	}

	graph.Snapshots = append(graph.Snapshots, GraphSnapshot{
		ChapterNumber: chapterNumber,
		Summary:       delta.ChapterSummary,
		NodesCount:    len(graph.Nodes),
		EdgesCount:    len(graph.Edges),
	})
}

func (s *GraphService) callLLM(ctx context.Context, deviceID uuid.UUID, prompt string) (string, error) {
	modelConfig, err := s.modelRepo.GetByPurpose(ctx, deviceID.String(), "architecture")
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
		Temperature: 0.3,
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

func (s *GraphService) parseGraphData(raw string) (*GraphData, error) {
	cleaned := cleanJSON(raw)
	var data GraphData
	if err := json.Unmarshal([]byte(cleaned), &data); err != nil {
		return nil, err
	}
	if data.Nodes == nil {
		data.Nodes = []GraphNode{}
	}
	if data.Edges == nil {
		data.Edges = []GraphEdge{}
	}
	return &data, nil
}

func (s *GraphService) parseChapterDelta(raw string) (*ChapterDelta, error) {
	cleaned := cleanJSON(raw)
	var delta ChapterDelta
	if err := json.Unmarshal([]byte(cleaned), &delta); err != nil {
		return nil, err
	}
	return &delta, nil
}

// cleanJSON 清理 LLM 返回中可能包含的 markdown 代码块标记
func cleanJSON(s string) string {
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "```json")
	s = strings.TrimPrefix(s, "```")
	s = strings.TrimSuffix(s, "```")
	return strings.TrimSpace(s)
}

func (s *GraphService) getMockGraphJSON(title string) string {
	data := s.getMockGraph(title)
	b, _ := json.Marshal(data)
	return string(b)
}

func (s *GraphService) getMockGraph(title string) *GraphData {
	return &GraphData{
		Nodes: []GraphNode{
			{ID: "liming", Name: "李明", Type: "protagonist", Description: "偶然获得超能力的程序员", Traits: []string{"聪明", "正义感强"}, Group: "主角阵营"},
			{ID: "suxiaoyu", Name: "苏小雨", Type: "supporting", Description: "神秘的短发女子", Traits: []string{"敏捷", "果断"}, Group: "主角阵营"},
			{ID: "heiyiren", Name: "黑衣人", Type: "antagonist", Description: "神秘组织的执行者", Traits: []string{"冷酷", "强大"}, Group: "神秘组织"},
			{ID: "laoshi", Name: "陈教授", Type: "supporting", Description: "李明的导师", Traits: []string{"睿智", "慈祥"}, Group: "学术界"},
		},
		Edges: []GraphEdge{
			{Source: "liming", Target: "suxiaoyu", Relation: "盟友", Description: "苏小雨救了李明，两人结为同盟", Weight: 8},
			{Source: "liming", Target: "heiyiren", Relation: "敌对", Description: "黑衣人追杀李明", Weight: 9},
			{Source: "suxiaoyu", Target: "heiyiren", Relation: "前同事", Description: "苏小雨曾是组织成员，后叛逃", Weight: 7},
			{Source: "liming", Target: "laoshi", Relation: "师生", Description: "陈教授是李明的大学导师", Weight: 5},
		},
	}
}

func (s *GraphService) applyMockDelta(graph *GraphData, chapterNumber int) *GraphData {
	graph.Snapshots = append(graph.Snapshots, GraphSnapshot{
		ChapterNumber: chapterNumber,
		Summary:       fmt.Sprintf("第%d章关系变化（模拟数据）", chapterNumber),
		NodesCount:    len(graph.Nodes),
		EdgesCount:    len(graph.Edges),
	})
	return graph
}
