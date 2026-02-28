import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import type { DeviceSettings } from '../types';

interface AppState {
  // 设备设置
  deviceSettings: DeviceSettings;
  setDeviceSettings: (settings: DeviceSettings) => void;

  // 主题
  theme: 'light' | 'dark';
  setTheme: (theme: 'light' | 'dark') => void;
  toggleTheme: () => void;

  // 当前项目
  currentProjectId: string | null;
  setCurrentProjectId: (id: string | null) => void;

  // 侧边栏状态
  sidebarCollapsed: boolean;
  toggleSidebar: () => void;
}

export const useAppStore = create<AppState>()(
  persist(
    (set, get) => ({
      // 设备设置
      deviceSettings: {
        id: '',
        device_id: '',
        theme: 'light',
        language: 'zh-CN',
        auto_save_enabled: true,
        auto_save_interval: 30000,
      },
      setDeviceSettings: (settings) => set({ deviceSettings: settings }),

      // 主题
      theme: 'light',
      setTheme: (theme) => set({ theme }),
      toggleTheme: () => {
        const currentTheme = get().theme;
        const newTheme = currentTheme === 'light' ? 'dark' : 'light';
        set({ theme: newTheme });

        // 应用主题到 document
        if (newTheme === 'dark') {
          document.documentElement.classList.add('dark');
        } else {
          document.documentElement.classList.remove('dark');
        }
      },

      // 当前项目
      currentProjectId: null,
      setCurrentProjectId: (id) => set({ currentProjectId: id }),

      // 侧边栏状态
      sidebarCollapsed: false,
      toggleSidebar: () => set((state) => ({ sidebarCollapsed: !state.sidebarCollapsed })),
    }),
    {
      name: 'x-novel-app',
      partialize: (state) => ({
        theme: state.theme,
        deviceSettings: state.deviceSettings,
        sidebarCollapsed: state.sidebarCollapsed,
      }),
    }
  )
);
