export const colors = {
  background: '#FFFFFF',
  text: '#000000',
  
  border: '#E0E0E0',
  borderDark: '#CCCCCC',
  
  tabBar: '#FFFFFF',
  tabBarActive: '#000000',
  tabBarInactive: '#999999',
} as const;

export type ColorKey = keyof typeof colors;
