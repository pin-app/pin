import React, { useState, useEffect, useCallback } from 'react';
import { View, StyleSheet, Text, SafeAreaView, TouchableOpacity, Alert, ScrollView } from 'react-native';
import { useFocusEffect } from '@react-navigation/native';
import FontAwesome6 from '@expo/vector-icons/FontAwesome6';
import { colors, spacing } from '@/theme';
import Button from '@/components/Button';
import UserProfile from '@/screens/Profile/components/UserProfile';
import ProfileHeader from './components/ProfileHeader';
import DevModeSettings from './components/DevModeSettings';
import DebugInfo from './components/DebugInfo';
import { useAuth } from '@/contexts/AuthContext';
import { useProfileRefresh } from '@/contexts/ProfileRefreshContext';
import { apiService } from '@/services/api';
import SidebarMenu, { MenuItem } from './components/sideBarMenu';
import { useNavigation } from '@react-navigation/native';

export default function ProfileScreen() {
  const { user, isDevMode } = useAuth();
  const { setRefreshFunction } = useProfileRefresh();
  const [showMenu, setShowMenu] = useState(false);
  const [showDevSettings, setShowDevSettings] = useState(false);
  const [followingCount, setFollowingCount] = useState(0);
  const [followersCount, setFollowersCount] = useState(0);
  const [postsCount, setPostsCount] = useState(0);
  const [isLoading, setIsLoading] = useState(true);
  const [refreshTrigger, setRefreshTrigger] = useState(0);

  const loadUserStats = useCallback(async () => {
    if (!user) return;

    try {
      setIsLoading(true);
      
      try {
        console.log('Loading user stats for:', user.id);
        const stats = await apiService.getUserStats(user.id);
        console.log('User stats loaded:', stats);
        console.log('Previous following count:', followingCount, 'New following count:', stats.following_count);
        setPostsCount(stats.posts_count);
        setFollowingCount(stats.following_count);
        setFollowersCount(stats.followers_count);
      } catch (error) {
        console.error('Failed to load user stats:', error);
        // fallback to individual calls if stats endpoint fails
        try {
          const posts = await apiService.getPostsByUser(user.id, 1, 0);
          setPostsCount(posts.length);
        } catch (postError) {
          console.error('Failed to load posts count:', postError);
          setPostsCount(0);
        }
        setFollowingCount(0);
        setFollowersCount(0);
      }
      
    } catch (error) {
      console.error('Failed to load user stats:', error);
      setFollowingCount(0);
      setFollowersCount(0);
      setPostsCount(0);
    } finally {
      setIsLoading(false);
    }
  }, [user]);

  useEffect(() => {
    if (user) {
      loadUserStats();
    }
  }, [user, loadUserStats]);

  // Reload stats when refreshTrigger changes
  useEffect(() => {
    if (user && refreshTrigger > 0) {
      console.log('Refresh trigger activated, reloading stats...');
      loadUserStats();
    }
  }, [refreshTrigger, user, loadUserStats]);

  const handleEditProfile = () => {
    console.log('Edit profile pressed');
  };

  const handleMenuPress = () => {
    setShowMenu(true);
  };

  const navigation = useNavigation<any>();

  const onEditProfile = () => {
    navigation.navigate('EditProfile');
  }
  const onShareProfile = () => console.log('Share profile');
  const onSaved = () => console.log('Saved posts');
  const onSettings = () => console.log('Settings');
  const onLogout = () => console.log('Logout');

  const sections = [
    {
      header: 'Profile',
      items: [
        { key: 'edit', label: 'Edit Profile', onPress: onEditProfile },
        { key: 'share', label: 'Share Profile', onPress: onShareProfile },
        { key: 'saved', label: 'Saved', subtitle: 'Your saved posts', onPress: onSaved },
      ],
    },
    {
      header: 'App',
      items: [
        { key: 'settings', label: 'Settings', onPress: onSettings },
        { key: 'dev', label: 'Developer Settings', subtitle: isDevMode ? 'Dev Mode ON' : 'Dev Mode OFF', onPress: () => setShowDevSettings(true) },
        { key: 'logout', label: 'Logout', destructive: true, onPress: onLogout },
      ],
    },
  ];

  // Function to refresh stats (can be called from other screens)
  const refreshStats = useCallback(() => {
    console.log('Manually refreshing stats...');
    setRefreshTrigger(prev => prev + 1);
  }, []);

  // Register the refresh function with the context
  useEffect(() => {
    setRefreshFunction(refreshStats);
  }, [setRefreshFunction, refreshStats]);


  if (!user) {
    return (
      <SafeAreaView style={styles.container}>
        <View style={styles.errorContainer}>
          <Text style={styles.errorText}>No user data available</Text>
        </View>
      </SafeAreaView>
    );
  }

  return (
    <SafeAreaView style={styles.container}>
      <ProfileHeader 
        username={user.username || user.display_name || 'user'} 
        onMenuPress={handleMenuPress} 
      />
      
      <ScrollView 
        style={styles.scrollContainer}
        contentContainerStyle={styles.content}
        showsVerticalScrollIndicator={false}
      >
        <UserProfile 
          user={user} 
          currentUserId={user.id}
          showFollowButton={false}
          onEditProfile={onEditProfile}
          onShareProfile={handleEditProfile}
          postsCount={postsCount}
          followingCount={followingCount}
          followersCount={followersCount}
        />
      </ScrollView>

      <SidebarMenu visible={showMenu} onClose={() => setShowMenu(false)} title="Profile Menu" sections={sections} />
        
      <DevModeSettings 
        visible={showDevSettings} 
        onClose={() => setShowDevSettings(false)} 
      />
    </SafeAreaView>
  );
}



const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: colors.background,
  },
  scrollContainer: {
    flex: 1,
  },
  content: {
    paddingHorizontal: spacing.md,
    paddingBottom: spacing.xl,
  },
  debugText: {
    fontSize: 12,
    color: colors.textSecondary,
    textAlign: 'center',
    marginBottom: spacing.sm,
  },


  divider: {
    width: StyleSheet.hairlineWidth,
    height: '60%',
    backgroundColor: '#E5E7B',
  },

  profilePicWrapper: {
    width: 150,
    height: 150,
    borderRadius: '100%',
    overflow: 'hidden',
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
    backgroundColor: 'gray',
  },
  buttonRow: {
    flexDirection: 'row',
    alignItems: 'center',
    marginTop: spacing.md,
    marginBottom: spacing.lg,
    gap: spacing.sm,
  },
  editButton: {
    flex: 1,
  },
  shareButton: {
    flex: 1,
  },
  followButton: {
    width: 44,
    minHeight: 32,
    borderRadius: 8,
    backgroundColor: colors.background,
    borderWidth: 1,
    borderColor: colors.border,
    justifyContent: 'center',
    alignItems: 'center',
  },

  profilePic: {
    width: '100%',
    height: '100%',
    objectFit: 'cover',
  },
  devModeIndicator: {
    backgroundColor: colors.background,
    paddingHorizontal: spacing.sm,
    paddingVertical: spacing.xs,
    borderRadius: 16,
    alignSelf: 'center',
    marginTop: spacing.md,
  },
  devModeText: {
    color: colors.background,
    fontSize: 12,
    fontWeight: '600',
  },
  errorContainer: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    paddingHorizontal: spacing.md,
  },
  errorText: {
    fontSize: 16,
    color: colors.textSecondary,
    textAlign: 'center',
  },
});
