import React, { useState, useEffect } from 'react';
import { 
  View, 
  StyleSheet, 
  Text, 
  SafeAreaView, 
  TouchableOpacity, 
  ScrollView, 
  TextInput,
  Alert,
  Image,
  Keyboard
} from 'react-native';
import { useNavigation } from '@react-navigation/native';
import FontAwesome6 from '@expo/vector-icons/FontAwesome6';
import { colors, spacing, typography } from '@/theme';
import { Post as PostType, Comment as CommentType } from '@/services/api';
import { apiService } from '@/services/api';
import Post from '@/shared/components/Post';
import { useAuth } from '@/contexts/AuthContext';

interface CommentsScreenProps {
  post: PostType;
  onBack: () => void;
}

type CommentWithReplies = CommentType & { 
  replies: CommentWithReplies[];
  isOptimistic?: boolean;
};

export default function CommentsScreen({ post, onBack }: CommentsScreenProps) {
  const navigation = useNavigation();
  const { user: currentUser } = useAuth();
  const [comments, setComments] = useState<CommentWithReplies[]>([]);
  const [newComment, setNewComment] = useState('');
  const [isLoading, setIsLoading] = useState(true);
  const [replyingTo, setReplyingTo] = useState<CommentWithReplies | null>(null);

  useEffect(() => {
    loadComments();
  }, [post.id]);

  const loadComments = async () => {
    try {
      setIsLoading(true);
      const commentsData = await apiService.getCommentsByPost(post.id);
      
      const commentMap = new Map<string, CommentWithReplies>();
      const rootComments: CommentWithReplies[] = [];
      
      commentsData.forEach(comment => {
        commentMap.set(comment.id, { ...comment, replies: [] });
      });
      
      commentsData.forEach(comment => {
        const commentWithReplies = commentMap.get(comment.id)!;
        if (comment.parent_id) {
          const parent = commentMap.get(comment.parent_id);
          if (parent) {
            parent.replies.push(commentWithReplies);
          }
        } else {
          rootComments.push(commentWithReplies);
        }
      });
      
      setComments(rootComments);
    } catch (error) {
      console.error('Failed to load comments:', error);
      Alert.alert('Error', 'Failed to load comments');
    } finally {
      setIsLoading(false);
    }
  };

  const handleAddComment = async () => {
    if (!newComment.trim() || !currentUser) return;

    const content = newComment.trim();
    const parentId = replyingTo?.id;
    
    // Create optimistic comment
    const optimisticComment: CommentWithReplies = {
      id: `optimistic-${Date.now()}`,
      post_id: post.id,
      user_id: currentUser.id,
      parent_id: parentId,
      content: content,
      created_at: new Date().toISOString(),
      updated_at: new Date().toISOString(),
      user: currentUser,
      replies: [],
      isOptimistic: true,
    };

    // Optimistically update UI
    if (parentId && replyingTo) {
      // Adding a reply
      setComments(prev => {
        const addReplyToComment = (comments: CommentWithReplies[]): CommentWithReplies[] => {
          return comments.map(comment => {
            if (comment.id === parentId) {
              return {
                ...comment,
                replies: [...comment.replies, optimisticComment],
              };
            }
            if (comment.replies.length > 0) {
              return {
                ...comment,
                replies: addReplyToComment(comment.replies),
              };
            }
            return comment;
          });
        };
        return addReplyToComment(prev);
      });
    } else {
      // Adding a root comment
      setComments(prev => [optimisticComment, ...prev]);
    }

    // Clear input and reset reply state
    setNewComment('');
    setReplyingTo(null);
    Keyboard.dismiss();

    try {
      // Make actual API call
      const newCommentData = await apiService.createComment(post.id, content, parentId);
      
      // Replace optimistic comment with real one
      setComments(prev => {
        const replaceOptimistic = (comments: CommentWithReplies[]): CommentWithReplies[] => {
          return comments.map(comment => {
            if (comment.id === optimisticComment.id) {
              return { ...newCommentData, replies: [] };
            }
            if (comment.replies.length > 0) {
              return {
                ...comment,
                replies: replaceOptimistic(comment.replies),
              };
            }
            return comment;
          });
        };
        return replaceOptimistic(prev);
      });
    } catch (error) {
      console.error('Failed to add comment:', error);
      
      // Remove optimistic comment on error
      setComments(prev => {
        const removeOptimistic = (comments: CommentWithReplies[]): CommentWithReplies[] => {
          return comments
            .filter(comment => comment.id !== optimisticComment.id)
            .map(comment => ({
              ...comment,
              replies: removeOptimistic(comment.replies),
            }));
        };
        return removeOptimistic(prev);
      });
      
      Alert.alert('Error', 'Failed to add comment. Please try again.');
      // Restore the comment text so user can retry
      setNewComment(content);
      setReplyingTo(parentId ? replyingTo : null);
    }
  };

  const handleReply = (comment: CommentWithReplies) => {
    setReplyingTo(comment);
    // Focus is handled by the TextInput when replyingTo changes
  };

  const handleCancelReply = () => {
    setReplyingTo(null);
  };

  const handleUserPress = (userId: string, username?: string) => {
    // If it's the current user's own profile, navigate to Profile tab
    if (currentUser && userId === currentUser.id) {
      (navigation as any).navigate('Profile');
    } else {
      // Otherwise, navigate to OtherUserProfile
      (navigation as any).navigate('OtherUserProfile', {
        userId,
        username,
      });
    }
  };

  const renderComment = (comment: CommentWithReplies, isReply = false) => {
    const timeAgo = new Date(comment.created_at).toLocaleDateString();
    
    return (
      <View key={comment.id} style={[styles.commentContainer, isReply && styles.replyContainer]}>
        <View style={styles.commentHeader}>
          <TouchableOpacity 
            style={styles.commentAvatar}
            onPress={() => handleUserPress(comment.user?.id || '', comment.user?.username)}
            disabled={!comment.user?.id}
          >
            {comment.user?.pfp_url ? (
              <Image source={{ uri: comment.user.pfp_url }} style={styles.avatarImage} />
            ) : (
              <FontAwesome6 name="user" size={16} color={colors.textSecondary} />
            )}
          </TouchableOpacity>
          <View style={styles.commentContent}>
            <View style={styles.commentUserInfo}>
              <TouchableOpacity 
                onPress={() => handleUserPress(comment.user?.id || '', comment.user?.username)}
                disabled={!comment.user?.id}
              >
                <Text style={[
                  styles.commentUsername,
                  comment.isOptimistic && styles.optimisticText
                ]}>
                  {comment.user?.display_name || comment.user?.username || 'Unknown User'}
                </Text>
              </TouchableOpacity>
              <Text style={styles.commentTime}>
                {comment.isOptimistic ? 'Posting...' : timeAgo}
              </Text>
            </View>
            <Text style={[
              styles.commentText,
              comment.isOptimistic && styles.optimisticText
            ]}>
              {comment.content}
            </Text>
            <TouchableOpacity 
              style={styles.replyButton}
              onPress={() => handleReply(comment)}
            >
              <Text style={styles.replyButtonText}>Reply</Text>
            </TouchableOpacity>
          </View>
        </View>
        
        {comment.replies && comment.replies.length > 0 && comment.replies.map((reply) => renderComment(reply as CommentWithReplies, true))}
      </View>
    );
  };

  return (
    <SafeAreaView style={styles.container}>
      <View style={styles.header}>
        <TouchableOpacity onPress={onBack} style={styles.backButton}>
          <FontAwesome6 name="arrow-left" size={20} color={colors.text} />
        </TouchableOpacity>
        <View style={styles.headerSpacer} />
      </View>

      <ScrollView style={styles.scrollContainer} showsVerticalScrollIndicator={false}>
        <View style={styles.postContainer}>
          <Post
            post={post}
            likes={0}
            isLiked={false}
            onLike={() => {}}
            onComment={() => {}}
            onRate={() => {}}
            onBookmark={() => {}}
            onUserPress={handleUserPress}
            showCommentsButton={false}
          />
        </View>

        {/* Comments Section */}
        <View style={styles.commentsSection}>
          <Text style={styles.commentsTitle}>{comments.length} Comments</Text>
          
          {isLoading ? (
            <View style={styles.loadingContainer}>
              <Text style={styles.loadingText}>Loading comments...</Text>
            </View>
          ) : (
            <View style={styles.commentsList}>
              {comments.map(comment => renderComment(comment))}
            </View>
          )}
        </View>
      </ScrollView>

      {/* Add Comment Input */}
      <View style={styles.addCommentContainer}>
        {replyingTo && (
          <View style={styles.replyingToBar}>
            <Text style={styles.replyingToText}>
              Replying to {replyingTo.user?.display_name || replyingTo.user?.username || 'Unknown User'}
            </Text>
            <TouchableOpacity onPress={handleCancelReply} style={styles.cancelReplyButton}>
              <FontAwesome6 name="xmark" size={14} color={colors.textSecondary} />
            </TouchableOpacity>
          </View>
        )}
        <View style={styles.inputRow}>
          <TextInput
            style={styles.commentInput}
            placeholder={replyingTo ? "Write a reply..." : "Add a comment..."}
            placeholderTextColor={colors.textSecondary}
            value={newComment}
            onChangeText={setNewComment}
            multiline
          />
          <TouchableOpacity 
            style={[styles.sendButton, !newComment.trim() && styles.sendButtonDisabled]}
            onPress={handleAddComment}
            disabled={!newComment.trim()}
          >
            <FontAwesome6 name="paper-plane" size={16} color={newComment.trim() ? colors.text : colors.textSecondary} />
          </TouchableOpacity>
        </View>
      </View>
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
    paddingVertical: spacing.sm,
    borderBottomWidth: 1,
    borderBottomColor: colors.border,
  },
  backButton: {
    padding: spacing.sm,
    marginLeft: -spacing.sm,
  },
  headerSpacer: {
    flex: 1,
  },
  scrollContainer: {
    flex: 1,
  },
  postContainer: {
    backgroundColor: colors.postBackground,
    marginBottom: spacing.sm,
  },
  commentsSection: {
    paddingHorizontal: spacing.md,
  },
  commentsTitle: {
    fontSize: typography.fontSize.lg,
    fontWeight: typography.fontWeight.semibold,
    color: colors.text,
    marginBottom: spacing.md,
  },
  loadingContainer: {
    padding: spacing.lg,
    alignItems: 'center',
  },
  loadingText: {
    color: colors.textSecondary,
    fontSize: typography.fontSize.sm,
  },
  commentsList: {
    paddingBottom: spacing.xl,
  },
  commentContainer: {
    marginBottom: spacing.md,
  },
  replyContainer: {
    marginLeft: spacing.lg,
    paddingLeft: spacing.md,
    borderLeftWidth: 2,
    borderLeftColor: colors.border,
    marginTop: spacing.md,
  },
  commentHeader: {
    flexDirection: 'row',
  },
  commentAvatar: {
    width: 32,
    height: 32,
    borderRadius: 16,
    backgroundColor: colors.border,
    borderWidth: 1,
    borderColor: colors.border,
    justifyContent: 'center',
    alignItems: 'center',
    marginRight: spacing.sm,
    overflow: 'hidden',
  },
  avatarImage: {
    width: '100%',
    height: '100%',
  },
  commentContent: {
    flex: 1,
  },
  commentUserInfo: {
    flexDirection: 'row',
    alignItems: 'center',
    marginBottom: spacing.xs,
  },
  commentUsername: {
    fontSize: typography.fontSize.sm,
    fontWeight: typography.fontWeight.semibold,
    color: colors.text,
    marginRight: spacing.sm,
  },
  commentTime: {
    fontSize: typography.fontSize.xs,
    color: colors.textSecondary,
  },
  commentText: {
    fontSize: typography.fontSize.sm,
    color: colors.text,
    lineHeight: typography.lineHeight.normal * typography.fontSize.sm,
    marginBottom: spacing.xs,
  },
  replyButton: {
    alignSelf: 'flex-start',
  },
  replyButtonText: {
    fontSize: typography.fontSize.xs,
    color: colors.textSecondary,
    fontWeight: typography.fontWeight.medium,
  },
  addCommentContainer: {
    borderTopWidth: 1,
    borderTopColor: colors.border,
    backgroundColor: colors.background,
  },
  replyingToBar: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    paddingHorizontal: spacing.md,
    paddingTop: spacing.sm,
    paddingBottom: spacing.xs,
    backgroundColor: colors.border,
  },
  replyingToText: {
    fontSize: typography.fontSize.xs,
    color: colors.text,
    fontWeight: typography.fontWeight.medium,
  },
  cancelReplyButton: {
    padding: spacing.xs,
  },
  inputRow: {
    flexDirection: 'row',
    alignItems: 'flex-end',
    paddingHorizontal: spacing.md,
    paddingVertical: spacing.sm,
  },
  commentInput: {
    flex: 1,
    borderWidth: 1,
    borderColor: colors.border,
    borderRadius: 20,
    paddingHorizontal: spacing.md,
    paddingVertical: spacing.sm,
    fontSize: typography.fontSize.sm,
    color: colors.text,
    backgroundColor: colors.postBackground,
    maxHeight: 100,
    marginRight: spacing.sm,
  },
  sendButton: {
    width: 40,
    height: 40,
    borderRadius: 20,
    backgroundColor: colors.border,
    justifyContent: 'center',
    alignItems: 'center',
  },
  sendButtonDisabled: {
    opacity: 0.5,
  },
  optimisticText: {
    opacity: 0.6,
  },
});
