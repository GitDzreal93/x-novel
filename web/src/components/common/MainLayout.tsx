import { Layout, theme } from 'antd';
import { Outlet, useNavigate, useLocation } from 'react-router-dom';
import {
  BookOutlined,
  SettingOutlined,
  SunOutlined,
  MoonOutlined,
  MenuFoldOutlined,
  MenuUnfoldOutlined,
} from '@ant-design/icons';
import { useAppStore } from '../../stores';
import styles from './MainLayout.module.css';

const { Header, Sider, Content } = Layout;

function MainLayout() {
  const navigate = useNavigate();
  const location = useLocation();
  const { token } = theme.useToken();
  const { theme: appTheme, toggleTheme, sidebarCollapsed, toggleSidebar } = useAppStore();

  const menuItems = [
    {
      key: '/projects',
      icon: <BookOutlined />,
      label: '我的项目',
      onClick: () => navigate('/projects'),
    },
    {
      key: '/settings',
      icon: <SettingOutlined />,
      label: '设置',
      onClick: () => navigate('/settings'),
    },
  ];

  // 根据路径确定选中的菜单项
  const getSelectedKey = () => {
    if (location.pathname.startsWith('/projects')) return '/projects';
    if (location.pathname.startsWith('/settings')) return '/settings';
    return location.pathname;
  };
  const selectedKey = getSelectedKey();

  return (
    <Layout className={styles.layout}>
      <Sider
        trigger={null}
        collapsible
        collapsed={sidebarCollapsed}
        className={styles.sider}
      >
        <div className={styles.logo}>
          <h1>{sidebarCollapsed ? 'XN' : 'X-Novel'}</h1>
        </div>
        <div className={styles.menu}>
          {menuItems.map((item) => (
            <div
              key={item.key}
              className={`${styles.menuItem} ${selectedKey === item.key ? styles.menuItemActive : ''}`}
              onClick={item.onClick}
            >
              <div className={styles.menuIcon}>{item.icon}</div>
              {!sidebarCollapsed && <span>{item.label}</span>}
            </div>
          ))}
        </div>
      </Sider>
      <Layout>
        <Header className={styles.header} style={{ background: token.colorBgContainer }}>
          <div className={styles.headerLeft}>
            {menuItems.find((item) => item.key === selectedKey)?.icon}
            <span className={styles.headerTitle}>
              {menuItems.find((item) => item.key === selectedKey)?.label}
            </span>
          </div>
          <div className={styles.headerRight}>
            <div className={styles.headerAction} onClick={toggleSidebar}>
              {sidebarCollapsed ? <MenuUnfoldOutlined /> : <MenuFoldOutlined />}
            </div>
            <div className={styles.headerAction} onClick={toggleTheme}>
              {appTheme === 'light' ? <MoonOutlined /> : <SunOutlined />}
            </div>
          </div>
        </Header>
        <Content className={styles.content}>
          <Outlet />
        </Content>
      </Layout>
    </Layout>
  );
}

export default MainLayout;
