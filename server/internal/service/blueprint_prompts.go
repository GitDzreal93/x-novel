package service

import (
	"fmt"
	"strings"
)

// BlueprintPromptParams 大纲生成参数
type BlueprintPromptParams struct {
	UserGuidance       string
	CoreSeed           string
	CharacterDynamics  string
	WorldBuilding      string
	PlotArchitecture   string
	ChapterCount       int
}

// BuildBlueprintPrompt 构建章节大纲提示词
func BuildBlueprintPrompt(params BlueprintPromptParams) string {
	// 构建小说架构字符串
	novelArchitecture := fmt.Sprintf(`核心种子：%s

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
		userGuidance = "无"
	}

	return fmt.Sprintf(`基于以下元素：
- 内容指导：%s
- 小说架构：
%s

设计%d章的节奏分布（根据小说类型调整风格）：
1. 章节集群划分：
- 每3-5章构成一个情节单元，包含完整的小高潮
- 单元之间合理安排情感节奏（张弛有度）
- 关键转折章需预留铺垫

2. 每章需明确：
- 章节定位（角色/事件/主题等）
- 核心内容（剧情推进/情感发展/角色成长等）
- 情感基调（符合小说类型的情感色彩）
- 伏笔操作（埋设/强化/回收）
- 情节张力（★☆☆☆☆ 到 ★★★★★）

输出格式示例：
第n章 - [标题]
本章定位：[角色/事件/主题/...]
核心作用：[推进/转折/发展/升温/...]
情感强度：[平缓/渐进/高潮/...]
伏笔操作：埋设(A线索)→强化(B关系)...
情节张力：★☆☆☆☆
本章简述：[一句话概括]

要求：
- 使用精炼语言描述，每章字数控制在100字以内。
- 合理安排节奏，确保整体情感曲线的连贯性。
- 在生成%d章前不要出现结局章节。
- **情节设计需符合小说类型的风格和情感基调**。

仅给出最终文本，不要解释任何内容。`,
		userGuidance, novelArchitecture, params.ChapterCount, params.ChapterCount)
}

// GenerateMockBlueprint 生成模拟章节大纲
func GenerateMockBlueprint(params BlueprintPromptParams) string {
	var blueprint strings.Builder

	// 根据章节数生成大纲
	for i := 1; i <= params.ChapterCount; i++ {
		var title, role, purpose, emotion, foreshadowing, tension, summary string

		// 第一幕（开端）- 约占20%
		if i <= params.ChapterCount/5 {
			switch i {
			case 1:
				title = "平凡的夜"
				role = "日常/事件"
				purpose = "展示日常状态，引出主线开端"
				emotion = "平缓"
				foreshadowing = "埋设(异常代码)"
				tension = "★☆☆☆☆"
				summary = "程序员李明在公司加班到深夜，处理完最后一个bug后准备回家，却发现代码中有一段神秘的异常。"
			case 2:
				title = "意外的发现"
				role = "事件/转折"
				purpose = "核心触发点，改变主角状态"
				emotion = "渐进"
				foreshadowing = "埋设(AI程序)"
				tension = "★★★☆☆"
				summary = "李明好奇心驱使下运行了那段代码，获得了超能力，整个世界在他眼中变成了数据流。"
			case 3:
				title = "初次相遇"
				role = "角色/感情"
				purpose = "引出女主角，建立初步联系"
				emotion = "平缓"
				foreshadowing = "埋设(苏小雨身份)"
				tension = "★☆☆☆☆"
				summary = "李明在咖啡店偶遇苏小雨，两人因一个小插曲产生交集，彼此留下印象。"
			case 4:
				title = "暗流涌动"
				role = "事件"
				purpose = "展示黑暗组织的存在"
				emotion = "渐进"
				foreshadowing = "强化(组织威胁)"
				tension = "★★☆☆☆"
				summary = "公司开始调查代码泄露事件，李明察觉到有人在使用他的超能力留下的痕迹。"
			case 5:
				title = "被迫卷入"
				role = "事件/转折"
				purpose = "主角被迫采取行动"
				emotion = "渐进"
				foreshadowing = "埋设(追杀开始)"
				tension = "★★★☆☆"
				summary = "李明被公司秘密部门盯上，第一次被迫使用超能力逃脱，意识到危险正在逼近。"
			default:
				title = fmt.Sprintf("第%d章 - 未命名", i)
				role = "事件"
				purpose = "剧情推进"
				emotion = "平缓"
				foreshadowing = ""
				tension = "★☆☆☆☆"
				summary = "剧情持续发展中..."
			}
		} else if i <= params.ChapterCount*2/5 {
			// 第二幕前段（发展）
			title = fmt.Sprintf("第%d章 - 追逃之路", i)
			role = "事件/角色"
			purpose = "主线与感情线交织发展"
			emotion = "渐进"
			foreshadowing = "强化(主线索)"
			tension = "★★★☆☆"
			summary = "李明开始逃亡生活，苏小雨主动接近，两人关系在危机中逐步升温。"
		} else if i <= params.ChapterCount*3/5 {
			// 第二幕中段（挑战）
			title = fmt.Sprintf("第%d章 - 能力觉醒", i)
			role = "角色/成长"
			purpose = "角色成长，能力提升"
			emotion = "高潮"
			foreshadowing = "埋设(最终对决)"
			tension = "★★★★☆"
			summary = "李明逐渐掌握超能力，苏小雨揭露身份，两人建立信任，共同对抗黑暗组织。"
		} else if i <= params.ChapterCount*4/5 {
			// 第二幕后段（转折）
			title = fmt.Sprintf("第%d章 - 信任危机", i)
			role = "转折/感情"
			purpose = "重要转折，关系发展"
			emotion = "高潮"
			foreshadowing = "强化(情感线)"
			tension = "★★★★★"
			summary = "黑暗组织启动重大计划，李明面临生死抉择，与苏小雨的关系也迎来考验。"
		} else {
			// 第三幕（高潮与结局）
			switch i {
			case params.ChapterCount - 2:
				title = "真相大白"
				role = "事件/高潮"
				purpose = "核心冲突爆发"
				emotion = "高潮"
				foreshadowing = "回收(所有伏笔)"
				tension = "★★★★★"
				summary = "黑暗组织的「数字永生」计划启动，李明发现真相，必须做出最终选择。"
			case params.ChapterCount - 1:
				title = "终极对决"
				role = "事件/高潮"
				purpose = "最高潮部分"
				emotion = "高潮"
				foreshadowing = ""
				tension = "★★★★★"
				summary = "李明与黑暗组织首领展开最终对决，苏小雨身受重伤，李明决定牺牲自己的超能力拯救世界。"
			case params.ChapterCount:
				title = "新生"
				role = "结局"
				purpose = "结局收尾"
				emotion = "平缓"
				foreshadowing = ""
				tension = "★☆☆☆☆"
				summary = "黑暗组织被瓦解，李明失去超能力回归普通生活，但内心已成长，与苏小雨开始了新的生活。"
			default:
				title = fmt.Sprintf("第%d章 - 终章前夕", i)
				role = "事件"
				purpose = "为结局做铺垫"
				emotion = "渐进"
				foreshadowing = "强化(主线索)"
				tension = "★★★★☆"
				summary = "最终对决前的准备，各方势力集结，大战一触即发。"
			}
		}

		blueprint.WriteString(fmt.Sprintf("第%d章 - %s\n", i, title))
		blueprint.WriteString(fmt.Sprintf("本章定位：%s\n", role))
		blueprint.WriteString(fmt.Sprintf("核心作用：%s\n", purpose))
		blueprint.WriteString(fmt.Sprintf("情感强度：%s\n", emotion))
		blueprint.WriteString(fmt.Sprintf("伏笔操作：%s\n", foreshadowing))
		blueprint.WriteString(fmt.Sprintf("情节张力：%s\n", tension))
		blueprint.WriteString(fmt.Sprintf("本章简述：%s\n\n", summary))
	}

	return blueprint.String()
}
