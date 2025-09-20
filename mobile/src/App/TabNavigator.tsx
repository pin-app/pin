import React from 'react';
import { createBottomTabNavigator } from '@react-navigation/bottom-tabs';
import { View, Image } from 'react-native';
import FontAwesome6 from '@expo/vector-icons/FontAwesome6';
import { colors, typography } from '@/theme';
import { FeedScreen, MapScreen, ProfileScreen } from '@/screens';
import { useAuth } from '@/contexts/AuthContext';

const Tab = createBottomTabNavigator();

export default function TabNavigator() {
  const { user } = useAuth();
  
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
        component={FeedScreen}
        options={{
          headerShown: false,
          tabBarLabel: '',
          tabBarIcon: ({ color, size }) => (
            <FontAwesome6 name="newspaper" size={size} color={color} />
          ),
        }}
      />
      <Tab.Screen 
        name="Map" 
        component={MapScreen}
        options={{
          headerShown: false,
          tabBarLabel: '',
          tabBarIcon: ({ color, size }) => (
            <FontAwesome6 name="globe" size={size} color={color} />
          ),
        }}
      />
      <Tab.Screen 
        name="Profile" 
        component={ProfileScreen}
        options={{
          headerShown: false,
          tabBarLabel: '',
          tabBarIcon: ({ color, size, focused }) => {
            const displayName = user?.display_name || user?.username || 'User';
            return (
              <View style={[
                { 
                  width: size, 
                  height: size, 
                  borderRadius: size/2, 
                  backgroundColor: colors.border,
                  borderWidth: 1,
                  borderColor: focused ? colors.text : colors.border,
                  overflow: 'hidden',
                  justifyContent: 'center',
                  alignItems: 'center'
                }
              ]}>
                {user?.pfp_url ? (
                  <Image 
                    source={{ uri: user.pfp_url }} 
                    style={{ width: '100%', height: '100%' }} 
                  />
                ) : (
                  <View style={{
                    width: '100%',
                    height: '100%',
                    backgroundColor: colors.border,
                    justifyContent: 'center',
                    alignItems: 'center'
                  }}>
                    <FontAwesome6 
                      name="user" 
                      size={size * 0.6} 
                      color={focused ? colors.text : colors.textSecondary} 
                    />
                  </View>
                )}
              </View>
            );
          },
        }}
      />
    </Tab.Navigator>
  );
}
