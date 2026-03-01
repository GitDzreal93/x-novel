import { lazy, Suspense } from 'react';
import { ConfigProvider, theme, App as AntApp, Spin, Flex } from 'antd';
import zhCN from 'antd/locale/zh_CN';
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { useAppStore } from './stores';
import MainLayout from './components/common/MainLayout';

const ProjectList = lazy(() => import('./pages/ProjectList'));
const ProjectDetail = lazy(() => import('./pages/ProjectDetail'));
const Chat = lazy(() => import('./pages/Chat'));
const Settings = lazy(() => import('./pages/Settings'));

function PageLoading() {
  return (
    <Flex align="center" justify="center" style={{ height: '60vh' }}>
      <Spin size="large" />
    </Flex>
  );
}

function App() {
  const { theme: appTheme } = useAppStore();

  return (
    <ConfigProvider
      locale={zhCN}
      theme={{
        algorithm: appTheme === 'dark' ? theme.darkAlgorithm : theme.defaultAlgorithm,
        token: {
          colorPrimary: '#ab372f',
          borderRadius: 8,
          colorBgLayout: '#faf8f6',
          fontFamily: '-apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif',
        },
      }}
    >
      <AntApp>
        <BrowserRouter>
          <Suspense fallback={<PageLoading />}>
            <Routes>
              <Route path="/" element={<MainLayout />}>
                <Route index element={<Navigate to="/projects" replace />} />
                <Route path="projects" element={<ProjectList />} />
                <Route path="projects/:id" element={<ProjectDetail />} />
                <Route path="chat" element={<Chat />} />
                <Route path="settings" element={<Settings />} />
              </Route>
            </Routes>
          </Suspense>
        </BrowserRouter>
      </AntApp>
    </ConfigProvider>
  );
}

export default App;
