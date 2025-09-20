export const config = {
  apiBaseUrl: process.env.EXPO_PUBLIC_API_BASE || 'http://localhost:8080',
  devMode: process.env.EXPO_PUBLIC_DEV_MODE === 'true' || __DEV__,
  googleClientId: process.env.EXPO_PUBLIC_GOOGLE_CLIENT_ID || '',
  appleClientId: process.env.EXPO_PUBLIC_APPLE_CLIENT_ID || '',
};
