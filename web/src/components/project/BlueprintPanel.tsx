import { useState, useEffect } from 'react';
import { Button, Form, Input, Modal, Space, App } from 'antd';
import { PlayCircleOutlined, SaveOutlined } from '@ant-design/icons';
import { projectApi } from '../../api';
import type { Project } from '../../types';
import { useMutation, useQueryClient } from '@tanstack/react-query';

const { TextArea } = Input;

interface BlueprintPanelProps {
  project: Project;
}

function BlueprintPanel({ project }: BlueprintPanelProps) {
  const { message, modal } = App.useApp();
  const queryClient = useQueryClient();
  const [form] = Form.useForm();
  const [generating, setGenerating] = useState(false);

  // 当 project 数据更新时，更新表单内容
  useEffect(() => {
    console.log('BlueprintPanel: project data updated', {
      chapter_blueprint_length: project.chapter_blueprint?.length || 0,
      hasData: !!project.chapter_blueprint,
    });
    form.setFieldsValue({
      chapter_blueprint: project.chapter_blueprint || '',
    });
  }, [project.chapter_blueprint, form]);

  // 更新项目
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

  // 生成大纲
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
    <div>
      <div className="flex justify-between items-center mb-4">
        <div>
          <h3 className="text-lg font-semibold">章节大纲</h3>
          <p className="text-gray-500 text-sm">
            为 {project.chapter_count} 个章节生成详细大纲，包含悬念节奏设计
          </p>
        </div>
        <Space>
          <Button
            type="primary"
            icon={<PlayCircleOutlined />}
            onClick={handleGenerate}
            loading={generating}
          >
            {project.blueprint_generated ? '重新生成' : '生成大纲'}
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
          chapter_blueprint: project.chapter_blueprint,
        }}
      >
        <Form.Item name="chapter_blueprint">
          <TextArea
            rows={20}
            placeholder={`章节大纲将以结构化格式生成，包含每个章节的：
- 章节号
- 章节标题
- 位置定位（情节阶段）
- 章节目的
- 悬念设置
- 伏笔铺垫
- 反转级别
- 内容摘要

点击"生成大纲"按钮开始生成...`}
          />
        </Form.Item>
      </Form>
    </div>
  );
}

export default BlueprintPanel;
