import React, { useState } from 'react';
import { View, Text, StyleSheet, FlatList, TouchableOpacity, Image, SafeAreaView } from 'react-native';
import Ionicons from '@expo/vector-icons/Ionicons';
import { ListRenderItemInfo } from 'react-native';
import { SearchBar } from '../../shared/components';
import Feather from '@expo/vector-icons/Feather';

// dummy data
type User = {
  id: string;
  name: string;
  profilePic: string;
};

const followersData: User[] = [
  { id: '1', name: 'Bob', profilePic: 'holder' },
  { id: '2', name: 'Billy', profilePic: 'holder' },
  { id: '3', name: 'Race', profilePic: 'holder' },
  { id: '4', name: 'Bob', profilePic: 'holder' },
  { id: '5', name: 'Billy', profilePic: 'holder' },
  { id: '6', name: 'Race', profilePic: 'holder' },
  { id: '7', name: 'Bob', profilePic: 'holder' },
  { id: '8', name: 'Billy', profilePic: 'holder' },
  { id: '9', name: 'Race', profilePic: 'holder' },
  { id: '10', name: 'Bob', profilePic: 'holder' },
  { id: '11', name: 'Billy', profilePic: 'holder' },
  { id: '12', name: 'Race', profilePic: 'holder' },
];

const followingData: User[] = [
  { id: '13', name: 'Tracy', profilePic: 'holder' },
  { id: '14', name: 'Ben', profilePic:'holder' },
  { id: '15', name: 'Anna', profilePic: 'holder' },
];

export default function FollowListScreen() {
  const [activeTab, setActiveTab] = useState('Followers');
  const [searchValue, setSearchValue] = useState('');
  const profileImagePlaceholder = 'holder'; 
  
  const renderItem = ({ item }: ListRenderItemInfo<User>) => (
    <View style={styles.listItem}>
      <Image 
        source={{ uri: profileImagePlaceholder }}
        style={styles.profilePic} 
      />
      <Text style={styles.username}>{item.name}</Text>
      
      {activeTab === 'Followers' ? (
        <TouchableOpacity style={styles.removeButton}>
          <Text style={styles.removeButtonText}>Remove</Text>
        </TouchableOpacity>
      ) : (
        <TouchableOpacity style={styles.followingButton}>
          <Text style={styles.followingButtonText}>Following</Text>
        </TouchableOpacity>
      )}
    </View>
  );

  return (
    <SafeAreaView style={styles.container}>
      {/* header */}
      <View style={styles.header}>
        <TouchableOpacity>
          <Feather name="x" style={styles.headerIcon}/>
        </TouchableOpacity>
        <Text style={styles.headerTitle}>Friends</Text>
        <TouchableOpacity>
          <Ionicons name="chevron-back" size={28} color="white" />  {/* white to make this invisible */}
        </TouchableOpacity>
      </View>

      {/* follower following tabs */}
      <View style={styles.tabContainer}>
        <TouchableOpacity 
          style={[styles.tab, activeTab === 'Followers' && styles.activeTab]} 
          onPress={() => setActiveTab('Followers')}
        >
          <Text style={[styles.tabText, activeTab === 'Followers' && styles.activeTabText]}>Followers</Text>
        </TouchableOpacity>

        <TouchableOpacity 
          style={[styles.tab, activeTab === 'Following' && styles.activeTab]} 
          onPress={() => setActiveTab('Following')}
        >
          <Text style={[styles.tabText, activeTab === 'Following' && styles.activeTabText]}>Following</Text>
        </TouchableOpacity>
      </View>

      {/* Search Bar */}
      <View style={{ marginTop: 15 }}>
        <SearchBar
          placeholder="Search for someone ..."
          value={searchValue}
          onInputChange={setSearchValue}
          onSearchPress={() => console.log('Search:', searchValue)}
          onClear={() => setSearchValue('')}
        />
      </View>

      <FlatList
        data={activeTab === 'Followers' ? followersData : followingData}
        renderItem={renderItem}
        keyExtractor={item => item.id}
        showsVerticalScrollIndicator={false}
      />
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: 'white',
  },
  header: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
    padding: 10,
    borderBottomWidth: StyleSheet.hairlineWidth,
    borderColor: '#CDCDCD',
  },
  headerTitle: {
    fontSize: 18,
    fontWeight: 'bold',
  },
  headerIcon: {
    fontSize: 28,
  },
  invisible: {
    opacity: 0,
  },
  
  tabContainer: {
    flexDirection: 'row',
    justifyContent: 'space-around',
    borderBottomWidth: StyleSheet.hairlineWidth,
    borderColor: '#CDCDCD',
  },
  tab: {
    flex: 1,
    paddingVertical: 15,
    alignItems: 'center',
  },
  tabText: {
    fontWeight: '600',
    color: 'gray',
  },
  activeTab: {
    borderBottomWidth: 2,
    borderColor: 'black',
  },
  activeTabText: {
    color: 'black',
  },
  
  listItem: {
    flexDirection: 'row',
    alignItems: 'center',
    padding: 15,
  },
  profilePic: {
    width: 50,
    height: 50,
    borderRadius: 25,
    backgroundColor: 'gray',
    marginRight: 10
  },
  username: {
    flex: 1,
    fontSize: 16,
    fontWeight: 'bold',
  },
  removeButton: {
    borderWidth: 1,
    borderColor: '#CDCDCD',
    borderRadius: 5,
    paddingVertical: 6,
    paddingHorizontal: 12,
  },
  removeButtonText: {
    fontWeight: '600',
    color: 'black',
  },
  followingButton: {
    borderWidth: 1,
    borderColor: '#CDCDCD',
    borderRadius: 5,
    paddingVertical: 6,
    paddingHorizontal: 12,
    backgroundColor: '#EFEFEF',
  },
  followingButtonText: {
    fontWeight: '600',
    color: 'black',
  },
});
