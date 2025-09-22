import React from 'react';
import { View, Text, StyleSheet, TouchableOpacity, Image, ScrollView } from 'react-native';
import FontAwesome6 from '@expo/vector-icons/FontAwesome6';
import { colors, spacing, typography } from '@/theme';
import { Post as PostType } from '@/services/api';

interface PostProps {
  post: PostType;
  likes?: number;
  isLiked?: boolean;
  onLike: (postId: string) => void;
  onComment: (postId: string) => void;
  onRate: (postId: string) => void;
  onBookmark: (postId: string) => void;
  onUserPress?: (userId: string, username?: string) => void;
  showCommentsButton?: boolean;
}

export default function Post({
  post,
  likes = 0,
  isLiked = false,
  onLike,
  onComment,
  onRate,
  onBookmark,
  onUserPress,
  showCommentsButton = true,
}: PostProps) {
  const getRatingColor = (rating: number) => {
    if (rating >= 7) return colors.ratingPositive;
    if (rating >= 4) return colors.ratingNeutral;
    return colors.ratingNegative;
  };

  const userName = post.user?.display_name || post.user?.username || 'Unknown User';
  const placeName = post.place?.name || 'Unknown Place';
  const description = post.description || '';
  const postImages = post.images || [];
  
  const postDate = new Date(post.created_at);
  const dayOfWeek = postDate.toLocaleDateString('en-US', { weekday: 'long' });
  
  // TODO: use real rating
  const rating = 8.2;
  const visits = 3;

  return (
    <View style={styles.container}>
      <View style={styles.header}>
        <TouchableOpacity 
          style={styles.avatar}
          onPress={() => {
            console.log('Avatar clicked, user ID:', post.user?.id, 'username:', post.user?.username);
            onUserPress?.(post.user?.id || '', post.user?.username);
          }}
          disabled={!onUserPress || !post.user?.id}
        >
          {post.user?.pfp_url ? (
            <Image source={{ uri: post.user.pfp_url }} style={styles.avatarImage} />
          ) : (
            <View style={styles.avatarPlaceholder}>
              <Text style={styles.avatarText}>
                {userName.charAt(0).toUpperCase()}
              </Text>
            </View>
          )}
        </TouchableOpacity>
        <View style={styles.userDetails}>
          <View style={styles.userNameContainer}>
            <TouchableOpacity 
              onPress={() => {
                console.log('Username clicked, user ID:', post.user?.id, 'username:', post.user?.username);
                onUserPress?.(post.user?.id || '', post.user?.username);
              }}
              disabled={!onUserPress || !post.user?.id}
            >
              <Text style={[styles.userName, styles.boldText]}>{userName}</Text>
            </TouchableOpacity>
            <Text style={styles.userName}>
              {' '}ranked <Text style={styles.boldText}>{placeName}</Text>
            </Text>
          </View>
          <Text style={styles.visitText}>{visits} visits • {dayOfWeek}</Text>
        </View>
        <View style={styles.ratingBadge}>
          <Text style={[styles.ratingText, { color: getRatingColor(rating) }]}>{rating}</Text>
        </View>
      </View>

      {description && (
        <View style={styles.textSection}>
          <Text style={styles.postText}>{description}</Text>
        </View>
      )}

      {postImages.length > 0 && (
        <ScrollView horizontal showsHorizontalScrollIndicator={false} style={styles.imagesContainer}>
          {postImages.map((image) => (
            <View key={image.id} style={styles.postImage}>
              <Image source={{ uri: image.image_url }} style={styles.postImageContent} />
            </View>
          ))}
        </ScrollView>
      )}

      <View style={styles.engagementSection}>
        <View style={styles.likesContainer}>
          <Text style={styles.likesText}>{likes} likes</Text>
        </View>
        <View style={styles.actionIcons}>
          <TouchableOpacity onPress={() => onLike(post.id)} style={styles.actionButton}>
            <FontAwesome6 
              name="heart" 
              size={20} 
              color={isLiked ? colors.ratingNegative : colors.iconDefault} 
              solid={isLiked}
            />
          </TouchableOpacity>
          {showCommentsButton && (
            <TouchableOpacity onPress={() => onComment(post.id)} style={styles.actionButton}>
              <FontAwesome6 name="comment" size={20} color={colors.iconDefault} />
            </TouchableOpacity>
          )}
          <TouchableOpacity onPress={() => onRate(post.id)} style={styles.actionButton}>
            <FontAwesome6 name="plus" size={20} color={colors.iconDefault} />
          </TouchableOpacity>
          <TouchableOpacity onPress={() => onBookmark(post.id)} style={styles.actionButton}>
            <FontAwesome6 name="bookmark" size={20} color={colors.iconDefault} />
          </TouchableOpacity>
        </View>
      </View>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    backgroundColor: colors.postBackground,
    marginBottom: spacing.lg,
    paddingHorizontal: spacing.md,
    paddingVertical: spacing.md,
  },
  header: {
    flexDirection: 'row',
    alignItems: 'center',
    marginBottom: spacing.sm,
    height: 48,
  },
  avatar: {
    width: 48,
    height: 48,
    borderRadius: 24,
    marginRight: 12,
    backgroundColor: colors.border,
    borderWidth: 1,
    borderColor: colors.border,
    overflow: 'hidden',
  },
  avatarImage: {
    width: '100%',
    height: '100%',
  },
  avatarPlaceholder: {
    width: '100%',
    height: '100%',
    backgroundColor: colors.border,
    justifyContent: 'center',
    alignItems: 'center',
  },
  avatarText: {
    fontSize: typography.fontSize.lg,
    fontWeight: typography.fontWeight.bold,
    color: colors.text,
  },
  userDetails: {
    flex: 1,
  },
  userNameContainer: {
    flexDirection: 'row',
    alignItems: 'center',
    marginBottom: 2,
  },
  userName: {
    fontSize: typography.fontSize.base,
    color: colors.text,
  },
  boldText: {
    fontWeight: typography.fontWeight.bold,
  },
  visitText: {
    fontSize: typography.fontSize.sm,
    color: colors.textTertiary,
  },
  ratingBadge: {
    width: 36,
    height: 36,
    borderRadius: 18,
    borderWidth: 1,
    borderColor: colors.border,
    justifyContent: 'center',
    alignItems: 'center',
    padding: 4,
  },
  ratingText: {
    fontSize: typography.fontSize.sm,
    fontWeight: typography.fontWeight.bold,
  },
  textSection: {
    marginBottom: 8,
  },
  postText: {
    fontSize: typography.fontSize.sm,
    color: colors.text,
    lineHeight: typography.lineHeight.normal * typography.fontSize.sm,
  },
  imagesContainer: {
    marginBottom: spacing.sm,
  },
  postImage: {
    width: 120,
    height: 120,
    borderRadius: 8,
    marginRight: 8,
    backgroundColor: colors.border,
    overflow: 'hidden',
  },
  postImageContent: {
    width: '100%',
    height: '100%',
  },
  engagementSection: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
  },
  likesContainer: {
    flexDirection: 'row',
    alignItems: 'center',
  },
  likesText: {
    fontSize: typography.fontSize.sm,
    color: colors.textSecondary,
    marginRight: spacing.xs,
  },
  actionIcons: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: spacing.md,
  },
  actionButton: {
    padding: spacing.xs,
  },
});
