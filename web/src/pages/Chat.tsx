import { useState, useRef, useEffect, useCallback } from 'react';
import {
  Layout, Button, Input, Typography, Flex, theme, List, Dropdown, Empty,
  Spin, Tag, App, Popconfirm, Avatar,
} from 'antd';
import {
  PlusOutlined,
  SendOutlined,
  DeleteOutlined,
  EditOutlined,
  BulbOutlined,
  GlobalOutlined,
  UserOutlined,
  RobotOutlined,
  ExperimentOutlined,
  TeamOutlined,
  MessageOutlined,
  CheckOutlined,
  CloseOutlined,
} from '@ant-design/icons';
import ReactMarkdown from 'react-markdown';
import remarkGfm from 'remark-gfm';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { chatApi } from '../api';
import type { Conversation, ChatMessage, ChatMode } from '../types';

const { Sider, Content } = Layout;
const { Text, Paragraph } = Typography;
const { TextArea } = Input;

const CHAT_MODES: { key: ChatMode; label: string; icon: React.ReactNode; color: string; desc: string }[] = [
  { key: 'creative', label: '创意启发', icon: <BulbOutlined />, color: '#f59e0b', desc: '激发灵感，突破瓶颈' },
  { key: 'building', label: '设定完善', icon: <GlobalOutlined />, color: '#3b82f6', desc: '完善世界观和体系' },
  { key: 'character', label: '角色塑造', icon: <TeamOutlined />, color: '#8b5cf6', desc: '创建立体人物' },
  { key: 'general', label: '通用助手', icon: <ExperimentOutlined />, color: '#10b981', desc: '全方位创作帮助' },
];

function Chat() {
  const { token } = theme.useToken();
  const { message: messageApi } = App.useApp();
  const queryClient = useQueryClient();

  const [activeConvId, setActiveConvId] = useState<string | null>(null);
  const [inputValue, setInputValue] = useState('');
  const [sending, setSending] = useState(false);
  const [streamingContent, setStreamingContent] = useState('');
  const [editingTitleId, setEditingTitleId] = useState<string | null>(null);
  const [editingTitle, setEditingTitle] = useState('');
  const messagesEndRef = useRef<HTMLDivElement>(null);
  const inputRef = useRef<any>(null);

  // 对话列表
  const { data: convListData, isLoading: loadingList } = useQuery({
    queryKey: ['conversations'],
    queryFn: async () => {
      const res = await chatApi.list({ page: 1, page_size: 100 });
      return res.data;
    },
  });

  const conversations = convListData?.conversations || [];

  // 当前对话详情
  const { data: convDetailData, isLoading: loadingDetail } = useQuery({
    queryKey: ['conversation', activeConvId],
    queryFn: async () => {
      if (!activeConvId) return null;
      const res = await chatApi.getById(activeConvId);
      return res.data;
    },
    enabled: !!activeConvId,
  });

  const messages: ChatMessage[] = convDetailData?.messages || [];

  // 创建对话
  const createMutation = useMutation({
    mutationFn: (mode: ChatMode) => chatApi.create({ mode }),
    onSuccess: (res) => {
      queryClient.invalidateQueries({ queryKey: ['conversations'] });
      setActiveConvId(res.data?.id || null);
    },
  });

  // 删除对话
  const deleteMutation = useMutation({
    mutationFn: (id: string) => chatApi.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['conversations'] });
      if (activeConvId) setActiveConvId(null);
    },
  });

  // 重命名对话
  const renameMutation = useMutation({
    mutationFn: ({ id, title }: { id: string; title: string }) => chatApi.update(id, { title }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['conversations'] });
      setEditingTitleId(null);
    },
  });

  // 滚动到底部
  const scrollToBottom = useCallback(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  }, []);

  useEffect(() => {
    scrollToBottom();
  }, [messages, streamingContent, scrollToBottom]);

  // SSE 发送消息
  const handleSend = async () => {
    if (!inputValue.trim() || !activeConvId || sending) return;

    const content = inputValue.trim();
    setInputValue('');
    setSending(true);
    setStreamingContent('');

    try {
      const deviceId = localStorage.getItem('x-novel-device-id') || '';

      const response = await fetch(`/api/v1/conversations/${activeConvId}/messages`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'X-Device-ID': deviceId,
        },
        body: JSON.stringify({ content, stream: true }),
      });

      if (!response.ok) {
        throw new Error(`HTTP ${response.status}`);
      }

      const reader = response.body?.getReader();
      const decoder = new TextDecoder();
      let accumulated = '';

      if (reader) {
        while (true) {
          const { done, value } = await reader.read();
          if (done) break;

          const text = decoder.decode(value, { stream: true });
          const lines = text.split('\n');

          for (const line of lines) {
            if (!line.startsWith('data: ')) continue;
            const jsonStr = line.slice(6).trim();
            if (!jsonStr) continue;

            try {
              const parsed = JSON.parse(jsonStr);
              if (parsed.done) {
                // 流结束，刷新对话数据
                queryClient.invalidateQueries({ queryKey: ['conversation', activeConvId] });
                queryClient.invalidateQueries({ queryKey: ['conversations'] });
                setStreamingContent('');
              } else if (parsed.error) {
                messageApi.error(parsed.error);
              } else if (parsed.content !== undefined) {
                accumulated += parsed.content;
                setStreamingContent(accumulated);
              }
            } catch {
              // 忽略解析失败的行
            }
          }
        }
      }
    } catch (err: any) {
      messageApi.error('发送失败：' + (err.message || '未知错误'));
      // 降级为非流式请求
      try {
        await chatApi.sendMessage(activeConvId, content);
        queryClient.invalidateQueries({ queryKey: ['conversation', activeConvId] });
        queryClient.invalidateQueries({ queryKey: ['conversations'] });
      } catch {
        messageApi.error('消息发送失败');
      }
    } finally {
      setSending(false);
      setStreamingContent('');
      inputRef.current?.focus();
    }
  };

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      handleSend();
    }
  };

  const getModeInfo = (mode: string) => {
    return CHAT_MODES.find(m => m.key === mode) || CHAT_MODES[3];
  };

  const activeConv = conversations.find(c => c.id === activeConvId);

  return (
    <Layout style={{ height: 'calc(100vh - 64px)', margin: -24 }}>
      {/* 左侧：对话列表 */}
      <Sider
        width={280}
        style={{
          background: token.colorBgContainer,
          borderRight: `1px solid ${token.colorBorderSecondary}`,
          display: 'flex',
          flexDirection: 'column',
          overflow: 'hidden',
        }}
      >
        <Flex vertical style={{ height: '100%' }}>
          {/* 新建对话 */}
          <div style={{ padding: '16px 12px 8px' }}>
            <Dropdown
              menu={{
                items: CHAT_MODES.map(m => ({
                  key: m.key,
                  icon: m.icon,
                  label: (
                    <Flex vertical gap={2}>
                      <Text strong style={{ fontSize: 13 }}>{m.label}</Text>
                      <Text type="secondary" style={{ fontSize: 11 }}>{m.desc}</Text>
                    </Flex>
                  ),
                  onClick: () => createMutation.mutate(m.key),
                })),
              }}
              trigger={['click']}
            >
              <Button
                type="dashed"
                icon={<PlusOutlined />}
                block
                loading={createMutation.isPending}
              >
                新建对话
              </Button>
            </Dropdown>
          </div>

          {/* 对话列表 */}
          <div style={{ flex: 1, overflow: 'auto', padding: '0 8px 8px' }}>
            {loadingList ? (
              <Flex justify="center" style={{ padding: 40 }}>
                <Spin size="small" />
              </Flex>
            ) : conversations.length === 0 ? (
              <Empty
                image={Empty.PRESENTED_IMAGE_SIMPLE}
                description="暂无对话"
                style={{ marginTop: 60 }}
              />
            ) : (
              <List
                dataSource={conversations}
                split={false}
                renderItem={(conv: Conversation) => {
                  const modeInfo = getModeInfo(conv.mode);
                  const isActive = conv.id === activeConvId;
                  const isEditing = editingTitleId === conv.id;

                  return (
                    <div
                      key={conv.id}
                      onClick={() => !isEditing && setActiveConvId(conv.id)}
                      style={{
                        padding: '10px 12px',
                        marginBottom: 2,
                        borderRadius: token.borderRadius,
                        cursor: 'pointer',
                        background: isActive ? token.colorPrimaryBg : 'transparent',
                        border: isActive ? `1px solid ${token.colorPrimaryBorder}` : '1px solid transparent',
                        transition: 'all 0.2s',
                      }}
                    >
                      <Flex justify="space-between" align="center">
                        <Flex vertical gap={4} style={{ flex: 1, minWidth: 0 }}>
                          {isEditing ? (
                            <Flex gap={4}>
                              <Input
                                size="small"
                                value={editingTitle}
                                onChange={(e) => setEditingTitle(e.target.value)}
                                onPressEnter={() => renameMutation.mutate({ id: conv.id, title: editingTitle })}
                                autoFocus
                                style={{ flex: 1 }}
                              />
                              <Button
                                size="small"
                                type="text"
                                icon={<CheckOutlined />}
                                onClick={(e) => {
                                  e.stopPropagation();
                                  renameMutation.mutate({ id: conv.id, title: editingTitle });
                                }}
                              />
                              <Button
                                size="small"
                                type="text"
                                icon={<CloseOutlined />}
                                onClick={(e) => {
                                  e.stopPropagation();
                                  setEditingTitleId(null);
                                }}
                              />
                            </Flex>
                          ) : (
                            <Text ellipsis style={{ fontSize: 13, fontWeight: isActive ? 600 : 400 }}>
                              {conv.title}
                            </Text>
                          )}
                          <Flex align="center" gap={6}>
                            <Tag
                              color={modeInfo.color}
                              style={{ fontSize: 10, lineHeight: '16px', margin: 0, padding: '0 4px' }}
                            >
                              {modeInfo.label}
                            </Tag>
                            <Text type="secondary" style={{ fontSize: 10 }}>
                              {new Date(conv.updated_at).toLocaleDateString()}
                            </Text>
                          </Flex>
                        </Flex>
                        {isActive && !isEditing && (
                          <Flex gap={2}>
                            <Button
                              size="small"
                              type="text"
                              icon={<EditOutlined style={{ fontSize: 12 }} />}
                              onClick={(e) => {
                                e.stopPropagation();
                                setEditingTitleId(conv.id);
                                setEditingTitle(conv.title);
                              }}
                            />
                            <Popconfirm
                              title="确定删除这个对话？"
                              onConfirm={(e) => {
                                e?.stopPropagation();
                                deleteMutation.mutate(conv.id);
                              }}
                              onCancel={(e) => e?.stopPropagation()}
                            >
                              <Button
                                size="small"
                                type="text"
                                danger
                                icon={<DeleteOutlined style={{ fontSize: 12 }} />}
                                onClick={(e) => e.stopPropagation()}
                              />
                            </Popconfirm>
                          </Flex>
                        )}
                      </Flex>
                    </div>
                  );
                }}
              />
            )}
          </div>
        </Flex>
      </Sider>

      {/* 右侧：聊天区域 */}
      <Content style={{ display: 'flex', flexDirection: 'column', background: token.colorBgLayout }}>
        {!activeConvId ? (
          // 空状态 - 欢迎界面
          <Flex
            vertical
            align="center"
            justify="center"
            gap={32}
            style={{ flex: 1, padding: 40 }}
          >
            <Flex vertical align="center" gap={8}>
              <div style={{
                width: 64,
                height: 64,
                borderRadius: 16,
                background: 'linear-gradient(135deg, #4f46e5, #7c3aed)',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
                fontSize: 28,
                color: '#fff',
              }}>
                <MessageOutlined />
              </div>
              <Text strong style={{ fontSize: 20, marginTop: 8 }}>AI 灵感工坊</Text>
              <Text type="secondary">选择一种模式开始对话，让 AI 为你的创作提供灵感</Text>
            </Flex>

            <Flex gap={16} wrap="wrap" justify="center" style={{ maxWidth: 640 }}>
              {CHAT_MODES.map(mode => (
                <div
                  key={mode.key}
                  onClick={() => createMutation.mutate(mode.key)}
                  style={{
                    width: 140,
                    padding: '20px 16px',
                    borderRadius: 12,
                    background: token.colorBgContainer,
                    border: `1px solid ${token.colorBorderSecondary}`,
                    cursor: 'pointer',
                    textAlign: 'center',
                    transition: 'all 0.2s',
                  }}
                  onMouseEnter={(e) => {
                    e.currentTarget.style.borderColor = mode.color;
                    e.currentTarget.style.transform = 'translateY(-2px)';
                    e.currentTarget.style.boxShadow = `0 4px 12px ${mode.color}22`;
                  }}
                  onMouseLeave={(e) => {
                    e.currentTarget.style.borderColor = token.colorBorderSecondary;
                    e.currentTarget.style.transform = 'none';
                    e.currentTarget.style.boxShadow = 'none';
                  }}
                >
                  <div style={{ fontSize: 28, color: mode.color, marginBottom: 8 }}>{mode.icon}</div>
                  <Text strong style={{ fontSize: 14 }}>{mode.label}</Text>
                  <br />
                  <Text type="secondary" style={{ fontSize: 11 }}>{mode.desc}</Text>
                </div>
              ))}
            </Flex>
          </Flex>
        ) : (
          <>
            {/* 消息头部 */}
            {activeConv && (
              <Flex
                align="center"
                gap={8}
                style={{
                  padding: '12px 24px',
                  borderBottom: `1px solid ${token.colorBorderSecondary}`,
                  background: token.colorBgContainer,
                }}
              >
                <Tag color={getModeInfo(activeConv.mode).color}>
                  {getModeInfo(activeConv.mode).icon} {getModeInfo(activeConv.mode).label}
                </Tag>
                <Text strong>{activeConv.title}</Text>
              </Flex>
            )}

            {/* 消息列表 */}
            <div style={{ flex: 1, overflow: 'auto', padding: '16px 0' }}>
              {loadingDetail ? (
                <Flex justify="center" style={{ padding: 60 }}><Spin /></Flex>
              ) : messages.length === 0 && !streamingContent ? (
                <Flex vertical align="center" justify="center" gap={12} style={{ padding: 60 }}>
                  <Text type="secondary" style={{ fontSize: 14 }}>
                    开始提问吧！试试：
                  </Text>
                  <Flex vertical gap={8} style={{ maxWidth: 400 }}>
                    {getQuickPrompts(activeConv?.mode || 'general').map((prompt, i) => (
                      <div
                        key={i}
                        onClick={() => {
                          setInputValue(prompt);
                          inputRef.current?.focus();
                        }}
                        style={{
                          padding: '10px 14px',
                          borderRadius: 8,
                          border: `1px solid ${token.colorBorderSecondary}`,
                          background: token.colorBgContainer,
                          cursor: 'pointer',
                          fontSize: 13,
                          transition: 'all 0.2s',
                        }}
                        onMouseEnter={(e) => {
                          e.currentTarget.style.borderColor = token.colorPrimary;
                        }}
                        onMouseLeave={(e) => {
                          e.currentTarget.style.borderColor = token.colorBorderSecondary;
                        }}
                      >
                        {prompt}
                      </div>
                    ))}
                  </Flex>
                </Flex>
              ) : (
                <div style={{ maxWidth: 800, margin: '0 auto', padding: '0 24px' }}>
                  {messages.map((msg: ChatMessage) => (
                    <MessageBubble key={msg.id} message={msg} token={token} />
                  ))}
                  {streamingContent && (
                    <MessageBubble
                      message={{
                        id: 'streaming',
                        role: 'assistant',
                        content: streamingContent,
                        created_at: new Date().toISOString(),
                      }}
                      token={token}
                      isStreaming
                    />
                  )}
                  <div ref={messagesEndRef} />
                </div>
              )}
            </div>

            {/* 输入区域 */}
            <div style={{
              padding: '12px 24px 16px',
              borderTop: `1px solid ${token.colorBorderSecondary}`,
              background: token.colorBgContainer,
            }}>
              <div style={{ maxWidth: 800, margin: '0 auto' }}>
                <Flex gap={8} align="end">
                  <TextArea
                    ref={inputRef}
                    value={inputValue}
                    onChange={(e) => setInputValue(e.target.value)}
                    onKeyDown={handleKeyDown}
                    placeholder={sending ? 'AI 正在思考...' : '输入你的问题，按 Enter 发送，Shift+Enter 换行'}
                    autoSize={{ minRows: 1, maxRows: 5 }}
                    disabled={sending}
                    style={{
                      flex: 1,
                      borderRadius: 12,
                      padding: '10px 14px',
                      fontSize: 14,
                      resize: 'none',
                    }}
                  />
                  <Button
                    type="primary"
                    icon={<SendOutlined />}
                    onClick={handleSend}
                    loading={sending}
                    disabled={!inputValue.trim()}
                    style={{ borderRadius: 12, height: 42, width: 42 }}
                  />
                </Flex>
                <Text type="secondary" style={{ fontSize: 11, marginTop: 4, display: 'block', textAlign: 'center' }}>
                  AI 生成内容仅供参考，请结合自身创作判断
                </Text>
              </div>
            </div>
          </>
        )}
      </Content>
    </Layout>
  );
}

function MessageBubble({
  message,
  token,
  isStreaming = false,
}: {
  message: ChatMessage;
  token: any;
  isStreaming?: boolean;
}) {
  const isUser = message.role === 'user';

  return (
    <Flex
      gap={12}
      style={{ marginBottom: 20 }}
      align="start"
    >
      <Avatar
        size={32}
        icon={isUser ? <UserOutlined /> : <RobotOutlined />}
        style={{
          background: isUser
            ? 'linear-gradient(135deg, #818cf8, #a855f7)'
            : 'linear-gradient(135deg, #4f46e5, #7c3aed)',
          flexShrink: 0,
          marginTop: 2,
        }}
      />
      <div style={{
        flex: 1,
        padding: '10px 16px',
        borderRadius: 12,
        background: isUser ? token.colorPrimaryBg : token.colorBgContainer,
        border: isUser ? 'none' : `1px solid ${token.colorBorderSecondary}`,
        minWidth: 0,
      }}>
        {isUser ? (
          <Paragraph
            style={{
              marginBottom: 0,
              whiteSpace: 'pre-wrap',
              wordBreak: 'break-word',
              fontSize: 14,
              lineHeight: 1.7,
            }}
          >
            {message.content}
          </Paragraph>
        ) : (
          <div className="chat-markdown" style={{ fontSize: 14, lineHeight: 1.7 }}>
            <ReactMarkdown remarkPlugins={[remarkGfm]}>
              {message.content}
            </ReactMarkdown>
            {isStreaming && <span style={{ animation: 'blink 1s steps(1) infinite' }}>▊</span>}
          </div>
        )}
      </div>
    </Flex>
  );
}

function getQuickPrompts(mode: string): string[] {
  switch (mode) {
    case 'creative':
      return [
        '我写了一个穿越题材的故事，但感觉情节太老套了，能给点新思路吗？',
        '帮我构思一个反派角色翻盘的情节转折',
        '如何让一个普通的校园故事变得更有吸引力？',
      ];
    case 'building':
      return [
        '帮我设计一个玄幻世界的修炼体系',
        '如何让科幻设定既有想象力又合理？',
        '一个末日后的城市应该有怎样的社会结构？',
      ];
    case 'character':
      return [
        '帮我塑造一个亦正亦邪的角色',
        '如何让配角也变得生动有趣？',
        '怎样设计角色间的关系让故事更有张力？',
      ];
    default:
      return [
        '帮我分析一下我的故事大纲有什么问题',
        '写小说的节奏该怎么把控？',
        '如何写出让读者沉浸的第一章？',
      ];
  }
}

export default Chat;
