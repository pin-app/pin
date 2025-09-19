import React from 'react';
import { View, Text, StyleSheet, TouchableOpacity } from 'react-native';
import FontAwesome6 from '@expo/vector-icons/FontAwesome6';
import { colors, spacing, typography } from '@/theme';

interface FeedHeaderProps {
  onMapPress: () => void;
  onNotificationPress: () => void;
  showSearchClose?: boolean;
  onSearchClose?: () => void;
}

export default function FeedHeader({ onMapPress, onNotificationPress, showSearchClose, onSearchClose }: FeedHeaderProps) {
  return (
    <View style={styles.container}>
      {/* TODO: make this the actual logo */}
      <View style={styles.logoContainer}>
        <Text style={styles.logo}>pin</Text>
      </View>
      
      <View style={styles.rightIcons}>
        {showSearchClose ? (
          <TouchableOpacity onPress={onSearchClose} style={styles.iconButton}>
            <FontAwesome6 name="xmark" size={20} color={colors.iconDefault} />
          </TouchableOpacity>
        ) : (
          <>
            <TouchableOpacity onPress={onMapPress} style={styles.iconButton}>
              <FontAwesome6 name="map" size={20} color={colors.iconDefault} />
            </TouchableOpacity>
            <TouchableOpacity onPress={onNotificationPress} style={styles.iconButton}>
              <FontAwesome6 name="bell" size={20} color={colors.iconDefault} />
            </TouchableOpacity>
          </>
        )}
      </View>
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
  logoContainer: {
    flexDirection: 'row',
    alignItems: 'center',
  },
  logo: {
    fontSize: typography.fontSize['2xl'],
    fontWeight: typography.fontWeight.bold,
    color: colors.text,
  },
  pinIcon: {
    marginLeft: 2,
  },
  rightIcons: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: spacing.md,
  },
  iconButton: {
    padding: spacing.xs,
  },
});
