import { useState } from 'react';
import { Card, Button, Empty, Modal, Form, Input, InputNumber, Select, App } from 'antd';
import { PlusOutlined, BookOutlined } from '@ant-design/icons';
import { useNavigate } from 'react-router-dom';
import { projectApi } from '../api';
import type { CreateProjectRequest } from '../types';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';

function ProjectList() {
  const { message } = App.useApp();
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const [page] = useState(1);
  const [createModalOpen, setCreateModalOpen] = useState(false);
  const [form] = Form.useForm();

  // 获取项目列表
  const { data: projectsData, isLoading } = useQuery({
    queryKey: ['projects', page],
    queryFn: () =>
      projectApi.list({ page, page_size: 10 }).then((res) => res.data),
  });

  const projects = projectsData?.projects || [];

  // 创建项目
  const createMutation = useMutation({
    mutationFn: (data: CreateProjectRequest) => projectApi.create(data),
    onSuccess: () => {
      message.success('项目创建成功');
      setCreateModalOpen(false);
      form.resetFields();
      queryClient.invalidateQueries({ queryKey: ['projects'] });
    },
    onError: () => {
      message.error('项目创建失败');
    },
  });

  const handleCreate = () => {
    form.validateFields().then((values) => {
      createMutation.mutate(values);
    });
  };

  return (
    <div>
      <div className="flex justify-between items-center mb-6">
        <h2 className="text-2xl font-bold">我的项目</h2>
        <Button
          type="primary"
          icon={<PlusOutlined />}
          onClick={() => setCreateModalOpen(true)}
        >
          新建项目
        </Button>
      </div>

      {isLoading ? (
        <div className="text-center py-12">加载中...</div>
      ) : projects.length === 0 ? (
        <Empty
          image={Empty.PRESENTED_IMAGE_SIMPLE}
          description="暂无项目"
          className="py-12"
        >
          <Button type="primary" icon={<PlusOutlined />} onClick={() => setCreateModalOpen(true)}>
            创建第一个项目
          </Button>
        </Empty>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {projects.map((project) => (
            <Card
              key={project.id}
              hoverable
              className="relative"
              onClick={() => navigate(`/projects/${project.id}`)}
            >
              <div className="flex items-start justify-between">
                <div className="flex-1">
                  <div className="flex items-center gap-2 mb-2">
                    <BookOutlined className="text-lg" />
                    <h3 className="text-lg font-semibold truncate">{project.title}</h3>
                  </div>
                  {project.topic && (
                    <p className="text-gray-500 text-sm mb-3 line-clamp-2">{project.topic}</p>
                  )}
                  <div className="flex flex-wrap gap-1 mb-3">
                    {project.genre?.map((g) => (
                      <span
                        key={g}
                        className="px-2 py-0.5 bg-blue-100 text-blue-600 text-xs rounded"
                      >
                        {g}
                      </span>
                    ))}
                  </div>
                  <div className="text-xs text-gray-400 space-y-1">
                    <div>进度: {project.completed_chapters}/{project.chapter_count} 章</div>
                    <div>字数: {project.total_words.toLocaleString()} 字</div>
                    <div>状态: {getStatusText(project.status)}</div>
                  </div>
                </div>
              </div>
            </Card>
          ))}
        </div>
      )}

      <Modal
        title="创建新项目"
        open={createModalOpen}
        onOk={handleCreate}
        onCancel={() => {
          setCreateModalOpen(false);
          form.resetFields();
        }}
        confirmLoading={createMutation.isPending}
      >
        <Form form={form} layout="vertical">
          <Form.Item
            name="title"
            label="项目名称"
            rules={[{ required: true, message: '请输入项目名称' }]}
          >
            <Input placeholder="例如：我的第一部小说" />
          </Form.Item>

          <Form.Item name="topic" label="故事梗概">
            <Input.TextArea
              rows={3}
              placeholder="简单描述你的故事创意..."
            />
          </Form.Item>

          <Form.Item name="genre" label="作品类型">
            <Select
              mode="tags"
              placeholder="选择或输入类型标签"
              options={[
                { label: '玄幻', value: '玄幻' },
                { label: '仙侠', value: '仙侠' },
                { label: '都市', value: '都市' },
                { label: '历史', value: '历史' },
                { label: '科幻', value: '科幻' },
                { label: '游戏', value: '游戏' },
                { label: '军事', value: '军事' },
              ]}
            />
          </Form.Item>

          <div className="grid grid-cols-2 gap-4">
            <Form.Item
              name="chapter_count"
              label="计划章节数"
              initialValue={100}
            >
              <InputNumber min={1} max={10000} className="w-full" />
            </Form.Item>

            <Form.Item
              name="words_per_chapter"
              label="每章字数"
              initialValue={3000}
            >
              <InputNumber min={100} max={50000} className="w-full" />
            </Form.Item>
          </div>

          <Form.Item name="user_guidance" label="创作要求">
            <Input.TextArea
              rows={3}
              placeholder="例如：风格轻松幽默，主角性格鲜明..."
            />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
}

function getStatusText(status: string): string {
  const statusMap: Record<string, string> = {
    draft: '草稿',
    writing: '写作中',
    completed: '已完成',
    published: '已发布',
  };
  return statusMap[status] || status;
}

export default ProjectList;
