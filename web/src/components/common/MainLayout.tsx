import { Outlet, useNavigate, useLocation } from 'react-router-dom';
import { Layout, Menu, Button, Avatar, Dropdown, Typography, theme, Flex, Space } from 'antd';
import {
  BookOutlined,
  SettingOutlined,
  SunOutlined,
  MoonOutlined,
  UserOutlined,
  MessageOutlined,
  GithubOutlined,
} from '@ant-design/icons';
import { useAppStore } from '../../stores';

const { Header, Content, Footer } = Layout;
const { Text } = Typography;

function MainLayout() {
  const navigate = useNavigate();
  const location = useLocation();
  const { theme: appTheme, toggleTheme } = useAppStore();
  const { token } = theme.useToken();

  const menuItems = [
    {
      key: '/projects',
      icon: <BookOutlined />,
      label: '我的项目',
    },
    {
      key: '/chat',
      icon: <MessageOutlined />,
      label: '灵感工坊',
    },
    {
      key: '/settings',
      icon: <SettingOutlined />,
      label: '系统设置',
    },
  ];

  const getSelectedKey = () => {
    if (location.pathname.startsWith('/projects')) return '/projects';
    if (location.pathname.startsWith('/chat')) return '/chat';
    if (location.pathname.startsWith('/settings')) return '/settings';
    return location.pathname;
  };

  return (
    <Layout style={{ minHeight: '100vh', background: token.colorBgLayout }}>
      {/* 顶部导航栏 */}
      <Header
        style={{
          background: token.colorBgContainer,
          borderBottom: `1px solid ${token.colorBorderSecondary}`,
          padding: 0,
          height: 56,
          lineHeight: '56px',
          position: 'sticky',
          top: 0,
          zIndex: 100,
          boxShadow: '0 1px 2px 0 rgba(0,0,0,0.03)',
        }}
      >
        <div
          style={{
            maxWidth: 1200,
            margin: '0 auto',
            padding: '0 32px',
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'space-between',
            height: '100%',
          }}
        >
          {/* 左：Logo + 导航 */}
          <Flex align="center" gap={32}>
            <Flex
              align="center"
              gap={10}
              style={{ cursor: 'pointer', flexShrink: 0 }}
              onClick={() => navigate('/projects')}
            >
              <div
                style={{
                  width: 30,
                  height: 30,
                  borderRadius: 8,
                  background: 'linear-gradient(135deg, #ab372f, #88100f)',
                  display: 'flex',
                  alignItems: 'center',
                  justifyContent: 'center',
                  color: '#fff',
                  fontWeight: 700,
                  fontSize: 13,
                }}
              >
                X
              </div>
              <Text
                strong
                style={{
                  fontSize: 16,
                  background: 'linear-gradient(90deg, #ab372f, #88100f)',
                  WebkitBackgroundClip: 'text',
                  WebkitTextFillColor: 'transparent',
                  letterSpacing: '-0.02em',
                }}
              >
                X-Novel
              </Text>
            </Flex>

            <Menu
              mode="horizontal"
              selectedKeys={[getSelectedKey()]}
              items={menuItems}
              onClick={({ key }) => navigate(key)}
              style={{
                border: 'none',
                background: 'transparent',
                lineHeight: '54px',
                fontSize: 14,
              }}
            />
          </Flex>

          {/* 右：操作区 */}
          <Flex align="center" gap={4}>
            <Button
              type="text"
              size="small"
              icon={appTheme === 'light' ? <MoonOutlined /> : <SunOutlined />}
              onClick={toggleTheme}
              title={appTheme === 'light' ? '切换暗色模式' : '切换亮色模式'}
              style={{ color: token.colorTextSecondary }}
            />

            <Dropdown
              menu={{
                items: [
                  { key: 'profile', label: '个人信息' },
                  { type: 'divider' },
                  { key: 'logout', label: '退出登录', danger: true },
                ],
              }}
              placement="bottomRight"
            >
              <Flex
                align="center"
                gap={8}
                style={{
                  cursor: 'pointer',
                  padding: '4px 12px',
                  borderRadius: token.borderRadius,
                  transition: 'background 0.2s',
                }}
              >
                <Avatar
                  size={28}
                  icon={<UserOutlined />}
                  style={{
                    background: 'linear-gradient(135deg, #e2695d, #ab372f)',
                  }}
                />
                <Text style={{ fontSize: 13, fontWeight: 500 }}>Admin</Text>
              </Flex>
            </Dropdown>
          </Flex>
        </div>
      </Header>

      {/* 内容区 */}
      <Content style={{ flex: 1 }}>
        <div
          style={{
            maxWidth: 1200,
            margin: '0 auto',
            padding: '28px 32px',
            minHeight: 'calc(100vh - 56px - 64px)',
          }}
        >
          <Outlet />
        </div>
      </Content>

      {/* 页脚 */}
      <Footer
        style={{
          textAlign: 'center',
          padding: '16px 32px',
          background: 'transparent',
        }}
      >
        <Space split={<span style={{ color: token.colorBorderSecondary }}>·</span>}>
          <Text type="secondary" style={{ fontSize: 12 }}>
            X-Novel v0.1.0
          </Text>
          <Text type="secondary" style={{ fontSize: 12 }}>
            AI 驱动的小说创作平台
          </Text>
          <a
            href="https://github.com"
            target="_blank"
            rel="noreferrer"
            style={{ color: token.colorTextQuaternary, fontSize: 12 }}
          >
            <GithubOutlined style={{ marginRight: 4 }} />
            GitHub
          </a>
        </Space>
      </Footer>
    </Layout>
  );
}

export default MainLayout;
