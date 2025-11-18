import React, { useState, useEffect, useCallback, useRef } from 'react';
import { View, StyleSheet, SafeAreaView, Alert } from 'react-native';
import { useNavigation, useFocusEffect } from '@react-navigation/native';
import { colors } from '@/theme';
import SearchBar from '@/components/SearchBar';
import SearchResults from '@/components/SearchResults';
import FeedHeader from './components/FeedHeader';
import Feed from './components/Feed';
import CommentsScreen from '../Comments';
import { apiService } from '@/services/api';
import { Post, Place, User } from '@/services/api';
import { useAuth } from '@/contexts/AuthContext';


interface SearchResult {
  id: string;
  type: 'place' | 'member' | 'recent';
  title: string;
  subtitle?: string;
  icon?: string;
  data?: any; // Store the actual data for navigation
}

export default function HomeScreen() {
  const navigation = useNavigation();
  const { user: currentUser } = useAuth();
  const [searchValue, setSearchValue] = useState('');
  const [posts, setPosts] = useState<Post[]>([]);
  const [isSearchFocused, setIsSearchFocused] = useState(false);
  const [searchResults, setSearchResults] = useState<SearchResult[]>([]);
  const [isSearching, setIsSearching] = useState(false);
  const [isLoading, setIsLoading] = useState(true);
  const [refreshing, setRefreshing] = useState(false);
  const [selectedPost, setSelectedPost] = useState<Post | null>(null);
  const [showComments, setShowComments] = useState(false);
  const hasFocusedOnce = useRef(false);

  const loadPosts = useCallback(async (options?: { useRefreshing?: boolean }) => {
    try {
      if (options?.useRefreshing) {
        setRefreshing(true);
      } else {
        setIsLoading(true);
      }
      const postsData = await apiService.getPosts(20, 0);
      setPosts(postsData);
    } catch (error) {
      console.error('Failed to load posts:', error);
    } finally {
      if (options?.useRefreshing) {
        setRefreshing(false);
      } else {
        setIsLoading(false);
      }
    }
  }, []);

  useEffect(() => {
    loadPosts();
  }, [loadPosts]);

  useFocusEffect(
    useCallback(() => {
      if (hasFocusedOnce.current) {
        loadPosts({ useRefreshing: true });
      } else {
        hasFocusedOnce.current = true;
      }
    }, [loadPosts])
  );

  const handleLike = async (postId: string) => {
    if (!currentUser) {
      Alert.alert('Sign in required', 'Log in to like posts.');
      return;
    }

    const post = posts.find(p => p.id === postId);
    if (!post) return;

    const wasLiked = post.liked_by_user;
    const previousState = { liked_by_user: post.liked_by_user, likes_count: post.likes_count };

    setPosts(prev =>
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
      if (wasLiked) {
        await apiService.unlikePost(postId);
      } else {
        await apiService.likePost(postId);
      }
    } catch (error) {
      console.error('Failed to update like', error);
      setPosts(prev =>
        prev.map(item =>
          item.id === postId ? { ...item, ...previousState } : item
        )
      );
      Alert.alert('Error', 'Could not update like.');
    }
  };

  const handleComment = (postId: string) => {
    const post = posts.find(p => p.id === postId);
    if (post) {
      setSelectedPost(post);
      setShowComments(true);
    }
  };

  const handleRate = (postId: string) => {
    console.log('Rate location for post:', postId);
  };

  const handleBookmark = (postId: string) => {
    console.log('Bookmark post:', postId);
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

  const handleMapPress = () => {
    console.log('Navigate to map');
  };

  const handleNotificationPress = () => {
    console.log('Open notifications');
  };

  const handleSearchFocus = () => {
    setIsSearchFocused(true);
  };

  const handleSearchBlur = () => {
    // idek what search blur is, do nothing for now
  };

  const handleSearch = async (query: string) => {
    if (!query.trim()) {
      setSearchResults([]);
      return;
    }

    try {
      setIsSearching(true);
      const [users, places] = await Promise.all([
        apiService.searchUsers(query, 5, 0),
        apiService.searchPlaces(query, 5, 0),
      ]);

      const results: SearchResult[] = [
        ...users.map(user => ({
          id: user.id,
          type: 'member' as const,
          title: user.display_name || user.username || 'Unknown User',
          subtitle: user.username ? `@${user.username}` : 'Member',
          icon: 'user',
          data: user,
        })),
        ...places.map(place => ({
          id: place.id,
          type: 'place' as const,
          title: place.name,
          subtitle: place.properties?.address || place.properties?.city || 'Place',
          icon: 'store',
          data: place,
        })),
      ];

      setSearchResults(results);
    } catch (error) {
      console.error('Search failed:', error);
      setSearchResults([]);
    } finally {
      setIsSearching(false);
    }
  };

  const resetSearch = () => {
    setIsSearchFocused(false);
    setSearchValue('');
    setSearchResults([]);
  };

  const handleSearchResultPress = (result: SearchResult) => {
    resetSearch();

    if (result.type === 'member' && result.data) {
      if (currentUser && result.data.id === currentUser.id) {
        (navigation as any).navigate('Profile');
      } else {
        (navigation as any).navigate('OtherUserProfile', {
          userId: result.data.id,
          username: result.data.username,
        });
      }
    } else if (result.type === 'place' && result.data) {
      (navigation as any).navigate('PlacePosts', {
        placeId: result.data.id,
        placeName: result.data.name,
      });
    }
  };

  const handleClearRecent = () => {
    setSearchResults(prev => prev.filter(result => result.type !== 'recent'));
  };

  const handleCloseSearch = () => {
    resetSearch();
  };

  const handleBackFromComments = () => {
    setShowComments(false);
    setSelectedPost(null);
  };

  if (showComments && selectedPost) {
    return (
      <CommentsScreen
        post={selectedPost}
        onBack={handleBackFromComments}
      />
    );
  }

  return (
    <SafeAreaView style={styles.container}>
      <FeedHeader 
        onMapPress={handleMapPress}
        onNotificationPress={handleNotificationPress}
        showSearchClose={isSearchFocused}
        onSearchClose={handleCloseSearch}
      />
      <SearchBar 
        placeholder="search a place, member, etc"
        value={searchValue}
        onInputChange={(value) => {
          setSearchValue(value);
          handleSearch(value);
        }}
        onSearchPress={() => handleSearch(searchValue)}
        onClear={() => {
          setSearchValue('');
          setSearchResults([]);
        }}
        onFocus={handleSearchFocus}
        onBlur={handleSearchBlur}
      />
      {isSearchFocused ? (
        <SearchResults
          results={searchResults}
          onResultPress={handleSearchResultPress}
          onClearRecent={handleClearRecent}
          onClose={resetSearch}
        />
      ) : (
        <Feed
          posts={posts}
          onLike={handleLike}
          onComment={handleComment}
          onRate={handleRate}
          onBookmark={handleBookmark}
          onUserPress={handleUserPress}
          refreshing={refreshing}
          onRefresh={() => loadPosts({ useRefreshing: true })}
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
});
