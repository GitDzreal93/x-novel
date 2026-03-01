import { useState, useEffect } from 'react';
import { Button, Form, Input, Modal, Space, Collapse, App } from 'antd';
import { PlayCircleOutlined, SaveOutlined } from '@ant-design/icons';
import { projectApi } from '../../api';
import type { Project } from '../../types';
import { useMutation, useQueryClient } from '@tanstack/react-query';

const { TextArea } = Input;

interface ArchitecturePanelProps {
  project: Project;
}

function ArchitecturePanel({ project }: ArchitecturePanelProps) {
  const { message } = App.useApp();
  const queryClient = useQueryClient();
  const [form] = Form.useForm();
  const [generating, setGenerating] = useState(false);

  // 当 project 架构数据更新时，更新表单内容
  useEffect(() => {
    form.setFieldsValue({
      core_seed: project.core_seed || '',
      character_dynamics: project.character_dynamics || '',
      world_building: project.world_building || '',
      plot_architecture: project.plot_architecture || '',
      character_state: project.character_state || '',
    });
  }, [project.core_seed, project.character_dynamics, project.world_building, project.plot_architecture, project.character_state, form]);

  // 更新项目
  const updateMutation = useMutation({
    mutationFn: (data: { core_seed?: string; character_dynamics?: string; world_building?: string; plot_architecture?: string; character_state?: string }) =>
      projectApi.update(project.id, data),
    onSuccess: () => {
      message.success('保存成功');
      queryClient.invalidateQueries({ queryKey: ['project', project.id] });
    },
    onError: () => {
      message.error('保存失败');
    },
  });

  // 生成架构
  const generateMutation = useMutation({
    mutationFn: (overwrite: boolean) =>
      projectApi.generateArchitecture(project.id, { overwrite }),
    onSuccess: () => {
      message.success('架构生成成功');
      queryClient.invalidateQueries({ queryKey: ['project', project.id] });
      setGenerating(false);
    },
    onError: () => {
      message.error('架构生成失败');
      setGenerating(false);
    },
  });

  const handleGenerate = () => {
    if (project.architecture_generated) {
      Modal.confirm({
        title: '重新生成',
        content: '架构已生成，确定要重新生成吗？当前内容将被覆盖。',
        onOk: () => {
          setGenerating(true);
          generateMutation.mutate(true);
        },
      });
    } else {
      setGenerating(true);
      generateMutation.mutate(false);
    }
  };

  const handleSave = () => {
    const values = form.getFieldsValue();
    updateMutation.mutate(values);
  };

  const collapseItems = [
    {
      key: 'core_seed',
      label: '1. 核心种子',
      children: (
        <Form.Item name="core_seed" className="mb-0">
          <TextArea
            rows={6}
            placeholder="核心种子是故事的创意火花，包含核心概念、创意前提、核心冲突、核心主题..."
          />
        </Form.Item>
      ),
    },
    {
      key: 'character_dynamics',
      label: '2. 角色动力学',
      children: (
        <Form.Item name="character_dynamics" className="mb-0">
          <TextArea
            rows={6}
            placeholder="角色动力学包含主要角色关系、动机三角形、角色弧光设计..."
          />
        </Form.Item>
      ),
    },
    {
      key: 'world_building',
      label: '3. 世界观构建',
      children: (
        <Form.Item name="world_building" className="mb-0">
          <TextArea
            rows={6}
            placeholder="世界观构建包含设定框架、时空体系、社会结构、能力体系..."
          />
        </Form.Item>
      ),
    },
    {
      key: 'plot_architecture',
      label: '4. 情节架构',
      children: (
        <Form.Item name="plot_architecture" className="mb-0">
          <TextArea
            rows={6}
            placeholder="情节架构包含主线结构、关键转折、高潮设计、结局走向..."
          />
        </Form.Item>
      ),
    },
    {
      key: 'character_state',
      label: '5. 角色状态',
      children: (
        <Form.Item name="character_state" className="mb-0">
          <TextArea
            rows={6}
            placeholder="角色状态包含人物档案、心理特质、外部关联、当前状态..."
          />
        </Form.Item>
      ),
    },
  ];

  return (
    <div>
      <div className="flex justify-between items-center mb-4">
        <div>
          <h3 className="text-lg font-semibold">小说架构</h3>
          <p className="text-gray-500 text-sm">
            基于雪花写作法，从 5 个维度构建完整小说架构
          </p>
        </div>
        <Space>
          <Button
            type="primary"
            icon={<PlayCircleOutlined />}
            onClick={handleGenerate}
            loading={generating}
          >
            {project.architecture_generated ? '重新生成' : '生成架构'}
          </Button>
          <Button
            icon={<SaveOutlined />}
            onClick={handleSave}
            loading={updateMutation.isPending}
          >
            保存
          </Button>
        </Space>
      </div>

      <Form
        form={form}
        layout="vertical"
        initialValues={{
          core_seed: project.core_seed,
          character_dynamics: project.character_dynamics,
          world_building: project.world_building,
          plot_architecture: project.plot_architecture,
          character_state: project.character_state,
        }}
      >
        <Collapse items={collapseItems} defaultActiveKey={['core_seed']} />
      </Form>
    </div>
  );
}

export default ArchitecturePanel;
