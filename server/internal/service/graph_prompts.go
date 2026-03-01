package service

import "fmt"

// GetExtractGraphPrompt 从小说架构中提取初始关系图谱
func GetExtractGraphPrompt(title, coreSeed, characterDynamics, worldBuilding string) string {
	return fmt.Sprintf(`你是一位专业的小说分析师，请从以下小说架构中提取角色关系图谱。

## 小说标题
%s

## 核心设定
%s

## 角色动态
%s

## 世界观
%s

## 输出要求
请严格按照以下 JSON 格式输出，不要添加任何其他文字或 markdown 标记：

{
  "nodes": [
    {
      "id": "唯一标识（使用角色名拼音或英文）",
      "name": "角色名称",
      "type": "protagonist|antagonist|supporting|minor",
      "description": "角色简短描述（20字以内）",
      "traits": ["性格特点1", "性格特点2"],
      "group": "所属阵营或组织"
    }
  ],
  "edges": [
    {
      "source": "角色A的id",
      "target": "角色B的id",
      "relation": "关系类型（如：师徒、恋人、对手、盟友、兄弟、上下级等）",
      "description": "关系的具体描述",
      "weight": 1到10的关系紧密度
    }
  ]
}

## 注意
1. 至少提取 3 个主要角色
2. 关系应双向考虑，但只需输出一条边
3. weight 越大表示关系越密切
4. type 分类：protagonist=主角，antagonist=反派，supporting=重要配角，minor=次要角色`, title, coreSeed, characterDynamics, worldBuilding)
}

// GetExtractChapterGraphPrompt 从单章内容中提取角色关系增量
func GetExtractChapterGraphPrompt(title string, chapterNumber int, chapterContent, existingGraphJSON string) string {
	return fmt.Sprintf(`你是一位专业的小说分析师，请分析以下章节内容，提取**新增或变化**的角色关系。

## 小说标题
%s

## 当前章节
第 %d 章

## 已有图谱
%s

## 章节内容
%s

## 输出要求
请严格按照以下 JSON 格式输出变化部分，不要添加任何其他文字或 markdown 标记：

{
  "new_nodes": [
    {
      "id": "唯一标识",
      "name": "新角色名称",
      "type": "protagonist|antagonist|supporting|minor",
      "description": "角色简短描述",
      "traits": ["性格特点"],
      "group": "所属阵营"
    }
  ],
  "new_edges": [
    {
      "source": "角色A的id",
      "target": "角色B的id",
      "relation": "关系类型",
      "description": "关系描述",
      "weight": 5
    }
  ],
  "updated_edges": [
    {
      "source": "角色A的id",
      "target": "角色B的id",
      "relation": "更新后的关系",
      "description": "关系变化描述",
      "weight": 7
    }
  ],
  "chapter_summary": "本章关系变化概述（30字以内）"
}

## 注意
1. 只输出本章**新增或变化**的内容
2. 如果没有新角色，new_nodes 为空数组
3. updated_edges 用于关系发生变化的情况（如从敌人变为盟友）
4. 已有角色如果出现新关系，放在 new_edges 中`, title, chapterNumber, existingGraphJSON, chapterContent)
}
