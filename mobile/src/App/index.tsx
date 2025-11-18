import React from "react";
import { NavigationContainer } from '@react-navigation/native';
import { createNativeStackNavigator } from '@react-navigation/native-stack';
import { SafeAreaProvider } from "react-native-safe-area-context";
import { AuthProvider, useAuth } from "../contexts/AuthContext";
import { ProfileRefreshProvider } from "../contexts/ProfileRefreshContext";
import TabNavigator from "./TabNavigator";
import AuthScreen from "../screens/Auth";
import { OtherUserProfileScreen, PlacePostsScreen } from "../screens";
import { EventProvider } from "react-native-outside-press";

type RootStackParamList = {
  MainTabs: undefined;
  OtherUserProfile: {
    userId: string;
    username?: string;
  };
  PlacePosts: {
    placeId: string;
    placeName: string;
  };
};

const Stack = createNativeStackNavigator<RootStackParamList>();

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
            <Stack.Screen 
              name="PlacePosts" 
              component={PlacePostsScreen}
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
        <ProfileRefreshProvider>
          <AppContent />
        </ProfileRefreshProvider>
      </AuthProvider>
    </SafeAreaProvider>
  );
}
