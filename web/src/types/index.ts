// 通用响应类型
export interface Response<T = any> {
  code: number;
  message: string;
  data?: T;
}

// 设备相关类型
export interface Device {
  id: string;
  device_id: string;
  device_info?: Record<string, any>;
  created_at: string;
  last_seen: string;
}

export interface DeviceSettings {
  id: string;
  device_id: string;
  theme: 'light' | 'dark';
  language: string;
  auto_save_enabled: boolean;
  auto_save_interval: number;
}

// 项目相关类型
export interface Project {
  id: string;
  title: string;
  topic?: string;
  genre?: string[];
  chapter_count: number;
  words_per_chapter: number;
  user_guidance?: string;

  // 架构数据
  core_seed?: string;
  character_dynamics?: string;
  world_building?: string;
  plot_architecture?: string;
  character_state?: string;
  architecture_generated: boolean;

  // 大纲数据
  chapter_blueprint?: string;
  blueprint_generated: boolean;

  // 统计
  global_summary?: string;
  total_chapters: number;
  completed_chapters: number;
  total_words: number;

  // 状态
  status: 'draft' | 'writing' | 'completed' | 'published';
  created_at: string;
  updated_at: string;
}

export interface CreateProjectRequest {
  title: string;
  topic?: string;
  genre?: string[];
  chapter_count?: number;
  words_per_chapter?: number;
  user_guidance?: string;
}

export interface UpdateProjectRequest {
  title?: string;
  topic?: string;
  genre?: string[];
  chapter_count?: number;
  words_per_chapter?: number;
  user_guidance?: string;
  status?: Project['status'];

  // 架构数据
  core_seed?: string;
  character_dynamics?: string;
  world_building?: string;
  plot_architecture?: string;
  character_state?: string;

  // 大纲数据
  chapter_blueprint?: string;

  // 全局摘要
  global_summary?: string;
}

// 章节相关类型
export interface Chapter {
  id: string;
  project_id: string;
  chapter_number: number;
  title?: string;

  // 大纲信息
  blueprint_position?: string;
  blueprint_purpose?: string;
  blueprint_suspense?: string;
  blueprint_foreshadowing?: string;
  blueprint_twist_level?: string;
  blueprint_summary?: string;

  // 内容
  content?: string;
  word_count: number;

  // 状态
  status: 'not_started' | 'draft' | 'completed';
  is_finalized: boolean;

  created_at: string;
  updated_at: string;
}

export interface CreateChapterRequest {
  chapter_number: number;
  title?: string;
  blueprint_summary?: string;
}

export interface UpdateChapterRequest {
  title?: string;
  content?: string;
  status?: Chapter['status'];
  is_finalized?: boolean;

  // 大纲信息
  blueprint_position?: string;
  blueprint_purpose?: string;
  blueprint_suspense?: string;
  blueprint_foreshadowing?: string;
  blueprint_twist_level?: string;
  blueprint_summary?: string;
}

// 列表响应类型
export interface ListResponse<T> {
  items: T[];
  total: number;
}

// 模型提供商类型
export interface ModelProvider {
  id: number;
  name: string;
  display_name: string;
  base_url?: string;
  auth_type?: string;
  is_active: boolean;
}

// 模型配置类型
export interface ModelConfig {
  id: string;
  provider_id: number;
  model_name: string;
  base_url?: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
  provider?: ModelProvider;
}

export interface CreateModelConfigRequest {
  provider_id: number;
  model_name: string;
  api_key: string;
  base_url?: string;
}

export interface UpdateModelConfigRequest {
  model_name?: string;
  api_key?: string;
  base_url?: string;
  is_active?: boolean;
}

export interface ValidateModelConfigRequest {
  provider_id: number;
  api_key: string;
  base_url?: string;
}

// 功能绑定类型
export type BindingPurpose = 'architecture' | 'chapter' | 'writing' | 'review' | 'general';

export interface ModelBinding {
  id: string;
  purpose: BindingPurpose;
  model_config_id: string;
  model_config?: ModelConfig;
  created_at: string;
  updated_at: string;
}

export interface UpsertModelBindingRequest {
  purpose: BindingPurpose;
  model_config_id: string;
}

// 关系图谱相关类型
export interface GraphNode {
  id: string;
  name: string;
  type: 'protagonist' | 'antagonist' | 'supporting' | 'minor';
  description: string;
  traits: string[];
  group: string;
}

export interface GraphEdge {
  source: string;
  target: string;
  relation: string;
  description: string;
  weight: number;
}

export interface GraphSnapshot {
  chapter_number: number;
  summary: string;
  nodes_count: number;
  edges_count: number;
}

export interface GraphData {
  nodes: GraphNode[];
  edges: GraphEdge[];
  snapshots?: GraphSnapshot[];
}

// 写作助手相关类型
export type WritingAction = 'polish' | 'continue' | 'suggestion';
export type PolishStyle = 'vivid' | 'concise' | 'literary' | 'dramatic';
export type SuggestionAspect = 'plot' | 'character' | 'dialogue' | 'description' | 'conflict';

export interface WritingAssistantRequest {
  action: WritingAction;
  content: string;
  project_id?: string;
  style?: PolishStyle;
  target_words?: number;
  aspect?: SuggestionAspect;
  stream?: boolean;
}

// 错误检测相关类型
export type DetectionType = 'typo' | 'grammar' | 'logic' | 'repetition';
export type IssueSeverity = 'error' | 'warning' | 'info';

export interface DetectionIssue {
  type: DetectionType;
  severity: IssueSeverity;
  position: string;
  original: string;
  suggestion: string;
  explanation: string;
}

export interface DetectionResult {
  issues: DetectionIssue[];
  summary: string;
  total_count: number;
  type_counts: Record<string, number>;
}

// AI 审阅相关类型
export type ReviewDimension = 'plot' | 'character' | 'writing' | 'pacing' | 'creativity' | 'readability';

export interface ReviewScore {
  dimension: ReviewDimension;
  score: number;
  comment: string;
}

export interface ReviewResult {
  scores: ReviewScore[];
  overall_score: number;
  highlights: string[];
  issues: string[];
  suggestions: string[];
  summary: string;
}

// 市场预测相关类型
export interface MarketTrendItem {
  trend: string;
  fit: number;
  analysis: string;
}

export interface ReaderAppealItem {
  dimension: string;
  score: number;
  comment: string;
}

export interface MonetizationAdvice {
  platforms: string[];
  pricing_model: string;
  ip_potential: number;
  suggestion: string;
}

export interface MarketPrediction {
  market_score: number;
  target_audience: string;
  competitive_edge: string;
  market_trends: MarketTrendItem[];
  reader_appeal: ReaderAppealItem[];
  monetization: MonetizationAdvice;
  risks: string[];
  recommendations: string[];
  summary: string;
}

// 数据备份相关类型
export interface BackupPreview {
  projects: number;
  chapters: number;
  total_words: number;
  conversations: number;
  messages: number;
}

export interface ImportResult {
  imported_projects: number;
  imported_chapters: number;
  imported_conversations: number;
  imported_messages: number;
  failed_projects: number;
  failed_chapters: number;
  failed_conversations: number;
  failed_messages: number;
}

// 对话相关类型
export type ChatMode = 'creative' | 'building' | 'character' | 'general';

export interface Conversation {
  id: string;
  title: string;
  mode: ChatMode;
  project_id?: string;
  messages?: ChatMessage[];
  created_at: string;
  updated_at: string;
}

export interface ChatMessage {
  id: string;
  role: 'user' | 'assistant';
  content: string;
  created_at: string;
}

export interface CreateConversationRequest {
  title?: string;
  mode?: ChatMode;
  project_id?: string;
}

export interface SendMessageRequest {
  content: string;
  stream?: boolean;
}

export interface SendMessageResponse {
  user_message: ChatMessage;
  assistant_message: ChatMessage;
}
