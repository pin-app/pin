import React, { useState, useEffect } from 'react';
import { View, Text, StyleSheet, TouchableOpacity, Image, Alert } from 'react-native';
import FontAwesome6 from '@expo/vector-icons/FontAwesome6';
import { colors, spacing, typography } from '@/theme';
import { User } from '@/services/api';
import { apiService } from '@/services/api';

interface UserProfileProps {
  user: User;
  currentUserId?: string;
  showFollowButton?: boolean;
  onFollowChange?: (isFollowing: boolean) => void;
  postsCount?: number;
  followingCount?: number;
  followersCount?: number;
}

export default function UserProfile({
  user,
  currentUserId,
  showFollowButton = true,
  onFollowChange,
  postsCount = 0,
  followingCount = 0,
  followersCount = 0,
}: UserProfileProps) {
  const [isFollowing, setIsFollowing] = useState(false);
  const [isLoading, setIsLoading] = useState(false);

  const isOwnProfile = currentUserId === user.id;

  useEffect(() => {
    if (!isOwnProfile && showFollowButton) {
      checkFollowStatus();
    }
  }, [user.id, isOwnProfile, showFollowButton]);

  const checkFollowStatus = async () => {
    try {
      const following = await apiService.isFollowing(user.id);
      setIsFollowing(following);
    } catch (error) {
      console.error('Failed to check follow status:', error);
    }
  };

  const handleFollowToggle = async () => {
    if (isLoading) return;

    try {
      setIsLoading(true);
      if (isFollowing) {
        await apiService.unfollowUser(user.id);
        setIsFollowing(false);
        onFollowChange?.(false);
      } else {
        await apiService.followUser(user.id);
        setIsFollowing(true);
        onFollowChange?.(true);
      }
    } catch (error) {
      console.error('Failed to toggle follow:', error);
      Alert.alert('Error', 'Failed to update follow status');
    } finally {
      setIsLoading(false);
    }
  };

  const displayName = user.display_name || user.username || 'Unknown User';
  const username = user.username ? `@${user.username}` : '';

  return (
    <View style={styles.container}>
      <View style={styles.topSection}>
        <View style={styles.avatar}>
          {user.pfp_url ? (
            <Image source={{ uri: user.pfp_url }} style={styles.avatarImage} />
          ) : (
            <View style={styles.avatarPlaceholder}>
              <Text style={styles.avatarText}>
                {displayName.charAt(0).toUpperCase()}
              </Text>
            </View>
          )}
        </View>
        <View style={styles.nameAndStats}>
          <View style={styles.nameSection}>
            <Text style={styles.displayName}>{displayName}</Text>
            {username && <Text style={styles.username}>{username}</Text>}
            <View style={styles.statsRow}>
              <View style={styles.statItem}>
                <Text style={styles.statValue}>{postsCount}</Text>
                <Text style={styles.statLabel}>posts</Text>
              </View>
              <View style={styles.statItem}>
                <Text style={styles.statValue}>{followingCount}</Text>
                <Text style={styles.statLabel}>following</Text>
              </View>
              <View style={styles.statItem}>
                <Text style={styles.statValue}>{followersCount}</Text>
                <Text style={styles.statLabel}>followers</Text>
              </View>
            </View>
          </View>
          {showFollowButton && !isOwnProfile && (
            <TouchableOpacity
              style={[styles.followButton, isFollowing && styles.followingButton]}
              onPress={handleFollowToggle}
              disabled={isLoading}
            >
              <Text style={[styles.followButtonText, isFollowing && styles.followingButtonText]}>
                {isFollowing ? 'Following' : 'Follow'}
              </Text>
            </TouchableOpacity>
          )}
        </View>
      </View>
      <View style={styles.bioSection}>
        {user.bio && <Text style={styles.bio}>{user.bio}</Text>}
        {user.location && (
          <View style={styles.locationContainer}>
            <FontAwesome6 name="location-dot" size={12} color={colors.textTertiary} style={styles.locationIcon} />
            <Text style={styles.location}>{user.location}</Text>
          </View>
        )}
      </View>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    backgroundColor: colors.postBackground,
    paddingHorizontal: spacing.md,
    marginTop: spacing.lg,
    marginBottom: spacing.sm,
  },
  topSection: {
    flexDirection: 'row',
    alignItems: 'flex-start',
    marginBottom: spacing.md,
  },
  avatar: {
    width: 100,
    height: 100,
    borderRadius: 50,
    backgroundColor: colors.border,
    borderWidth: 2,
    borderColor: colors.border,
    overflow: 'hidden',
    shadowColor: '#000',
    shadowOpacity: 0.1,
    shadowRadius: 8,
    shadowOffset: { width: 0, height: 4 },
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
    fontSize: 36,
    fontWeight: typography.fontWeight.bold,
    color: colors.text,
  },
  nameAndStats: {
    flex: 1,
    marginLeft: spacing.lg,
    justifyContent: 'space-between',
  },
  nameSection: {
    marginBottom: spacing.sm,
  },
  displayName: {
    fontSize: typography.fontSize.xl,
    fontWeight: typography.fontWeight.bold,
    color: colors.text,
    marginBottom: 2,
  },
  username: {
    fontSize: typography.fontSize.base,
    color: colors.textSecondary,
  },
  bioSection: {
    marginTop: spacing.sm,
  },
  bio: {
    fontSize: typography.fontSize.base,
    color: colors.text,
    lineHeight: typography.lineHeight.normal * typography.fontSize.base,
    marginBottom: spacing.sm,
  },
  locationContainer: {
    flexDirection: 'row',
    alignItems: 'center',
  },
  locationIcon: {
    marginRight: 4,
  },
  location: {
    fontSize: typography.fontSize.sm,
    color: colors.textTertiary,
  },
  followButton: {
    backgroundColor: colors.text,
    paddingHorizontal: spacing.md,
    paddingVertical: spacing.sm,
    borderRadius: 20,
    minWidth: 80,
    alignItems: 'center',
    alignSelf: 'flex-start',
  },
  followingButton: {
    backgroundColor: 'transparent',
    borderWidth: 1,
    borderColor: colors.border,
  },
  followButtonText: {
    color: colors.background,
    fontSize: typography.fontSize.sm,
    fontWeight: typography.fontWeight.bold,
  },
  followingButtonText: {
    color: colors.text,
  },
  statsRow: {
    flexDirection: 'row',
    marginTop: spacing.sm,
    gap: spacing.md,
  },
  statItem: {
    alignItems: 'flex-start',
  },
  statValue: {
    fontSize: typography.fontSize.lg,
    fontWeight: typography.fontWeight.bold,
    color: colors.text,
  },
  statLabel: {
    fontSize: typography.fontSize.sm,
    color: colors.textSecondary,
    marginTop: 2,
  },
});
