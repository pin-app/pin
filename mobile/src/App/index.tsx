import React from "react";
import { NavigationContainer } from "@react-navigation/native";
import { createStackNavigator } from "@react-navigation/stack";
import { SafeAreaProvider } from "react-native-safe-area-context";
import { AuthProvider, useAuth } from "../contexts/AuthContext";
import TabNavigator from "./TabNavigator";
import AuthScreen from "../screens/Auth";
import { OtherUserProfileScreen } from "../screens";
import { EventProvider } from "react-native-outside-press";

type RootStackParamList = {
  MainTabs: undefined;
  OtherUserProfile: {
    userId: string;
    username?: string;
  };
};

const Stack = createStackNavigator<RootStackParamList>();

export { default as TabNavigator } from "./TabNavigator";

function AppContent() {
  const { isAuthenticated, isLoading } = useAuth();

  if (isLoading) {
    // could add a loading screen here
    return null;
  }

  return (
    <NavigationContainer>
      <EventProvider>
        {isAuthenticated ? (
          <Stack.Navigator
            screenOptions={{
              headerShown: false,
            }}
          >
            <Stack.Screen name="MainTabs" component={TabNavigator} />
            <Stack.Screen 
              name="OtherUserProfile" 
              component={OtherUserProfileScreen}
              options={{
                presentation: 'card',
                gestureEnabled: true,
                headerShown: false,
              }}
            />
          </Stack.Navigator>
        ) : (
          <AuthScreen />
        )}
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
