import React from 'react';
import { FlatList, StyleSheet, View } from 'react-native';
import { colors } from '@/theme/colors';
import { spacing } from '@/theme/spacing';
import Post from '@/components/Post';
import { Post as PostType } from '@/services/api';

interface FeedProps {
  posts: PostType[];
  onLike: (postId: string) => void;
  onComment: (postId: string) => void;
  onRate: (postId: string) => void;
  onBookmark: (postId: string) => void;
  onUserPress?: (userId: string, username?: string) => void;
  onLoadMore?: () => void;
  refreshing?: boolean;
  onRefresh?: () => void;
}

export default function Feed({
  posts,
  onLike,
  onComment,
  onRate,
  onBookmark,
  onUserPress,
  onLoadMore,
  refreshing = false,
  onRefresh,
}: FeedProps) {
  const renderPost = ({ item, index }: { item: PostType; index: number }) => (
    <>
      <Post
        post={item}
        likes={item.likes_count}
        isLiked={item.liked_by_user}
        onLike={onLike}
        onComment={onComment}
        onRate={onRate}
        onBookmark={onBookmark}
        onUserPress={onUserPress}
      />
      {index < posts.length - 1 && <View style={styles.separator} />}
    </>
  );

  return (
    <View style={styles.container}>
      <FlatList
        data={posts}
        renderItem={renderPost}
        keyExtractor={(item) => item.id}
        showsVerticalScrollIndicator={false}
        onEndReached={onLoadMore}
        onEndReachedThreshold={0.5}
        refreshing={refreshing}
        onRefresh={onRefresh}
        contentContainerStyle={styles.contentContainer}
      />
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: colors.background,
  },
  contentContainer: {
    paddingBottom: spacing.xl,
  },
  separator: {
    height: 1,
    backgroundColor: colors.border,
    marginHorizontal: spacing.md,
  },
});
