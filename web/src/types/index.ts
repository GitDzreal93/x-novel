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
  purpose: 'architecture' | 'chapter' | 'writing' | 'review' | 'general';
  base_url?: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
  provider?: ModelProvider;
}

export interface CreateModelConfigRequest {
  provider_id: number;
  model_name: string;
  purpose: ModelConfig['purpose'];
  api_key: string;
  base_url?: string;
}

export interface UpdateModelConfigRequest {
  model_name?: string;
  purpose?: ModelConfig['purpose'];
  api_key?: string;
  base_url?: string;
  is_active?: boolean;
}

export interface ValidateModelConfigRequest {
  provider_id: number;
  api_key: string;
  base_url?: string;
}
