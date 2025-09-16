import React from 'react';
import { View, StyleSheet } from 'react-native';
import { colors, spacing } from '../../theme';
import { PageTitle, Button } from '../../shared/components';

export default function ProfileScreen() {
  const handleEditProfile = () => {
    // useless for now
    console.log('Edit profile pressed');
  };

  return (
    <View style={styles.container}>
      <PageTitle title="Profile" />
      <Button 
        title="Edit Profile" 
        onPress={handleEditProfile}
        variant="primary"
        size="md"
        style={styles.button}
      />
      <Button 
        title="Share" 
        onPress={handleEditProfile}
        variant="secondary"
        size="md"
        style={styles.button}
      />
      <Button 
        title="Share but with a border around it" 
        onPress={handleEditProfile}
        variant="outline"
        size="md"
        style={styles.button}
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
  button: {
    marginTop: spacing.lg,
  },
});
