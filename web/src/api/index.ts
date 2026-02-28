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
