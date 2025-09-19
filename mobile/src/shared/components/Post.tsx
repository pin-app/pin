import React from 'react';
import { View, Text, StyleSheet, TouchableOpacity, Image, ScrollView } from 'react-native';
import FontAwesome6 from '@expo/vector-icons/FontAwesome6';
import { colors, spacing, typography } from '@/theme';

interface PostImage {
  id: string;
  uri: string;
}

interface PostProps {
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
  images: PostImage[];
  likes: number;
  isLiked: boolean;
  onLike: (postId: string) => void;
  onComment: (postId: string) => void;
  onRate: (postId: string) => void;
  onBookmark: (postId: string) => void;
}

export default function Post({
  id,
  user,
  location,
  visits,
  dayOfWeek,
  rating,
  text,
  images,
  likes,
  isLiked,
  onLike,
  onComment,
  onRate,
  onBookmark,
}: PostProps) {
  const getRatingColor = (rating: number) => {
    if (rating >= 7) return colors.ratingPositive;
    if (rating >= 4) return colors.ratingNeutral;
    return colors.ratingNegative;
  };

  return (
    <View style={styles.container}>
      {/* User info and rating */}
      <View style={styles.header}>
        <View style={styles.avatar} />
        <View style={styles.userDetails}>
          <Text style={styles.userName}>
            <Text style={styles.boldText}>{user.name}</Text> ranked <Text style={styles.boldText}>{location}</Text>
          </Text>
          <Text style={styles.visitText}>{visits} visits • {dayOfWeek}</Text>
        </View>
        <View style={styles.ratingBadge}>
          <Text style={[styles.ratingText, { color: getRatingColor(rating) }]}>{rating}</Text>
        </View>
      </View>

      {/* Post text */}
      <View style={styles.textSection}>
        <Text style={styles.postText}>{text}</Text>
      </View>

      {/* Images */}
      {images.length > 0 && (
        <ScrollView horizontal showsHorizontalScrollIndicator={false} style={styles.imagesContainer}>
          {images.map((image) => (
            <View key={image.id} style={styles.postImage} />
          ))}
        </ScrollView>
      )}

      {/* Engagement icons */}
      <View style={styles.engagementSection}>
        <View style={styles.likesContainer}>
          <Text style={styles.likesText}>{likes} likes</Text>
        </View>
        <View style={styles.actionIcons}>
          <TouchableOpacity onPress={() => onLike(id)} style={styles.actionButton}>
            <FontAwesome6 
              name="heart" 
              size={20} 
              color={isLiked ? colors.ratingNegative : colors.iconDefault} 
              solid={isLiked}
            />
          </TouchableOpacity>
          <TouchableOpacity onPress={() => onComment(id)} style={styles.actionButton}>
            <FontAwesome6 name="comment" size={20} color={colors.iconDefault} />
          </TouchableOpacity>
          <TouchableOpacity onPress={() => onRate(id)} style={styles.actionButton}>
            <FontAwesome6 name="plus" size={20} color={colors.iconDefault} />
          </TouchableOpacity>
          <TouchableOpacity onPress={() => onBookmark(id)} style={styles.actionButton}>
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
    height: 40,
  },
  avatar: {
    width: 40,
    height: 40,
    borderRadius: 20,
    marginRight: 12,
    backgroundColor: colors.border,
  },
  userDetails: {
    flex: 1,
  },
  userName: {
    fontSize: typography.fontSize.base,
    color: colors.text,
    marginBottom: 2,
  },
  boldText: {
    fontWeight: typography.fontWeight.bold,
  },
  visitText: {
    fontSize: typography.fontSize.sm,
    color: colors.textTertiary,
  },
  ratingBadge: {
    width: 32,
    height: 32,
    borderRadius: 16,
    borderWidth: 2,
    borderColor: colors.text,
    justifyContent: 'center',
    alignItems: 'center',
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
    width: 100,
    height: 100,
    borderRadius: 8,
    marginRight: 8,
    backgroundColor: colors.border,
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
