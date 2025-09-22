import React from 'react';
import { createStackNavigator } from '@react-navigation/stack';
import { ProfileScreen } from '@/screens';

const Stack = createStackNavigator();

export default function ProfileStack() {
  return (
    <Stack.Navigator
      screenOptions={{
        headerShown: false,
      }}
    >
      <Stack.Screen name="ProfileMain" component={ProfileScreen} />
    </Stack.Navigator>
  );
}
