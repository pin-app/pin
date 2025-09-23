import React, { useState, useEffect } from 'react';
import { View, StyleSheet, Text, SafeAreaView, TouchableOpacity, ScrollView, Alert } from 'react-native';
import FontAwesome6 from '@expo/vector-icons/FontAwesome6';
import { colors, spacing } from '@/theme';
import UserProfile from '@/screens/Profile/components/UserProfile';
import ProfileHeader from '@/screens/Profile/components/ProfileHeader';
import { useAuth } from '@/contexts/AuthContext';
import { useProfileRefresh } from '@/contexts/ProfileRefreshContext';
import { apiService } from '@/services/api';

interface OtherUserProfileScreenProps {
  route: {
    params: {
      userId: string;
      username?: string;
    };
  };
  navigation: any;
}

export default function OtherUserProfileScreen({ route, navigation }: OtherUserProfileScreenProps) {
  const { user: currentUser } = useAuth();
  const { refreshProfile } = useProfileRefresh();
  const { userId, username } = route.params;
  const [otherUser, setOtherUser] = useState<any>(null);
  const [followingCount, setFollowingCount] = useState(0);
  const [followersCount, setFollowersCount] = useState(0);
  const [postsCount, setPostsCount] = useState(0);
  const [isLoading, setIsLoading] = useState(true);
  const [isFollowing, setIsFollowing] = useState(false);

  useEffect(() => {
    loadUserProfile();
  }, [userId]);

  const loadUserProfile = async () => {
    try {
      setIsLoading(true);
      
      console.log('Loading user profile for ID:', userId);
      
      // Load user data
      const userData = await apiService.getUser(userId);
      console.log('User data loaded:', userData);
      setOtherUser(userData);
      
      // Load user stats
      try {
        const stats = await apiService.getUserStats(userId);
        setPostsCount(stats.posts_count);
        setFollowingCount(stats.following_count);
        setFollowersCount(stats.followers_count);
      } catch (error) {
        console.error('Failed to load user stats:', error);
        setPostsCount(0);
        setFollowingCount(0);
        setFollowersCount(0);
      }
      
      // Check if current user is following this user
      if (currentUser) {
        try {
          const following = await apiService.isFollowing(userId);
          setIsFollowing(following);
        } catch (error) {
          console.error('Failed to check follow status:', error);
          // Default to not following if check fails
          setIsFollowing(false);
        }
      }
      
    } catch (error) {
      console.error('Failed to load user profile:', error);
      console.error('Error details:', error);
      Alert.alert('Error', `Failed to load user profile: ${error instanceof Error ? error.message : 'Unknown error'}`);
      navigation.goBack();
    } finally {
      setIsLoading(false);
    }
  };


  const handleBack = () => {
    navigation.goBack();
  };

  if (isLoading) {
    return (
      <SafeAreaView style={styles.container}>
        <View style={styles.loadingContainer}>
          <Text style={styles.loadingText}>Loading...</Text>
        </View>
      </SafeAreaView>
    );
  }

  if (!otherUser) {
    return (
      <SafeAreaView style={styles.container}>
        <ProfileHeader 
          username={username || 'Unknown User'} 
          onMenuPress={handleBack}
          showBackButton={true}
        />
        <View style={styles.errorContainer}>
          <Text style={styles.errorText}>User not found</Text>
        </View>
      </SafeAreaView>
    );
  }

  return (
    <SafeAreaView style={styles.container}>
      <ProfileHeader 
        username={otherUser.username || otherUser.display_name || 'user'} 
        onMenuPress={handleBack}
        showBackButton={true}
      />
      
      <ScrollView 
        style={styles.scrollContainer}
        contentContainerStyle={styles.content}
        showsVerticalScrollIndicator={false}
      >
        <UserProfile 
          user={otherUser} 
          currentUserId={currentUser?.id}
          showFollowButton={true}
          onFollowChange={setIsFollowing}
          onFollowAction={() => {
            // Update follower count immediately
            setFollowersCount(prev => isFollowing ? Math.max(0, prev - 1) : prev + 1);
            // Refresh the current user's profile stats
            refreshProfile();
          }}
          postsCount={postsCount}
          followingCount={followingCount}
          followersCount={followersCount}
        />
      </ScrollView>
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
  loadingContainer: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    paddingHorizontal: spacing.md,
  },
  loadingText: {
    fontSize: 16,
    color: colors.textSecondary,
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
