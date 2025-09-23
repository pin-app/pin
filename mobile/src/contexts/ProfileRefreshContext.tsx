import React, { createContext, useContext, useRef } from 'react';

interface ProfileRefreshContextType {
  refreshProfile: () => void;
  setRefreshFunction: (fn: () => void) => void;
}

const ProfileRefreshContext = createContext<ProfileRefreshContextType | null>(null);

export function ProfileRefreshProvider({ children }: { children: React.ReactNode }) {
  const refreshFunctionRef = useRef<(() => void) | null>(null);

  const refreshProfile = () => {
    if (refreshFunctionRef.current) {
      refreshFunctionRef.current();
    }
  };

  const setRefreshFunction = (fn: () => void) => {
    refreshFunctionRef.current = fn;
  };

  return (
    <ProfileRefreshContext.Provider value={{ refreshProfile, setRefreshFunction }}>
      {children}
    </ProfileRefreshContext.Provider>
  );
}

export function useProfileRefresh() {
  const context = useContext(ProfileRefreshContext);
  if (!context) {
    throw new Error('useProfileRefresh must be used within a ProfileRefreshProvider');
  }
  return context;
}
