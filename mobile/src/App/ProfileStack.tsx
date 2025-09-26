import React from 'react';
import { NavigationContainer } from '@react-navigation/native';
import { createNativeStackNavigator } from '@react-navigation/native-stack';
import { ProfileScreen } from '@/screens';
import EditProfileScreen from '@/screens/EditProfile';
import FollowListScreen from '@/screens/Following';

const Stack = createNativeStackNavigator();

export default function ProfileStack() {
  return (
    <Stack.Navigator
      screenOptions={{
        headerShown: false,
      }}
    >
      <Stack.Screen name="ProfileMain" component={ProfileScreen} />
      <Stack.Screen name="EditProfile" component={EditProfileScreen} />
      <Stack.Screen name="FollowList" component={FollowListScreen} />
    </Stack.Navigator>
  );
}
