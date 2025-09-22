import React, { createContext, useContext, useState, useEffect, ReactNode } from 'react';
import AsyncStorage from '@react-native-async-storage/async-storage';
import { apiService } from '../services/api';

export interface User {
  id: string;
  email: string;
  username?: string;
  display_name?: string;
  bio?: string;
  location?: string;
  pfp_url?: string;
  created_at: string;
  updated_at: string;
}

export interface AuthState {
  user: User | null;
  sessionToken: string | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  isDevMode: boolean;
}

interface AuthContextType extends AuthState {
  login: (sessionToken: string, user: User) => Promise<void>;
  logout: () => Promise<void>;
  setDevMode: (enabled: boolean) => void;
  setDevUser: (userId: string) => void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

const SESSION_TOKEN_KEY = '@pin_session_token';
const USER_KEY = '@pin_user';
const DEV_MODE_KEY = '@pin_dev_mode';
const DEV_USER_KEY = '@pin_dev_user';

interface AuthProviderProps {
  children: ReactNode;
}

export function AuthProvider({ children }: AuthProviderProps) {
  const [user, setUser] = useState<User | null>(null);
  const [sessionToken, setSessionToken] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [isDevMode, setIsDevMode] = useState(false);
  const [devUserId, setDevUserId] = useState<string | null>(null);

  const isAuthenticated = !!(user && sessionToken);

  useEffect(() => {
    loadAuthState();
  }, []);

  const loadAuthState = async () => {
    try {
      const [storedToken, storedUser, storedDevMode, storedDevUser] = await Promise.all([
        AsyncStorage.getItem(SESSION_TOKEN_KEY),
        AsyncStorage.getItem(USER_KEY),
        AsyncStorage.getItem(DEV_MODE_KEY),
        AsyncStorage.getItem(DEV_USER_KEY),
      ]);

      const devMode = storedDevMode === 'true';
      setIsDevMode(devMode);
      
      // Update API service with dev mode state
      apiService.setDevMode(devMode);

      if (devMode && storedDevUser) {
        setDevUserId(storedDevUser);
        apiService.setDevUserId(storedDevUser);
        // In dev mode, we don't need a real session token
        setSessionToken('dev-token');
        // Try to fetch user data or create a mock user
        try {
          const userData = await apiService.getUser(storedDevUser);
          setUser(userData);
        } catch (error) {
          // If user doesn't exist, create a mock user
          const mockUser: User = {
            id: storedDevUser,
            email: 'dev@localhost',
            username: 'devuser',
            display_name: 'Dev User',
            created_at: new Date().toISOString(),
            updated_at: new Date().toISOString(),
          };
          setUser(mockUser);
        }
      } else if (storedToken && storedUser) {
        setSessionToken(storedToken);
        apiService.setSessionToken(storedToken);
        setUser(JSON.parse(storedUser));
      }
    } catch (error) {
      console.error('Failed to load auth state:', error);
    } finally {
      setIsLoading(false);
    }
  };

  const login = async (token: string, userData: User) => {
    try {
      await Promise.all([
        AsyncStorage.setItem(SESSION_TOKEN_KEY, token),
        AsyncStorage.setItem(USER_KEY, JSON.stringify(userData)),
      ]);
      setSessionToken(token);
      apiService.setSessionToken(token);
      setUser(userData);
    } catch (error) {
      console.error('Failed to save auth state:', error);
      throw error;
    }
  };

  const logout = async () => {
    try {
      if (sessionToken && !isDevMode) {
        await apiService.logout(sessionToken);
      }
      await Promise.all([
        AsyncStorage.removeItem(SESSION_TOKEN_KEY),
        AsyncStorage.removeItem(USER_KEY),
      ]);
      setSessionToken(null);
      apiService.setSessionToken(null);
      setUser(null);
    } catch (error) {
      console.error('Failed to logout:', error);
      // Still clear local state even if API call fails
      setSessionToken(null);
      apiService.setSessionToken(null);
      setUser(null);
    }
  };

  const setDevMode = async (enabled: boolean) => {
    try {
      console.log('🔧 Setting dev mode:', enabled);
      await AsyncStorage.setItem(DEV_MODE_KEY, enabled.toString());
      setIsDevMode(enabled);
      
      // Update API service
      apiService.setDevMode(enabled);
      
      if (!enabled) {
        console.log('🔧 Disabling dev mode, clearing data');
        // Clear dev mode data
        await AsyncStorage.removeItem(DEV_USER_KEY);
        setDevUserId(null);
        apiService.setDevUserId(null);
        if (sessionToken === 'dev-token') {
          await logout();
        }
      } else {
        console.log('🔧 Dev mode enabled');
      }
    } catch (error) {
      console.error('Failed to set dev mode:', error);
    }
  };

  const setDevUser = async (userId: string) => {
    try {
      console.log('🔧 Setting dev user:', userId, 'Dev mode:', isDevMode);
      await AsyncStorage.setItem(DEV_USER_KEY, userId);
      setDevUserId(userId);
      
      // Update API service
      apiService.setDevUserId(userId);
      
      // If we're in dev mode, update the user and set session token
      if (isDevMode) {
        console.log('🔧 Dev mode enabled, setting session token and user');
        setSessionToken('dev-token');
        apiService.setSessionToken('dev-token');
        
        try {
          console.log('🔧 Fetching user data from API...');
          const userData = await apiService.getUser(userId);
          console.log('🔧 User data fetched:', userData);
          setUser(userData);
        } catch (error) {
          console.log('🔧 API failed, creating mock user:', error);
          // Create mock user if not found
          const mockUser: User = {
            id: userId,
            email: 'dev@localhost',
            username: 'devuser',
            display_name: 'Dev User',
            created_at: new Date().toISOString(),
            updated_at: new Date().toISOString(),
          };
          setUser(mockUser);
        }
      } else {
        console.log('⚠️ Dev mode not enabled, cannot set dev user');
      }
    } catch (error) {
      console.error('Failed to set dev user:', error);
    }
  };

  const value: AuthContextType = {
    user,
    sessionToken,
    isAuthenticated,
    isLoading,
    isDevMode,
    login,
    logout,
    setDevMode,
    setDevUser,
  };

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
}
