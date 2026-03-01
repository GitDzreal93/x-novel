import { Outlet, useNavigate, useLocation } from 'react-router-dom';
import { Layout, Menu, Button, Avatar, Dropdown, Typography, theme, Flex } from 'antd';
import {
  BookOutlined,
  SettingOutlined,
  SunOutlined,
  MoonOutlined,
  MenuFoldOutlined,
  MenuUnfoldOutlined,
  UserOutlined,
  MessageOutlined,
} from '@ant-design/icons';
import { useAppStore } from '../../stores';

const { Sider, Header, Content } = Layout;
const { Text } = Typography;

function MainLayout() {
  const navigate = useNavigate();
  const location = useLocation();
  const { theme: appTheme, toggleTheme, sidebarCollapsed, toggleSidebar } = useAppStore();
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
    <Layout style={{ minHeight: '100vh' }}>
      <Sider
        collapsible
        collapsed={sidebarCollapsed}
        onCollapse={toggleSidebar}
        trigger={null}
        width={220}
        collapsedWidth={64}
        style={{
          background: token.colorBgContainer,
          borderRight: `1px solid ${token.colorBorderSecondary}`,
          overflow: 'auto',
          height: '100vh',
          position: 'fixed',
          left: 0,
          top: 0,
          bottom: 0,
          zIndex: 10,
        }}
      >
        {/* Logo */}
        <Flex
          align="center"
          justify={sidebarCollapsed ? 'center' : 'flex-start'}
          gap={10}
          style={{
            height: 64,
            padding: sidebarCollapsed ? '0' : '0 20px',
            borderBottom: `1px solid ${token.colorBorderSecondary}`,
          }}
        >
          <div
            style={{
              width: 32,
              height: 32,
              borderRadius: 10,
              background: 'linear-gradient(135deg, #4f46e5, #7c3aed)',
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
              color: '#fff',
              fontWeight: 700,
              fontSize: 14,
              flexShrink: 0,
            }}
          >
            X
          </div>
          {!sidebarCollapsed && (
            <Text
              strong
              style={{
                fontSize: 17,
                background: 'linear-gradient(90deg, #4f46e5, #7c3aed)',
                WebkitBackgroundClip: 'text',
                WebkitTextFillColor: 'transparent',
                whiteSpace: 'nowrap',
                letterSpacing: '-0.02em',
              }}
            >
              X-Novel
            </Text>
          )}
        </Flex>

        {/* Menu */}
        <Menu
          mode="inline"
          selectedKeys={[getSelectedKey()]}
          items={menuItems}
          onClick={({ key }) => navigate(key)}
          style={{
            border: 'none',
            padding: '8px',
          }}
        />

        {/* 底部版本 */}
        {!sidebarCollapsed && (
          <div
            style={{
              position: 'absolute',
              bottom: 0,
              left: 0,
              right: 0,
              padding: '12px 16px',
              borderTop: `1px solid ${token.colorBorderSecondary}`,
              textAlign: 'center',
            }}
          >
            <Text type="secondary" style={{ fontSize: 11 }}>
              X-Novel v0.1.0
            </Text>
          </div>
        )}
      </Sider>

      <Layout
        style={{
          marginLeft: sidebarCollapsed ? 64 : 220,
          transition: 'margin-left 0.2s',
        }}
      >
        {/* Header */}
        <Header
          style={{
            background: token.colorBgContainer,
            borderBottom: `1px solid ${token.colorBorderSecondary}`,
            padding: '0 24px',
            height: 64,
            lineHeight: '64px',
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'space-between',
            position: 'sticky',
            top: 0,
            zIndex: 9,
          }}
        >
          <Flex align="center" gap={12}>
            <Button
              type="text"
              icon={sidebarCollapsed ? <MenuUnfoldOutlined /> : <MenuFoldOutlined />}
              onClick={toggleSidebar}
              style={{ fontSize: 16 }}
            />
          </Flex>

          <Flex align="center" gap={8}>
            <Button
              type="text"
              icon={appTheme === 'light' ? <MoonOutlined /> : <SunOutlined />}
              onClick={toggleTheme}
              title={appTheme === 'light' ? '切换暗色模式' : '切换亮色模式'}
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
                  padding: '4px 10px',
                  borderRadius: token.borderRadius,
                }}
              >
                <Avatar
                  size={30}
                  icon={<UserOutlined />}
                  style={{
                    background: 'linear-gradient(135deg, #818cf8, #a855f7)',
                  }}
                />
                <Text style={{ fontSize: 14, fontWeight: 500 }}>Admin</Text>
              </Flex>
            </Dropdown>
          </Flex>
        </Header>

        {/* Content */}
        <Content
          style={{
            padding: 24,
            minHeight: 'calc(100vh - 64px)',
            overflow: 'auto',
          }}
        >
          <div style={{ maxWidth: 1400, margin: '0 auto' }}>
            <Outlet />
          </div>
        </Content>
      </Layout>
    </Layout>
  );
}

export default MainLayout;
