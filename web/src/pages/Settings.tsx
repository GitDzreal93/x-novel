import { useState } from 'react';
import {
  Card,
  Form,
  Input,
  Select,
  Button,
  Table,
  Modal,
  Space,
  Tag,
  Switch,
  Popconfirm,
  Divider,
  App,
  Tabs,
} from 'antd';
import {
  PlusOutlined,
  DeleteOutlined,
  EditOutlined,
  CheckCircleOutlined,
  ApiOutlined,
  SettingOutlined,
} from '@ant-design/icons';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { modelConfigApi } from '../api';
import { useAppStore } from '../stores';
import type { ModelConfig, ModelProvider, CreateModelConfigRequest } from '../types';

// 用途选项
const PURPOSE_OPTIONS = [
  { label: '架构生成', value: 'architecture' },
  { label: '章节生成', value: 'chapter' },
  { label: '写作辅助', value: 'writing' },
  { label: 'AI 审阅', value: 'review' },
  { label: '通用', value: 'general' },
];

// 常用模型选项
const MODEL_OPTIONS: Record<string, string[]> = {
  openai: ['gpt-4o', 'gpt-4o-mini', 'gpt-4-turbo', 'gpt-4', 'gpt-3.5-turbo'],
  anthropic: ['claude-3-5-sonnet-20241022', 'claude-3-opus-20240229', 'claude-3-haiku-20240307'],
  deepseek: ['deepseek-chat', 'deepseek-coder'],
  qwen: ['qwen-max', 'qwen-plus', 'qwen-turbo'],
  zhipu: ['glm-4', 'glm-4-flash', 'glm-3-turbo'],
  custom: [],
};

function Settings() {
  const { message } = App.useApp();
  const queryClient = useQueryClient();
  const { theme, toggleTheme } = useAppStore();
  const [modalOpen, setModalOpen] = useState(false);
  const [editingConfig, setEditingConfig] = useState<ModelConfig | null>(null);
  const [form] = Form.useForm();
  const [selectedProvider, setSelectedProvider] = useState<number | null>(null);

  // 获取提供商列表
  const { data: providersRes } = useQuery({
    queryKey: ['model-providers'],
    queryFn: () => modelConfigApi.listProviders().then((res) => res.data?.data || []),
  });
  const providers = providersRes || [];

  // 获取模型配置列表
  const { data: configsRes, isLoading } = useQuery({
    queryKey: ['model-configs'],
    queryFn: () => modelConfigApi.list({ page: 1, page_size: 100 }).then((res) => res.data?.data),
  });
  const configs = configsRes?.configs || [];

  // 创建模型配置
  const createMutation = useMutation({
    mutationFn: (data: CreateModelConfigRequest) => modelConfigApi.create(data),
    onSuccess: () => {
      message.success('配置创建成功');
      setModalOpen(false);
      form.resetFields();
      queryClient.invalidateQueries({ queryKey: ['model-configs'] });
    },
    onError: (err: any) => {
      message.error(err?.response?.data?.message || '创建失败');
    },
  });

  // 更新模型配置
  const updateMutation = useMutation({
    mutationFn: ({ id, data }: { id: string; data: any }) => modelConfigApi.update(id, data),
    onSuccess: () => {
      message.success('配置更新成功');
      setModalOpen(false);
      setEditingConfig(null);
      form.resetFields();
      queryClient.invalidateQueries({ queryKey: ['model-configs'] });
    },
    onError: (err: any) => {
      message.error(err?.response?.data?.message || '更新失败');
    },
  });

  // 删除模型配置
  const deleteMutation = useMutation({
    mutationFn: (id: string) => modelConfigApi.delete(id),
    onSuccess: () => {
      message.success('配置已删除');
      queryClient.invalidateQueries({ queryKey: ['model-configs'] });
    },
    onError: () => {
      message.error('删除失败');
    },
  });

  // 验证模型配置
  const validateMutation = useMutation({
    mutationFn: modelConfigApi.validate,
    onSuccess: () => {
      message.success('验证成功，API Key 有效');
    },
    onError: (err: any) => {
      message.error(err?.response?.data?.message || '验证失败');
    },
  });

  const handleSubmit = () => {
    form.validateFields().then((values) => {
      if (editingConfig) {
        updateMutation.mutate({
          id: editingConfig.id,
          data: values,
        });
      } else {
        createMutation.mutate(values);
      }
    });
  };

  const handleEdit = (config: ModelConfig) => {
    setEditingConfig(config);
    setSelectedProvider(config.provider_id);
    form.setFieldsValue({
      provider_id: config.provider_id,
      model_name: config.model_name,
      purpose: config.purpose,
      base_url: config.base_url,
      is_active: config.is_active,
    });
    setModalOpen(true);
  };

  const handleValidate = () => {
    const values = form.getFieldsValue();
    if (!values.provider_id || !values.api_key) {
      message.warning('请先填写提供商和 API Key');
      return;
    }
    validateMutation.mutate({
      provider_id: values.provider_id,
      api_key: values.api_key,
      base_url: values.base_url,
    });
  };

  const getPurposeTag = (purpose: string) => {
    const colors: Record<string, string> = {
      architecture: 'blue',
      chapter: 'green',
      writing: 'purple',
      review: 'orange',
      general: 'default',
    };
    const labels: Record<string, string> = {
      architecture: '架构生成',
      chapter: '章节生成',
      writing: '写作辅助',
      review: 'AI 审阅',
      general: '通用',
    };
    return <Tag color={colors[purpose]}>{labels[purpose] || purpose}</Tag>;
  };

  const columns = [
    {
      title: '提供商',
      dataIndex: ['provider', 'display_name'],
      key: 'provider',
      render: (_: any, record: ModelConfig) => record.provider?.display_name || '-',
    },
    {
      title: '模型',
      dataIndex: 'model_name',
      key: 'model_name',
    },
    {
      title: '用途',
      dataIndex: 'purpose',
      key: 'purpose',
      render: (purpose: string) => getPurposeTag(purpose),
    },
    {
      title: '状态',
      dataIndex: 'is_active',
      key: 'is_active',
      render: (isActive: boolean) => (
        isActive ? (
          <Tag icon={<CheckCircleOutlined />} color="success">启用</Tag>
        ) : (
          <Tag color="default">禁用</Tag>
        )
      ),
    },
    {
      title: '操作',
      key: 'action',
      render: (_: any, record: ModelConfig) => (
        <Space>
          <Button
            type="link"
            size="small"
            icon={<EditOutlined />}
            onClick={() => handleEdit(record)}
          >
            编辑
          </Button>
          <Popconfirm
            title="确定要删除这个配置吗？"
            onConfirm={() => deleteMutation.mutate(record.id)}
          >
            <Button type="link" size="small" danger icon={<DeleteOutlined />}>
              删除
            </Button>
          </Popconfirm>
        </Space>
      ),
    },
  ];

  const selectedProviderName = providers.find((p) => p.id === selectedProvider)?.name;
  const modelOptions = selectedProviderName ? MODEL_OPTIONS[selectedProviderName] || [] : [];

  const tabItems = [
    {
      key: 'models',
      label: (
        <span>
          <ApiOutlined />
          模型配置
        </span>
      ),
      children: (
        <div>
          <div className="flex justify-between items-center mb-4">
            <div>
              <h3 className="text-lg font-semibold">模型配置</h3>
              <p className="text-gray-500 text-sm">
                配置不同用途的 AI 模型，支持多个提供商
              </p>
            </div>
            <Button
              type="primary"
              icon={<PlusOutlined />}
              onClick={() => {
                setEditingConfig(null);
                form.resetFields();
                setSelectedProvider(null);
                setModalOpen(true);
              }}
            >
              添加配置
            </Button>
          </div>

          <Table
            columns={columns}
            dataSource={configs}
            rowKey="id"
            loading={isLoading}
            pagination={false}
          />

          {configs.length === 0 && !isLoading && (
            <div className="text-center py-8 text-gray-500">
              <p className="mb-4">还没有配置任何模型</p>
              <p className="text-sm">添加模型配置后，AI 功能才能正常使用</p>
            </div>
          )}
        </div>
      ),
    },
    {
      key: 'general',
      label: (
        <span>
          <SettingOutlined />
          通用设置
        </span>
      ),
      children: (
        <div className="max-w-xl">
          <h3 className="text-lg font-semibold mb-4">通用设置</h3>
          
          <div className="space-y-6">
            <div className="flex justify-between items-center p-4 bg-gray-50 dark:bg-gray-800 rounded-lg">
              <div>
                <div className="font-medium">深色模式</div>
                <div className="text-gray-500 text-sm">切换应用的显示主题</div>
              </div>
              <Switch
                checked={theme === 'dark'}
                onChange={toggleTheme}
              />
            </div>

            <Divider />

            <div className="text-gray-500 text-sm">
              <p>更多设置选项即将推出...</p>
            </div>
          </div>
        </div>
      ),
    },
  ];

  return (
    <div>
      <h2 className="text-2xl font-bold mb-6">设置</h2>

      <Card>
        <Tabs items={tabItems} />
      </Card>

      {/* 模型配置弹窗 */}
      <Modal
        title={editingConfig ? '编辑模型配置' : '添加模型配置'}
        open={modalOpen}
        onCancel={() => {
          setModalOpen(false);
          setEditingConfig(null);
          form.resetFields();
        }}
        footer={[
          <Button key="cancel" onClick={() => setModalOpen(false)}>
            取消
          </Button>,
          !editingConfig && (
            <Button
              key="validate"
              onClick={handleValidate}
              loading={validateMutation.isPending}
            >
              验证 API
            </Button>
          ),
          <Button
            key="submit"
            type="primary"
            onClick={handleSubmit}
            loading={createMutation.isPending || updateMutation.isPending}
          >
            {editingConfig ? '保存' : '创建'}
          </Button>,
        ].filter(Boolean)}
        width={600}
      >
        <Form form={form} layout="vertical" className="mt-4">
          <Form.Item
            name="provider_id"
            label="提供商"
            rules={[{ required: true, message: '请选择提供商' }]}
          >
            <Select
              placeholder="选择 AI 提供商"
              options={providers.map((p) => ({
                label: p.display_name,
                value: p.id,
              }))}
              onChange={(value) => {
                setSelectedProvider(value);
                form.setFieldValue('model_name', undefined);
              }}
            />
          </Form.Item>

          <Form.Item
            name="model_name"
            label="模型名称"
            rules={[{ required: true, message: '请输入或选择模型名称' }]}
          >
            {modelOptions.length > 0 ? (
              <Select
                placeholder="选择或输入模型名称"
                showSearch
                allowClear
                options={modelOptions.map((m) => ({ label: m, value: m }))}
              />
            ) : (
              <Input placeholder="输入模型名称，如 gpt-4o" />
            )}
          </Form.Item>

          <Form.Item
            name="purpose"
            label="用途"
            rules={[{ required: true, message: '请选择用途' }]}
          >
            <Select
              placeholder="选择此配置的用途"
              options={PURPOSE_OPTIONS}
            />
          </Form.Item>

          <Form.Item
            name="api_key"
            label="API Key"
            rules={[{ required: !editingConfig, message: '请输入 API Key' }]}
            extra={editingConfig ? '留空表示不修改' : undefined}
          >
            <Input.Password
              placeholder={editingConfig ? '留空表示不修改' : '输入 API Key'}
            />
          </Form.Item>

          <Form.Item
            name="base_url"
            label="自定义 Base URL（可选）"
            extra="如果使用代理或自定义端点，请填写完整的 API 地址"
          >
            <Input placeholder="例如: https://api.example.com/v1" />
          </Form.Item>

          {editingConfig && (
            <Form.Item
              name="is_active"
              label="启用状态"
              valuePropName="checked"
            >
              <Switch />
            </Form.Item>
          )}
        </Form>
      </Modal>
    </div>
  );
}

export default Settings;
