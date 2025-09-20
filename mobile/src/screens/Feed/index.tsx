import React, { useState, useEffect } from 'react';
import { View, StyleSheet, SafeAreaView, Alert } from 'react-native';
import { colors } from '@/theme';
import SearchBar from '@/components/SearchBar';
import SearchResults from '@/components/SearchResults';
import FeedHeader from './components/FeedHeader';
import Feed from './components/Feed';
import CommentsScreen from '../Comments';
import { apiService } from '@/services/api';
import { Post, Place, User } from '@/services/api';

// TODO: merge the db tables pr and setup some type of dev mode data in the db by default so we
// dont have to hardcode this
const mockPosts = [
  {
    id: '1',
    user: {
      name: 'Pablo',
      avatar: 'none'
    },
    location: 'Doraville',
    visits: 3,
    dayOfWeek: 'Sunday',
    rating: 8.2,
    text: 'Great for bigbacking, grocery. It\'s calm. LanZhou Ramen, El Rey Del Taco, Lees Bakery',
    images: [
      { id: '1', uri: '' },
      { id: '2', uri: '' },
      { id: '3', uri: '' },
    ],
    likes: 9,
    isLiked: true,
  },
  {
    id: '2',
    user: {
      name: 'Alain',
      avatar: 'none'
    },
    location: 'Sandy Springs',
    visits: 2,
    dayOfWeek: 'Friday',
    rating: 0.3,
    text: 'Bro this place is 💩. Food is trash, nothing to do, streets are empty, boring asl',
    images: [
      { id: '1', uri: '' },
      { id: '2', uri: '' },
      { id: '3', uri: '' },
    ],
    likes: 2,
    isLiked: false,
  },
];

interface SearchResult {
  id: string;
  type: 'place' | 'member' | 'recent';
  title: string;
  subtitle?: string;
  icon?: string;
  data?: any; // Store the actual data for navigation
}

export default function HomeScreen() {
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
      // Fallback to mock data if API fails
      setPosts(mockPosts as any);
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
