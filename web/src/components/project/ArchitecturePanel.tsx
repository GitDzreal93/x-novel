import { useState, useEffect } from 'react';
import { Button, Form, Input, Space, Collapse, App, Badge, Typography, Flex, theme } from 'antd';
import { PlayCircleOutlined, SaveOutlined, ThunderboltOutlined } from '@ant-design/icons';
import { projectApi } from '../../api';
import type { Project } from '../../types';
import { useMutation, useQueryClient } from '@tanstack/react-query';

const { TextArea } = Input;
const { Title, Text } = Typography;

interface ArchitecturePanelProps {
  project: Project;
}

function ArchitecturePanel({ project }: ArchitecturePanelProps) {
  const { message, modal } = App.useApp();
  const queryClient = useQueryClient();
  const { token } = theme.useToken();
  const [form] = Form.useForm();
  const [generating, setGenerating] = useState(false);

  useEffect(() => {
    form.setFieldsValue({
      core_seed: project.core_seed || '',
      character_dynamics: project.character_dynamics || '',
      world_building: project.world_building || '',
      plot_architecture: project.plot_architecture || '',
      character_state: project.character_state || '',
    });
  }, [project.core_seed, project.character_dynamics, project.world_building, project.plot_architecture, project.character_state, form]);

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
      modal.confirm({
        title: '重新生成',
        content: '架构已生成，确定要重新生成吗？当前内容将被覆盖。',
        centered: true,
        okButtonProps: { danger: true },
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

  const createHeader = (title: string, hasContent: boolean, color: string) => (
    <Flex align="center" gap={8} style={{ padding: '4px 0' }}>
      <Badge color={hasContent ? color : '#d9d9d9'} />
      <Text strong={hasContent} type={hasContent ? undefined : 'secondary'} style={{ fontSize: 15 }}>
        {title}
      </Text>
    </Flex>
  );

  const textAreaStyle: React.CSSProperties = {
    fontFamily: 'monospace',
    fontSize: 13,
    lineHeight: 1.8,
    padding: 16,
    minHeight: 200,
    resize: 'vertical' as const,
  };

  const collapseItems = [
    {
      key: 'core_seed',
      label: createHeader('1. 核心种子', !!project.core_seed, '#3b82f6'),
      children: (
        <Form.Item name="core_seed" style={{ marginBottom: 0 }}>
          <TextArea
            rows={10}
            style={textAreaStyle}
            placeholder="核心种子是故事的创意火花，包含核心概念、创意前提、核心冲突、核心主题..."
          />
        </Form.Item>
      ),
    },
    {
      key: 'character_dynamics',
      label: createHeader('2. 角色动力学', !!project.character_dynamics, '#ec4899'),
      children: (
        <Form.Item name="character_dynamics" style={{ marginBottom: 0 }}>
          <TextArea
            rows={12}
            style={textAreaStyle}
            placeholder="角色动力学包含主要角色关系、动机三角形、角色弧光设计..."
          />
        </Form.Item>
      ),
    },
    {
      key: 'world_building',
      label: createHeader('3. 世界观构建', !!project.world_building, '#10b981'),
      children: (
        <Form.Item name="world_building" style={{ marginBottom: 0 }}>
          <TextArea
            rows={12}
            style={textAreaStyle}
            placeholder="世界观构建包含设定框架、时空体系、社会结构、能力体系..."
          />
        </Form.Item>
      ),
    },
    {
      key: 'plot_architecture',
      label: createHeader('4. 情节架构', !!project.plot_architecture, '#f59e0b'),
      children: (
        <Form.Item name="plot_architecture" style={{ marginBottom: 0 }}>
          <TextArea
            rows={12}
            style={textAreaStyle}
            placeholder="情节架构包含主线结构、关键转折、高潮设计、结局走向..."
          />
        </Form.Item>
      ),
    },
    {
      key: 'character_state',
      label: createHeader('5. 角色状态', !!project.character_state, '#8b5cf6'),
      children: (
        <Form.Item name="character_state" style={{ marginBottom: 0 }}>
          <TextArea
            rows={12}
            style={textAreaStyle}
            placeholder="角色状态包含人物档案、心理特质、外部关联、当前状态..."
          />
        </Form.Item>
      ),
    },
  ];

  return (
    <div style={{ padding: '0 24px 24px' }}>
      <Flex
        justify="space-between"
        align="center"
        wrap
        gap={16}
        style={{ marginBottom: 20 }}
      >
        <div>
          <Flex align="center" gap={8}>
            <Title level={5} style={{ margin: 0 }}>小说架构设计</Title>
            <ThunderboltOutlined style={{ color: token.colorPrimary, fontSize: 16 }} />
          </Flex>
          <Text type="secondary" style={{ fontSize: 13 }}>
            基于雪花写作法，从 5 个维度智能构建完整的小说底层架构
          </Text>
        </div>
        <Space>
          <Button
            type={project.architecture_generated ? 'default' : 'primary'}
            icon={project.architecture_generated ? <PlayCircleOutlined /> : <ThunderboltOutlined />}
            onClick={handleGenerate}
            loading={generating || generateMutation.isPending}
            size="large"
          >
            {project.architecture_generated ? '重新生成' : 'AI 一键生成架构'}
          </Button>
          <Button
            type="primary"
            icon={<SaveOutlined />}
            onClick={handleSave}
            loading={updateMutation.isPending}
            size="large"
            style={{ background: '#16a34a' }}
          >
            保存修改
          </Button>
        </Space>
      </Flex>

      <div
        style={{
          padding: 8,
          borderRadius: token.borderRadiusLG,
          border: `1px solid ${token.colorBorderSecondary}`,
          background: token.colorBgLayout,
        }}
      >
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
          <Collapse
            items={collapseItems}
            defaultActiveKey={['core_seed', 'character_dynamics', 'world_building', 'plot_architecture', 'character_state']}
            expandIconPlacement="end"
            style={{
              background: 'transparent',
              border: 'none',
            }}
          />
        </Form>
      </div>
    </div>
  );
}

export default ArchitecturePanel;
