import { useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { Tabs, Button, Card, Spin } from 'antd';
import { ArrowLeftOutlined } from '@ant-design/icons';
import { projectApi } from '../api';
import { useQuery } from '@tanstack/react-query';
import ArchitecturePanel from '../components/project/ArchitecturePanel';
import BlueprintPanel from '../components/project/BlueprintPanel';
import ChapterPanel from '../components/project/ChapterPanel';

function ProjectDetail() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [activeTab, setActiveTab] = useState('architecture');

  // 获取项目详情
  const { data: projectRes, isLoading } = useQuery({
    queryKey: ['project', id],
    queryFn: () => projectApi.getById(id!).then((res) => res.data.data),
    enabled: !!id,
  });

  const project = projectRes;

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
      children: <ArchitecturePanel project={project} />,
    },
    {
      key: 'blueprint',
      label: '章节大纲',
      children: <BlueprintPanel project={project} />,
      disabled: !project.architecture_generated,
    },
    {
      key: 'chapters',
      label: '章节写作',
      children: <ChapterPanel project={project} />,
      disabled: !project.blueprint_generated,
    },
  ];

  return (
    <div>
      <div className="flex items-center gap-4 mb-6">
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
