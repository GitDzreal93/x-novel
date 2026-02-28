import { useState } from 'react';
import { Button, Card, List, Tag, message, Modal, InputNumber, Spin, Space, Input, Form, Divider, Alert } from 'antd';
import { PlayCircleOutlined, CheckOutlined, PlusOutlined, EditOutlined, SaveOutlined, ExpandOutlined, LockOutlined } from '@ant-design/icons';
import { useNavigate } from 'react-router-dom';
import { chapterApi } from '../../api';
import type { Project, Chapter } from '../../types';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';

const { TextArea } = Input;

interface ChapterPanelProps {
  project: Project;
}

function ChapterPanel({ project }: ChapterPanelProps) {
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const [page, setPage] = useState(1);
  const [selectedChapter, setSelectedChapter] = useState<Chapter | null>(null);
  const [detailModalOpen, setDetailModalOpen] = useState(false);
  const [createModalOpen, setCreateModalOpen] = useState(false);
  const [createForm] = Form.useForm();
  const [editForm] = Form.useForm();

  // 获取章节列表
  const { data: chaptersData, isLoading, refetch } = useQuery({
    queryKey: ['chapters', project.id, page],
    queryFn: () =>
      chapterApi.list(project.id, { page, page_size: 20 }).then((res) => res.data.data),
    enabled: !!project.id,
  });

  const chapters = chaptersData?.chapters || [];
  const total = chaptersData?.total || 0;

  // 创建章节
  const createMutation = useMutation({
    mutationFn: (data: { chapter_number: number; title: string; blueprint_summary: string }) =>
      chapterApi.create(project.id, data),
    onSuccess: () => {
      message.success('章节创建成功');
      setCreateModalOpen(false);
      createForm.resetFields();
      queryClient.invalidateQueries({ queryKey: ['chapters', project.id] });
    },
    onError: () => {
      message.error('创建失败');
    },
  });

  // 生成章节内容
  const generateMutation = useMutation({
    mutationFn: (chapterNumber: number) =>
      chapterApi.generateContent(project.id, chapterNumber, { overwrite: false }),
    onSuccess: () => {
      message.success('章节内容生成成功');
      queryClient.invalidateQueries({ queryKey: ['chapters', project.id] });
      if (selectedChapter) {
        // 重新获取章节详情
        chapterApi.getByNumber(project.id, selectedChapter.chapter_number).then((res) => {
          setSelectedChapter(res.data.data);
        });
      }
    },
    onError: (err: any) => {
      if (err?.response?.data?.message?.includes('已有内容')) {
        message.error('章节已有内容，请先清空或选择覆盖');
      } else {
        message.error('生成失败');
      }
    },
  });

  // 扩写章节
  const enrichMutation = useMutation({
    mutationFn: (chapterNumber: number) =>
      chapterApi.enrich(project.id, chapterNumber, { target_words: 3000 }),
    onSuccess: () => {
      message.success('章节扩写成功');
      queryClient.invalidateQueries({ queryKey: ['chapters', project.id] });
      if (selectedChapter) {
        chapterApi.getByNumber(project.id, selectedChapter.chapter_number).then((res) => {
          setSelectedChapter(res.data.data);
        });
      }
    },
    onError: () => {
      message.error('扩写失败');
    },
  });

  // 定稿章节
  const finalizeMutation = useMutation({
    mutationFn: (chapterNumber: number) =>
      chapterApi.finalize(project.id, chapterNumber, { update_summary: true }),
    onSuccess: () => {
      message.success('章节定稿成功');
      queryClient.invalidateQueries({ queryKey: ['chapters', project.id] });
      if (selectedChapter) {
        chapterApi.getByNumber(project.id, selectedChapter.chapter_number).then((res) => {
          setSelectedChapter(res.data.data);
        });
      }
    },
    onError: () => {
      message.error('定稿失败');
    },
  });

  const handleCreate = () => {
    createForm.validateFields().then((values) => {
      createMutation.mutate(values);
    });
  };

  const handleGenerate = () => {
    if (!selectedChapter) return;
    Modal.confirm({
      title: '生成章节内容',
      content: `确定要生成第 ${selectedChapter.chapter_number} 章的内容吗？`,
      onOk: () => {
        generateMutation.mutate(selectedChapter.chapter_number);
      },
    });
  };

  const handleEnrich = () => {
    if (!selectedChapter) return;
    enrichMutation.mutate(selectedChapter.chapter_number);
  };

  const handleFinalize = () => {
    if (!selectedChapter) return;
    Modal.confirm({
      title: '定稿章节',
      content: '定稿后将无法再修改内容，确定要定稿吗？',
      onOk: () => {
        finalizeMutation.mutate(selectedChapter.chapter_number);
      },
    });
  };

  const handleChapterClick = (chapter: Chapter) => {
    setSelectedChapter(chapter);
    setDetailModalOpen(true);
    editForm.setFieldsValue({
      title: chapter.title,
      content: chapter.content,
    });
  };

  const getStatusTag = (status: string) => {
    const statusMap: Record<string, { color: string; text: string }> = {
      not_started: { color: 'default', text: '未开始' },
      draft: { color: 'processing', text: '草稿' },
      completed: { color: 'success', text: '已完成' },
    };
    const { color, text } = statusMap[status] || { color: 'default', text: status };
    return <Tag color={color}>{text}</Tag>;
  };

  return (
    <div>
      <div className="flex justify-between items-center mb-4">
        <div>
          <h3 className="text-lg font-semibold">章节写作</h3>
          <p className="text-gray-500 text-sm">
            进度: {project.completed_chapters || 0}/{project.chapter_count} 章 |
            总字数: {project.total_words?.toLocaleString() || 0} 字
          </p>
        </div>
        <Button
          type="primary"
          icon={<PlusOutlined />}
          onClick={() => setCreateModalOpen(true)}
        >
          创建章节
        </Button>
      </div>

      {isLoading ? (
        <div className="text-center py-12">
          <Spin />
        </div>
      ) : chapters.length === 0 ? (
        <Card className="text-center py-8">
          <p className="mb-4 text-gray-500">暂无章节</p>
          <div className="space-x-2">
            <Button type="primary" onClick={() => setCreateModalOpen(true)}>
              创建章节
            </Button>
            <Button onClick={() => navigate(`/projects/${project.id}`)}>
              前往生成大纲
            </Button>
          </div>
        </Card>
      ) : (
        <List
          grid={{ gutter: 16, xs: 1, sm: 2, md: 3, lg: 4, xl: 4, xxl: 6 }}
          dataSource={chapters}
          pagination={{
            total,
            pageSize: 20,
            current: page,
            onChange: setPage,
          }}
          renderItem={(chapter) => (
            <List.Item>
              <Card
                size="small"
                hoverable
                onClick={() => handleChapterClick(chapter)}
                className={chapter.is_finalized ? 'border-green-500' : ''}
              >
                <div className="space-y-2">
                  <div className="flex justify-between items-start">
                    <span className="font-semibold">第 {chapter.chapter_number} 章</span>
                    {getStatusTag(chapter.status)}
                  </div>
                  {chapter.title && (
                    <div className="text-sm font-medium truncate">{chapter.title}</div>
                  )}
                  {chapter.word_count > 0 && (
                    <div className="text-xs text-gray-500">
                      {chapter.word_count} 字
                    </div>
                  )}
                  <div className="flex gap-1">
                    {chapter.status === 'not_started' && (
                      <Button
                        size="small"
                        type="primary"
                        icon={<PlayCircleOutlined />}
                        onClick={(e) => {
                          e.stopPropagation();
                          setSelectedChapter(chapter);
                          handleGenerate();
                        }}
                        loading={generateMutation.isPending}
                      >
                        生成
                      </Button>
                    )}
                    {chapter.is_finalized && (
                      <Tag icon={<CheckOutlined />} color="success">
                        已定稿
                      </Tag>
                    )}
                  </div>
                </div>
              </Card>
            </List.Item>
          )}
        />
      )}

      {/* 创建章节弹窗 */}
      <Modal
        title="创建章节"
        open={createModalOpen}
        onCancel={() => {
          setCreateModalOpen(false);
          createForm.resetFields();
        }}
        onOk={handleCreate}
        confirmLoading={createMutation.isPending}
      >
        <Form form={createForm} layout="vertical">
          <Form.Item
            name="chapter_number"
            label="章节号"
            rules={[{ required: true, message: '请输入章节号' }]}
          >
            <InputNumber type="number" min={1} placeholder="章节号" className="w-full" />
          </Form.Item>
          <Form.Item
            name="title"
            label="章节标题"
            rules={[{ required: true, message: '请输入章节标题' }]}
          >
            <Input placeholder="章节标题" />
          </Form.Item>
          <Form.Item
            name="blueprint_summary"
            label="章节摘要"
            rules={[{ required: true, message: '请输入章节摘要' }]}
          >
            <TextArea rows={4} placeholder="章节摘要（可从大纲中复制）" />
          </Form.Item>
        </Form>
      </Modal>

      {/* 章节详情弹窗 */}
      <Modal
        title={`第 ${selectedChapter?.chapter_number} 章${selectedChapter?.title ? ' - ' + selectedChapter.title : ''}`}
        open={detailModalOpen}
        onCancel={() => {
          setDetailModalOpen(false);
          setSelectedChapter(null);
        }}
        footer={null}
        width={900}
      >
        {selectedChapter && (
          <div className="space-y-4">
            <div className="flex justify-between items-center">
              <span className="text-sm text-gray-500">
                状态: {getStatusTag(selectedChapter.status)} |
                字数: {selectedChapter.word_count} |
                {selectedChapter.is_finalized && <Tag icon={<LockOutlined />} color="success">已定稿</Tag>}
              </span>
              <Space>
                {!selectedChapter.is_finalized && (
                  <>
                    {selectedChapter.status !== 'not_started' && (
                      <>
                        <Button
                          size="small"
                          icon={<ExpandOutlined />}
                          onClick={handleEnrich}
                          loading={enrichMutation.isPending}
                        >
                          扩写
                        </Button>
                        <Button
                          size="small"
                          type="primary"
                          icon={<CheckOutlined />}
                          onClick={handleFinalize}
                          loading={finalizeMutation.isPending}
                        >
                          定稿
                        </Button>
                      </>
                    )}
                    {selectedChapter.status === 'not_started' && (
                      <Button
                        type="primary"
                        icon={<PlayCircleOutlined />}
                        onClick={handleGenerate}
                        loading={generateMutation.isPending}
                      >
                        生成内容
                      </Button>
                    )}
                  </>
                )}
              </Space>
            </div>

            <Divider />

            <Form form={editForm} layout="vertical">
              <Form.Item name="content">
                <TextArea
                  rows={20}
                  placeholder="章节内容..."
                  disabled={selectedChapter.is_finalized}
                  style={{ fontFamily: 'monospace', lineHeight: '1.8' }}
                />
              </Form.Item>
            </Form>

            {selectedChapter.is_finalized && (
              <Alert
                message="此章节已定稿，无法再修改内容"
                type="info"
                showIcon
              />
            )}
          </div>
        )}
      </Modal>
    </div>
  );
}

export default ChapterPanel;
