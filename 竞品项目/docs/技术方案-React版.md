# 🛠️ 技术方案 - React + Ant Design

## 技术栈选型

### 核心框架
```javascript
{
  "框架": "React 18",
  "语言": "TypeScript",
  "构建工具": "Vite",
  "UI 组件库": "Ant Design 5.x",
  "状态管理": "Zustand + React Query",
  "路由": "React Router v6",
  "HTTP 客户端": "Axios + Fetch",
  "富文本编辑器": "Tiptap / Slate",
  "图表可视化": "ECharts / Recharts",
  "关系图谱": "ECharts-Graph / D3.js",
  "Markdown": "react-markdown",
  "样式方案": "CSS Modules + Tailwind CSS",
  "表单处理": "React Hook Form + Zod",
  "国际化": "i18next"
}
```

---

## 📁 项目目录结构

```
x-novel/
├── src/
│   ├── App.tsx                      # 根组件
│   ├── main.tsx                     # 入口文件
│   ├── vite-env.d.ts
│   │
│   ├── pages/                       # 页面组件
│   │   ├── home/                    # 首页
│   │   │   ├── HomePage.tsx
│   │   │   ├── ProjectList.tsx
│   │   │   ├── ProjectCard.tsx
│   │   │   └── CreateProjectModal.tsx
│   │   │
│   │   ├── project/                 # 项目详情页
│   │   │   ├── ProjectPage.tsx
│   │   │   ├── ProjectHeader.tsx
│   │   │   └── tabs/
│   │   │       ├── ArchitectureTab.tsx
│   │   │       ├── BlueprintTab.tsx
│   │   │       ├── WritingTab.tsx
│   │   │       ├── CompassTab.tsx
│   │   │       ├── ReviewTab.tsx
│   │   │       └── ExportTab.tsx
│   │   │
│   │   ├── inspiration/             # 灵感激发页（新增）
│   │   │   ├── InspirationPage.tsx
│   │   │   ├── ChatInterface.tsx
│   │   │   ├── CreativeCards.tsx
│   │   │   └── TrendAnalysis.tsx
│   │   │
│   │   └── settings/                # 设置页
│   │       ├── SettingsPage.tsx
│   │       ├── APIConfig.tsx
│   │       └── ModelConfig.tsx
│   │
│   ├── components/                  # 业务组件
│   │   ├── layout/
│   │   │   ├── AppLayout.tsx
│   │   │   ├── AppHeader.tsx
│   │   │   └── AppSider.tsx
│   │   │
│   │   ├── writing/                 # 写作相关组件
│   │   │   ├── WritingEditor.tsx    # 写作编辑器
│   │   │   ├── ChapterList.tsx      # 章节列表
│   │   │   ├── WritingAssistant.tsx # 写作助手（右侧边栏）
│   │   │   ├── PolishPanel.tsx      # 润色面板
│   │   │   ├── ExpandPanel.tsx      # 扩写面板
│   │   │   └── ErrorDetection.tsx   # 错误检测
│   │   │
│   │   ├── architecture/            # 架构相关组件
│   │   │   ├── ArchitectureCollapse.tsx
│   │   │   ├── CoreSeedEditor.tsx
│   │   │   ├── CharacterEditor.tsx
│   │   │   └── WorldEditor.tsx
│   │   │
│   │   ├── compass/                 # 关系图谱组件
│   │   │   ├── CompassGraph.tsx
│   │   │   ├── CompassTimeline.tsx
│   │   │   ├── NodeDetailPanel.tsx
│   │   │   └── RelationPopover.tsx
│   │   │
│   │   ├── review/                  # 审阅相关组件（新增）
│   │   │   ├── ExpertReview.tsx     # AI 专家审阅
│   │   │   ├── QualityRadar.tsx     # 质量雷达图
│   │   │   ├── PopularityPredict.tsx # 流行度预测
│   │   │   ├── PublishModal.tsx     # 发布多平台
│   │   │   └── ScoreCard.tsx        # 评分卡片
│   │   │
│   │   └── promotion/               # 推广相关组件（新增）
│   │       ├── SliceGenerator.tsx   # 切片生成器
│   │       ├── SocialGenerator.tsx  # 社交内容生成
│   │       ├── XiaohongshuCard.tsx  # 小红书卡片
│   │       ├── VideoScript.tsx      # 视频脚本
│   │       └── PromoAnalytics.tsx   # 推广数据分析
│   │
│   ├── hooks/                       # 自定义 Hooks
│   │   ├── useNovel.ts              # 小说项目管理
│   │   ├── useSettings.ts           # 设置管理
│   │   ├── useChat.ts               # 聊天对话
│   │   ├── useAI.ts                 # AI 调用
│   │   ├── useEditor.ts             # 编辑器相关
│   │   ├── useErrorDetection.ts     # 错误检测
│   │   └── useLocalStorage.ts       # LocalStorage
│   │
│   ├── stores/                      # 状态管理
│   │   ├── novelStore.ts            # 小说状态
│   │   ├── settingsStore.ts         # 设置状态
│   │   ├── chatStore.ts             # 聊天状态
│   │   └── editorStore.ts           # 编辑器状态
│   │
│   ├── services/                    # API 服务
│   │   ├── api/
│   │   │   ├── novel.ts             # 小说相关 API
│   │   │   ├── ai.ts                # AI 相关 API
│   │   │   ├── chat.ts              # 聊天 API
│   │   │   ├── compass.ts           # 图谱 API
│   │   │   ├── review.ts            # 审阅 API
│   │   │   ├── publish.ts           # 发布 API
│   │   │   └── promotion.ts         # 推广 API
│   │   └── llm/
│   │       ├── openai.ts            # OpenAI 兼容接口
│   │       ├── stream.ts            # 流式处理
│   │       └── config.ts            # 配置管理
│   │
│   ├── prompts/                     # AI 提示词
│   │   ├── architecture.ts
│   │   ├── chapter.ts
│   │   ├── compass.ts
│   │   ├── chat.ts                  # 新增：聊天提示词
│   │   ├── polish.ts                # 新增：润色提示词
│   │   ├── review.ts                # 新增：审阅提示词
│   │   └── promotion.ts             # 新增：推广提示词
│   │
│   ├── types/                       # TypeScript 类型
│   │   ├── novel.ts                 # 小说相关类型
│   │   ├── chat.ts                  # 聊天相关类型
│   │   ├── compass.ts               # 图谱相关类型
│   │   └── index.ts
│   │
│   ├── utils/                       # 工具函数
│   │   ├── storage.ts               # 存储工具
│   │   ├── format.ts                # 格式化工具
│   │   ├── parse.ts                 # 解析工具
│   │   ├── graph-helpers.ts         # 图谱工具
│   │   ├── text-analyze.ts          # 新增：文本分析
│   │   └── export.ts                # 导出工具
│   │
│   ├── constants/                   # 常量
│   │   ├── genres.ts                # 小说类型
│   │   ├── platforms.ts             # 新增：发布平台
│   │   └── errors.ts                # 错误码
│   │
│   └── assets/                      # 静态资源
│       ├── styles/
│       │   ├── global.css
│       │   ├── variables.css        # CSS 变量
│       │   └── theme.ts             # Ant Design 主题配置
│       └── images/
│
├── public/                          # 公共资源
├── index.html
├── package.json
├── tsconfig.json
├── tsconfig.node.json
├── vite.config.ts
├── tailwind.config.js
└── README.md
```

---

## 🔧 核心技术实现

### 1. 状态管理方案

#### 使用 Zustand + React Query

```typescript
// src/stores/novelStore.ts
import { create } from 'zustand'
import { persist } from 'zustand/middleware'

interface NovelProject {
  id: string
  title: string
  topic: string
  genre: string[]
  numberOfChapters: number
  wordNumber: number
  userGuidance: string

  // 架构数据
  coreSeed: string
  characterDynamics: string
  worldBuilding: string
  plotArchitecture: string
  characterState: string

  // 大纲数据
  chapterBlueprint: string

  // 章节内容
  chapters: Record<number, string>

  // 上下文数据
  globalSummary: string

  // 关系图谱
  graphData: GraphData

  // 状态标记
  architectureGenerated: boolean
  blueprintGenerated: boolean

  // 元数据
  createdAt: string
  updatedAt: string
}

interface NovelStore {
  projects: NovelProject[]
  currentProject: NovelProject | null

  // Actions
  createProject: (data: Partial<NovelProject>) => NovelProject
  updateProject: (id: string, updates: Partial<NovelProject>) => void
  deleteProject: (id: string) => void
  setCurrentProject: (id: string) => void
  getCurrentProject: () => NovelProject | null
}

export const useNovelStore = create<NovelStore>()(
  persist(
    (set, get) => ({
      projects: [],
      currentProject: null,

      createProject: (data) => {
        const newProject: NovelProject = {
          id: Date.now().toString(),
          createdAt: new Date().toISOString(),
          updatedAt: new Date().toISOString(),
          coreSeed: '',
          characterDynamics: '',
          worldBuilding: '',
          plotArchitecture: '',
          characterState: '',
          chapterBlueprint: '',
          chapters: {},
          globalSummary: '',
          graphData: {
            version: 1,
            generatedAt: null,
            snapshots: {},
            audit: { inconsistencies: [], lastAuditAt: null },
            graphGenerated: false
          },
          architectureGenerated: false,
          blueprintGenerated: false,
          ...data
        }

        set((state) => ({
          projects: [newProject, ...state.projects]
        }))

        return newProject
      },

      updateProject: (id, updates) => {
        set((state) => ({
          projects: state.projects.map(p =>
            p.id === id
              ? { ...p, ...updates, updatedAt: new Date().toISOString() }
              : p
          ),
          currentProject: state.currentProject?.id === id
            ? { ...state.currentProject, ...updates, updatedAt: new Date().toISOString() }
            : state.currentProject
        }))
      },

      deleteProject: (id) => {
        set((state) => ({
          projects: state.projects.filter(p => p.id !== id),
          currentProject: state.currentProject?.id === id ? null : state.currentProject
        }))
      },

      setCurrentProject: (id) => {
        const project = get().projects.find(p => p.id === id)
        set({ currentProject: project || null })
      },

      getCurrentProject: () => {
        return get().currentProject
      }
    }),
    {
      name: 'novel-storage',
      partialize: (state) => ({
        projects: state.projects
      })
    }
  )
)
```

```typescript
// src/stores/chatStore.ts - 新增聊天状态
import { create } from 'zustand'

interface ChatMessage {
  id: string
  role: 'user' | 'assistant' | 'system'
  content: string
  timestamp: number
}

interface ChatStore {
  currentChat: ChatMessage[]
  chatHistory: Record<string, ChatMessage[]>
  contextMode: 'creative' | 'building' | 'character' | 'general'

  // Actions
  addMessage: (message: Omit<ChatMessage, 'id' | 'timestamp'>) => void
  clearCurrentChat: () => void
  saveChat: (projectId: string) => void
  loadChat: (projectId: string) => void
  setContextMode: (mode: ChatStore['contextMode']) => void
}

export const useChatStore = create<ChatStore>((set, get) => ({
  currentChat: [{
    id: 'welcome',
    role: 'assistant',
    content: '你好！我是你的 AI 创作助手。你可以告诉我你想写什么类型的小说，或者有什么创意点子，我来帮你完善。',
    timestamp: Date.now()
  }],
  chatHistory: {},
  contextMode: 'general',

  addMessage: (message) => {
    const newMessage: ChatMessage = {
      ...message,
      id: Date.now().toString(),
      timestamp: Date.now()
    }

    set((state) => ({
      currentChat: [...state.currentChat, newMessage]
    }))
  },

  clearCurrentChat: () => {
    set({ currentChat: [] })
  },

  saveChat: (projectId) => {
    const { currentChat, chatHistory } = get()
    set({
      chatHistory: {
        ...chatHistory,
        [projectId]: currentChat
      }
    })
  },

  loadChat: (projectId) => {
    const { chatHistory } = get()
    set({
      currentChat: chatHistory[projectId] || []
    })
  },

  setContextMode: (mode) => {
    set({ contextMode: mode })
  }
}))
```

#### 使用 React Query 管理服务端状态

```typescript
// src/services/api/chat.ts
import { useMutation, useQuery } from '@tanstack/react-query'
import { chatAPI } from '../llm/openai'

export const useChatMutation = () => {
  return useMutation({
    mutationFn: async (messages: ChatMessage[]) => {
      return await chatAPI.completions(messages)
    }
  })
}

export const usePolishMutation = () => {
  return useMutation({
    mutationFn: async (params: { text: string; option: string }) => {
      return await chatAPI.polish(params.text, params.option)
    }
  })
}

export const useExpandMutation = () => {
  return useMutation({
    mutationFn: async (params: { text: string; length: number; direction: string }) => {
      return await chatAPI.expand(params.text, params.length, params.direction)
    }
  })
}
```

---

### 2. AI 调用层

```typescript
// src/services/llm/openai.ts
import axios from 'axios'
import { getSettings } from '@/stores/settingsStore'

export class LLMService {
  private baseURL: string
  private apiKey: string
  private model: string
  private timeout: number

  constructor() {
    const settings = getSettings()
    this.baseURL = settings.apiConfig.baseUrl
    this.apiKey = settings.apiConfig.apiKey
    this.model = settings.apiConfig.model
    this.timeout = settings.apiConfig.timeout * 1000
  }

  // 非流式补全
  async chatCompletion(messages: ChatMessage[], options?: ChatOptions): Promise<string> {
    try {
      const response = await axios.post(
        `${this.baseURL}/chat/completions`,
        {
          model: this.model,
          messages,
          temperature: options?.temperature ?? 0.7,
          max_tokens: options?.maxTokens ?? 8192,
          stream: false
        },
        {
          headers: {
            'Authorization': `Bearer ${this.apiKey}`,
            'Content-Type': 'application/json'
          },
          timeout: this.timeout
        }
      )

      return response.data.choices[0].message.content
    } catch (error) {
      throw new Error(`AI 调用失败: ${error.message}`)
    }
  }

  // 流式补全
  async streamCompletion(
    messages: ChatMessage[],
    onChunk: (chunk: string, full: string) => void,
    options?: ChatOptions
  ): Promise<string> {
    const response = await fetch(`${this.baseURL}/chat/completions`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${this.apiKey}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        model: this.model,
        messages,
        temperature: options?.temperature ?? 0.7,
        max_tokens: options?.maxTokens ?? 8192,
        stream: true
      })
    })

    if (!response.ok) {
      throw new Error(`API 请求失败: ${response.status}`)
    }

    const reader = response.body!.getReader()
    const decoder = new TextDecoder()
    let fullContent = ''

    while (true) {
      const { done, value } = await reader.read()
      if (done) break

      const chunk = decoder.decode(value)
      const lines = chunk.split('\n').filter(line => line.trim() !== '')

      for (const line of lines) {
        if (line.startsWith('data: ')) {
          const data = line.slice(6)
          if (data === '[DONE]') continue

          try {
            const parsed = JSON.parse(data)
            const content = parsed.choices?.[0]?.delta?.content || ''
            if (content) {
              fullContent += content
              onChunk(content, fullContent)
            }
          } catch (e) {
            // 跳过无效 JSON
          }
        }
      }
    }

    return fullContent
  }
}

export const llmService = new LLMService()
```

---

### 3. 自定义 Hooks

```typescript
// src/hooks/useChat.ts
import { useChatStore } from '@/stores/chatStore'
import { useChatMutation } from '@/services/api/chat'
import { useChatPrompts } from '@/prompts/chat'

export const useChat = () => {
  const { currentChat, addMessage, contextMode } = useChatStore()
  const chatMutation = useChatMutation()
  const { getSystemPrompt } = useChatPrompts()

  const sendMessage = async (userMessage: string) => {
    // 添加用户消息
    addMessage({
      role: 'user',
      content: userMessage
    })

    // 添加临时助手消息
    const tempId = Date.now().toString()
    addMessage({
      role: 'assistant',
      content: '思考中...'
    })

    try {
      // 获取系统提示词
      const systemPrompt = getSystemPrompt(contextMode)

      // 构建消息历史
      const messages = [
        { role: 'system', content: systemPrompt },
        ...currentChat.slice(-10).map(m => ({
          role: m.role,
          content: m.content
        }))
      ]

      // 调用 AI
      const response = await chatMutation.mutateAsync(messages)

      // 更新助手消息
      // 实现中...
    } catch (error) {
      addMessage({
        role: 'assistant',
        content: `抱歉，出现了错误：${error.message}`
      })
    }
  }

  return {
    currentChat,
    sendMessage,
    isLoading: chatMutation.isPending
  }
}
```

```typescript
// src/hooks/useWritingAssistant.ts - 新增
import { useState } from 'react'
import { usePolishMutation, useExpandMutation, useSuggestMutation } from '@/services/api/writing'
import { usePolishPrompts, useExpandPrompts, useSuggestPrompts } from '@/prompts/writing'

export const useWritingAssistant = () => {
  const [selectedText, setSelectedText] = useState('')
  const [polishResult, setPolishResult] = useState('')
  const [expandResult, setExpandResult] = useState('')
  const [suggestions, setSuggestions] = useState<string[]>([])

  const polishMutation = usePolishMutation()
  const expandMutation = useExpandMutation()
  const suggestMutation = useSuggestMutation()

  const polish = async (text: string, option: string) => {
    setSelectedText(text)

    const prompt = usePolishPrompts(text, option)
    const result = await polishMutation.mutateAsync({ text, option })
    setPolishResult(result)
  }

  const expand = async (text: string, length: number, direction: string) => {
    setSelectedText(text)

    const result = await expandMutation.mutateAsync({ text, length, direction })
    setExpandResult(result)
  }

  const getSuggestions = async (context: string) => {
    const prompt = useSuggestPrompts(context)
    const result = await suggestMutation.mutateAsync({ context })
    setSuggestions(result)
  }

  return {
    selectedText,
    polishResult,
    expandResult,
    suggestions,
    polish,
    expand,
    getSuggestions,
    isPolishing: polishMutation.isPending,
    isExpanding: expandMutation.isPending,
    isSuggesting: suggestMutation.isPending
  }
}
```

```typescript
// src/hooks/useErrorDetection.ts - 新增
import { useMemo } from 'react'
import { detectTypos, detectGrammarIssues, detectRepetition } from '@/utils/text-analyze'

export const useErrorDetection = (text: string) => {
  const errors = useMemo(() => {
    const typos = detectTypos(text)
    const grammar = detectGrammarIssues(text)
    const repetition = detectRepetition(text)

    return {
      typos,      // 错别字
      grammar,    // 病句
      repetition  // 重复问题
    }
  }, [text])

  const errorCount = errors.typos.length + errors.grammar.length + errors.repetition.length

  return {
    errors,
    errorCount,
    hasErrors: errorCount > 0
  }
}
```

---

### 4. 富文本编辑器

使用 **Tiptap** 作为编辑器（更适合写作场景）

```typescript
// src/components/writing/WritingEditor.tsx
import { useEditor, EditorContent } from '@tiptap/react'
import StarterKit from '@tiptap/starter-kit'
import Placeholder from '@tiptap/extension-placeholder'
import CharacterCount from '@tiptap/extension-character-count'

interface WritingEditorProps {
  content: string
  onChange: (content: string) => void
  placeholder?: string
  readOnly?: boolean
}

export const WritingEditor: React.FC<WritingEditorProps> = ({
  content,
  onChange,
  placeholder = '开始写作...',
  readOnly = false
}) => {
  const editor = useEditor({
    extensions: [
      StarterKit,
      Placeholder.configure({
        placeholder
      }),
      CharacterCount
    ],
    content,
    onUpdate: ({ editor }) => {
      onChange(editor.getHTML())
    },
    editable: !readOnly
  })

  if (!editor) {
    return null
  }

  return (
    <div className="writing-editor">
      <EditorContent editor={editor} />

      {/* 工具栏 */}
      {!readOnly && (
        <div className="editor-toolbar">
          <button onClick={() => editor.chain().focus().toggleBold().run()}>
            Bold
          </button>
          <button onClick={() => editor.chain().focus().toggleItalic().run()}>
            Italic
          </button>
          {/* 更多格式按钮... */}
        </div>
      )}

      {/* 字数统计 */}
      <div className="editor-footer">
        {editor.storage.characterCount.characters()} 字符
      </div>
    </div>
  )
}
```

---

### 5. 写作助手组件（右侧边栏）

```typescript
// src/components/writing/WritingAssistant.tsx
import { Drawer, Tabs, Button, Space, Select, Slider } from 'antd'
import { useWritingAssistant } from '@/hooks/useWritingAssistant'
import { PolishPanel } from './PolishPanel'
import { ExpandPanel } from './ExpandPanel'
import { SuggestionPanel } from './SuggestionPanel'

interface WritingAssistantProps {
  visible: boolean
  onClose: () => void
  selectedText: string
  contextText: string
  onApplyPolish: (text: string) => void
  onApplyExpand: (text: string) => void
}

export const WritingAssistant: React.FC<WritingAssistantProps> = ({
  visible,
  onClose,
  selectedText,
  contextText,
  onApplyPolish,
  onApplyExpand
}) => {
  const {
    polishResult,
    expandResult,
    suggestions,
    polish,
    expand,
    getSuggestions,
    isPolishing,
    isExpanding,
    isSuggesting
  } = useWritingAssistant()

  const [polishOption, setPolishOption] = useState('vivid')
  const [expandLength, setExpandLength] = useState(500)
  const [expandDirection, setExpandDirection] = useState('comprehensive')

  return (
    <Drawer
      title="✨ 写作助手"
      placement="right"
      width={400}
      open={visible}
      onClose={onClose}
    >
      <Tabs
        defaultActiveKey="polish"
        items={[
          {
            key: 'polish',
            label: '润色',
            children: (
              <PolishPanel
                selectedText={selectedText}
                result={polishResult}
                option={polishOption}
                onOptionChange={setPolishOption}
                onPolish={() => polish(selectedText, polishOption)}
                onApply={onApplyPolish}
                loading={isPolishing}
              />
            )
          },
          {
            key: 'expand',
            label: '扩写',
            children: (
              <ExpandPanel
                selectedText={selectedText}
                result={expandResult}
                length={expandLength}
                direction={expandDirection}
                onLengthChange={setExpandLength}
                onDirectionChange={setExpandDirection}
                onExpand={() => expand(selectedText, expandLength, expandDirection)}
                onApply={onApplyExpand}
                loading={isExpanding}
              />
            )
          },
          {
            key: 'suggest',
            label: '灵感',
            children: (
              <SuggestionPanel
                contextText={contextText}
                suggestions={suggestions}
                onGetSuggestions={() => getSuggestions(contextText)}
                loading={isSuggesting}
              />
            )
          }
        ]}
      />
    </Drawer>
  )
}
```

---

### 6. 聊天界面组件

```typescript
// src/pages/inspiration/ChatInterface.tsx
import { useState, useRef, useEffect } from 'react'
import { Card, Input, Button, Space, Select, Tag } from 'antd'
import { SendOutlined, RobotOutlined, UserOutlined } from '@ant-design/icons'
import { useChat } from '@/hooks/useChat'
import ReactMarkdown from 'react-markdown'

export const ChatInterface: React.FC = () => {
  const { currentChat, sendMessage, isLoading } = useChat()
  const [input, setInput] = useState('')
  const [contextMode, setContextMode] = useState<'creative' | 'building' | 'character' | 'general'>('general')
  const messagesEndRef = useRef<HTMLDivElement>(null)

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' })
  }

  useEffect(() => {
    scrollToBottom()
  }, [currentChat])

  const handleSend = () => {
    if (!input.trim()) return

    sendMessage(input)
    setInput('')
  }

  return (
    <div className="chat-interface">
      {/* 模式选择 */}
      <Space className="chat-mode-selector">
        <span>对话模式：</span>
        <Select
          value={contextMode}
          onChange={setContextMode}
          style={{ width: 120 }}
          options={[
            { label: '创意启发', value: 'creative' },
            { label: '设定完善', value: 'building' },
            { label: '角色塑造', value: 'character' },
            { label: '通用', value: 'general' }
          ]}
        />
      </Space>

      {/* 消息列表 */}
      <div className="chat-messages">
        {currentChat.map((message) => (
          <div
            key={message.id}
            className={`chat-message ${message.role}`}
          >
            {message.role === 'assistant' && (
              <RobotOutlined className="message-avatar" />
            )}
            <Card className="message-content">
              {message.role === 'assistant' ? (
                <ReactMarkdown>{message.content}</ReactMarkdown>
              ) : (
                message.content
              )}
            </Card>
            {message.role === 'user' && (
              <UserOutlined className="message-avatar" />
            )}
          </div>
        ))}
        {isLoading && (
          <div className="chat-message assistant">
            <RobotOutlined className="message-avatar" />
            <Card className="message-content">
              <span className="typing-indicator">思考中...</span>
            </Card>
          </div>
        )}
        <div ref={messagesEndRef} />
      </div>

      {/* 输入框 */}
      <div className="chat-input">
        <Space.Compact style={{ width: '100%' }}>
          <Input.TextArea
            value={input}
            onChange={(e) => setInput(e.target.value)}
            onPressEnter={(e) => {
              if (e.shiftKey) return
              e.preventDefault()
              handleSend()
            }}
            placeholder="输入你的问题...（Shift + Enter 换行）"
            autoSize={{ minRows: 2, maxRows: 6 }}
          />
          <Button
            type="primary"
            icon={<SendOutlined />}
            onClick={handleSend}
            loading={isLoading}
          >
            发送
          </Button>
        </Space.Compact>
      </div>
    </div>
  )
}
```

---

### 7. Ant Design 主题配置

```typescript
// src/assets/styles/theme.ts
import { ConfigTheme, theme } from 'antd'

export const darkTheme: ConfigTheme = {
  algorithm: theme.darkAlgorithm,
  token: {
    colorPrimary: '#6366f1',      // 靛蓝色
    colorSuccess: '#10b981',
    colorWarning: '#f59e0b',
    colorError: '#ef4444',
    colorInfo: '#3b82f6',
    borderRadius: 8,
    fontSize: 14,
    fontFamily: '-apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif'
  },
  components: {
    Layout: {
      headerBg: '#1f1f23',
      siderBg: '#1f1f23'
    },
    Input: {
      colorBgContainer: '#2a2a2e',
      colorBorder: '#3f3f46'
    },
    Card: {
      colorBgContainer: '#1f1f23',
      colorBorderSecondary: '#3f3f46'
    }
  }
}

export const lightTheme: ConfigTheme = {
  algorithm: theme.defaultAlgorithm,
  token: {
    colorPrimary: '#6366f1',
    borderRadius: 8,
    fontSize: 14
  }
}
```

---

### 8. 路由配置

```typescript
// src/App.tsx
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom'
import { ConfigProvider, App as AntdApp } from 'antd'
import { useSettingsStore } from './stores/settingsStore'
import { lightTheme, darkTheme } from './assets/styles/theme'
import AppLayout from './components/layout/AppLayout'

// Pages
import HomePage from './pages/home/HomePage'
import ProjectPage from './pages/project/ProjectPage'
import InspirationPage from './pages/inspiration/InspirationPage'
import SettingsPage from './pages/settings/SettingsPage'

const App: React.FC = () => {
  const { isDark } = useSettingsStore()

  return (
    <ConfigProvider theme={isDark ? darkTheme : lightTheme}>
      <AntdApp>
        <BrowserRouter>
          <Routes>
            <Route path="/" element={<AppLayout />}>
              <Route index element={<HomePage />} />
              <Route path="project/:id" element={<ProjectPage />} />
              <Route path="inspiration" element={<InspirationPage />} />
              <Route path="settings" element={<SettingsPage />} />
            </Route>
            <Route path="*" element={<Navigate to="/" replace />} />
          </Routes>
        </BrowserRouter>
      </AntdApp>
    </ConfigProvider>
  )
}

export default App
```

---

## 📦 核心依赖安装

```bash
# 核心框架
npm install react react-dom
npm install -D @types/react @types/react-dom

# 构建工具
npm install -D vite @vitejs/plugin-react typescript

# UI 组件
npm install antd
npm install @ant-design/icons

# 状态管理
npm install zustand
npm install @tanstack/react-query

# 路由
npm install react-router-dom

# HTTP
npm install axios

# 编辑器
npm install @tiptap/react @tiptap/starter-kit @tiptap/extension-placeholder @tiptap/extension-character-count

# Markdown
npm install react-markdown

# 图表
npm install echarts echarts-for-react
# 或
npm install recharts

# 表单
npm install react-hook-form zod

# 工具
npm install dayjs
npm install lodash-es
npm install -D @types/lodash

# 样式
npm install -D tailwindcss postcss autoprefixer
npm install -D sass

# 持久化
npm install zustand.persist

# 国际化（可选）
npm install react-i18next i18next
```

---

## 🎯 开发优先级建议

### Phase 1: 基础框架（第 1-2 周）
- [x] 项目初始化（Vite + React + TS）
- [x] Ant Design 配置
- [x] 路由配置
- [x] 状态管理（Zustand）
- [x] 基础布局组件

### Phase 2: 竞品功能复刻（第 3-6 周）
- [x] 项目管理（CRUD）
- [x] 小说架构生成
- [x] 章节大纲生成
- [x] 章节写作面板
- [x] 关系图谱（简化版）

### Phase 3: 核心超越功能（第 7-10 周）
- [x] 智能对话激发灵感
- [x] 写作助手（润色/扩写/灵感）
- [x] 错误检测系统
- [x] AI 专家审阅
- [x] 流行度预测

### Phase 4: 推广功能（第 11-12 周）
- [x] 小说切片生成
- [x] 社交媒体内容生成
- [x] 一键发布多平台

---

## 📝 关键技术点说明

### 1. 为什么选择 Zustand 而不是 Redux？
- 更简洁的 API
- 更小的包体积
- 内置 TypeScript 支持
- 不需要 boilerplate 代码
- 支持持久化中间件

### 2. 为什么选择 Tiptap 而不是 Slate？
- 更好的文档
- 更容易扩展
- 内置常用扩展
- 更好的性能
- ProseMirror 底层，功能强大

### 3. 为什么选择 React Query？
- 自动缓存和重新验证
- 乐观更新支持
- 更好的开发者体验
- 减少 boilerplate

### 4. CSS 方案
- **Tailwind CSS**: 快速开发，原子化 CSS
- **CSS Modules**: 组件隔离，避免冲突
- **Ant Design 主题**: 统一设计语言

---

**文档生成时间**: 2026-02-27
**技术栈**: React 18 + Ant Design 5 + Vite + TypeScript
