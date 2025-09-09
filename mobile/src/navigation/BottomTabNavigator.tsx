import React from 'react';
import { createBottomTabNavigator } from '@react-navigation/bottom-tabs';
import { colors, typography } from '../theme';
import { FeedPage, MapPage, ProfilePage } from '../pages';

const Tab = createBottomTabNavigator();

export default function BottomTabNavigator() {
  return (
    <Tab.Navigator
      screenOptions={{
        tabBarStyle: {
          backgroundColor: colors.tabBar,
          borderTopColor: colors.border,
          borderTopWidth: 1,
        },
        tabBarActiveTintColor: colors.tabBarActive,
        tabBarInactiveTintColor: colors.tabBarInactive,
        tabBarLabelStyle: {
          fontSize: typography.fontSize.sm,
          fontWeight: typography.fontWeight.medium,
        },
        headerStyle: {
          backgroundColor: colors.background,
          borderBottomColor: colors.border,
          borderBottomWidth: 1,
        },
        headerTitleStyle: {
          color: colors.text,
          fontSize: typography.fontSize.lg,
          fontWeight: typography.fontWeight.semibold,
        },
      }}
    >
      <Tab.Screen 
        name="Feed" 
        component={FeedPage}
        options={{
          title: 'Feed',
        }}
      />
      <Tab.Screen 
        name="Map" 
        component={MapPage}
        options={{
          title: 'Map',
        }}
      />
      <Tab.Screen 
        name="Profile" 
        component={ProfilePage}
        options={{
          title: 'Profile',
        }}
      />
    </Tab.Navigator>
  );
}
