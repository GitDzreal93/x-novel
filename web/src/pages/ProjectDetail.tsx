import { useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { Tabs, Button, Card, Spin, Space, Dropdown, App } from 'antd';
import { ArrowLeftOutlined, DownloadOutlined, FileTextOutlined, FileOutlined } from '@ant-design/icons';
import { projectApi } from '../api';
import { useQuery } from '@tanstack/react-query';
import ArchitecturePanel from '../components/project/ArchitecturePanel';
import BlueprintPanel from '../components/project/BlueprintPanel';
import ChapterPanel from '../components/project/ChapterPanel';

function ProjectDetail() {
  const { message } = App.useApp();
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [activeTab, setActiveTab] = useState('architecture');

  // 获取项目详情
  const { data: projectRes, isLoading } = useQuery({
    queryKey: ['project', id],
    queryFn: () => projectApi.getById(id!).then((res) => {
      if (!res?.data) {
        throw new Error(res?.message || '获取项目失败');
      }
      return res.data;
    }),
    enabled: !!id,
  });

  const project = projectRes;

  // 导出处理
  const handleExport = async (format: 'txt' | 'md') => {
    try {
      const res = await projectApi.export(id!, format);
      if (!res?.data?.data?.download_url) {
        throw new Error(res?.data?.message || '导出失败');
      }
      const content = res.data.data.download_url;

      // 创建 Blob 并下载
      const blob = new Blob([content], {
        type: format === 'txt' ? 'text/plain' : 'text/markdown'
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
    } catch (error) {
      message.error('导出失败');
    }
  };

  const exportMenuItems = [
    {
      key: 'txt',
      label: '导出为 TXT',
      icon: <FileTextOutlined />,
      onClick: () => handleExport('txt'),
    },
    {
      key: 'md',
      label: '导出为 Markdown',
      icon: <FileOutlined />,
      onClick: () => handleExport('md'),
    },
  ];

  if (isLoading) {
    return (
      <div className="flex justify-center items-center h-64">
        <Spin size="large" />
      </div>
    );
  }

  if (!project) {
    return (
      <div className="text-center py-12">
        <p className="mb-4">项目不存在</p>
        <Button onClick={() => navigate('/projects')}>返回项目列表</Button>
      </div>
    );
  }

  const tabItems = [
    {
      key: 'architecture',
      label: '小说架构',
      children: <ArchitecturePanel key={`arch-${project.updated_at}`} project={project} />,
    },
    {
      key: 'blueprint',
      label: '章节大纲',
      children: <BlueprintPanel key={`blueprint-${project.updated_at}`} project={project} />,
      disabled: !project.architecture_generated,
    },
    {
      key: 'chapters',
      label: '章节写作',
      children: <ChapterPanel key={`chapters-${project.updated_at}`} project={project} />,
      disabled: !project.blueprint_generated,
    },
  ];

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <div className="flex items-center gap-4">
          <Button
            icon={<ArrowLeftOutlined />}
            onClick={() => navigate('/projects')}
          >
            返回
          </Button>
          <div>
            <h1 className="text-2xl font-bold">{project.title}</h1>
            {project.topic && (
              <p className="text-gray-500 text-sm mt-1">{project.topic}</p>
            )}
          </div>
        </div>
        <Space>
          <Dropdown menu={{ items: exportMenuItems }} placement="bottomRight">
            <Button icon={<DownloadOutlined />} loading={isLoading}>
              导出
            </Button>
          </Dropdown>
        </Space>
      </div>

      <Card>
        <Tabs
          activeKey={activeTab}
          onChange={setActiveTab}
          items={tabItems}
        />
      </Card>
    </div>
  );
}

export default ProjectDetail;
