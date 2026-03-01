import { useState } from 'react';
import {
  Button, Empty, Modal, Form, Input, InputNumber, Select, App,
  Tag, Progress, Spin, Typography, Card, Row, Col, Flex, theme,
} from 'antd';
import { PlusOutlined, EditOutlined, FileTextOutlined, ClockCircleOutlined } from '@ant-design/icons';
import { useNavigate } from 'react-router-dom';
import { projectApi } from '../api';
import type { CreateProjectRequest } from '../types';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';

const { Title, Text, Paragraph } = Typography;

function ProjectList() {
  const { message } = App.useApp();
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const { token } = theme.useToken();
  const [page] = useState(1);
  const [createModalOpen, setCreateModalOpen] = useState(false);
  const [form] = Form.useForm();

  const { data: projectsData, isLoading } = useQuery({
    queryKey: ['projects', page],
    queryFn: () => projectApi.list({ page, page_size: 10 }).then((res) => res.data),
  });

  const projects = projectsData?.projects || [];

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
    <>
      <Flex justify="space-between" align="center" style={{ marginBottom: 24 }}>
        <div>
          <Title level={4} style={{ margin: 0 }}>我的项目</Title>
          <Text type="secondary" style={{ fontSize: 13 }}>
            管理你的所有小说创作项目
          </Text>
        </div>
        <Button
          type="primary"
          size="large"
          icon={<PlusOutlined />}
          onClick={() => setCreateModalOpen(true)}
        >
          新建项目
        </Button>
      </Flex>

      {isLoading ? (
        <Flex justify="center" align="center" style={{ height: 256 }}>
          <Spin size="large" />
        </Flex>
      ) : projects.length === 0 ? (
        <Card
          style={{
            textAlign: 'center',
            borderStyle: 'dashed',
            padding: '60px 0',
          }}
        >
          <Empty
            image={<FileTextOutlined style={{ fontSize: 48, color: token.colorTextQuaternary }} />}
            imageStyle={{ height: 60 }}
            description={
              <div>
                <Text style={{ fontSize: 15 }}>暂无项目</Text>
                <br />
                <Text type="secondary" style={{ fontSize: 13 }}>
                  创建一个项目，开始你的创作之旅
                </Text>
              </div>
            }
          >
            <Button type="primary" icon={<PlusOutlined />} onClick={() => setCreateModalOpen(true)}>
              创建第一个项目
            </Button>
          </Empty>
        </Card>
      ) : (
        <Row gutter={[20, 20]}>
          {projects.map((project) => {
            const progress = project.chapter_count > 0
              ? Math.round((project.completed_chapters / project.chapter_count) * 100)
              : 0;

            return (
              <Col key={project.id} xs={24} md={12} xl={8}>
                <Card
                  hoverable
                  onClick={() => navigate(`/projects/${project.id}`)}
                  styles={{ body: { padding: 20 } }}
                >
                  <Flex justify="space-between" align="flex-start" style={{ marginBottom: 8 }}>
                    <Tag color={getStatusColor(project.status)} bordered={false}>
                      {getStatusText(project.status)}
                    </Tag>
                    <Flex align="center" gap={4}>
                      <ClockCircleOutlined style={{ fontSize: 11, color: token.colorTextQuaternary }} />
                      <Text type="secondary" style={{ fontSize: 11 }}>
                        {formatDate(project.updated_at)}
                      </Text>
                    </Flex>
                  </Flex>

                  <Title
                    level={5}
                    ellipsis
                    style={{ margin: '8px 0', fontSize: 17 }}
                  >
                    {project.title}
                  </Title>

                  <Paragraph
                    type="secondary"
                    ellipsis={{ rows: 2 }}
                    style={{ fontSize: 13, marginBottom: 12, minHeight: 40 }}
                  >
                    {project.topic || '暂无故事梗概'}
                  </Paragraph>

                  {project.genre && project.genre.length > 0 && (
                    <Flex wrap gap={6} style={{ marginBottom: 12 }}>
                      {project.genre.slice(0, 3).map((g) => (
                        <Tag key={g} color="purple" bordered={false} style={{ fontSize: 11 }}>
                          {g}
                        </Tag>
                      ))}
                      {project.genre.length > 3 && (
                        <Tag bordered={false} style={{ fontSize: 11 }}>
                          +{project.genre.length - 3}
                        </Tag>
                      )}
                    </Flex>
                  )}

                  <div
                    style={{
                      paddingTop: 12,
                      borderTop: `1px solid ${token.colorBorderSecondary}`,
                    }}
                  >
                    <Flex justify="space-between" style={{ marginBottom: 6, fontSize: 12 }}>
                      <Flex align="center" gap={4}>
                        <EditOutlined style={{ color: token.colorTextQuaternary }} />
                        <Text type="secondary" style={{ fontSize: 12 }}>
                          {project.completed_chapters}/{project.chapter_count} 章
                        </Text>
                      </Flex>
                      <Flex align="center" gap={4}>
                        <FileTextOutlined style={{ color: token.colorTextQuaternary }} />
                        <Text type="secondary" style={{ fontSize: 12 }}>
                          {project.total_words > 10000
                            ? `${(project.total_words / 10000).toFixed(1)}万字`
                            : `${project.total_words}字`}
                        </Text>
                      </Flex>
                    </Flex>
                    <Progress
                      percent={progress}
                      size="small"
                      showInfo={false}
                      strokeColor={{ '0%': '#818cf8', '100%': '#6366f1' }}
                    />
                  </div>
                </Card>
              </Col>
            );
          })}
        </Row>
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
        okText="创建"
        cancelText="取消"
        width={580}
        centered
      >
        <Form form={form} layout="vertical" style={{ marginTop: 16 }}>
          <Form.Item
            name="title"
            label="项目名称"
            rules={[{ required: true, message: '请输入项目名称' }]}
          >
            <Input size="large" placeholder="例如：我的第一部小说" />
          </Form.Item>

          <Form.Item name="topic" label="故事梗概">
            <Input.TextArea rows={3} placeholder="简单描述你的故事创意..." />
          </Form.Item>

          <Form.Item name="genre" label="作品类型">
            <Select
              mode="tags"
              size="large"
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

          <Row gutter={16}>
            <Col span={12}>
              <Form.Item name="chapter_count" label="计划章节数" initialValue={100}>
                <InputNumber size="large" min={1} max={10000} style={{ width: '100%' }} />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="words_per_chapter" label="每章字数" initialValue={3000}>
                <InputNumber size="large" min={100} max={50000} style={{ width: '100%' }} />
              </Form.Item>
            </Col>
          </Row>

          <Form.Item name="user_guidance" label="创作要求（选填）">
            <Input.TextArea rows={2} placeholder="例如：风格轻松幽默，主角性格鲜明..." />
          </Form.Item>
        </Form>
      </Modal>
    </>
  );
}

function getStatusColor(status: string) {
  switch (status) {
    case 'writing': return 'processing';
    case 'completed': return 'success';
    case 'published': return 'purple';
    default: return 'default';
  }
}

function getStatusText(status: string): string {
  const statusMap: Record<string, string> = {
    draft: '构思中',
    writing: '连载中',
    completed: '已完结',
    published: '已发布',
  };
  return statusMap[status] || status;
}

function formatDate(dateStr: string): string {
  if (!dateStr) return '';
  const d = new Date(dateStr);
  const now = new Date();
  const diff = now.getTime() - d.getTime();
  const days = Math.floor(diff / (1000 * 60 * 60 * 24));
  if (days === 0) return '今天';
  if (days === 1) return '昨天';
  if (days < 7) return `${days}天前`;
  return `${d.getMonth() + 1}/${d.getDate()}`;
}

export default ProjectList;
