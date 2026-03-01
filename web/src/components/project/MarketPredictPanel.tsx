import { useState } from 'react';
import {
  Button, Card, Typography, Flex, Space, Spin, Empty, App, theme,
  Progress, Tag, List, Statistic, Row, Col, Descriptions,
} from 'antd';
import {
  StockOutlined, RiseOutlined, TeamOutlined, DollarOutlined,
  WarningOutlined, BulbOutlined, ThunderboltOutlined,
} from '@ant-design/icons';
import ReactECharts from 'echarts-for-react';
import { reviewApi } from '../../api';
import type { Project, MarketPrediction } from '../../types';

const { Text, Paragraph, Title } = Typography;

interface MarketPredictPanelProps {
  project: Project;
}

function MarketPredictPanel({ project }: MarketPredictPanelProps) {
  const { message } = App.useApp();
  const { token } = theme.useToken();
  const [loading, setLoading] = useState(false);
  const [result, setResult] = useState<MarketPrediction | null>(null);

  const handlePredict = async () => {
    setLoading(true);
    try {
      const res = await reviewApi.marketPredict(project.id);
      if (res?.data) {
        setResult(res.data);
      }
    } catch {
      message.error('市场预测失败');
    } finally {
      setLoading(false);
    }
  };

  const getScoreColor = (score: number) => {
    if (score >= 80) return token.colorSuccess;
    if (score >= 60) return token.colorWarning;
    return token.colorError;
  };

  const getRadarOption = () => {
    if (!result) return {};
    const indicators = result.reader_appeal.map(item => ({
      name: item.dimension,
      max: 100,
    }));
    const values = result.reader_appeal.map(item => item.score);

    return {
      radar: {
        indicator: indicators,
        shape: 'circle',
        radius: '60%',
        axisName: {
          color: token.colorText,
          fontSize: 12,
        },
        splitArea: {
          areaStyle: { color: ['rgba(24,144,255,0.02)', 'rgba(24,144,255,0.05)'] },
        },
      },
      series: [{
        type: 'radar',
        data: [{
          value: values,
          name: '读者吸引力',
          areaStyle: { color: 'rgba(24,144,255,0.2)' },
          lineStyle: { color: token.colorPrimary },
          itemStyle: { color: token.colorPrimary },
        }],
      }],
      tooltip: { trigger: 'item' },
    };
  };

  const getTrendBarOption = () => {
    if (!result) return {};
    const trends = result.market_trends;
    return {
      tooltip: { trigger: 'axis' },
      grid: { left: 20, right: 20, top: 10, bottom: 30, containLabel: true },
      xAxis: {
        type: 'category',
        data: trends.map(t => t.trend),
        axisLabel: { fontSize: 11, interval: 0, rotate: trends.length > 4 ? 15 : 0 },
      },
      yAxis: { type: 'value', max: 100, axisLabel: { fontSize: 11 } },
      series: [{
        type: 'bar',
        data: trends.map(t => ({
          value: t.fit,
          itemStyle: { color: getScoreColor(t.fit), borderRadius: [4, 4, 0, 0] },
        })),
        barWidth: '40%',
      }],
    };
  };

  return (
    <Flex vertical gap={16}>
      <div>
        <Title level={5} style={{ marginBottom: 4 }}>
          <StockOutlined style={{ marginRight: 8 }} />
          市场预测
        </Title>
        <Text type="secondary">
          基于 AI 分析作品的市场潜力、目标读者和变现建议
        </Text>
      </div>

      <Button
        type="primary"
        icon={<ThunderboltOutlined />}
        onClick={handlePredict}
        loading={loading}
        size="large"
        block
      >
        开始市场预测分析
      </Button>

      {loading && (
        <Flex justify="center" align="center" style={{ padding: 60 }}>
          <Spin tip="AI 正在分析市场数据..." size="large" />
        </Flex>
      )}

      {!loading && !result && (
        <Empty description="点击上方按钮开始市场预测" />
      )}

      {!loading && result && (
        <Flex vertical gap={16}>
          {/* 核心评分 */}
          <Card size="small">
            <Row gutter={16} align="middle">
              <Col span={8}>
                <Flex vertical align="center">
                  <Progress
                    type="dashboard"
                    percent={result.market_score}
                    size={100}
                    strokeColor={getScoreColor(result.market_score)}
                    format={(p) => <span style={{ fontSize: 20, fontWeight: 700 }}>{p}</span>}
                  />
                  <Text strong style={{ marginTop: 4 }}>市场适配度</Text>
                </Flex>
              </Col>
              <Col span={16}>
                <Space direction="vertical" style={{ width: '100%' }} size={4}>
                  <Flex align="center" gap={6}>
                    <TeamOutlined style={{ color: token.colorPrimary }} />
                    <Text strong>目标读者</Text>
                  </Flex>
                  <Text type="secondary">{result.target_audience}</Text>
                  <Flex align="center" gap={6} style={{ marginTop: 8 }}>
                    <RiseOutlined style={{ color: token.colorSuccess }} />
                    <Text strong>竞争优势</Text>
                  </Flex>
                  <Text type="secondary">{result.competitive_edge}</Text>
                </Space>
              </Col>
            </Row>
          </Card>

          {/* 读者吸引力雷达图 */}
          <Card
            size="small"
            title={<><TeamOutlined style={{ marginRight: 6 }} />读者吸引力分析</>}
          >
            <ReactECharts option={getRadarOption()} style={{ height: 280 }} />
            <List
              size="small"
              dataSource={result.reader_appeal}
              renderItem={(item) => (
                <List.Item>
                  <Flex align="center" justify="space-between" style={{ width: '100%' }}>
                    <Text style={{ width: 90 }}>{item.dimension}</Text>
                    <Progress
                      percent={item.score}
                      size="small"
                      strokeColor={getScoreColor(item.score)}
                      style={{ flex: 1, margin: '0 12px' }}
                    />
                    <Text type="secondary" style={{ fontSize: 12, width: 120, textAlign: 'right' }}>
                      {item.comment}
                    </Text>
                  </Flex>
                </List.Item>
              )}
            />
          </Card>

          {/* 市场趋势 */}
          <Card
            size="small"
            title={<><RiseOutlined style={{ marginRight: 6 }} />市场趋势契合度</>}
          >
            <ReactECharts option={getTrendBarOption()} style={{ height: 200 }} />
            <List
              size="small"
              dataSource={result.market_trends}
              renderItem={(item) => (
                <List.Item>
                  <Flex vertical style={{ width: '100%' }}>
                    <Flex justify="space-between" align="center">
                      <Text strong>{item.trend}</Text>
                      <Tag color={item.fit >= 70 ? 'green' : item.fit >= 50 ? 'orange' : 'red'}>
                        契合度 {item.fit}%
                      </Tag>
                    </Flex>
                    <Text type="secondary" style={{ fontSize: 12 }}>{item.analysis}</Text>
                  </Flex>
                </List.Item>
              )}
            />
          </Card>

          {/* 变现建议 */}
          <Card
            size="small"
            title={<><DollarOutlined style={{ marginRight: 6 }} />变现建议</>}
          >
            <Descriptions column={1} size="small" bordered>
              <Descriptions.Item label="推荐平台">
                <Space>
                  {result.monetization.platforms.map((p, i) => (
                    <Tag key={i} color="blue">{p}</Tag>
                  ))}
                </Space>
              </Descriptions.Item>
              <Descriptions.Item label="定价模式">
                {result.monetization.pricing_model}
              </Descriptions.Item>
              <Descriptions.Item label="IP 潜力">
                <Flex align="center" gap={8}>
                  <Progress
                    percent={result.monetization.ip_potential}
                    size="small"
                    strokeColor={getScoreColor(result.monetization.ip_potential)}
                    style={{ width: 120 }}
                  />
                  <Statistic
                    value={result.monetization.ip_potential}
                    suffix="/ 100"
                    valueStyle={{ fontSize: 14, color: getScoreColor(result.monetization.ip_potential) }}
                  />
                </Flex>
              </Descriptions.Item>
              <Descriptions.Item label="核心建议">
                <Text>{result.monetization.suggestion}</Text>
              </Descriptions.Item>
            </Descriptions>
          </Card>

          {/* 风险与建议 */}
          <Row gutter={12}>
            <Col span={12}>
              <Card
                size="small"
                title={<><WarningOutlined style={{ color: token.colorError, marginRight: 6 }} />风险提示</>}
              >
                <List
                  size="small"
                  dataSource={result.risks}
                  renderItem={(item) => (
                    <List.Item style={{ padding: '6px 0' }}>
                      <Text style={{ fontSize: 13 }}>
                        <WarningOutlined style={{ color: token.colorError, marginRight: 6, fontSize: 11 }} />
                        {item}
                      </Text>
                    </List.Item>
                  )}
                />
              </Card>
            </Col>
            <Col span={12}>
              <Card
                size="small"
                title={<><BulbOutlined style={{ color: token.colorSuccess, marginRight: 6 }} />改进建议</>}
              >
                <List
                  size="small"
                  dataSource={result.recommendations}
                  renderItem={(item) => (
                    <List.Item style={{ padding: '6px 0' }}>
                      <Text style={{ fontSize: 13 }}>
                        <BulbOutlined style={{ color: token.colorSuccess, marginRight: 6, fontSize: 11 }} />
                        {item}
                      </Text>
                    </List.Item>
                  )}
                />
              </Card>
            </Col>
          </Row>

          {/* 总结 */}
          <Card size="small" style={{ background: token.colorBgLayout }}>
            <Paragraph style={{ margin: 0, fontStyle: 'italic' }}>
              {result.summary}
            </Paragraph>
          </Card>
        </Flex>
      )}
    </Flex>
  );
}

export default MarketPredictPanel;
