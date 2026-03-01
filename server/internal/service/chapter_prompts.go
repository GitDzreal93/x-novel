package service

import (
	"fmt"
	"strings"
)

// ChapterPromptParams 章节生成参数
type ChapterPromptParams struct {
	// 项目信息
	Title            string
	Topic            string
	Genre            []string
	UserGuidance     string
	WordsPerChapter  int

	// 架构信息
	CoreSeed          string
	CharacterDynamics string
	WorldBuilding     string
	PlotArchitecture  string
	CharacterState    string

	// 章节信息
	ChapterNumber     int
	ChapterTitle      string
	BlueprintSummary  string

	// 上下文
	GlobalSummary     string
	PreviousSummary   string // 前一章摘要

	// 当前章节内容（用于扩写）
	CurrentContent    string
	TargetWords       int
}

// GetFirstDraftPrompt 获取第一章草稿提示词
func GetFirstDraftPrompt(params ChapterPromptParams) string {
	genreStr := strings.Join(params.Genre, "、")

	return fmt.Sprintf(`你是一位专业的小说作家，现在需要为一部【%s】类型的小说撰写第一章。

## 小说基本信息
- 标题：%s
- 主题：%s
- 类型：%s
- 每章目标字数：约 %d 字

## 核心设定
%s

## 角色体系
%s

## 世界观
%s

## 情节架构
%s

## 角色状态
%s

## 本章大纲
章节号：第 %d 章
章节标题：%s
章节摘要：%s

## 写作要求
1. **严格遵循【%s】类型的写作风格和情感基调**
2. 以生动的场景描写开篇，迅速吸引读者
3. 自然地引入主要角色和背景设定
4. 在章节末尾设置悬念或引子，吸引读者继续阅读
5. 目标字数约 %d 字，确保内容充实但不拖沓
6. 使用第三人称视角
7. 对话要自然流畅，符合人物性格
8. 注重细节描写，让场景具有画面感

## 输出要求
直接输出章节正文内容，不要包含章节标题、作者注释或任何额外说明。`,
		genreStr, params.Title, params.Topic, genreStr, params.WordsPerChapter,
		params.CoreSeed, params.CharacterDynamics, params.WorldBuilding,
		params.PlotArchitecture, params.CharacterState,
		params.ChapterNumber, params.ChapterTitle, params.BlueprintSummary,
		genreStr, params.WordsPerChapter)
}

// GetNextDraftPrompt 获取后续章节草稿提示词
func GetNextDraftPrompt(params ChapterPromptParams) string {
	genreStr := strings.Join(params.Genre, "、")

	return fmt.Sprintf(`你是一位专业的小说作家，现在需要为一部【%s】类型的小说撰写第 %d 章。

## 小说基本信息
- 标题：%s
- 类型：%s
- 每章目标字数：约 %d 字

## 核心设定
%s

## 角色状态（当前）
%s

## 前文摘要
%s

## 本章大纲
章节号：第 %d 章
章节标题：%s
章节摘要：%s

## 写作要求
1. **严格遵循【%s】类型的写作风格和情感基调**
2. 承接上一章的剧情，保持故事连贯性
3. 按照本章大纲推进剧情
4. 保持人物性格和行为的一致性
5. 目标字数约 %d 字
6. 使用第三人称视角
7. 在章节末尾适当设置悬念或铺垫
8. 对话要自然，符合人物性格特点

## 输出要求
直接输出章节正文内容，不要包含章节标题、作者注释或任何额外说明。`,
		genreStr, params.ChapterNumber,
		params.Title, genreStr, params.WordsPerChapter,
		params.CoreSeed, params.CharacterState,
		params.GlobalSummary,
		params.ChapterNumber, params.ChapterTitle, params.BlueprintSummary,
		genreStr, params.WordsPerChapter)
}

// GetEnrichPrompt 获取扩写提示词
func GetEnrichPrompt(params ChapterPromptParams) string {
	genreStr := strings.Join(params.Genre, "、")
	currentWordCount := len([]rune(params.CurrentContent))

	return fmt.Sprintf(`你是一位专业的小说编辑，现在需要对一段【%s】类型小说的章节内容进行扩写。

## 当前内容
%s

## 扩写要求
1. 当前字数：约 %d 字
2. 目标字数：约 %d 字
3. 需要增加：约 %d 字

## 扩写方向
1. 增加更多的环境描写和氛围渲染
2. 深化人物的心理活动描写
3. 丰富对话内容，增加人物互动
4. 补充细节，让场景更有画面感
5. 适当增加过渡段落，使情节更流畅

## 扩写原则
1. 保持原有剧情走向不变
2. 保持人物性格一致
3. 不要改变原有的情节逻辑
4. 扩写的内容要自然融入原文
5. 符合【%s】类型的风格特点

## 输出要求
直接输出扩写后的完整章节内容，不要包含任何额外说明。`,
		genreStr, params.CurrentContent,
		currentWordCount, params.TargetWords, params.TargetWords-currentWordCount,
		genreStr)
}

// GetSummaryPrompt 获取章节摘要生成提示词
func GetSummaryPrompt(chapterContent string, chapterNumber int) string {
	return fmt.Sprintf(`请为以下第 %d 章的内容生成一个简洁的摘要，用于帮助后续章节保持剧情连贯性。

## 章节内容
%s

## 摘要要求
1. 概括本章的主要事件和情节发展
2. 记录重要的人物出场和关系变化
3. 标注关键的伏笔和悬念
4. 总结人物的心理变化或成长
5. 控制在 200-300 字以内

## 输出格式
【主要事件】xxx
【人物动态】xxx
【伏笔/悬念】xxx
【情感变化】xxx`,
		chapterNumber, chapterContent)
}

// GetUpdateCharacterStatePrompt 获取更新角色状态的提示词
func GetUpdateCharacterStatePrompt(currentState string, chapterContent string, chapterNumber int) string {
	return fmt.Sprintf(`根据第 %d 章的内容，更新角色状态文档。

## 当前角色状态
%s

## 本章内容
%s

## 更新要求
1. 根据本章发生的事件，更新相关角色的状态
2. 包括：物品变化、能力变化、身心状态、关系网变化
3. 新增触发的重要事件
4. 保持原有格式结构
5. 只更新有变化的部分，未变化的保持原样

## 输出要求
输出更新后的完整角色状态文档，使用原有的树形格式。`,
		chapterNumber, currentState, chapterContent)
}

// GetChapterPrompt 根据章节号获取对应的提示词
func GetChapterPrompt(chapterNumber int, params ChapterPromptParams) string {
	if chapterNumber == 1 {
		return GetFirstDraftPrompt(params)
	}
	return GetNextDraftPrompt(params)
}
