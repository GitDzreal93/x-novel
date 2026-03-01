import { useState, useRef } from 'react';
import {
  Button, Card, Tag, Modal, InputNumber, Spin, Space, Input, Form,
  Alert, App, Pagination, Typography, Flex, Row, Col, theme, Collapse,
} from 'antd';
import {
  PlayCircleOutlined, CheckOutlined, PlusOutlined, ExpandOutlined,
  LockOutlined, EditOutlined, FileSearchOutlined, RobotOutlined,
  MenuFoldOutlined, MenuUnfoldOutlined, BugOutlined,
} from '@ant-design/icons';
import { useNavigate } from 'react-router-dom';
import { chapterApi } from '../../api';
import type { Project, Chapter } from '../../types';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import WritingAssistant from './WritingAssistant';
import ErrorDetectionPanel from './ErrorDetectionPanel';
import RichEditor, { type RichEditorRef, type ErrorMark } from '../common/RichEditor';
import type { DetectionIssue } from '../../types';

const { TextArea } = Input;
const { Title, Text } = Typography;

interface ChapterPanelProps {
  project: Project;
}

function ChapterPanel({ project }: ChapterPanelProps) {
  const { message, modal } = App.useApp();
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const { token } = theme.useToken();
  const [page, setPage] = useState(1);
  const [selectedChapter, setSelectedChapter] = useState<Chapter | null>(null);
  const [detailModalOpen, setDetailModalOpen] = useState(false);
  const [createModalOpen, setCreateModalOpen] = useState(false);
  const [assistantOpen, setAssistantOpen] = useState(false);
  const [detectionOpen, setDetectionOpen] = useState(false);
  const [editorContent, setEditorContent] = useState('');
  const [errorMarks, setErrorMarks] = useState<ErrorMark[]>([]);
  const editorRef = useRef<RichEditorRef>(null);
  const [createForm] = Form.useForm();
  const [editForm] = Form.useForm();

  const { data: chaptersData, isLoading } = useQuery({
    queryKey: ['chapters', project.id, page],
    queryFn: () =>
      chapterApi.list(project.id, { page, page_size: 20 }).then((res) => {
        if (!res?.data) {
          throw new Error(res?.message || '获取章节列表失败');
        }
        return res.data;
      }),
    enabled: !!project.id,
  });

  const chapters = chaptersData?.chapters || [];
  const total = chaptersData?.total || 0;

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

  const updateMutation = useMutation({
    mutationFn: (data: { title?: string; content?: string }) =>
      chapterApi.update(project.id, selectedChapter!.chapter_number, data),
    onSuccess: () => {
      message.success('保存成功');
      queryClient.invalidateQueries({ queryKey: ['chapters', project.id] });
    },
    onError: () => {
      message.error('保存失败');
    },
  });

  const generateMutation = useMutation({
    mutationFn: (chapterNumber: number) =>
      chapterApi.generateContent(project.id, chapterNumber, { overwrite: false }),
    onSuccess: () => {
      message.success('章节内容生成成功');
      queryClient.invalidateQueries({ queryKey: ['chapters', project.id] });
      if (selectedChapter) {
        chapterApi.getByNumber(project.id, selectedChapter.chapter_number).then((res) => {
          if (res?.data) {
            setSelectedChapter(res.data);
            setEditorContent(res.data.content || '');
          }
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

  const enrichMutation = useMutation({
    mutationFn: (chapterNumber: number) =>
      chapterApi.enrich(project.id, chapterNumber, { target_words: 3000 }),
    onSuccess: () => {
      message.success('章节扩写成功');
      queryClient.invalidateQueries({ queryKey: ['chapters', project.id] });
      if (selectedChapter) {
        chapterApi.getByNumber(project.id, selectedChapter.chapter_number).then((res) => {
          if (res?.data) {
            setSelectedChapter(res.data);
            setEditorContent(res.data.content || '');
          }
        });
      }
    },
    onError: () => {
      message.error('扩写失败');
    },
  });

  const finalizeMutation = useMutation({
    mutationFn: (chapterNumber: number) =>
      chapterApi.finalize(project.id, chapterNumber, { update_summary: true }),
    onSuccess: () => {
      message.success('章节定稿成功');
      queryClient.invalidateQueries({ queryKey: ['chapters', project.id] });
      if (selectedChapter) {
        chapterApi.getByNumber(project.id, selectedChapter.chapter_number).then((res) => {
          if (res?.data) {
            setSelectedChapter(res.data);
          }
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
    modal.confirm({
      title: '生成章节内容',
      content: `确定要生成第 ${selectedChapter.chapter_number} 章的内容吗？`,
      centered: true,
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
    modal.confirm({
      title: '定稿章节',
      content: '定稿后将无法再修改内容，确定要定稿吗？',
      centered: true,
      okButtonProps: { danger: true },
      onOk: () => {
        finalizeMutation.mutate(selectedChapter.chapter_number);
      },
    });
  };

  const handleSaveContent = () => {
    if (!selectedChapter) return;
    const content = editorRef.current?.getText() || '';
    updateMutation.mutate({ content });
  };

  const handleChapterClick = (chapter: Chapter) => {
    setSelectedChapter(chapter);
    setDetailModalOpen(true);
    setEditorContent(chapter.content || '');
    editForm.setFieldsValue({
      title: chapter.title,
    });
  };

  const getStatusTag = (status: string) => {
    switch (status) {
      case 'not_started': return <Tag color="default" bordered={false}>未开始</Tag>;
      case 'draft': return <Tag color="processing" bordered={false}>草稿</Tag>;
      case 'completed': return <Tag color="success" bordered={false}>已完成</Tag>;
      default: return <Tag color="default" bordered={false}>{status}</Tag>;
    }
  };

  return (
    <div style={{ padding: '0 24px 24px' }}>
      <Flex
        justify="space-between"
        align="center"
        wrap
        gap={16}
        style={{ marginBottom: 20 }}
      >
        <div>
          <Flex align="center" gap={8}>
            <Title level={5} style={{ margin: 0 }}>章节创作区</Title>
            <EditOutlined style={{ color: token.colorPrimary, fontSize: 16 }} />
          </Flex>
          <Text type="secondary" style={{ fontSize: 13 }}>
            进度：<Text strong>{project.completed_chapters || 0}</Text> / {project.chapter_count} 章
            <span style={{
              display: 'inline-block',
              width: 4, height: 4, borderRadius: '50%',
              background: token.colorTextQuaternary,
              margin: '0 10px',
              verticalAlign: 'middle',
            }} />
            总字数：<Text strong>{(project.total_words || 0).toLocaleString()}</Text> 字
          </Text>
        </div>
        <Button
          type="primary"
          icon={<PlusOutlined />}
          onClick={() => setCreateModalOpen(true)}
          size="large"
        >
          创建章节
        </Button>
      </Flex>

      {isLoading ? (
        <Flex justify="center" align="center" style={{ height: 160 }}>
          <Spin size="large" />
        </Flex>
      ) : chapters.length === 0 ? (
        <Card
          style={{
            textAlign: 'center',
            padding: '60px 0',
            borderStyle: 'dashed',
          }}
        >
          <Text style={{ fontSize: 16, display: 'block', marginBottom: 16 }}>暂无章节内容</Text>
          <Space>
            <Button type="primary" onClick={() => setCreateModalOpen(true)}>
              创建章节
            </Button>
            <Button onClick={() => navigate(`/projects/${project.id}?tab=blueprint`)}>
              前往生成大纲
            </Button>
          </Space>
        </Card>
      ) : (
        <>
          <Row gutter={[16, 16]}>
            {chapters.map((chapter) => (
              <Col key={chapter.id} xs={24} sm={12} md={8} lg={6} xl={4} xxl={4}>
                <Card
                  size="small"
                  hoverable
                  onClick={() => handleChapterClick(chapter)}
                  style={{
                    borderColor: chapter.is_finalized ? token.colorSuccess : undefined,
                    background: chapter.is_finalized ? token.colorSuccessBg : undefined,
                  }}
                  styles={{ body: { padding: 16 } }}
                >
                  <Flex vertical style={{ height: '100%' }}>
                    <Flex justify="space-between" align="flex-start" style={{ marginBottom: 8 }}>
                      <Text
                        strong
                        style={{
                          fontSize: 16,
                          color: chapter.is_finalized ? token.colorSuccess : undefined,
                        }}
                      >
                        第 {chapter.chapter_number} 章
                      </Text>
                      {getStatusTag(chapter.status)}
                    </Flex>

                    {chapter.title && (
                      <Text
                        ellipsis
                        style={{ fontSize: 13, fontWeight: 500, marginBottom: 8 }}
                      >
                        {chapter.title}
                      </Text>
                    )}

                    <Flex
                      justify="space-between"
                      align="center"
                      style={{
                        marginTop: 'auto',
                        paddingTop: 10,
                        borderTop: `1px solid ${token.colorBorderSecondary}`,
                      }}
                    >
                      <Text type="secondary" style={{ fontSize: 12 }}>
                        {chapter.word_count > 0 ? `${chapter.word_count.toLocaleString()} 字` : '-- 字'}
                      </Text>

                      <div onClick={(e) => e.stopPropagation()}>
                        {chapter.status === 'not_started' && (
                          <Button
                            type="text"
                            size="small"
                            icon={<PlayCircleOutlined />}
                            onClick={() => {
                              setSelectedChapter(chapter);
                              handleGenerate();
                            }}
                            loading={generateMutation.isPending && selectedChapter?.id === chapter.id}
                          >
                            生成
                          </Button>
                        )}
                        {chapter.is_finalized && (
                          <LockOutlined
                            style={{ color: token.colorSuccess }}
                            title="已定稿"
                          />
                        )}
                      </div>
                    </Flex>
                  </Flex>
                </Card>
              </Col>
            ))}
          </Row>
          <Flex justify="center" style={{ marginTop: 32 }}>
            <Pagination
              total={total}
              pageSize={20}
              current={page}
              onChange={setPage}
              showSizeChanger={false}
              showQuickJumper
            />
          </Flex>
        </>
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
        centered
      >
        <Form form={createForm} layout="vertical" style={{ marginTop: 16 }}>
          <Form.Item
            name="chapter_number"
            label="章节号"
            rules={[{ required: true, message: '请输入章节号' }]}
            initialValue={chapters.length > 0 ? Math.max(...chapters.map(c => c.chapter_number)) + 1 : 1}
          >
            <InputNumber size="large" min={1} placeholder="章节号" style={{ width: '100%' }} />
          </Form.Item>
          <Form.Item
            name="title"
            label="章节标题"
            rules={[{ required: true, message: '请输入章节标题' }]}
          >
            <Input size="large" placeholder="章节标题" />
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

      {/* 章节详情弹窗/编辑器 */}
      <Modal
        title={
          <Flex align="center" gap={8}>
            <span style={{ fontSize: 16, fontWeight: 600 }}>
              第 {selectedChapter?.chapter_number} 章{selectedChapter?.title ? `：${selectedChapter.title}` : ''}
            </span>
            {selectedChapter?.is_finalized && (
              <Tag icon={<LockOutlined />} color="success" bordered={false}>已定稿</Tag>
            )}
          </Flex>
        }
        open={detailModalOpen}
        onCancel={() => {
          setDetailModalOpen(false);
          setSelectedChapter(null);
          setAssistantOpen(false);
          setDetectionOpen(false);
        }}
        footer={null}
        width={(assistantOpen || detectionOpen) ? 1300 : 1000}
        centered
        style={{ top: 20, transition: 'width 0.3s' }}
        styles={{ body: { padding: 0 } }}
      >
        {selectedChapter && (
          <Flex vertical style={{ height: '80vh', maxHeight: 800 }}>
            {/* 工具栏 */}
            <Flex
              justify="space-between"
              align="center"
              style={{
                padding: '10px 24px',
                borderBottom: `1px solid ${token.colorBorderSecondary}`,
                background: token.colorBgLayout,
                flexShrink: 0,
              }}
            >
              <Flex align="center" gap={16} style={{ fontSize: 13 }}>
                <Text type="secondary">状态: {getStatusTag(selectedChapter.status)}</Text>
                <Text type="secondary">
                  字数: <Text strong>{selectedChapter.word_count.toLocaleString()}</Text>
                </Text>
              </Flex>
              <Space>
                {!selectedChapter.is_finalized && (
                  <>
                    <Button onClick={handleSaveContent} loading={updateMutation.isPending}>
                      保存草稿
                    </Button>

                    {selectedChapter.status !== 'not_started' ? (
                      <>
                        <Button
                          icon={<ExpandOutlined />}
                          onClick={handleEnrich}
                          loading={enrichMutation.isPending}
                        >
                          智能扩写
                        </Button>
                        <Button
                          type="primary"
                          icon={<CheckOutlined />}
                          onClick={handleFinalize}
                          loading={finalizeMutation.isPending}
                          style={{ background: '#16a34a' }}
                        >
                          完成定稿
                        </Button>
                      </>
                    ) : (
                      <Button
                        type="primary"
                        icon={<PlayCircleOutlined />}
                        onClick={handleGenerate}
                        loading={generateMutation.isPending}
                      >
                        AI 生成内容
                      </Button>
                    )}
                  </>
                )}
                <Button
                  type={detectionOpen ? 'primary' : 'default'}
                  icon={<BugOutlined />}
                  onClick={() => {
                    const next = !detectionOpen;
                    setDetectionOpen(next);
                    if (next) setAssistantOpen(false);
                    if (!next) setErrorMarks([]);
                  }}
                  title={detectionOpen ? '收起错误检测' : '展开错误检测'}
                >
                  错误检测
                </Button>
                <Button
                  type={assistantOpen ? 'primary' : 'default'}
                  icon={assistantOpen ? <MenuFoldOutlined /> : <RobotOutlined />}
                  onClick={() => {
                    setAssistantOpen(!assistantOpen);
                    if (!assistantOpen) setDetectionOpen(false);
                  }}
                  title={assistantOpen ? '收起写作助手' : '展开写作助手'}
                >
                  写作助手
                </Button>
              </Space>
            </Flex>

            {/* 编辑器 + 写作助手 */}
            <Flex style={{ flex: 1, minHeight: 0 }}>
              {/* 编辑器区域 */}
              <div style={{ flex: 1, padding: 24, overflow: 'auto', minWidth: 0 }}>
                {selectedChapter.is_finalized && (
                  <Alert
                    message="此章节已定稿，无法再修改内容"
                    type="info"
                    showIcon
                    style={{ marginBottom: 16 }}
                  />
                )}

                {/* 前文摘要 */}
                {project.global_summary && selectedChapter.chapter_number > 1 && (
                  <Collapse
                    size="small"
                    style={{ marginBottom: 16 }}
                    items={[{
                      key: 'summary',
                      label: (
                        <Flex align="center" gap={6}>
                          <FileSearchOutlined style={{ color: token.colorPrimary }} />
                          <Text strong style={{ fontSize: 13 }}>前文摘要</Text>
                          <Text type="secondary" style={{ fontSize: 12 }}>
                            （帮助保持剧情连贯性）
                          </Text>
                        </Flex>
                      ),
                      children: (
                        <pre style={{
                          margin: 0,
                          whiteSpace: 'pre-wrap',
                          wordBreak: 'break-word',
                          fontFamily: 'monospace',
                          fontSize: 12,
                          lineHeight: 1.7,
                          color: token.colorTextSecondary,
                          maxHeight: 200,
                          overflow: 'auto',
                        }}>
                          {project.global_summary}
                        </pre>
                      ),
                    }]}
                  />
                )}

                <RichEditor
                  ref={editorRef}
                  content={editorContent}
                  placeholder={selectedChapter.is_finalized ? '内容为空' : '在这里开始写作，或使用 AI 生成...'}
                  disabled={selectedChapter.is_finalized}
                  onChange={(text) => {
                    setEditorContent(text);
                    if (errorMarks.length > 0) setErrorMarks([]);
                  }}
                  minHeight={400}
                  showToolbar={!selectedChapter.is_finalized}
                  errorMarks={errorMarks}
                />
              </div>

              {/* 错误检测侧边栏 */}
              {detectionOpen && (
                <div style={{
                  width: 300,
                  borderLeft: `1px solid ${token.colorBorderSecondary}`,
                  padding: 16,
                  overflow: 'auto',
                  display: 'flex',
                  flexDirection: 'column',
                  flexShrink: 0,
                }}>
                  <ErrorDetectionPanel
                    content={editorContent}
                    disabled={!editorContent?.trim()}
                    onApplySuggestion={(original, suggestion) => {
                      const updated = editorContent.replace(original, suggestion);
                      setEditorContent(updated);
                      editorRef.current?.setContent(updated);
                      setErrorMarks([]);
                    }}
                    onErrorsDetected={(issues: DetectionIssue[]) => {
                      // 在编辑器文本中查找每个 issue 的位置并标记
                      const text = editorRef.current?.getText() || '';
                      const marks: ErrorMark[] = [];
                      for (const issue of issues) {
                        if (!issue.original) continue;
                        const idx = text.indexOf(issue.original);
                        if (idx >= 0) {
                          marks.push({
                            from: idx + 1, // ProseMirror offset: +1 for doc node
                            to: idx + 1 + issue.original.length,
                            severity: issue.severity,
                            message: `${issue.suggestion || issue.explanation}`,
                          });
                        }
                      }
                      setErrorMarks(marks);
                    }}
                  />
                </div>
              )}

              {/* 写作助手侧边栏 */}
              {assistantOpen && (
                <div style={{
                  width: 300,
                  borderLeft: `1px solid ${token.colorBorderSecondary}`,
                  padding: 16,
                  overflow: 'auto',
                  display: 'flex',
                  flexDirection: 'column',
                  flexShrink: 0,
                }}>
                  <WritingAssistant
                    content={editorContent}
                    projectId={project.id}
                    disabled={selectedChapter.is_finalized}
                    onApplyResult={(result) => {
                      setEditorContent(result);
                      editorRef.current?.setContent(result);
                    }}
                    onAppendResult={(result) => {
                      const updated = editorContent + '\n\n' + result;
                      setEditorContent(updated);
                      editorRef.current?.setContent(updated);
                    }}
                  />
                </div>
              )}
            </Flex>
          </Flex>
        )}
      </Modal>
    </div>
  );
}

export default ChapterPanel;
