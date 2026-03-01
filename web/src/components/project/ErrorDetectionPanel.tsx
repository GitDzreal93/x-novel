import { useState } from 'react';
import {
  Button, Card, Tag, Typography, Flex, Space, Checkbox, Spin, Empty,
  App, theme, Badge,
} from 'antd';
import {
  BugOutlined, WarningOutlined, InfoCircleOutlined,
  CheckCircleOutlined, SyncOutlined, CopyOutlined,
} from '@ant-design/icons';
import { reviewApi } from '../../api';
import type { DetectionResult, DetectionIssue, DetectionType, IssueSeverity } from '../../types';

const { Text, Paragraph } = Typography;

export interface ErrorDetectionPanelProps {
  content: string;
  disabled?: boolean;
  onApplySuggestion?: (original: string, suggestion: string) => void;
  onErrorsDetected?: (issues: DetectionIssue[]) => void;
}

const TYPE_LABELS: Record<DetectionType, { label: string; color: string }> = {
  typo: { label: '错别字', color: 'red' },
  grammar: { label: '语法', color: 'orange' },
  logic: { label: '逻辑', color: 'purple' },
  repetition: { label: '重复', color: 'blue' },
};

const SEVERITY_CONFIG: Record<IssueSeverity, { icon: React.ReactNode; color: string }> = {
  error: { icon: <BugOutlined />, color: '#ff4d4f' },
  warning: { icon: <WarningOutlined />, color: '#faad14' },
  info: { icon: <InfoCircleOutlined />, color: '#1677ff' },
};

function ErrorDetectionPanel({ content, disabled, onApplySuggestion, onErrorsDetected }: ErrorDetectionPanelProps) {
  const { message: messageApi } = App.useApp();
  const { token } = theme.useToken();
  const [loading, setLoading] = useState(false);
  const [result, setResult] = useState<DetectionResult | null>(null);
  const [selectedTypes, setSelectedTypes] = useState<DetectionType[]>(['typo', 'grammar', 'logic', 'repetition']);

  const handleDetect = async () => {
    if (!content?.trim()) {
      messageApi.warning('章节内容为空，无法检测');
      return;
    }

    setLoading(true);
    try {
      const res = await reviewApi.detect({ content, types: selectedTypes });
      if (res?.data) {
        setResult(res.data);
        onErrorsDetected?.(res.data.issues);
      }
    } catch {
      messageApi.error('检测失败，请稍后重试');
    } finally {
      setLoading(false);
    }
  };

  const handleApply = (issue: DetectionIssue) => {
    if (onApplySuggestion) {
      onApplySuggestion(issue.original, issue.suggestion);
      messageApi.success('已应用修改建议');
    }
  };

  const typeOptions = Object.entries(TYPE_LABELS).map(([key, val]) => ({
    label: val.label,
    value: key,
  }));

  return (
    <Flex vertical gap={12} style={{ height: '100%' }}>
      <Text strong style={{ fontSize: 14 }}>
        <BugOutlined style={{ marginRight: 6, color: token.colorPrimary }} />
        错误检测
      </Text>

      <Checkbox.Group
        options={typeOptions}
        value={selectedTypes}
        onChange={(vals) => setSelectedTypes(vals as DetectionType[])}
        style={{ fontSize: 13 }}
      />

      <Button
        type="primary"
        icon={<SyncOutlined />}
        onClick={handleDetect}
        loading={loading}
        disabled={disabled || !content?.trim()}
        block
      >
        开始检测
      </Button>

      {loading && (
        <Flex justify="center" align="center" style={{ flex: 1 }}>
          <Spin tip="AI 正在检测中..." />
        </Flex>
      )}

      {!loading && result && (
        <Flex vertical gap={8} style={{ flex: 1, overflow: 'auto' }}>
          {/* 概要 */}
          <Card size="small" style={{ background: token.colorBgLayout, flexShrink: 0 }}>
            <Flex justify="space-between" align="center">
              <Space size={12}>
                {result.total_count === 0 ? (
                  <Tag icon={<CheckCircleOutlined />} color="success">无问题</Tag>
                ) : (
                  <Tag color="warning">发现 {result.total_count} 个问题</Tag>
                )}
              </Space>
              <Space size={4}>
                {Object.entries(result.type_counts || {}).map(([type, count]) => (
                  <Badge
                    key={type}
                    count={count}
                    size="small"
                    color={TYPE_LABELS[type as DetectionType]?.color}
                    title={TYPE_LABELS[type as DetectionType]?.label}
                  />
                ))}
              </Space>
            </Flex>
            {result.summary && (
              <Text type="secondary" style={{ fontSize: 12, display: 'block', marginTop: 6 }}>
                {result.summary}
              </Text>
            )}
          </Card>

          {/* 问题列表 */}
          {result.issues.length === 0 ? (
            <Flex justify="center" align="center" style={{ flex: 1 }}>
              <Empty
                image={Empty.PRESENTED_IMAGE_SIMPLE}
                description="恭喜！未发现任何问题"
              />
            </Flex>
          ) : (
            result.issues.map((issue, idx) => (
              <IssueCard
                key={idx}
                issue={issue}
                token={token}
                onApply={() => handleApply(issue)}
              />
            ))
          )}
        </Flex>
      )}

      {!loading && !result && (
        <Flex justify="center" align="center" style={{ flex: 1 }}>
          <Empty
            image={Empty.PRESENTED_IMAGE_SIMPLE}
            description="点击检测按钮开始分析"
          />
        </Flex>
      )}
    </Flex>
  );
}

function IssueCard({
  issue,
  token,
  onApply,
}: {
  issue: DetectionIssue;
  token: any;
  onApply: () => void;
}) {
  const severityCfg = SEVERITY_CONFIG[issue.severity];
  const typeCfg = TYPE_LABELS[issue.type];

  return (
    <Card
      size="small"
      style={{
        borderLeft: `3px solid ${severityCfg.color}`,
        flexShrink: 0,
      }}
      styles={{ body: { padding: '10px 12px' } }}
    >
      <Flex justify="space-between" align="flex-start" style={{ marginBottom: 6 }}>
        <Space size={6}>
          <Tag
            color={typeCfg.color}
            bordered={false}
            style={{ fontSize: 11, lineHeight: '18px', padding: '0 6px' }}
          >
            {typeCfg.label}
          </Tag>
          <Tag
            bordered={false}
            style={{
              fontSize: 11,
              lineHeight: '18px',
              padding: '0 6px',
              color: severityCfg.color,
              background: `${severityCfg.color}15`,
            }}
          >
            {issue.severity === 'error' ? '严重' : issue.severity === 'warning' ? '建议' : '提示'}
          </Tag>
        </Space>
        <Text type="secondary" style={{ fontSize: 11 }}>{issue.position}</Text>
      </Flex>

      <div style={{
        background: token.colorErrorBg,
        padding: '4px 8px',
        borderRadius: 4,
        marginBottom: 6,
        fontSize: 13,
      }}>
        <Text delete type="danger" style={{ fontSize: 13 }}>{issue.original}</Text>
      </div>

      {issue.suggestion && (
        <div style={{
          background: token.colorSuccessBg,
          padding: '4px 8px',
          borderRadius: 4,
          marginBottom: 6,
          fontSize: 13,
        }}>
          <Text type="success" style={{ fontSize: 13 }}>{issue.suggestion}</Text>
        </div>
      )}

      <Paragraph
        type="secondary"
        style={{ fontSize: 12, margin: 0, marginBottom: 6 }}
        ellipsis={{ rows: 2, expandable: true, symbol: '展开' }}
      >
        {issue.explanation}
      </Paragraph>

      <Flex justify="flex-end" gap={6}>
        <Button
          size="small"
          type="text"
          icon={<CopyOutlined />}
          onClick={() => {
            navigator.clipboard.writeText(issue.suggestion);
          }}
          style={{ fontSize: 12 }}
        >
          复制
        </Button>
        <Button
          size="small"
          type="link"
          onClick={onApply}
          style={{ fontSize: 12 }}
        >
          应用修改
        </Button>
      </Flex>
    </Card>
  );
}

export default ErrorDetectionPanel;
