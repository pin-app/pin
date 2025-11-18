import React, { useCallback, useEffect, useState } from 'react';
import { SafeAreaView, View, Text, StyleSheet, TouchableOpacity, Alert, ActivityIndicator } from 'react-native';
import FontAwesome6 from '@expo/vector-icons/FontAwesome6';
import { colors, spacing, typography } from '@/theme';
import { apiService, Post as PostType } from '@/services/api';
import Feed from '@/screens/Feed/components/Feed';
import CommentsScreen from '@/screens/Comments';
import { useAuth } from '@/contexts/AuthContext';

interface PlacePostsScreenProps {
  route: {
    params: {
      placeId: string;
      placeName: string;
    };
  };
  navigation: any;
}

export default function PlacePostsScreen({ route, navigation }: PlacePostsScreenProps) {
  const { placeId, placeName } = route.params;
  const { user: currentUser } = useAuth();

  const [posts, setPosts] = useState<PostType[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [refreshing, setRefreshing] = useState(false);
  const [selectedPost, setSelectedPost] = useState<PostType | null>(null);
  const [showComments, setShowComments] = useState(false);

  const loadPosts = useCallback(
    async (opts?: { refreshing?: boolean }) => {
      try {
        if (opts?.refreshing) {
          setRefreshing(true);
        } else {
          setIsLoading(true);
        }

        const placePosts = await apiService.getPostsByPlace(placeId, 20, 0);
        setPosts(placePosts);
      } catch (error) {
        console.error('Failed to load place posts:', error);
      } finally {
        if (opts?.refreshing) {
          setRefreshing(false);
        } else {
          setIsLoading(false);
        }
      }
    },
    [placeId]
  );

  useEffect(() => {
    loadPosts();
  }, [loadPosts]);

  const handleLike = async (postId: string) => {
    if (!currentUser) {
      Alert.alert('Sign in required', 'Log in to like posts.');
      return;
    }

    const target = posts.find((post) => post.id === postId);
    if (!target) return;

    const previousState = {
      liked_by_user: target.liked_by_user,
      likes_count: target.likes_count,
    };

    setPosts((prev) =>
      prev.map((post) => {
        if (post.id !== postId) return post;
        const delta = post.liked_by_user ? -1 : 1;
        return {
          ...post,
          liked_by_user: !post.liked_by_user,
          likes_count: Math.max(0, post.likes_count + delta),
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
      console.error('Failed to update like for place posts:', error);
      setPosts((prev) =>
        prev.map((post) => (post.id === postId ? { ...post, ...previousState } : post))
      );
      Alert.alert('Error', 'Could not update like.');
    }
  };

  const handleComment = (postId: string) => {
    const post = posts.find((p) => p.id === postId);
    if (post) {
      setSelectedPost(post);
      setShowComments(true);
    }
  };

  const handleBackFromComments = () => {
    setShowComments(false);
    setSelectedPost(null);
  };

  const handleUserPress = (userId: string, username?: string) => {
    if (!userId) return;

    if (currentUser && userId === currentUser.id) {
      navigation.navigate('MainTabs', {
        screen: 'Profile',
      });
    } else {
      navigation.navigate('OtherUserProfile', {
        userId,
        username,
      });
    }
  };

  const handleBookmark = () => {};
  const handleRate = () => {};

  const handleBack = () => {
    navigation.goBack();
  };

  if (showComments && selectedPost) {
    return <CommentsScreen post={selectedPost} onBack={handleBackFromComments} />;
  }

  return (
    <SafeAreaView style={styles.container}>
      <View style={styles.header}>
        <TouchableOpacity onPress={handleBack} style={styles.backButton}>
          <FontAwesome6 name="arrow-left" size={18} color={colors.text} />
        </TouchableOpacity>
        <View style={styles.headerContent}>
          <Text style={styles.placeName}>{placeName}</Text>
          <Text style={styles.subtitle}>posts at this spot</Text>
        </View>
      </View>

      {isLoading ? (
        <View style={styles.loading}>
          <ActivityIndicator color={colors.textSecondary} />
        </View>
      ) : (
        <Feed
          posts={posts}
          onLike={handleLike}
          onComment={handleComment}
          onRate={handleRate}
          onBookmark={handleBookmark}
          onUserPress={handleUserPress}
          refreshing={refreshing}
          onRefresh={() => loadPosts({ refreshing: true })}
        />
      )}
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: colors.background,
  },
  header: {
    flexDirection: 'row',
    alignItems: 'center',
    paddingHorizontal: spacing.md,
    paddingVertical: spacing.md,
    borderBottomWidth: 1,
    borderBottomColor: colors.border,
  },
  backButton: {
    padding: spacing.sm,
    marginRight: spacing.sm,
  },
  headerContent: {
    flex: 1,
  },
  placeName: {
    fontSize: typography.fontSize.lg,
    color: colors.text,
    fontWeight: typography.fontWeight.semibold,
  },
  subtitle: {
    fontSize: typography.fontSize.sm,
    color: colors.textSecondary,
  },
  loading: {
    flex: 1,
    alignItems: 'center',
    justifyContent: 'center',
  },
});

