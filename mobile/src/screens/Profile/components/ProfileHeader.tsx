import React from 'react';
import { View, Text, StyleSheet, TouchableOpacity } from 'react-native';
import FontAwesome6 from '@expo/vector-icons/FontAwesome6';
import { colors } from '@/theme/colors';
import { typography } from '@/theme/typography';
import { spacing } from '@/theme/spacing';

interface ProfileHeaderProps {
  username: string;
  onMenuPress: () => void;
  showBackButton?: boolean;
}

export default function ProfileHeader({ username, onMenuPress, showBackButton = false }: ProfileHeaderProps) {
  return (
    <View style={styles.container}>
      {showBackButton ? (
        <TouchableOpacity onPress={onMenuPress} style={styles.backButton}>
          <FontAwesome6 name="arrow-left" size={20} color={colors.iconDefault} />
        </TouchableOpacity>
      ) : (
        <View style={styles.placeholder} />
      )}
      
      <Text style={styles.username}>{username}</Text>
      
      {!showBackButton && (
        <TouchableOpacity onPress={onMenuPress} style={styles.menuButton}>
          <FontAwesome6 name="bars" size={20} color={colors.iconDefault} />
        </TouchableOpacity>
      )}
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    paddingHorizontal: spacing.md,
    paddingTop: spacing.sm,
    paddingBottom: spacing.sm,
    backgroundColor: colors.background,
  },
  username: {
    fontSize: typography.fontSize['2xl'],
    fontWeight: typography.fontWeight.bold,
    color: colors.text,
    position: 'absolute',
    left: 0,
    right: 0,
    textAlign: 'center',
    zIndex: -1,
  },
  backButton: {
    padding: spacing.xs,
    zIndex: 1,
  },
  menuButton: {
    padding: spacing.xs,
  },
  placeholder: {
    width: 32,
    height: 32,
  },
});
