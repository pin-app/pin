import React from 'react';
import { View, Text, StyleSheet, TouchableOpacity } from 'react-native';
import FontAwesome6 from '@expo/vector-icons/FontAwesome6';
import { colors } from '@/theme/colors';
import { typography } from '@/theme/typography';
import { spacing } from '@/theme/spacing';

interface ProfileHeaderProps {
  username: string;
  onMenuPress: () => void;
}

export default function ProfileHeader({ username, onMenuPress }: ProfileHeaderProps) {
  return (
    <View style={styles.container}>
      {/* Username */}
      <Text style={styles.username}>{username}</Text>
      
      {/* Menu button */}
      <TouchableOpacity onPress={onMenuPress} style={styles.menuButton}>
        <FontAwesome6 name="bars" size={20} color={colors.iconDefault} />
      </TouchableOpacity>
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
  },
  menuButton: {
    padding: spacing.xs,
  },
});
