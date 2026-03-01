import { useState, useRef } from 'react';
import {
  Button, Tabs, Select, InputNumber, Flex, Typography, Space, App, theme, Tooltip,
} from 'antd';
import {
  HighlightOutlined,
  EditOutlined,
  BulbOutlined,
  CopyOutlined,
  CheckOutlined,
  SwapOutlined,
} from '@ant-design/icons';
import ReactMarkdown from 'react-markdown';
import remarkGfm from 'remark-gfm';
import type { PolishStyle, SuggestionAspect } from '../../types';

const { Text } = Typography;

interface WritingAssistantProps {
  content: string;
  projectId: string;
  onApplyResult: (result: string) => void;
  onAppendResult: (result: string) => void;
  disabled?: boolean;
}

const POLISH_STYLES: { value: PolishStyle; label: string }[] = [
  { value: 'vivid', label: '生动形象' },
  { value: 'concise', label: '精炼简洁' },
  { value: 'literary', label: '文学性' },
  { value: 'dramatic', label: '戏剧张力' },
];

const SUGGESTION_ASPECTS: { value: SuggestionAspect; label: string }[] = [
  { value: 'plot', label: '情节发展' },
  { value: 'character', label: '角色刻画' },
  { value: 'dialogue', label: '对话设计' },
  { value: 'description', label: '环境描写' },
  { value: 'conflict', label: '冲突悬念' },
];

function WritingAssistant({ content, projectId, onApplyResult, onAppendResult, disabled }: WritingAssistantProps) {
  const { token } = theme.useToken();
  const { message: messageApi } = App.useApp();

  const [activeTab, setActiveTab] = useState('polish');
  const [loading, setLoading] = useState(false);
  const [result, setResult] = useState('');
  const [streamingContent, setStreamingContent] = useState('');
  const [copied, setCopied] = useState(false);

  // 润色配置
  const [polishStyle, setPolishStyle] = useState<PolishStyle>('vivid');
  // 续写配置
  const [targetWords, setTargetWords] = useState(500);
  // 建议配置
  const [suggestionAspect, setSuggestionAspect] = useState<SuggestionAspect>('plot');

  const abortRef = useRef<AbortController | null>(null);

  const handleExecute = async () => {
    if (!content.trim()) {
      messageApi.warning('请先在编辑器中输入内容');
      return;
    }

    setLoading(true);
    setResult('');
    setStreamingContent('');

    const controller = new AbortController();
    abortRef.current = controller;

    try {
      const deviceId = localStorage.getItem('x-novel-device-id') || '';

      const body: Record<string, unknown> = {
        action: activeTab,
        content: content.trim(),
        project_id: projectId,
        stream: true,
      };

      if (activeTab === 'polish') body.style = polishStyle;
      if (activeTab === 'continue') body.target_words = targetWords;
      if (activeTab === 'suggestion') body.aspect = suggestionAspect;

      const response = await fetch('/api/v1/writing/assist', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'X-Device-ID': deviceId,
        },
        body: JSON.stringify(body),
        signal: controller.signal,
      });

      if (!response.ok) throw new Error(`HTTP ${response.status}`);

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
                setResult(parsed.result || accumulated);
                setStreamingContent('');
              } else if (parsed.error) {
                messageApi.error(parsed.error);
              } else if (parsed.content !== undefined) {
                accumulated += parsed.content;
                setStreamingContent(accumulated);
              }
            } catch {
              // ignore
            }
          }
        }
      }

      if (accumulated && !result) {
        setResult(accumulated);
        setStreamingContent('');
      }
    } catch (err: any) {
      if (err.name !== 'AbortError') {
        messageApi.error('处理失败：' + (err.message || '未知错误'));
      }
    } finally {
      setLoading(false);
      abortRef.current = null;
    }
  };

  const handleStop = () => {
    abortRef.current?.abort();
    if (streamingContent) {
      setResult(streamingContent);
      setStreamingContent('');
    }
    setLoading(false);
  };

  const handleCopy = () => {
    navigator.clipboard.writeText(result);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  const displayContent = streamingContent || result;
  const actionLabel = activeTab === 'polish' ? '润色' : activeTab === 'continue' ? '续写' : '建议';

  const tabItems = [
    {
      key: 'polish',
      label: (
        <Flex align="center" gap={4}>
          <HighlightOutlined />
          <span>润色</span>
        </Flex>
      ),
      children: (
        <Flex vertical gap={8}>
          <Text type="secondary" style={{ fontSize: 12 }}>选择润色风格，AI 将对编辑器中的内容进行优化</Text>
          <Select
            value={polishStyle}
            onChange={setPolishStyle}
            options={POLISH_STYLES}
            style={{ width: '100%' }}
            size="small"
          />
        </Flex>
      ),
    },
    {
      key: 'continue',
      label: (
        <Flex align="center" gap={4}>
          <EditOutlined />
          <span>续写</span>
        </Flex>
      ),
      children: (
        <Flex vertical gap={8}>
          <Text type="secondary" style={{ fontSize: 12 }}>基于编辑器内容继续写作</Text>
          <Flex align="center" gap={8}>
            <Text style={{ fontSize: 12, whiteSpace: 'nowrap' }}>目标字数</Text>
            <InputNumber
              value={targetWords}
              onChange={(v) => setTargetWords(v || 500)}
              min={100}
              max={3000}
              step={100}
              size="small"
              style={{ width: 100 }}
            />
          </Flex>
        </Flex>
      ),
    },
    {
      key: 'suggestion',
      label: (
        <Flex align="center" gap={4}>
          <BulbOutlined />
          <span>灵感</span>
        </Flex>
      ),
      children: (
        <Flex vertical gap={8}>
          <Text type="secondary" style={{ fontSize: 12 }}>基于当前内容提供创作灵感和建议</Text>
          <Select
            value={suggestionAspect}
            onChange={setSuggestionAspect}
            options={SUGGESTION_ASPECTS}
            style={{ width: '100%' }}
            size="small"
          />
        </Flex>
      ),
    },
  ];

  return (
    <Flex vertical style={{ height: '100%' }}>
      <Tabs
        activeKey={activeTab}
        onChange={setActiveTab}
        items={tabItems}
        size="small"
        style={{ marginBottom: 0 }}
      />

      <Button
        type="primary"
        onClick={loading ? handleStop : handleExecute}
        loading={loading}
        disabled={disabled || !content.trim()}
        block
        style={{ marginTop: 8, marginBottom: 12 }}
      >
        {loading ? '停止' : `AI ${actionLabel}`}
      </Button>

      {/* 结果区域 */}
      {displayContent && (
        <Flex vertical style={{ flex: 1, minHeight: 0 }}>
          <Flex justify="space-between" align="center" style={{ marginBottom: 6 }}>
            <Text strong style={{ fontSize: 12 }}>
              {streamingContent ? '生成中...' : `${actionLabel}结果`}
            </Text>
            {result && !streamingContent && (
              <Space size={4}>
                <Tooltip title="复制结果">
                  <Button
                    type="text"
                    size="small"
                    icon={copied ? <CheckOutlined style={{ color: token.colorSuccess }} /> : <CopyOutlined />}
                    onClick={handleCopy}
                  />
                </Tooltip>
                {activeTab === 'polish' && (
                  <Tooltip title="替换编辑器内容">
                    <Button
                      type="text"
                      size="small"
                      icon={<SwapOutlined />}
                      onClick={() => {
                        onApplyResult(result);
                        messageApi.success('已替换');
                      }}
                    />
                  </Tooltip>
                )}
                {activeTab === 'continue' && (
                  <Tooltip title="追加到编辑器末尾">
                    <Button
                      type="text"
                      size="small"
                      icon={<EditOutlined />}
                      onClick={() => {
                        onAppendResult(result);
                        messageApi.success('已追加');
                      }}
                    />
                  </Tooltip>
                )}
              </Space>
            )}
          </Flex>
          <div
            style={{
              flex: 1,
              overflow: 'auto',
              padding: 10,
              borderRadius: 8,
              border: `1px solid ${token.colorBorderSecondary}`,
              background: token.colorBgLayout,
              fontSize: 13,
              lineHeight: 1.7,
            }}
          >
            <div className="chat-markdown">
              <ReactMarkdown remarkPlugins={[remarkGfm]}>
                {displayContent}
              </ReactMarkdown>
              {streamingContent && <span style={{ animation: 'blink 1s steps(1) infinite' }}>▊</span>}
            </div>
          </div>
        </Flex>
      )}
    </Flex>
  );
}

export default WritingAssistant;
