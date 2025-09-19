import React, { useState } from 'react';
import { View, StyleSheet, SafeAreaView } from 'react-native';
import { colors } from '@/theme';
import SearchBar from '@/components/SearchBar';
import SearchResults from '@/components/SearchResults';
import FeedHeader from './components/FeedHeader';
import Feed from './components/Feed';

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

// TODO: same for search results. These should be stored in local storage on the app
const mockSearchResults = [
  { id: '1', type: 'recent' as const, title: 'Doraville', subtitle: 'Recent search' },
  { id: '2', type: 'recent' as const, title: 'Sandy Springs', subtitle: 'Recent search' },
  { id: '3', type: 'place' as const, title: 'Buford Highway Farmers Market', subtitle: 'Atlanta, GA', icon: 'store' },
  { id: '4', type: 'member' as const, title: 'Pablo', subtitle: 'Member', icon: 'user' },
  { id: '5', type: 'member' as const, title: 'Alain', subtitle: 'Member', icon: 'user' },
];

export default function HomeScreen() {
  const [searchValue, setSearchValue] = useState('');
  const [posts, setPosts] = useState(mockPosts);
  const [isSearchFocused, setIsSearchFocused] = useState(false);
  const [searchResults, setSearchResults] = useState(mockSearchResults);

  const handleLike = (postId: string) => {
    setPosts(prevPosts =>
      prevPosts.map(post =>
        post.id === postId
          ? {
              ...post,
              isLiked: !post.isLiked,
              likes: post.isLiked ? post.likes - 1 : post.likes + 1,
            }
          : post
      )
    );
  };

  const handleComment = (postId: string) => {
    console.log('Comment on post:', postId);
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

  const handleSearchResultPress = (result: any) => {
    console.log('Search result pressed:', result);
  };

  const handleClearRecent = () => {
    setSearchResults(prev => prev.filter(result => result.type !== 'recent'));
  };

  const handleCloseSearch = () => {
    setIsSearchFocused(false);
    setSearchValue(''); // clear the search value when closing
  };

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
        onInputChange={setSearchValue}
        onSearchPress={() => console.log('Search:', searchValue)}
        onClear={() => setSearchValue('')}
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
