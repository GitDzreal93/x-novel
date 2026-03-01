import { useState, useEffect } from 'react';
import { useParams, useNavigate, useLocation } from 'react-router-dom';
import {
  Tabs, Button, Spin, Dropdown, App, Tag, Tooltip,
  Typography, Card, Flex, theme,
} from 'antd';
import {
  ArrowLeftOutlined,
  DownloadOutlined,
  FileTextOutlined,
  FileOutlined,
  BlockOutlined,
  ReadOutlined,
  EditOutlined,
  CheckCircleOutlined,
  ApartmentOutlined,
  AuditOutlined,
  StockOutlined,
} from '@ant-design/icons';
import { projectApi } from '../api';
import { useQuery } from '@tanstack/react-query';
import ArchitecturePanel from '../components/project/ArchitecturePanel';
import BlueprintPanel from '../components/project/BlueprintPanel';
import ChapterPanel from '../components/project/ChapterPanel';
import GraphPanel from '../components/project/GraphPanel';
import ReviewPanel from '../components/project/ReviewPanel';
import MarketPredictPanel from '../components/project/MarketPredictPanel';

const { Title, Text, Paragraph } = Typography;

function ProjectDetail() {
  const { message } = App.useApp();
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const location = useLocation();
  const { token } = theme.useToken();
  const [activeTab, setActiveTab] = useState('architecture');

  useEffect(() => {
    const params = new URLSearchParams(location.search);
    const tab = params.get('tab');
    if (tab && ['architecture', 'blueprint', 'chapters', 'graph', 'review', 'market'].includes(tab)) {
      setActiveTab(tab);
    }
  }, [location]);

  const handleTabChange = (key: string) => {
    setActiveTab(key);
    navigate(`/projects/${id}?tab=${key}`, { replace: true });
  };

  const { data: projectRes, isLoading } = useQuery({
    queryKey: ['project', id],
    queryFn: () =>
      projectApi.getById(id!).then((res) => {
        if (!res?.data) {
          throw new Error(res?.message || '获取项目失败');
        }
        return res.data;
      }),
    enabled: !!id,
  });

  const project = projectRes;

  const handleExport = async (format: 'txt' | 'md') => {
    try {
      const res = await projectApi.export(id!, format);
      if (!res?.data?.download_url) {
        throw new Error(res?.message || '导出失败');
      }
      const content = res.data.download_url;

      const blob = new Blob([content], {
        type: format === 'txt' ? 'text/plain' : 'text/markdown',
      });
      const url = URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = `${project?.title || 'novel'}.${format}`;
      document.body.appendChild(a);
      a.click();
      document.body.removeChild(a);
      URL.revokeObjectURL(url);

      message.success(`${format === 'txt' ? 'TXT' : 'Markdown'} 导出成功`);
    } catch {
      message.error('导出失败');
    }
  };

  const exportMenuItems = [
    {
      key: 'txt',
      label: '导出为 TXT',
      icon: <FileTextOutlined style={{ color: token.colorPrimary }} />,
      onClick: () => handleExport('txt'),
    },
    {
      key: 'md',
      label: '导出为 Markdown',
      icon: <FileOutlined style={{ color: token.colorTextSecondary }} />,
      onClick: () => handleExport('md'),
    },
  ];

  if (isLoading) {
    return (
      <Flex justify="center" align="center" style={{ height: 'calc(100vh - 200px)' }}>
        <Spin size="large" />
      </Flex>
    );
  }

  if (!project) {
    return (
      <Flex vertical align="center" justify="center" style={{ padding: '80px 0' }}>
        <Text type="secondary" style={{ fontSize: 16, marginBottom: 16 }}>
          项目不存在或已被删除
        </Text>
        <Button type="primary" onClick={() => navigate('/projects')}>
          返回项目列表
        </Button>
      </Flex>
    );
  }

  const tabItems = [
    {
      key: 'architecture',
      label: (
        <Flex align="center" gap={6}>
          <BlockOutlined />
          小说架构
          {project.architecture_generated && (
            <CheckCircleOutlined style={{ color: token.colorSuccess, fontSize: 12 }} />
          )}
        </Flex>
      ),
      children: <ArchitecturePanel key={`arch-${project.updated_at}`} project={project} />,
    },
    {
      key: 'blueprint',
      label: (
        <Flex align="center" gap={6}>
          <ReadOutlined />
          章节大纲
          {project.blueprint_generated && (
            <CheckCircleOutlined style={{ color: token.colorSuccess, fontSize: 12 }} />
          )}
        </Flex>
      ),
      children: <BlueprintPanel key={`blueprint-${project.updated_at}`} project={project} />,
      disabled: !project.architecture_generated,
    },
    {
      key: 'chapters',
      label: (
        <Flex align="center" gap={6}>
          <EditOutlined />
          章节写作
          {project.completed_chapters > 0 && (
            <Tag color="processing" bordered={false} style={{ marginLeft: 4, fontSize: 12 }}>
              {project.completed_chapters}
            </Tag>
          )}
        </Flex>
      ),
      children: <ChapterPanel key={`chapters-${project.updated_at}`} project={project} />,
      disabled: !project.blueprint_generated,
    },
    {
      key: 'graph',
      label: (
        <Flex align="center" gap={6}>
          <ApartmentOutlined />
          关系图谱
        </Flex>
      ),
      children: <GraphPanel key={`graph-${project.updated_at}`} project={project} />,
      disabled: !project.architecture_generated,
    },
    {
      key: 'review',
      label: (
        <Flex align="center" gap={6}>
          <AuditOutlined />
          AI 审阅
        </Flex>
      ),
      children: <ReviewPanel key={`review-${project.updated_at}`} project={project} />,
      disabled: project.completed_chapters === 0,
    },
    {
      key: 'market',
      label: (
        <Flex align="center" gap={6}>
          <StockOutlined />
          市场预测
        </Flex>
      ),
      children: <MarketPredictPanel key={`market-${project.updated_at}`} project={project} />,
      disabled: !project.architecture_generated,
    },
  ];

  return (
    <div style={{ paddingBottom: 24 }}>
      {/* 页面头部 */}
      <Card style={{ marginBottom: 20 }}>
        <Flex justify="space-between" align="flex-start" wrap gap={16}>
          <Flex gap={16}>
            <Tooltip title="返回列表">
              <Button
                type="text"
                icon={<ArrowLeftOutlined />}
                onClick={() => navigate('/projects')}
                style={{ marginTop: 2 }}
              />
            </Tooltip>

            <div style={{ minWidth: 0, flex: 1 }}>
              <Flex align="center" gap={12} style={{ marginBottom: 6 }}>
                <Title level={4} ellipsis style={{ margin: 0 }}>
                  {project.title}
                </Title>
                <StatusTag status={project.status} />
              </Flex>

              {project.topic && (
                <Paragraph
                  type="secondary"
                  style={{ fontSize: 13, marginBottom: 8, maxWidth: 700 }}
                >
                  {project.topic}
                </Paragraph>
              )}

              {project.genre && project.genre.length > 0 && (
                <Flex wrap gap={6} style={{ marginBottom: 10 }}>
                  {project.genre.map((g: string) => (
                    <Tag key={g} color="purple" bordered={false} style={{ fontSize: 12 }}>
                      {g}
                    </Tag>
                  ))}
                </Flex>
              )}

              <Flex
                wrap
                gap={24}
                style={{
                  paddingTop: 10,
                  borderTop: `1px solid ${token.colorBorderSecondary}`,
                  fontSize: 12,
                }}
              >
                <Text type="secondary">
                  <span style={{ display: 'inline-block', width: 6, height: 6, borderRadius: '50%', background: '#3b82f6', marginRight: 6, verticalAlign: 'middle' }} />
                  总计：{project.chapter_count} 章
                </Text>
                <Text type="secondary">
                  <span style={{ display: 'inline-block', width: 6, height: 6, borderRadius: '50%', background: '#22c55e', marginRight: 6, verticalAlign: 'middle' }} />
                  完成：{project.completed_chapters} 章
                </Text>
                <Text type="secondary">
                  <span style={{ display: 'inline-block', width: 6, height: 6, borderRadius: '50%', background: '#a855f7', marginRight: 6, verticalAlign: 'middle' }} />
                  字数：{project.total_words > 10000
                    ? `${(project.total_words / 10000).toFixed(1)}万`
                    : `${project.total_words}`} 字
                </Text>
              </Flex>
            </div>
          </Flex>

          <div style={{ flexShrink: 0 }}>
            <Dropdown menu={{ items: exportMenuItems }} placement="bottomRight" trigger={['click']}>
              <Button type="primary" icon={<DownloadOutlined />} size="large">
                导出作品
              </Button>
            </Dropdown>
          </div>
        </Flex>
      </Card>

      {/* Tabs */}
      <Card styles={{ body: { padding: 0 } }}>
        <Tabs
          activeKey={activeTab}
          onChange={handleTabChange}
          items={tabItems}
          size="large"
          tabBarStyle={{
            padding: '0 24px',
            marginBottom: 0,
          }}
        />
      </Card>
    </div>
  );
}

function StatusTag({ status }: { status: string }) {
  switch (status) {
    case 'writing':
      return <Tag color="processing" bordered={false}>连载中</Tag>;
    case 'completed':
      return <Tag color="success" bordered={false}>已完结</Tag>;
    case 'published':
      return <Tag color="purple" bordered={false}>已发布</Tag>;
    default:
      return <Tag color="default" bordered={false}>构思中</Tag>;
  }
}

export default ProjectDetail;
