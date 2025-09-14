import React, { useState } from 'react';
import { View, StyleSheet } from 'react-native';
import { colors, spacing } from '../../theme';
import { PageTitle, SearchBar } from '../../shared/components';

export default function HomeScreen() {
  const [searchValue, setSearchValue] = useState('');

  return (
    <View style={styles.container}>
      <PageTitle title="Feed" />
      <SearchBar 
        placeholder="search a place, member, etc"
        value={searchValue}
        onInputChange={setSearchValue}
        onSearchPress={() => console.log('Search:', searchValue)}
        onClear={() => setSearchValue('')}
      />
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: colors.background,
    alignItems: 'center',
    justifyContent: 'center',
    padding: spacing.md,
  },
});
