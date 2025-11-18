import React, { useState, useEffect, useCallback } from 'react';
import { View, StyleSheet, Text, SafeAreaView, TouchableOpacity, ScrollView, Alert, ActivityIndicator } from 'react-native';
import FontAwesome6 from '@expo/vector-icons/FontAwesome6';
import { colors, spacing } from '@/theme';
import UserProfile from '@/screens/Profile/components/UserProfile';
import ProfileHeader from '@/screens/Profile/components/ProfileHeader';
import { useAuth } from '@/contexts/AuthContext';
import { useProfileRefresh } from '@/contexts/ProfileRefreshContext';
import { apiService, Post as PostType } from '@/services/api';
import PostCard from '@/components/Post';
import CommentsScreen from '@/screens/Comments';

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
  const [userPosts, setUserPosts] = useState<PostType[]>([]);
  const [postsLoading, setPostsLoading] = useState(false);
  const [showComments, setShowComments] = useState(false);
  const [selectedPost, setSelectedPost] = useState<PostType | null>(null);

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

      await loadPostsForUser(userId);
      
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
  const loadPostsForUser = useCallback(async (id: string) => {
    try {
      setPostsLoading(true);
      const posts = await apiService.getPostsByUser(id, 20, 0);
      setUserPosts(posts);
    } catch (error) {
      console.error('Failed to load user posts:', error);
      setUserPosts([]);
    } finally {
      setPostsLoading(false);
    }
  }, []);

  const handleLikePost = async (postId: string) => {
    const post = userPosts.find(p => p.id === postId);
    if (!post) return;

    const previousState = { liked_by_user: post.liked_by_user, likes_count: post.likes_count };

    setUserPosts(prev =>
      prev.map(item => {
        if (item.id !== postId) return item;
        const delta = item.liked_by_user ? -1 : 1;
        return {
          ...item,
          liked_by_user: !item.liked_by_user,
          likes_count: Math.max(0, item.likes_count + delta),
        };
      })
    );

    try {
      if (previousState.liked_by_user) {
        await apiService.unlikePost(postId);
      } else {
        await apiService.likePost(postId);
      }
    } catch (error) {
      console.error('Failed to update like on other profile:', error);
      setUserPosts(prev =>
        prev.map(item =>
          item.id === postId ? { ...item, ...previousState } : item
        )
      );
      Alert.alert('Error', 'Could not update like.');
    }
  };

  const handleCommentPost = (postId: string) => {
    const post = userPosts.find(p => p.id === postId);
    if (post) {
      setSelectedPost(post);
      setShowComments(true);
    }
  };

  const handleBackFromComments = () => {
    setShowComments(false);
    setSelectedPost(null);
  };

  const renderPosts = () => {
    if (postsLoading) {
      return <ActivityIndicator color={colors.textSecondary} style={styles.postsLoader} />;
    }
    if (userPosts.length === 0) {
      return <Text style={styles.emptyPosts}>no posts yet</Text>;
    }
    return userPosts.map(post => (
      <View key={post.id} style={styles.postSpacing}>
        <PostCard
          post={post}
          likes={post.likes_count}
          isLiked={post.liked_by_user}
          onLike={() => handleLikePost(post.id)}
          onComment={() => handleCommentPost(post.id)}
          onRate={() => {}}
          onBookmark={() => {}}
          onUserPress={(id) => {
            if (id && id !== userId) {
              navigation.push('OtherUserProfile', { userId: id });
            }
          }}
        />
      </View>
    ));
  };

  if (showComments && selectedPost) {
    return (
      <CommentsScreen
        post={selectedPost}
        onBack={handleBackFromComments}
      />
    );
  }



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
        <View style={styles.postsSection}>
          {renderPosts()}
        </View>
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
  postsSection: {
    marginTop: spacing.lg,
  },
  postSpacing: {
    marginBottom: spacing.lg,
  },
  emptyPosts: {
    color: colors.textSecondary,
  },
  postsLoader: {
    marginTop: spacing.sm,
  },
});
