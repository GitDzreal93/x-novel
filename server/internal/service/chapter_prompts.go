package service

import (
	"fmt"
)

// ChapterDraftParams 章节草稿生成参数
type ChapterDraftParams struct {
	ChapterNumber       int
	ChapterTitle        string
	ChapterRole         string
	ChapterPurpose      string
	SuspenseLevel       string
	Foreshadowing       string
	PlotTwistLevel      string
	ChapterSummary      string
	WordNumber          int
	CharactersInvolved  string
	KeyItems            string
	SceneLocation       string
	TimeConstraint      string

	// 架构信息
	CoreSeed           string
	CharacterDynamics  string
	WorldBuilding      string
	PlotArchitecture   string
	CharacterState     string

	// 上下文信息
	GlobalSummary         string
	PreviousChapterExcerpt string
	IsFirstChapter        bool

	UserGuidance       string
}

// BuildFirstDraftPrompt 构建第一章草稿提示词
func BuildFirstDraftPrompt(params ChapterDraftParams) string {
	novelSetting := fmt.Sprintf(`核心种子：%s

角色动力学：
%s

世界观：
%s

情节架构：
%s`,
		params.CoreSeed,
		params.CharacterDynamics,
		params.WorldBuilding,
		params.PlotArchitecture)

	userGuidance := params.UserGuidance
	if userGuidance == "" {
		userGuidance = "（无）"
	}

	return fmt.Sprintf(`即将创作：第 %d 章《%s》
本章定位：%s
核心作用：%s
悬念密度：%s
伏笔操作：%s
认知颠覆：%s
本章简述：%s

可用元素：
- 核心人物：%s
- 关键道具：%s
- 空间坐标：%s
- 时间压力：%s

参考文档：
- 小说设定：
%s

完成第 %d 章的正文，字数要求%d字，根据小说类型设计2个或以上的场景：
1. 对话场景：
   - 体现人物性格和关系
   - 推动剧情或情感发展

2. 动作/互动场景：
   - 环境交互细节（感官描写）
   - 节奏控制（根据情节需要调整）
   - 通过行动展现人物特质

3. 心理/情感场景：
   - 人物内心活动描写
   - 情感变化的细腻刻画
   - 符合小说类型的情感基调

4. 环境场景：
   - 场景氛围营造
   - 环境与人物心情的呼应
   - 符合小说类型的整体风格

格式要求：
- 仅返回章节正文文本；
- 不使用分章节小标题；
- 不要使用markdown格式。

额外指导：%s`,
		params.ChapterNumber,
		params.ChapterTitle,
		params.ChapterRole,
		params.ChapterPurpose,
		params.SuspenseLevel,
		params.Foreshadowing,
		params.PlotTwistLevel,
		params.ChapterSummary,
		orDefault(params.CharactersInvolved, "（未指定）"),
		orDefault(params.KeyItems, "（未指定）"),
		orDefault(params.SceneLocation, "（未指定）"),
		orDefault(params.TimeConstraint, "（未指定）"),
		novelSetting,
		params.ChapterNumber,
		params.WordNumber,
		userGuidance)
}

// BuildNextDraftPrompt 构建后续章节草稿提示词
func BuildNextDraftPrompt(params ChapterDraftParams) string {
	userGuidance := params.UserGuidance
	if userGuidance == "" {
		userGuidance = "（无）"
	}

	return fmt.Sprintf(`参考文档：
└── 前文摘要：
    %s

└── 前章结尾段：
    %s

└── 用户指导：
    %s

└── 角色状态：
    %s

└── 当前章节摘要：
    %s

当前章节信息：
第%d章《%s》：
├── 章节定位：%s
├── 核心作用：%s
├── 悬念密度：%s
├── 伏笔设计：%s
├── 转折程度：%s
├── 章节简述：%s
├── 字数要求：%d字
├── 核心人物：%s
├── 关键道具：%s
├── 场景地点：%s
└── 时间压力：%s

依据前面所有设定，开始完成第 %d 章的正文，字数要求%d字，
内容生成严格遵循：
- 用户指导
- 当前章节摘要
- 无逻辑漏洞
确保章节内容与前文摘要、前章结尾段衔接流畅，

格式要求：
- 仅返回章节正文文本；
- 不使用分章节小标题；
- 不要使用markdown格式。

根据小说类型设计2个或以上的场景：
1. 对话场景：体现人物性格和关系，推动剧情或情感发展
2. 动作/互动场景：环境交互细节（感官描写），节奏控制
3. 心理/情感场景：人物内心活动描写，情感变化细腻刻画
4. 环境场景：场景氛围营造，环境与人物心情呼应`,
		orDefault(params.GlobalSummary, "（无）"),
		orDefault(params.PreviousChapterExcerpt, "（无）"),
		userGuidance,
		orDefault(params.CharacterState, "（无）"),
		orDefault(params.ChapterSummary, "（无）"),
		params.ChapterNumber,
		params.ChapterTitle,
		params.ChapterRole,
		params.ChapterPurpose,
		params.SuspenseLevel,
		params.Foreshadowing,
		params.PlotTwistLevel,
		params.ChapterSummary,
		params.WordNumber,
		orDefault(params.CharactersInvolved, "（未指定）"),
		orDefault(params.KeyItems, "（未指定）"),
		orDefault(params.SceneLocation, "（未指定）"),
		orDefault(params.TimeConstraint, "（未指定）"),
		params.ChapterNumber,
		params.WordNumber)
}

// ChapterEnrichParams 章节扩写参数
type ChapterEnrichParams struct {
	ChapterText string
	WordNumber  int
}

// BuildEnrichPrompt 构建章节扩写提示词
func BuildEnrichPrompt(params ChapterEnrichParams) string {
	return fmt.Sprintf(`以下章节文本较短，请在保持剧情连贯的前提下进行扩写，使其更充实，接近 %d 字左右。

原内容：
%s

要求：
- 保持剧情连贯性
- 增加环境描写和心理描写
- 丰富对话和场景细节
- 仅给出最终文本，不要解释任何内容`,
		params.WordNumber,
		params.ChapterText)
}

// orDefault 返回默认值
func orDefault(s, defaultVal string) string {
	if s == "" {
		return defaultVal
	}
	return s
}
