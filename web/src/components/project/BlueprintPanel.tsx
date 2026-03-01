import { useState, useEffect } from 'react';
import { Button, Form, Input, Space, App, Typography, Flex, theme } from 'antd';
import { PlayCircleOutlined, SaveOutlined, UnorderedListOutlined } from '@ant-design/icons';
import { projectApi } from '../../api';
import type { Project } from '../../types';
import { useMutation, useQueryClient } from '@tanstack/react-query';

const { TextArea } = Input;
const { Title, Text } = Typography;

interface BlueprintPanelProps {
  project: Project;
}

function BlueprintPanel({ project }: BlueprintPanelProps) {
  const { message, modal } = App.useApp();
  const queryClient = useQueryClient();
  const { token } = theme.useToken();
  const [form] = Form.useForm();
  const [generating, setGenerating] = useState(false);

  useEffect(() => {
    form.setFieldsValue({
      chapter_blueprint: project.chapter_blueprint || '',
    });
  }, [project.chapter_blueprint, form]);

  const updateMutation = useMutation({
    mutationFn: (data: { chapter_blueprint?: string }) =>
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
      projectApi.generateBlueprint(project.id, { overwrite }),
    onSuccess: () => {
      message.success('大纲生成成功');
      queryClient.invalidateQueries({ queryKey: ['project', project.id] });
      setGenerating(false);
    },
    onError: () => {
      message.error('大纲生成失败');
      setGenerating(false);
    },
  });

  const handleGenerate = () => {
    if (project.blueprint_generated) {
      modal.confirm({
        title: '重新生成',
        content: '大纲已生成，确定要重新生成吗？当前内容将被覆盖。',
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
            <Title level={5} style={{ margin: 0 }}>章节大纲设计</Title>
            <UnorderedListOutlined style={{ color: token.colorPrimary, fontSize: 16 }} />
          </Flex>
          <Text type="secondary" style={{ fontSize: 13 }}>
            为规划的{' '}
            <Text strong style={{ color: token.colorPrimary }}>{project.chapter_count}</Text>
            {' '}个章节智能生成详细大纲与悬念节奏
          </Text>
        </div>
        <Space>
          <Button
            type={project.blueprint_generated ? 'default' : 'primary'}
            icon={project.blueprint_generated ? <PlayCircleOutlined /> : <UnorderedListOutlined />}
            onClick={handleGenerate}
            loading={generating || generateMutation.isPending}
            size="large"
          >
            {project.blueprint_generated ? '重新生成' : 'AI 一键生成大纲'}
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
          padding: 16,
          borderRadius: token.borderRadiusLG,
          border: `1px solid ${token.colorBorderSecondary}`,
          background: token.colorBgLayout,
        }}
      >
        <Form
          form={form}
          layout="vertical"
          initialValues={{
            chapter_blueprint: project.chapter_blueprint,
          }}
        >
          <Form.Item name="chapter_blueprint" style={{ marginBottom: 0 }}>
            <TextArea
              rows={25}
              style={{
                fontFamily: 'monospace',
                fontSize: 13,
                lineHeight: 1.8,
                padding: 20,
                resize: 'vertical',
              }}
              placeholder={`章节大纲将以结构化格式生成，包含每个章节的：
- 章节号
- 章节标题
- 位置定位（情节阶段）
- 章节目的
- 悬念设置
- 伏笔铺垫
- 反转级别
- 内容摘要

点击右上角"AI 一键生成大纲"按钮开始生成...`}
            />
          </Form.Item>
        </Form>
      </div>
    </div>
  );
}

export default BlueprintPanel;
