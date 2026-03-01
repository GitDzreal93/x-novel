package service

import "fmt"

// WritingAssistantAction 写作助手动作类型
type WritingAssistantAction string

const (
	ActionPolish      WritingAssistantAction = "polish"      // 润色
	ActionContinue    WritingAssistantAction = "continue"     // 续写
	ActionSuggestion  WritingAssistantAction = "suggestion"   // 灵感建议
)

// GetPolishPrompt 润色提示词
func GetPolishPrompt(content, style string) string {
	styleDesc := "保持原有风格"
	switch style {
	case "vivid":
		styleDesc = "让文字更加生动形象，增加修辞和细节描写"
	case "concise":
		styleDesc = "让文字更加精炼简洁，去除冗余表达"
	case "literary":
		styleDesc = "增强文学性，提升语言的美感和节奏感"
	case "dramatic":
		styleDesc = "增强戏剧张力，让情节更加紧凑吸引人"
	}

	return fmt.Sprintf(`你是一位资深的文学编辑，请对以下小说段落进行润色修改。

## 润色要求
%s

## 规则
1. 保持原文的核心剧情和人物不变
2. 不要改变叙事视角和时态
3. 润色后的文字应该自然流畅
4. 直接输出润色后的完整文本，不要加任何说明

## 原文
%s`, styleDesc, content)
}

// GetContinuePrompt 续写提示词
func GetContinuePrompt(content string, targetWords int, context string) string {
	contextSection := ""
	if context != "" {
		contextSection = fmt.Sprintf("\n## 小说背景\n%s\n", context)
	}

	return fmt.Sprintf(`你是一位专业的小说作家，请基于以下已有内容继续写作。
%s
## 续写要求
1. 续写约 %d 字
2. 保持与前文一致的风格、语气和叙事视角
3. 情节发展要自然合理，不能出现突兀转折
4. 注意人物性格和语言习惯的一致性
5. 直接续写，不要重复已有内容，不要加说明

## 已有内容（请从此处继续）
%s`, contextSection, targetWords, content)
}

// GetSuggestionPrompt 灵感建议提示词
func GetSuggestionPrompt(content, aspect string, context string) string {
	aspectDesc := "故事后续发展"
	switch aspect {
	case "plot":
		aspectDesc = "接下来的情节发展方向"
	case "character":
		aspectDesc = "角色行为和心理刻画"
	case "dialogue":
		aspectDesc = "人物对话设计"
	case "description":
		aspectDesc = "环境和氛围描写"
	case "conflict":
		aspectDesc = "冲突和悬念设计"
	}

	contextSection := ""
	if context != "" {
		contextSection = fmt.Sprintf("\n## 小说背景\n%s\n", context)
	}

	return fmt.Sprintf(`你是一位经验丰富的小说创作顾问，请根据以下内容提供创作建议。
%s
## 建议方向
重点分析和建议：**%s**

## 输出格式
请提供 3-5 条具体可行的建议，每条建议包含：
- 建议标题（简短概括）
- 详细说明（具体的实施方式和示例）

## 当前内容
%s`, contextSection, aspectDesc, content)
}
