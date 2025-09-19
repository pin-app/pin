import React from 'react';
import { FlatList, StyleSheet, View } from 'react-native';
import { colors } from '@/theme/colors';
import { spacing } from '@/theme/spacing';
import Post from '@/components/Post';

interface PostData {
  id: string;
  user: {
    name: string;
    avatar: string;
  };
  location: string;
  visits: number;
  dayOfWeek: string;
  rating: number;
  text: string;
  images: Array<{
    id: string;
    uri: string;
  }>;
  likes: number;
  isLiked: boolean;
}

interface FeedProps {
  posts: PostData[];
  onLike: (postId: string) => void;
  onComment: (postId: string) => void;
  onRate: (postId: string) => void;
  onBookmark: (postId: string) => void;
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
  onLoadMore,
  refreshing = false,
  onRefresh,
}: FeedProps) {
  const renderPost = ({ item, index }: { item: PostData; index: number }) => (
    <>
      <Post
        id={item.id}
        user={item.user}
        location={item.location}
        visits={item.visits}
        dayOfWeek={item.dayOfWeek}
        rating={item.rating}
        text={item.text}
        images={item.images}
        likes={item.likes}
        isLiked={item.isLiked}
        onLike={onLike}
        onComment={onComment}
        onRate={onRate}
        onBookmark={onBookmark}
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
