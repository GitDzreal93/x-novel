import { useState } from 'react';
import {
  Button, Card, Typography, Flex, Space, Spin, Empty, Select, App,
  theme, Progress, Tag, List, Segmented,
} from 'antd';
import {
  AuditOutlined, StarOutlined, BulbOutlined,
  WarningOutlined, ThunderboltOutlined,
} from '@ant-design/icons';
import ReactECharts from 'echarts-for-react';
import { reviewApi, chapterApi } from '../../api';
import { useQuery } from '@tanstack/react-query';
import type { Project, ReviewResult, ReviewDimension } from '../../types';

const { Text, Title, Paragraph } = Typography;

interface ReviewPanelProps {
  project: Project;
}

const DIMENSION_LABELS: Record<ReviewDimension, string> = {
  plot: '情节',
  character: '人物',
  writing: '文笔',
  pacing: '节奏',
  creativity: '创意',
  readability: '可读性',
};

const SCORE_COLORS: [number, string][] = [
  [60, '#ff4d4f'],
  [70, '#faad14'],
  [80, '#52c41a'],
  [90, '#1677ff'],
  [100, '#722ed1'],
];

function getScoreColor(score: number): string {
  for (const [threshold, color] of SCORE_COLORS) {
    if (score <= threshold) return color;
  }
  return '#722ed1';
}

function ReviewPanel({ project }: ReviewPanelProps) {
  const { message: messageApi } = App.useApp();
  const { token } = theme.useToken();
  const [loading, setLoading] = useState(false);
  const [result, setResult] = useState<ReviewResult | null>(null);
  const [reviewMode, setReviewMode] = useState<'project' | 'chapter'>('project');
  const [selectedChapter, setSelectedChapter] = useState<number | undefined>(undefined);

  const { data: chaptersData } = useQuery({
    queryKey: ['chapters', project.id, 'all'],
    queryFn: () =>
      chapterApi.list(project.id, { page: 1, page_size: 100 }).then((res) => {
        if (!res?.data) throw new Error('获取章节失败');
        return res.data;
      }),
    enabled: !!project.id,
  });

  const chapters = chaptersData?.chapters || [];
  const chapterOptions = chapters
    .filter((c) => c.content && c.word_count > 0)
    .map((c) => ({
      value: c.chapter_number,
      label: `第${c.chapter_number}章${c.title ? ` ${c.title}` : ''}`,
    }));

  const handleReview = async () => {
    setLoading(true);
    try {
      let res;
      if (reviewMode === 'project') {
        res = await reviewApi.reviewProject(project.id);
      } else {
        if (!selectedChapter) {
          messageApi.warning('请选择要审阅的章节');
          setLoading(false);
          return;
        }
        res = await reviewApi.reviewChapter(project.id, selectedChapter);
      }
      if (res?.data) {
        setResult(res.data);
      }
    } catch {
      messageApi.error('审阅失败，请稍后重试');
    } finally {
      setLoading(false);
    }
  };

  const getRadarOption = () => {
    if (!result) return {};

    const indicator = result.scores.map((s) => ({
      name: DIMENSION_LABELS[s.dimension] || s.dimension,
      max: 100,
    }));

    const values = result.scores.map((s) => s.score);

    return {
      radar: {
        indicator,
        shape: 'polygon',
        radius: '65%',
        splitNumber: 5,
        axisName: {
          color: token.colorText,
          fontSize: 13,
        },
        splitArea: {
          areaStyle: {
            color: [
              'rgba(22, 119, 255, 0.02)',
              'rgba(22, 119, 255, 0.04)',
              'rgba(22, 119, 255, 0.06)',
              'rgba(22, 119, 255, 0.08)',
              'rgba(22, 119, 255, 0.10)',
            ],
          },
        },
        splitLine: {
          lineStyle: { color: token.colorBorderSecondary },
        },
        axisLine: {
          lineStyle: { color: token.colorBorderSecondary },
        },
      },
      series: [
        {
          type: 'radar',
          data: [
            {
              value: values,
              name: '评分',
              areaStyle: {
                color: 'rgba(22, 119, 255, 0.15)',
              },
              lineStyle: {
                color: '#1677ff',
                width: 2,
              },
              itemStyle: {
                color: '#1677ff',
              },
            },
          ],
        },
      ],
      tooltip: {
        trigger: 'item',
      },
    };
  };

  return (
    <div style={{ padding: '0 24px 24px' }}>
      <Flex justify="space-between" align="center" wrap gap={16} style={{ marginBottom: 20 }}>
        <div>
          <Flex align="center" gap={8}>
            <Title level={5} style={{ margin: 0 }}>AI 审阅</Title>
            <AuditOutlined style={{ color: token.colorPrimary, fontSize: 16 }} />
          </Flex>
          <Text type="secondary" style={{ fontSize: 13 }}>
            AI 将从多维度评估作品质量，提供专业修改建议
          </Text>
        </div>
      </Flex>

      {/* 审阅配置 */}
      <Card style={{ marginBottom: 20 }}>
        <Flex vertical gap={16}>
          <Segmented
            value={reviewMode}
            onChange={(val) => setReviewMode(val as 'project' | 'chapter')}
            options={[
              { label: '整体审阅', value: 'project' },
              { label: '章节审阅', value: 'chapter' },
            ]}
            block
          />

          {reviewMode === 'chapter' && (
            <Select
              placeholder="选择要审阅的章节"
              value={selectedChapter}
              onChange={setSelectedChapter}
              options={chapterOptions}
              style={{ width: '100%' }}
              size="large"
              notFoundContent={<Empty image={Empty.PRESENTED_IMAGE_SIMPLE} description="暂无可审阅的章节" />}
            />
          )}

          <Button
            type="primary"
            icon={<ThunderboltOutlined />}
            onClick={handleReview}
            loading={loading}
            size="large"
            block
            disabled={reviewMode === 'chapter' && !selectedChapter}
          >
            {loading ? 'AI 审阅中...' : '开始审阅'}
          </Button>
        </Flex>
      </Card>

      {loading && (
        <Flex justify="center" align="center" style={{ height: 300 }}>
          <Spin size="large" tip="AI 正在认真审阅作品..." />
        </Flex>
      )}

      {!loading && result && (
        <Flex vertical gap={20}>
          {/* 总分 + 雷达图 */}
          <Flex gap={20} wrap>
            <Card style={{ flex: 1, minWidth: 280, textAlign: 'center' }}>
              <Flex vertical align="center" gap={12}>
                <Text type="secondary" style={{ fontSize: 13 }}>综合评分</Text>
                <Progress
                  type="dashboard"
                  percent={result.overall_score}
                  size={160}
                  strokeColor={getScoreColor(result.overall_score)}
                  format={(percent) => (
                    <Flex vertical align="center">
                      <span style={{ fontSize: 36, fontWeight: 700, color: getScoreColor(percent || 0) }}>
                        {percent}
                      </span>
                      <Text type="secondary" style={{ fontSize: 12 }}>
                        {(percent || 0) >= 90 ? '优秀' : (percent || 0) >= 80 ? '良好' : (percent || 0) >= 70 ? '中等' : '待提升'}
                      </Text>
                    </Flex>
                  )}
                />
              </Flex>
            </Card>

            <Card title="多维度评分" style={{ flex: 2, minWidth: 350 }}>
              <ReactECharts
                option={getRadarOption()}
                style={{ height: 260 }}
                opts={{ renderer: 'svg' }}
              />
            </Card>
          </Flex>

          {/* 各维度详情 */}
          <Card title="维度详情">
            <Flex vertical gap={12}>
              {result.scores.map((score) => (
                <Flex
                  key={score.dimension}
                  justify="space-between"
                  align="center"
                  style={{
                    padding: '10px 12px',
                    background: token.colorBgLayout,
                    borderRadius: token.borderRadius,
                  }}
                >
                  <Flex align="center" gap={12} style={{ flex: 1 }}>
                    <Tag
                      color={getScoreColor(score.score)}
                      bordered={false}
                      style={{ minWidth: 60, textAlign: 'center' }}
                    >
                      {DIMENSION_LABELS[score.dimension]}
                    </Tag>
                    <Progress
                      percent={score.score}
                      size="small"
                      strokeColor={getScoreColor(score.score)}
                      style={{ flex: 1, maxWidth: 200 }}
                      format={(p) => <Text strong>{p}</Text>}
                    />
                  </Flex>
                  <Text type="secondary" style={{ fontSize: 13, maxWidth: 300 }} ellipsis>
                    {score.comment}
                  </Text>
                </Flex>
              ))}
            </Flex>
          </Card>

          {/* 亮点 */}
          {result.highlights.length > 0 && (
            <Card
              title={
                <Flex align="center" gap={6}>
                  <StarOutlined style={{ color: '#faad14' }} />
                  <span>亮点</span>
                </Flex>
              }
            >
              <List
                dataSource={result.highlights}
                renderItem={(item) => (
                  <List.Item style={{ padding: '8px 0' }}>
                    <Flex align="flex-start" gap={8}>
                      <Tag color="gold" bordered={false} style={{ flexShrink: 0 }}>亮点</Tag>
                      <Text>{item}</Text>
                    </Flex>
                  </List.Item>
                )}
              />
            </Card>
          )}

          {/* 问题 */}
          {result.issues.length > 0 && (
            <Card
              title={
                <Flex align="center" gap={6}>
                  <WarningOutlined style={{ color: '#ff4d4f' }} />
                  <span>问题诊断</span>
                </Flex>
              }
            >
              <List
                dataSource={result.issues}
                renderItem={(item) => (
                  <List.Item style={{ padding: '8px 0' }}>
                    <Flex align="flex-start" gap={8}>
                      <Tag color="red" bordered={false} style={{ flexShrink: 0 }}>问题</Tag>
                      <Text>{item}</Text>
                    </Flex>
                  </List.Item>
                )}
              />
            </Card>
          )}

          {/* 修改建议 */}
          {result.suggestions.length > 0 && (
            <Card
              title={
                <Flex align="center" gap={6}>
                  <BulbOutlined style={{ color: '#1677ff' }} />
                  <span>修改建议</span>
                </Flex>
              }
            >
              <List
                dataSource={result.suggestions}
                renderItem={(item, idx) => (
                  <List.Item style={{ padding: '8px 0' }}>
                    <Flex align="flex-start" gap={8}>
                      <Tag color="blue" bordered={false} style={{ flexShrink: 0 }}>建议 {idx + 1}</Tag>
                      <Text>{item}</Text>
                    </Flex>
                  </List.Item>
                )}
              />
            </Card>
          )}

          {/* 总评 */}
          <Card
            title="总体评价"
            style={{ borderColor: token.colorPrimary }}
          >
            <Paragraph style={{ fontSize: 15, lineHeight: 1.8, margin: 0 }}>
              {result.summary}
            </Paragraph>
          </Card>
        </Flex>
      )}

      {!loading && !result && (
        <Card style={{ textAlign: 'center', padding: '60px 0' }}>
          <Empty
            image={Empty.PRESENTED_IMAGE_SIMPLE}
            description="选择审阅模式后点击开始审阅"
          />
        </Card>
      )}
    </div>
  );
}

export default ReviewPanel;
