import request from './request';
import type {
  Response,
  Device,
  DeviceSettings,
  Project,
  CreateProjectRequest,
  UpdateProjectRequest,
  Chapter,
  CreateChapterRequest,
  UpdateChapterRequest,
  ModelConfig,
  ModelProvider,
  CreateModelConfigRequest,
  UpdateModelConfigRequest,
  ValidateModelConfigRequest,
  Conversation,
  CreateConversationRequest,
  SendMessageResponse,
  WritingAssistantRequest,
  GraphData,
  DetectionResult,
  ReviewResult,
  MarketPrediction,
  BackupPreview,
  ImportResult,
} from '../types';

// ========== 设备相关 API ==========

export const deviceApi = {
  // 获取设备信息
  getInfo: () => {
    return request.get<Response<Device>>('/api/v1/device/info');
  },

  // 获取设备设置
  getSettings: () => {
    return request.get<Response<DeviceSettings>>('/api/v1/device/settings');
  },

  // 更新设备设置
  updateSettings: (data: Partial<DeviceSettings>) => {
    return request.put<Response<DeviceSettings>>('/api/v1/device/settings', data);
  },
};

// ========== 项目相关 API ==========

export const projectApi = {
  // 创建项目
  create: (data: CreateProjectRequest) => {
    return request.post<Response<Project>>('/api/v1/projects', data);
  },

  // 获取项目列表
  list: (params: { page?: number; page_size?: number }) => {
    return request.get<Response<{ projects: Project[]; total: number }>>(
      '/api/v1/projects',
      { params }
    );
  },

  // 获取项目详情
  getById: (id: string) => {
    return request.get<Response<Project>>(`/api/v1/projects/${id}`);
  },

  // 更新项目
  update: (id: string, data: UpdateProjectRequest) => {
    return request.put<Response<Project>>(`/api/v1/projects/${id}`, data);
  },

  // 删除项目
  delete: (id: string) => {
    return request.delete<Response<Project>>(`/api/v1/projects/${id}`);
  },

  // 生成小说架构
  generateArchitecture: (id: string, data: { overwrite?: boolean }) => {
    return request.post<Response<Project>>(
      `/api/v1/projects/${id}/architecture/generate`,
      data
    );
  },

  // 生成章节大纲
  generateBlueprint: (id: string, data: { overwrite?: boolean }) => {
    return request.post<Response<Project>>(
      `/api/v1/projects/${id}/blueprint/generate`,
      data
    );
  },

  // 导出项目
  export: (id: string, format: 'txt' | 'md') => {
    return request.get<Response<{ download_url: string }>>(
      `/api/v1/projects/${id}/export/${format}`
    );
  },
};

// ========== 章节相关 API ==========

export const chapterApi = {
  // 创建章节
  create: (projectId: string, data: CreateChapterRequest) => {
    return request.post<Response<Chapter>>(
      `/api/v1/projects/${projectId}/chapters`,
      data
    );
  },

  // 获取章节列表
  list: (projectId: string, params: { page?: number; page_size?: number }) => {
    return request.get<Response<{ chapters: Chapter[]; total: number }>>(
      `/api/v1/projects/${projectId}/chapters`,
      { params }
    );
  },

  // 获取章节详情
  getByNumber: (projectId: string, chapterNumber: number) => {
    return request.get<Response<Chapter>>(
      `/api/v1/projects/${projectId}/chapters/${chapterNumber}`
    );
  },

  // 更新章节
  update: (
    projectId: string,
    chapterNumber: number,
    data: UpdateChapterRequest
  ) => {
    return request.put<Response<Chapter>>(
      `/api/v1/projects/${projectId}/chapters/${chapterNumber}`,
      data
    );
  },

  // 生成章节内容
  generateContent: (
    projectId: string,
    chapterNumber: number,
    data: { overwrite?: boolean }
  ) => {
    return request.post<Response<Chapter>>(
      `/api/v1/projects/${projectId}/chapters/${chapterNumber}/generate`,
      data
    );
  },

  // 定稿章节
  finalize: (projectId: string, chapterNumber: number, data: { update_summary?: boolean }) => {
    return request.post<Response<Chapter>>(
      `/api/v1/projects/${projectId}/chapters/${chapterNumber}/finalize`,
      data
    );
  },

  // 扩写章节
  enrich: (projectId: string, chapterNumber: number, data: { target_words?: number }) => {
    return request.post<Response<Chapter>>(
      `/api/v1/projects/${projectId}/chapters/${chapterNumber}/enrich`,
      data
    );
  },
};

// ========== 模型配置相关 API ==========

export const modelConfigApi = {
  // 获取提供商列表
  listProviders: () => {
    return request.get<Response<ModelProvider[]>>('/api/v1/models/providers');
  },

  // 获取模型配置列表
  list: (params?: { page?: number; page_size?: number }) => {
    return request.get<Response<{ configs: ModelConfig[]; total: number }>>(
      '/api/v1/models',
      { params }
    );
  },

  // 创建模型配置
  create: (data: CreateModelConfigRequest) => {
    return request.post<Response<ModelConfig>>('/api/v1/models', data);
  },

  // 获取模型配置详情
  getById: (id: string) => {
    return request.get<Response<ModelConfig>>(`/api/v1/models/${id}`);
  },

  // 更新模型配置
  update: (id: string, data: UpdateModelConfigRequest) => {
    return request.put<Response<ModelConfig>>(`/api/v1/models/${id}`, data);
  },

  // 删除模型配置
  delete: (id: string) => {
    return request.delete<Response<void>>(`/api/v1/models/${id}`);
  },

  // 验证模型配置
  validate: (data: ValidateModelConfigRequest) => {
    return request.post<Response<void>>('/api/v1/models/validate', data);
  },
};

// ========== 关系图谱 API ==========

export const graphApi = {
  // 获取项目图谱
  get: (projectId: string) => {
    return request.get<Response<GraphData>>(`/api/v1/projects/${projectId}/graph`);
  },

  // 生成图谱（从架构提取）
  generate: (projectId: string) => {
    return request.post<Response<GraphData>>(`/api/v1/projects/${projectId}/graph/generate`);
  },

  // 从章节更新图谱
  updateFromChapter: (projectId: string, chapterNumber: number) => {
    return request.post<Response<GraphData>>(`/api/v1/projects/${projectId}/graph/chapters/${chapterNumber}`);
  },

  // 获取章节快照
  getChapterSnapshot: (projectId: string, chapterNumber: number) => {
    return request.get<Response<GraphData>>(`/api/v1/projects/${projectId}/graph/chapters/${chapterNumber}`);
  },
};

// ========== 数据备份 API ==========

export const backupApi = {
  preview: () => {
    return request.get<Response<BackupPreview>>('/api/v1/backup/preview');
  },

  exportData: () => {
    return request.get<Blob>('/api/v1/backup/export', { responseType: 'blob' });
  },

  importData: (file: File) => {
    const formData = new FormData();
    formData.append('file', file);
    return request.post<Response<ImportResult>>('/api/v1/backup/import', formData, {
      headers: { 'Content-Type': 'multipart/form-data' },
    });
  },
};

// ========== 错误检测 & AI 审阅 API ==========

export const reviewApi = {
  detect: (data: { content: string; types?: string[] }) => {
    return request.post<Response<DetectionResult>>('/api/v1/review/detect', data);
  },

  reviewChapter: (projectId: string, chapterNumber: number) => {
    return request.post<Response<ReviewResult>>(
      `/api/v1/projects/${projectId}/review/chapters/${chapterNumber}`
    );
  },

  reviewProject: (projectId: string) => {
    return request.post<Response<ReviewResult>>(
      `/api/v1/projects/${projectId}/review`
    );
  },

  marketPredict: (projectId: string) => {
    return request.post<Response<MarketPrediction>>(
      `/api/v1/projects/${projectId}/market-predict`
    );
  },
};

// ========== 写作助手 API ==========

export const writingApi = {
  assist: (data: WritingAssistantRequest) => {
    return request.post<Response<{ result: string }>>('/api/v1/writing/assist', data);
  },
};

// ========== 对话相关 API ==========

export const chatApi = {
  // 创建对话
  create: (data: CreateConversationRequest) => {
    return request.post<Response<Conversation>>('/api/v1/conversations', data);
  },

  // 获取对话列表
  list: (params?: { page?: number; page_size?: number }) => {
    return request.get<Response<{ conversations: Conversation[]; total: number }>>(
      '/api/v1/conversations',
      { params }
    );
  },

  // 获取对话详情（含消息）
  getById: (id: string) => {
    return request.get<Response<Conversation>>(`/api/v1/conversations/${id}`);
  },

  // 更新对话标题
  update: (id: string, data: { title: string }) => {
    return request.put<Response<void>>(`/api/v1/conversations/${id}`, data);
  },

  // 删除对话
  delete: (id: string) => {
    return request.delete<Response<void>>(`/api/v1/conversations/${id}`);
  },

  // 发送消息（非流式）
  sendMessage: (conversationId: string, content: string) => {
    return request.post<Response<SendMessageResponse>>(
      `/api/v1/conversations/${conversationId}/messages`,
      { content, stream: false }
    );
  },
};
