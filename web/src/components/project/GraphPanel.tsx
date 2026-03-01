import { useState, useMemo, useRef } from 'react';
import {
  Button, Empty, Spin, Typography, Flex, theme, Tag, Card, Select, App, Space, Tooltip,
} from 'antd';
import {
  ApartmentOutlined,
  ReloadOutlined,
  NodeIndexOutlined,
  HistoryOutlined,
  DownloadOutlined,
} from '@ant-design/icons';
import ReactECharts from 'echarts-for-react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { graphApi } from '../../api';
import type { Project, GraphData, GraphNode, GraphEdge } from '../../types';

const { Title, Text } = Typography;

const NODE_TYPE_CONFIG: Record<string, { color: string; size: number; label: string }> = {
  protagonist: { color: '#ab372f', size: 50, label: '主角' },
  antagonist: { color: '#dc2626', size: 45, label: '反派' },
  supporting: { color: '#0891b2', size: 38, label: '配角' },
  minor: { color: '#6b7280', size: 28, label: '次要' },
};

interface GraphPanelProps {
  project: Project;
}

function GraphPanel({ project }: GraphPanelProps) {
  const { token } = theme.useToken();
  const { message: messageApi } = App.useApp();
  const queryClient = useQueryClient();
  const [selectedNode, setSelectedNode] = useState<GraphNode | null>(null);
  const [snapshotChapter, setSnapshotChapter] = useState<number | null>(null);
  const chartRef = useRef<ReactECharts>(null);

  const { data: graphData, isLoading } = useQuery({
    queryKey: ['graph', project.id, snapshotChapter],
    queryFn: async () => {
      if (snapshotChapter !== null) {
        const res = await graphApi.getChapterSnapshot(project.id, snapshotChapter);
        return res.data;
      }
      const res = await graphApi.get(project.id);
      return res.data;
    },
    enabled: !!project.id,
  });

  const generateMutation = useMutation({
    mutationFn: () => graphApi.generate(project.id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['graph', project.id] });
      messageApi.success('关系图谱生成成功');
    },
    onError: () => {
      messageApi.error('生成失败');
    },
  });

  const hasGraph = graphData && graphData.nodes && graphData.nodes.length > 0;

  const chartOption = useMemo(() => {
    if (!hasGraph) return {};
    return buildChartOption(graphData!, token, selectedNode?.id);
  }, [graphData, token, selectedNode, hasGraph]);

  const handleChartClick = (params: any) => {
    if (params.dataType === 'node') {
      const node = graphData?.nodes.find(n => n.id === params.data.id);
      setSelectedNode(node || null);
    } else {
      setSelectedNode(null);
    }
  };

  const handleExportPNG = () => {
    const instance = chartRef.current?.getEchartsInstance();
    if (!instance) {
      messageApi.error('图谱未就绪');
      return;
    }
    const url = instance.getDataURL({
      type: 'png',
      pixelRatio: 2,
      backgroundColor: token.colorBgContainer,
    });
    const a = document.createElement('a');
    a.href = url;
    a.download = `${project.title}-关系图谱.png`;
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    messageApi.success('图谱已导出为 PNG');
  };

  const snapshotOptions = (graphData?.snapshots || []).map(s => ({
    value: s.chapter_number,
    label: `第${s.chapter_number}章 (${s.nodes_count}角色, ${s.edges_count}关系)`,
  }));

  // 找到当前选中节点的关系
  const selectedEdges = selectedNode
    ? (graphData?.edges || []).filter(e => e.source === selectedNode.id || e.target === selectedNode.id)
    : [];

  return (
    <div style={{ padding: 24 }}>
      <Flex justify="space-between" align="center" wrap gap={16} style={{ marginBottom: 20 }}>
        <div>
          <Flex align="center" gap={8}>
            <Title level={5} style={{ margin: 0 }}>角色关系图谱</Title>
            <ApartmentOutlined style={{ color: token.colorPrimary, fontSize: 16 }} />
          </Flex>
          <Text type="secondary" style={{ fontSize: 13 }}>
            {hasGraph
              ? `${graphData!.nodes.length} 个角色，${graphData!.edges.length} 条关系`
              : '从小说架构中提取角色关系网络'}
          </Text>
        </div>
        <Space>
          {snapshotOptions.length > 0 && (
            <Select
              placeholder="查看章节快照"
              allowClear
              style={{ width: 220 }}
              options={snapshotOptions}
              value={snapshotChapter}
              onChange={(v) => {
                setSnapshotChapter(v ?? null);
                setSelectedNode(null);
              }}
              suffixIcon={<HistoryOutlined />}
            />
          )}
          {hasGraph && (
            <Button
              icon={<DownloadOutlined />}
              onClick={handleExportPNG}
            >
              导出 PNG
            </Button>
          )}
          <Button
            type="primary"
            icon={hasGraph ? <ReloadOutlined /> : <NodeIndexOutlined />}
            onClick={() => generateMutation.mutate()}
            loading={generateMutation.isPending}
            disabled={!project.architecture_generated}
          >
            {hasGraph ? '重新生成' : '生成图谱'}
          </Button>
        </Space>
      </Flex>

      {!project.architecture_generated ? (
        <Card style={{ textAlign: 'center', padding: '60px 0', borderStyle: 'dashed' }}>
          <Text style={{ fontSize: 16, display: 'block', marginBottom: 8 }}>请先生成小说架构</Text>
          <Text type="secondary">关系图谱将从架构中的角色设定自动提取</Text>
        </Card>
      ) : isLoading ? (
        <Flex justify="center" align="center" style={{ height: 400 }}>
          <Spin size="large" />
        </Flex>
      ) : !hasGraph ? (
        <Card style={{ textAlign: 'center', padding: '60px 0', borderStyle: 'dashed' }}>
          <ApartmentOutlined style={{ fontSize: 48, color: token.colorTextQuaternary, marginBottom: 16 }} />
          <Text style={{ fontSize: 16, display: 'block', marginBottom: 16 }}>暂无图谱数据</Text>
          <Button type="primary" onClick={() => generateMutation.mutate()} loading={generateMutation.isPending}>
            从架构生成图谱
          </Button>
        </Card>
      ) : (
        <Flex gap={16} style={{ height: 560 }}>
          {/* 图谱画布 */}
          <div style={{
            flex: 1,
            border: `1px solid ${token.colorBorderSecondary}`,
            borderRadius: token.borderRadius,
            overflow: 'hidden',
            background: token.colorBgContainer,
            minWidth: 0,
          }}>
            <ReactECharts
              ref={chartRef}
              option={chartOption}
              style={{ height: '100%', width: '100%' }}
              onEvents={{ click: handleChartClick }}
              notMerge
            />
          </div>

          {/* 右侧详情 */}
          <div style={{
            width: 280,
            flexShrink: 0,
            overflow: 'auto',
          }}>
            {selectedNode ? (
              <Card size="small" styles={{ body: { padding: 16 } }}>
                <Flex vertical gap={12}>
                  <Flex align="center" gap={8}>
                    <div style={{
                      width: 12, height: 12, borderRadius: '50%',
                      background: NODE_TYPE_CONFIG[selectedNode.type]?.color || '#999',
                    }} />
                    <Text strong style={{ fontSize: 16 }}>{selectedNode.name}</Text>
                    <Tag color={NODE_TYPE_CONFIG[selectedNode.type]?.color} style={{ margin: 0 }}>
                      {NODE_TYPE_CONFIG[selectedNode.type]?.label || selectedNode.type}
                    </Tag>
                  </Flex>

                  <Text type="secondary">{selectedNode.description}</Text>

                  {selectedNode.traits.length > 0 && (
                    <div>
                      <Text strong style={{ fontSize: 12 }}>性格特点</Text>
                      <Flex wrap gap={4} style={{ marginTop: 4 }}>
                        {selectedNode.traits.map(t => (
                          <Tag key={t} bordered={false}>{t}</Tag>
                        ))}
                      </Flex>
                    </div>
                  )}

                  {selectedNode.group && (
                    <div>
                      <Text strong style={{ fontSize: 12 }}>阵营</Text>
                      <div style={{ marginTop: 4 }}>
                        <Tag color="processing" bordered={false}>{selectedNode.group}</Tag>
                      </div>
                    </div>
                  )}

                  {selectedEdges.length > 0 && (
                    <div>
                      <Text strong style={{ fontSize: 12 }}>关系（{selectedEdges.length}）</Text>
                      <Flex vertical gap={6} style={{ marginTop: 6 }}>
                        {selectedEdges.map((edge, i) => {
                          const otherId = edge.source === selectedNode.id ? edge.target : edge.source;
                          const otherNode = graphData!.nodes.find(n => n.id === otherId);
                          return (
                            <div
                              key={i}
                              style={{
                                padding: '6px 10px',
                                borderRadius: 6,
                                background: token.colorBgLayout,
                                fontSize: 12,
                              }}
                            >
                              <Flex justify="space-between" align="center">
                                <Text strong>{otherNode?.name || otherId}</Text>
                                <Tag bordered={false} style={{ margin: 0, fontSize: 11 }}>{edge.relation}</Tag>
                              </Flex>
                              {edge.description && (
                                <Text type="secondary" style={{ fontSize: 11 }}>{edge.description}</Text>
                              )}
                            </div>
                          );
                        })}
                      </Flex>
                    </div>
                  )}
                </Flex>
              </Card>
            ) : (
              <Card size="small" style={{ textAlign: 'center', padding: '40px 16px' }}>
                <ApartmentOutlined style={{ fontSize: 32, color: token.colorTextQuaternary, marginBottom: 8 }} />
                <br />
                <Text type="secondary" style={{ fontSize: 12 }}>
                  点击图谱中的角色节点查看详情
                </Text>

                {/* 图例 */}
                <Flex vertical gap={6} style={{ marginTop: 24 }}>
                  <Text strong style={{ fontSize: 12, textAlign: 'left' }}>图例</Text>
                  {Object.entries(NODE_TYPE_CONFIG).map(([key, cfg]) => (
                    <Flex key={key} align="center" gap={8}>
                      <div style={{
                        width: 10, height: 10, borderRadius: '50%',
                        background: cfg.color,
                      }} />
                      <Text style={{ fontSize: 12 }}>{cfg.label}</Text>
                    </Flex>
                  ))}
                </Flex>
              </Card>
            )}

            {/* 快照时间轴 */}
            {(graphData?.snapshots || []).length > 0 && (
              <Card size="small" style={{ marginTop: 12 }} styles={{ body: { padding: 12 } }}>
                <Text strong style={{ fontSize: 12, display: 'block', marginBottom: 8 }}>
                  <HistoryOutlined /> 更新记录
                </Text>
                <Flex vertical gap={4}>
                  {(graphData?.snapshots || []).slice(-5).reverse().map((snap, i) => (
                    <Tooltip key={i} title={snap.summary}>
                      <div
                        onClick={() => setSnapshotChapter(snap.chapter_number)}
                        style={{
                          padding: '4px 8px',
                          borderRadius: 4,
                          background: snapshotChapter === snap.chapter_number ? token.colorPrimaryBg : 'transparent',
                          cursor: 'pointer',
                          fontSize: 11,
                        }}
                      >
                        <Text style={{ fontSize: 11 }}>
                          第{snap.chapter_number}章 · {snap.nodes_count}角色 · {snap.edges_count}关系
                        </Text>
                      </div>
                    </Tooltip>
                  ))}
                </Flex>
              </Card>
            )}
          </div>
        </Flex>
      )}
    </div>
  );
}

function buildChartOption(graph: GraphData, token: any, selectedNodeId?: string) {
  const categories = Object.entries(NODE_TYPE_CONFIG).map(([, cfg]) => ({
    name: cfg.label,
    itemStyle: { color: cfg.color },
  }));

  const categoryIndex: Record<string, number> = {
    protagonist: 0,
    antagonist: 1,
    supporting: 2,
    minor: 3,
  };

  const nodes = graph.nodes.map(node => {
    const cfg = NODE_TYPE_CONFIG[node.type] || NODE_TYPE_CONFIG.minor;
    const isSelected = node.id === selectedNodeId;
    return {
      id: node.id,
      name: node.name,
      symbolSize: isSelected ? cfg.size * 1.3 : cfg.size,
      category: categoryIndex[node.type] ?? 3,
      itemStyle: {
        color: cfg.color,
        borderColor: isSelected ? '#fff' : 'transparent',
        borderWidth: isSelected ? 3 : 0,
        shadowBlur: isSelected ? 15 : 0,
        shadowColor: isSelected ? cfg.color : 'transparent',
      },
      label: {
        show: true,
        fontSize: isSelected ? 14 : 12,
        fontWeight: isSelected ? 'bold' as const : 'normal' as const,
        color: token.colorText,
      },
    };
  });

  const edges = graph.edges.map(edge => ({
    source: edge.source,
    target: edge.target,
    value: edge.weight,
    lineStyle: {
      width: Math.max(1, edge.weight / 3),
      color: '#aaa',
      opacity: selectedNodeId
        ? (edge.source === selectedNodeId || edge.target === selectedNodeId ? 0.8 : 0.15)
        : 0.5,
    },
    label: {
      show: !selectedNodeId || edge.source === selectedNodeId || edge.target === selectedNodeId,
      formatter: edge.relation,
      fontSize: 10,
      color: token.colorTextSecondary,
    },
  }));

  return {
    tooltip: {
      trigger: 'item' as const,
      formatter: (params: any) => {
        if (params.dataType === 'node') {
          const node = graph.nodes.find(n => n.id === params.data.id);
          if (!node) return '';
          return `<b>${node.name}</b><br/>${node.description}<br/>阵营：${node.group}`;
        }
        if (params.dataType === 'edge') {
          const edge = graph.edges.find(
            e => e.source === params.data.source && e.target === params.data.target
          );
          if (!edge) return '';
          return `<b>${edge.relation}</b><br/>${edge.description}`;
        }
        return '';
      },
    },
    legend: {
      data: categories.map(c => c.name),
      bottom: 10,
      textStyle: { color: token.colorText, fontSize: 12 },
    },
    animationDuration: 800,
    animationEasingUpdate: 'quinticInOut' as const,
    series: [{
      type: 'graph' as const,
      layout: 'force' as const,
      data: nodes,
      links: edges,
      categories,
      roam: true,
      draggable: true,
      force: {
        repulsion: 300,
        edgeLength: [100, 250],
        gravity: 0.1,
      },
      emphasis: {
        focus: 'adjacency' as const,
        lineStyle: { width: 3 },
      },
      label: {
        position: 'bottom' as const,
        distance: 5,
      },
      edgeLabel: {
        show: true,
        position: 'middle' as const,
      },
    }],
  };
}

export default GraphPanel;
