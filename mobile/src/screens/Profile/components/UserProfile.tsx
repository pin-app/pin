import React, { useState, useEffect } from 'react';
import { View, Text, StyleSheet, TouchableOpacity, Image, Alert } from 'react-native';
import FontAwesome6 from '@expo/vector-icons/FontAwesome6';
import { colors, spacing, typography } from '@/theme';
import { User } from '@/services/api';
import { apiService } from '@/services/api';
import Button from '@/shared/components/Button';

interface UserProfileProps {
  user: User;
  currentUserId?: string;
  showFollowButton?: boolean;
  onFollowChange?: (isFollowing: boolean) => void;
  onEditProfile?: () => void;
  onShareProfile?: () => void;
  postsCount?: number;
  followingCount?: number;
  followersCount?: number;
}

export default function UserProfile({
  user,
  currentUserId,
  showFollowButton = true,
  onFollowChange,
  onEditProfile,
  onShareProfile,
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
      <View style={styles.avatarContainer}>
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
      </View>

      <View style={styles.userInfoSection}>
        <Text style={styles.displayName}>{displayName}</Text>

        <View style={styles.usernameLocationRow}>
          {username && <Text style={styles.username}>{username}</Text>}
          {username && user.location && <Text style={styles.separator}>•</Text>}
          {user.location && (
            <View style={styles.locationContainer}>
              <FontAwesome6 name="location-dot" size={12} color={colors.textTertiary} style={styles.locationIcon} />
              <Text style={styles.location}>{user.location}</Text>
            </View>
          )}
        </View>

        <View style={styles.bioSection}>
          {user.bio && <Text style={styles.bio}>{user.bio}</Text>}
        </View>

        {/* buttons */}
        <View style={styles.actionButtonsContainer}>
          {isOwnProfile && (
            <>
              <Button 
                title="Edit Profile" 
                onPress={onEditProfile || (() => {})}
                variant="secondary"
                size="sm"
                style={styles.actionButton}
              />
              <Button 
                title="Share Profile" 
                onPress={onShareProfile || (() => {})}
                variant="secondary"
                size="sm"
                style={styles.actionButton}
              />
            </>
          )}
          {!isOwnProfile && showFollowButton && (
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

        {/* stats */}
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
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    backgroundColor: colors.postBackground,
    paddingHorizontal: spacing.md,
    marginTop: spacing.lg,
    marginBottom: spacing.sm,
    alignItems: 'center',
  },
  avatarContainer: {
    alignItems: 'center',
    marginBottom: spacing.lg,
  },
  avatar: {
    width: 150,
    height: 150,
    borderRadius: 75,
    backgroundColor: colors.border,
    borderWidth: 3,
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
    fontSize: 48,
    fontWeight: typography.fontWeight.bold,
    color: colors.text,
  },
  userInfoSection: {
    width: '100%',
    alignItems: 'center',
  },
  displayName: {
    fontSize: typography.fontSize.xl,
    fontWeight: typography.fontWeight.bold,
    color: colors.text,
    marginBottom: spacing.sm,
    textAlign: 'center',
  },
  usernameLocationRow: {
    flexDirection: 'row',
    alignItems: 'center',
    marginBottom: spacing.md,
    justifyContent: 'center',
  },
  username: {
    fontSize: typography.fontSize.base,
    color: colors.textSecondary,
  },
  separator: {
    fontSize: typography.fontSize.base,
    color: colors.textSecondary,
    marginHorizontal: spacing.sm,
  },
  bioSection: {
    width: '100%',
    alignItems: 'center',
    marginBottom: spacing.md,
  },
  bio: {
    fontSize: typography.fontSize.base,
    color: colors.text,
    lineHeight: typography.lineHeight.normal * typography.fontSize.base,
    textAlign: 'center',
  },
  locationContainer: {
    flexDirection: 'row',
    alignItems: 'center',
  },
  locationIcon: {
    marginRight: 4,
  },
  location: {
    fontSize: typography.fontSize.base,
    color: colors.textTertiary,
  },
  actionButtonsContainer: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'center',
    marginBottom: spacing.lg,
    gap: spacing.sm,
  },
  actionButton: {
    flex: 1,
  },
  followButton: {
    backgroundColor: colors.text,
    paddingHorizontal: spacing.lg,
    paddingVertical: spacing.sm,
    borderRadius: 20,
    minWidth: 100,
    alignItems: 'center',
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
    gap: spacing.xl,
    justifyContent: 'center',
  },
  statItem: {
    alignItems: 'center',
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
