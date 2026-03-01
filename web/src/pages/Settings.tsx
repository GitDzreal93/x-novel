import { useState, useRef } from 'react';
import {
  Form, Input, Select, Button, Table, Modal, Space, Tag, Switch,
  Popconfirm, App, Tabs, Typography, Card, Flex, theme, Statistic,
  Upload, Spin, Result, Divider,
} from 'antd';
import {
  PlusOutlined,
  DeleteOutlined,
  EditOutlined,
  CheckCircleOutlined,
  ApiOutlined,
  SettingOutlined,
  KeyOutlined,
  CloudDownloadOutlined,
  CloudUploadOutlined,
  DatabaseOutlined,
  FileTextOutlined,
  MessageOutlined,
  BookOutlined,
  LinkOutlined,
} from '@ant-design/icons';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { modelConfigApi, backupApi } from '../api';
import { useAppStore } from '../stores';
import type { ModelConfig, ModelProvider, ModelBinding, CreateModelConfigRequest, BindingPurpose, ImportResult } from '../types';

const { Title, Text } = Typography;

const BINDING_PURPOSES: { key: BindingPurpose; label: string; description: string }[] = [
  { key: 'architecture', label: '架构生成', description: '生成小说世界观、人物关系等整体架构' },
  { key: 'chapter', label: '章节生成', description: '生成章节内容和大纲' },
  { key: 'writing', label: '写作助手', description: '文本润色、续写、建议等辅助功能' },
  { key: 'review', label: 'AI 审阅', description: '错误检测、质量评审、市场预测' },
  { key: 'general', label: '通用 / 对话', description: '灵感对话等通用 AI 交互' },
];

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
  const { theme: appTheme, toggleTheme } = useAppStore();
  const { token } = theme.useToken();
  const [modalOpen, setModalOpen] = useState(false);
  const [editingConfig, setEditingConfig] = useState<ModelConfig | null>(null);
  const [form] = Form.useForm();
  const [selectedProvider, setSelectedProvider] = useState<number | null>(null);

  const { data: providersRes } = useQuery({
    queryKey: ['model-providers'],
    queryFn: () => modelConfigApi.listProviders().then((res) => res.data || []),
  });
  const providers = providersRes || [];

  const { data: configsRes, isLoading } = useQuery({
    queryKey: ['model-configs'],
    queryFn: () => modelConfigApi.list({ page: 1, page_size: 100 }).then((res) => res.data),
  });
  const configs = configsRes?.configs || [];

  const { data: bindingsRes } = useQuery({
    queryKey: ['model-bindings'],
    queryFn: () => modelConfigApi.listBindings().then((res) => res.data || []),
  });
  const bindings: ModelBinding[] = bindingsRes || [];

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

  const deleteMutation = useMutation({
    mutationFn: (id: string) => modelConfigApi.delete(id),
    onSuccess: () => {
      message.success('配置已删除');
      queryClient.invalidateQueries({ queryKey: ['model-configs'] });
      queryClient.invalidateQueries({ queryKey: ['model-bindings'] });
    },
    onError: () => {
      message.error('删除失败');
    },
  });

  const validateMutation = useMutation({
    mutationFn: modelConfigApi.validate,
    onSuccess: () => {
      message.success('验证成功，API Key 有效');
    },
    onError: (err: any) => {
      message.error(err?.response?.data?.message || '验证失败');
    },
  });

  const bindingMutation = useMutation({
    mutationFn: (data: { purpose: BindingPurpose; model_config_id: string }) =>
      modelConfigApi.upsertBinding(data),
    onSuccess: () => {
      message.success('绑定已更新');
      queryClient.invalidateQueries({ queryKey: ['model-bindings'] });
    },
    onError: () => {
      message.error('绑定更新失败');
    },
  });

  const unbindMutation = useMutation({
    mutationFn: (purpose: string) => modelConfigApi.deleteBinding(purpose),
    onSuccess: () => {
      message.success('已取消绑定');
      queryClient.invalidateQueries({ queryKey: ['model-bindings'] });
    },
  });

  const handleSubmit = () => {
    form.validateFields().then((values) => {
      if (editingConfig) {
        updateMutation.mutate({ id: editingConfig.id, data: values });
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

  const getBindingForPurpose = (purpose: BindingPurpose): ModelBinding | undefined => {
    return bindings.find((b) => b.purpose === purpose);
  };

  const configSelectOptions = configs.map((c: ModelConfig) => ({
    label: `${c.provider?.display_name || '未知'} / ${c.model_name}`,
    value: c.id,
  }));

  const columns = [
    {
      title: '提供商',
      dataIndex: ['provider', 'display_name'],
      key: 'provider',
      render: (_: any, record: ModelConfig) => (
        <Text strong>{record.provider?.display_name || '-'}</Text>
      ),
    },
    {
      title: '模型',
      dataIndex: 'model_name',
      key: 'model_name',
      render: (text: string) => (
        <Text type="secondary" code style={{ fontSize: 13 }}>{text}</Text>
      ),
    },
    {
      title: '状态',
      dataIndex: 'is_active',
      key: 'is_active',
      render: (isActive: boolean) =>
        isActive ? (
          <Tag icon={<CheckCircleOutlined />} color="success" bordered={false}>
            已启用
          </Tag>
        ) : (
          <Tag color="default" bordered={false}>已禁用</Tag>
        ),
    },
    {
      title: '操作',
      key: 'action',
      render: (_: any, record: ModelConfig) => (
        <Space size="middle">
          <Button type="text" size="small" icon={<EditOutlined />} onClick={() => handleEdit(record)}>
            编辑
          </Button>
          <Popconfirm
            title="删除配置"
            description="确定要删除这个模型配置吗？关联的功能绑定也会被清除。"
            okText="确定"
            cancelText="取消"
            onConfirm={() => deleteMutation.mutate(record.id)}
            okButtonProps={{ danger: true }}
          >
            <Button type="text" size="small" danger icon={<DeleteOutlined />}>
              删除
            </Button>
          </Popconfirm>
        </Space>
      ),
    },
  ];

  const selectedProviderName = providers.find((p: ModelProvider) => p.id === selectedProvider)?.name;
  const modelOptions = selectedProviderName ? MODEL_OPTIONS[selectedProviderName] || [] : [];

  const tabItems = [
    {
      key: 'models',
      label: (
        <Flex align="center" gap={6}>
          <ApiOutlined />
          模型配置
        </Flex>
      ),
      children: (
        <div style={{ padding: '8px 0' }}>
          {/* 模型配置列表 */}
          <Flex
            justify="space-between"
            align="center"
            wrap
            gap={16}
            style={{ marginBottom: 20 }}
          >
            <div>
              <Flex align="center" gap={8}>
                <Title level={5} style={{ margin: 0 }}>AI 模型配置</Title>
                <KeyOutlined style={{ color: token.colorPrimary }} />
              </Flex>
              <Text type="secondary" style={{ fontSize: 13 }}>
                添加 AI 模型的连接信息（提供商、模型、API Key），然后在下方为各功能选择使用哪个模型
              </Text>
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
              size="large"
            >
              添加配置
            </Button>
          </Flex>

          <Table
            columns={columns}
            dataSource={configs}
            rowKey="id"
            loading={isLoading}
            pagination={false}
            locale={{
              emptyText: (
                <div style={{ padding: '48px 0' }}>
                  <ApiOutlined style={{ fontSize: 36, color: token.colorTextQuaternary, marginBottom: 12, display: 'block' }} />
                  <Text strong>还没有配置任何模型</Text>
                  <br />
                  <Text type="secondary" style={{ fontSize: 13 }}>
                    添加模型配置后，才能在下方为各个功能绑定模型
                  </Text>
                </div>
              ),
            }}
          />

          {/* 功能绑定 */}
          <Divider />
          <div style={{ marginBottom: 20 }}>
            <Flex align="center" gap={8}>
              <Title level={5} style={{ margin: 0 }}>功能绑定</Title>
              <LinkOutlined style={{ color: token.colorPrimary }} />
            </Flex>
            <Text type="secondary" style={{ fontSize: 13 }}>
              为每个 AI 功能选择使用哪个模型配置，未绑定的功能将自动使用"通用"配置
            </Text>
          </div>

          <Flex vertical gap={12}>
            {BINDING_PURPOSES.map(({ key, label, description }) => {
              const binding = getBindingForPurpose(key);
              return (
                <Flex
                  key={key}
                  justify="space-between"
                  align="center"
                  style={{
                    padding: '14px 20px',
                    background: token.colorBgLayout,
                    border: `1px solid ${token.colorBorderSecondary}`,
                    borderRadius: token.borderRadius,
                  }}
                >
                  <div style={{ minWidth: 160 }}>
                    <Text strong>{label}</Text>
                    <br />
                    <Text type="secondary" style={{ fontSize: 12 }}>{description}</Text>
                  </div>
                  <Flex align="center" gap={8}>
                    <Select
                      style={{ width: 280 }}
                      placeholder="选择模型配置"
                      allowClear
                      value={binding?.model_config_id}
                      options={configSelectOptions}
                      onChange={(value) => {
                        if (value) {
                          bindingMutation.mutate({ purpose: key, model_config_id: value });
                        } else {
                          unbindMutation.mutate(key);
                        }
                      }}
                      notFoundContent={
                        <Text type="secondary" style={{ fontSize: 13 }}>
                          请先在上方添加模型配置
                        </Text>
                      }
                    />
                  </Flex>
                </Flex>
              );
            })}
          </Flex>
        </div>
      ),
    },
    {
      key: 'backup',
      label: (
        <Flex align="center" gap={6}>
          <DatabaseOutlined />
          数据管理
        </Flex>
      ),
      children: <BackupTab />,
    },
    {
      key: 'general',
      label: (
        <Flex align="center" gap={6}>
          <SettingOutlined />
          通用设置
        </Flex>
      ),
      children: (
        <div style={{ maxWidth: 640, padding: '8px 0' }}>
          <div style={{ marginBottom: 20 }}>
            <Title level={5} style={{ margin: 0 }}>系统偏好</Title>
            <Text type="secondary" style={{ fontSize: 13 }}>自定义你的使用体验</Text>
          </div>

          <Flex vertical gap={12}>
            <Flex
              justify="space-between"
              align="center"
              style={{
                padding: 20,
                background: token.colorBgLayout,
                border: `1px solid ${token.colorBorderSecondary}`,
                borderRadius: token.borderRadius,
              }}
            >
              <div>
                <Text strong>深色模式</Text>
                <br />
                <Text type="secondary" style={{ fontSize: 13 }}>切换应用的显示主题外观</Text>
              </div>
              <Switch checked={appTheme === 'dark'} onChange={toggleTheme} />
            </Flex>

            <Flex
              justify="space-between"
              align="center"
              style={{
                padding: 20,
                background: token.colorBgLayout,
                border: `1px solid ${token.colorBorderSecondary}`,
                borderRadius: token.borderRadius,
                opacity: 0.5,
                cursor: 'not-allowed',
              }}
            >
              <div>
                <Flex align="center" gap={8}>
                  <Text strong>自动保存</Text>
                  <Tag color="blue" bordered={false} style={{ fontSize: 11 }}>即将推出</Tag>
                </Flex>
                <Text type="secondary" style={{ fontSize: 13 }}>写作内容自动保存的间隔时间</Text>
              </div>
              <Switch disabled />
            </Flex>
          </Flex>
        </div>
      ),
    },
  ];

  return (
    <>
      <div style={{ marginBottom: 24 }}>
        <Title level={4} style={{ margin: 0 }}>系统设置</Title>
      </div>

      <Card styles={{ body: { padding: 0 } }}>
        <Tabs
          items={tabItems}
          size="large"
          tabBarStyle={{
            padding: '0 24px',
            marginBottom: 0,
          }}
          style={{ padding: '0 24px 24px' }}
        />
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
            <Button key="validate" onClick={handleValidate} loading={validateMutation.isPending}>
              测试连接
            </Button>
          ),
          <Button
            key="submit"
            type="primary"
            onClick={handleSubmit}
            loading={createMutation.isPending || updateMutation.isPending}
          >
            {editingConfig ? '保存修改' : '确认创建'}
          </Button>,
        ].filter(Boolean)}
        width={580}
        centered
      >
        <Form form={form} layout="vertical" style={{ marginTop: 16 }}>
          <Form.Item
            name="provider_id"
            label="AI 提供商"
            rules={[{ required: true, message: '请选择提供商' }]}
          >
            <Select
              size="large"
              placeholder="选择 AI 服务提供商"
              options={providers.map((p: ModelProvider) => ({ label: p.display_name, value: p.id }))}
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
                size="large"
                placeholder="选择或输入模型名称"
                showSearch
                allowClear
                options={modelOptions.map((m) => ({ label: m, value: m }))}
              />
            ) : (
              <Input size="large" placeholder="输入模型名称，如 gpt-4o" />
            )}
          </Form.Item>

          <Form.Item
            name="api_key"
            label="API Key"
            rules={[{ required: !editingConfig, message: '请输入 API Key' }]}
            extra={
              <Text type="secondary" style={{ fontSize: 12 }}>
                {editingConfig ? '留空表示不修改现有 Key' : '提供商分配的身份验证凭证'}
              </Text>
            }
          >
            <Input.Password size="large" placeholder={editingConfig ? '留空表示不修改' : 'sk-...'} />
          </Form.Item>

          <Form.Item
            name="base_url"
            label="自定义接口地址（可选）"
            extra={
              <Text type="secondary" style={{ fontSize: 12 }}>
                如果使用 API 代理或本地模型服务，请填写完整的 Base URL
              </Text>
            }
          >
            <Input
              size="large"
              placeholder="例如: https://api.example.com/v1"
              style={{ fontFamily: 'monospace', fontSize: 13 }}
            />
          </Form.Item>

          {editingConfig && (
            <Form.Item name="is_active" label="启用状态" valuePropName="checked">
              <Switch />
            </Form.Item>
          )}
        </Form>
      </Modal>
    </>
  );
}

function BackupTab() {
  const { message: messageApi, modal } = App.useApp();
  const { token } = theme.useToken();
  const fileInputRef = useRef<HTMLInputElement>(null);
  const [exporting, setExporting] = useState(false);
  const [importing, setImporting] = useState(false);
  const [importResult, setImportResult] = useState<ImportResult | null>(null);

  const { data: preview, isLoading: previewLoading } = useQuery({
    queryKey: ['backup-preview'],
    queryFn: () => backupApi.preview().then((res) => res.data),
  });

  const handleExport = async () => {
    setExporting(true);
    try {
      const res = await backupApi.exportData();
      const blob = res instanceof Blob ? res : new Blob([JSON.stringify(res)], { type: 'application/json' });
      const url = URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = `x-novel-backup-${new Date().toISOString().slice(0, 10)}.json`;
      document.body.appendChild(a);
      a.click();
      document.body.removeChild(a);
      URL.revokeObjectURL(url);
      messageApi.success('数据导出成功');
    } catch {
      messageApi.error('导出失败，请稍后重试');
    } finally {
      setExporting(false);
    }
  };

  const handleImportClick = () => {
    fileInputRef.current?.click();
  };

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;

    if (!file.name.endsWith('.json')) {
      messageApi.error('请选择 JSON 格式的备份文件');
      return;
    }

    modal.confirm({
      title: '确认导入数据',
      content: `将从文件 "${file.name}" 导入数据。导入的项目和对话将作为新数据添加，不会覆盖现有数据。`,
      okText: '确认导入',
      cancelText: '取消',
      centered: true,
      onOk: () => doImport(file),
    });

    if (fileInputRef.current) {
      fileInputRef.current.value = '';
    }
  };

  const doImport = async (file: File) => {
    setImporting(true);
    setImportResult(null);
    try {
      const res = await backupApi.importData(file);
      if (res?.data) {
        setImportResult(res.data);
        messageApi.success('数据导入完成');
      }
    } catch {
      messageApi.error('导入失败，请检查文件格式');
    } finally {
      setImporting(false);
    }
  };

  return (
    <div style={{ maxWidth: 800, padding: '8px 0' }}>
      <div style={{ marginBottom: 20 }}>
        <Title level={5} style={{ margin: 0 }}>数据管理</Title>
        <Text type="secondary" style={{ fontSize: 13 }}>
          导出和导入你的所有创作数据，支持跨设备迁移
        </Text>
      </div>

      {/* 当前数据概览 */}
      <Card title="当前数据概览" style={{ marginBottom: 20 }}>
        {previewLoading ? (
          <Flex justify="center" style={{ padding: 24 }}><Spin /></Flex>
        ) : preview ? (
          <Flex wrap gap={32}>
            <Statistic
              title={<Flex align="center" gap={4}><BookOutlined />项目</Flex>}
              value={preview.projects}
              suffix="个"
            />
            <Statistic
              title={<Flex align="center" gap={4}><FileTextOutlined />章节</Flex>}
              value={preview.chapters}
              suffix="章"
            />
            <Statistic
              title="总字数"
              value={preview.total_words}
              suffix="字"
            />
            <Statistic
              title={<Flex align="center" gap={4}><MessageOutlined />对话</Flex>}
              value={preview.conversations}
              suffix="条"
            />
            <Statistic
              title="消息"
              value={preview.messages}
              suffix="条"
            />
          </Flex>
        ) : (
          <Text type="secondary">暂无数据</Text>
        )}
      </Card>

      {/* 导出 */}
      <Card
        style={{ marginBottom: 20 }}
        styles={{ body: { padding: '20px 24px' } }}
      >
        <Flex justify="space-between" align="center">
          <div>
            <Flex align="center" gap={8} style={{ marginBottom: 4 }}>
              <CloudDownloadOutlined style={{ color: token.colorPrimary, fontSize: 18 }} />
              <Text strong style={{ fontSize: 15 }}>导出数据</Text>
            </Flex>
            <Text type="secondary" style={{ fontSize: 13 }}>
              将所有项目、章节、对话数据导出为 JSON 文件
            </Text>
          </div>
          <Button
            type="primary"
            icon={<CloudDownloadOutlined />}
            onClick={handleExport}
            loading={exporting}
            size="large"
          >
            导出
          </Button>
        </Flex>
      </Card>

      {/* 导入 */}
      <Card styles={{ body: { padding: '20px 24px' } }}>
        <Flex justify="space-between" align="center" style={{ marginBottom: importResult ? 16 : 0 }}>
          <div>
            <Flex align="center" gap={8} style={{ marginBottom: 4 }}>
              <CloudUploadOutlined style={{ color: '#52c41a', fontSize: 18 }} />
              <Text strong style={{ fontSize: 15 }}>导入数据</Text>
            </Flex>
            <Text type="secondary" style={{ fontSize: 13 }}>
              从 JSON 备份文件中恢复数据（不会覆盖现有数据）
            </Text>
          </div>
          <Button
            icon={<CloudUploadOutlined />}
            onClick={handleImportClick}
            loading={importing}
            size="large"
          >
            选择文件
          </Button>
          <input
            ref={fileInputRef}
            type="file"
            accept=".json"
            style={{ display: 'none' }}
            onChange={handleFileChange}
          />
        </Flex>

        {importing && (
          <Flex justify="center" style={{ padding: 24 }}>
            <Spin tip="正在导入数据..." />
          </Flex>
        )}

        {importResult && (
          <Result
            status="success"
            title="导入完成"
            subTitle={
              <Flex vertical gap={4} style={{ textAlign: 'left', display: 'inline-block' }}>
                <Text>导入项目：{importResult.imported_projects} 个</Text>
                <Text>导入章节：{importResult.imported_chapters} 章</Text>
                <Text>导入对话：{importResult.imported_conversations} 条</Text>
                <Text>导入消息：{importResult.imported_messages} 条</Text>
                {(importResult.failed_projects > 0 || importResult.failed_chapters > 0) && (
                  <Text type="danger">
                    失败：{importResult.failed_projects} 项目 / {importResult.failed_chapters} 章节
                  </Text>
                )}
              </Flex>
            }
            extra={
              <Button onClick={() => setImportResult(null)}>关闭</Button>
            }
            style={{
              padding: '16px 0',
              background: token.colorBgLayout,
              borderRadius: token.borderRadius,
            }}
          />
        )}
      </Card>
    </div>
  );
}

export default Settings;
