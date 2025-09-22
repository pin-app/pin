import React from "react";
import { NavigationContainer } from "@react-navigation/native";
import { SafeAreaProvider } from "react-native-safe-area-context";
import { AuthProvider, useAuth } from "../contexts/AuthContext";
import TabNavigator from "./TabNavigator";
import AuthScreen from "../screens/Auth";
import { EventProvider } from "react-native-outside-press";

export { default as TabNavigator } from "./TabNavigator";

function AppContent() {
  const { isAuthenticated, isLoading } = useAuth();

  if (isLoading) {
    // You could add a loading screen here
    return null;
  }

  return (
    <NavigationContainer>
      <EventProvider>
        {isAuthenticated ? <TabNavigator /> : <AuthScreen />}
      </EventProvider>
    </NavigationContainer>
  );
}

export default function App() {
  return (
    <SafeAreaProvider>
      <AuthProvider>
        <AppContent />
      </AuthProvider>
    </SafeAreaProvider>
  );
}
