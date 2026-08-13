import type { GlobalThemeOverrides } from 'naive-ui'

export const lightThemeOverrides: GlobalThemeOverrides = {
  common: {
    primaryColor: '#00ADD8FF',
    primaryColorHover: '#26B6DEFF',
    primaryColorPressed: '#008CB0FF',
    primaryColorSuppl: '#26B6DEFF',
    borderRadius: '4px',
  },
  Layout: {
    color: '#f1f5f9',
  },
  Menu: {
    itemHeight: '32px',
  },
}

export const darkThemeOverrides: GlobalThemeOverrides = {
  common: {
    primaryColor: '#00ADD8FF',
    primaryColorHover: '#26B6DEFF',
    primaryColorPressed: '#008CB0FF',
    primaryColorSuppl: '#26B6DEFF',
    borderRadius: '4px',
  },
  Layout: {
    color: '#1c1c1c',
    siderColor: '#232323',
    siderBorderColor: 'transparent',
  },
  Card: {
    color: '#232323',
    borderColor: '#282828',
  },
  Menu: {
    itemHeight: '32px',
  },
}
