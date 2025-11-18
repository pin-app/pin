import React, { useEffect, useState, useCallback } from 'react';
import { SafeAreaView, View, Text, StyleSheet, TouchableOpacity, FlatList, RefreshControl } from 'react-native';
import FontAwesome6 from '@expo/vector-icons/FontAwesome6';
import { colors, spacing, typography } from '@/theme';
import { apiService, Notification } from '@/services/api';

interface NotificationsScreenProps {
  onBack: () => void;
  onOpenPost?: (postId: string) => void;
}

export default function NotificationsScreen({ onBack, onOpenPost }: NotificationsScreenProps) {
  const [notifications, setNotifications] = useState<Notification[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [isClearing, setIsClearing] = useState(false);
  const [refreshing, setRefreshing] = useState(false);

  const loadNotifications = useCallback(async (opts?: { refreshing?: boolean }) => {
    try {
      if (opts?.refreshing) {
        setRefreshing(true);
      } else {
        setIsLoading(true);
      }
      const data = await apiService.getNotifications(50, 0);
      setNotifications(data);
    } catch (error) {
      console.error('Failed to load notifications', error);
    } finally {
      if (opts?.refreshing) {
        setRefreshing(false);
      } else {
        setIsLoading(false);
      }
    }
  }, []);

  useEffect(() => {
    loadNotifications();
  }, [loadNotifications]);

  const handleClear = async () => {
    try {
      setIsClearing(true);
      await apiService.clearNotifications();
      await loadNotifications();
    } catch (error) {
      console.error('Failed to clear notifications', error);
    } finally {
      setIsClearing(false);
    }
  };

  const handlePressNotification = (notification: Notification) => {
    if (notification.post_id && onOpenPost) {
      onOpenPost(notification.post_id);
    }
  };

  const renderNotification = ({ item }: { item: Notification }) => {
    const actorName = item.actor?.display_name || item.actor?.username || 'Someone';

    let message = '';
    switch (item.type) {
      case 'like_post':
        message = `${actorName} liked your post`;
        break;
      case 'comment_post':
        message = `${actorName} commented on your post`;
        break;
      case 'comment_reply':
        message = `${actorName} replied to your comment`;
        break;
      default:
        message = `${actorName} interacted with you`;
    }

    return (
      <TouchableOpacity
        style={styles.notificationItem}
        onPress={() => handlePressNotification(item)}
        disabled={!item.post_id}
      >
        <View style={styles.notificationContent}>
          <Text style={styles.notificationText}>{message}</Text>
          <Text style={styles.notificationTime}>{new Date(item.created_at).toLocaleString()}</Text>
        </View>
        {item.post_id && (
          <FontAwesome6 name="chevron-right" size={14} color={colors.textSecondary} />
        )}
      </TouchableOpacity>
    );
  };

  return (
    <SafeAreaView style={styles.container}>
      <View style={styles.header}>
        <TouchableOpacity onPress={onBack} style={styles.backButton}>
          <FontAwesome6 name="arrow-left" size={18} color={colors.text} />
        </TouchableOpacity>
        <Text style={styles.headerTitle}>notifications</Text>
        <TouchableOpacity onPress={handleClear} style={styles.clearButton} disabled={isClearing}>
          <Text style={[styles.clearText, isClearing && styles.clearTextDisabled]}>
            {isClearing ? 'clearing...' : 'clear'}
          </Text>
        </TouchableOpacity>
      </View>

      <FlatList
        data={notifications}
        keyExtractor={(item) => item.id}
        renderItem={renderNotification}
        contentContainerStyle={styles.listContent}
        refreshControl={
          <RefreshControl refreshing={refreshing} onRefresh={() => loadNotifications({ refreshing: true })} />
        }
        ListEmptyComponent={
          !isLoading ? (
            <View style={styles.emptyState}>
              <Text style={styles.emptyText}>no notifications yet</Text>
            </View>
          ) : null
        }
      />
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
    padding: spacing.xs,
    marginRight: spacing.sm,
  },
  headerTitle: {
    flex: 1,
    fontSize: typography.fontSize.lg,
    fontWeight: typography.fontWeight.semibold,
    color: colors.text,
  },
  clearButton: {
    paddingHorizontal: spacing.sm,
    paddingVertical: spacing.xs,
  },
  clearText: {
    color: colors.textSecondary,
    fontSize: typography.fontSize.sm,
  },
  clearTextDisabled: {
    opacity: 0.5,
  },
  listContent: {
    paddingHorizontal: spacing.md,
  },
  notificationItem: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
    paddingVertical: spacing.md,
    borderBottomWidth: 1,
    borderBottomColor: colors.border,
  },
  notificationContent: {
    flex: 1,
    marginRight: spacing.md,
  },
  notificationText: {
    color: colors.text,
    fontSize: typography.fontSize.base,
    marginBottom: spacing.xs,
  },
  notificationTime: {
    color: colors.textSecondary,
    fontSize: typography.fontSize.xs,
  },
  emptyState: {
    padding: spacing.lg,
    alignItems: 'center',
  },
  emptyText: {
    color: colors.textSecondary,
  },
});

