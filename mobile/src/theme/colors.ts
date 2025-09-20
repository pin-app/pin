export const colors = {
  background: '#FFFFFF',
  text: '#000000',
  
  border: '#E0E0E0',
  borderDark: '#CCCCCC',
  
  tabBar: '#FFFFFF',
  tabBarActive: '#000000',
  tabBarInactive: '#999999',
  
  searchBackground: '#F5F5F5',
  postBackground: '#FFFFFF',
  ratingPositive: '#22C55E',
  ratingNegative: '#EF4444',
  ratingNeutral: '#6B7280',
  
  iconDefault: '#000000',
  iconInactive: '#9CA3AF',
  iconActive: '#000000',
  
  textSecondary: '#6B7280',
  textTertiary: '#9CA3AF',
} as const;

export type ColorKey = keyof typeof colors;
