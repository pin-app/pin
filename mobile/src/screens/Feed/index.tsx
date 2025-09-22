import React, { useState, useEffect } from 'react';
import { View, StyleSheet, SafeAreaView, Alert } from 'react-native';
import { useNavigation } from '@react-navigation/native';
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
  const [selectedPost, setSelectedPost] = useState<Post | null>(null);
  const [showComments, setShowComments] = useState(false);

  useEffect(() => {
    loadPosts();
  }, []);

  const loadPosts = async () => {
    try {
      setIsLoading(true);
      const postsData = await apiService.getPosts(20, 0);
      setPosts(postsData);
    } catch (error) {
      console.error('Failed to load posts:', error);
    } finally {
      setIsLoading(false);
    }
  };

  const handleLike = (postId: string) => {
    // TODO: Implement real like functionality with backend
    console.log('Like post:', postId);
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

  const handleSearchResultPress = (result: SearchResult) => {
    console.log('Search result pressed:', result);
    // TODO: Navigate to user profile or place details
    if (result.type === 'member') {
      // Navigate to user profile
      console.log('Navigate to user profile:', result.data);
    } else if (result.type === 'place') {
      // Navigate to place details or show posts for this place
      console.log('Navigate to place:', result.data);
    }
  };

  const handleClearRecent = () => {
    setSearchResults(prev => prev.filter(result => result.type !== 'recent'));
  };

  const handleCloseSearch = () => {
    setIsSearchFocused(false);
    setSearchValue(''); // clear the search value when closing
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
        />
      ) : (
        <Feed
          posts={posts}
          onLike={handleLike}
          onComment={handleComment}
          onRate={handleRate}
          onBookmark={handleBookmark}
          onUserPress={handleUserPress}
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
