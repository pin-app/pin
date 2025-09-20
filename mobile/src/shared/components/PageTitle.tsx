import React from 'react';
import { Text, StyleSheet } from 'react-native';
import { colors, spacing, typography } from '@/theme';

interface PageTitleProps {
  title: string;
  style?: any;
}

export default function PageTitle({ title, style }: PageTitleProps) {
  return (
    <Text style={[styles.title, style]}>
      {title}
    </Text>
  );
}

const styles = StyleSheet.create({
  title: {
    fontSize: typography.fontSize['2xl'],
    fontWeight: typography.fontWeight.bold,
    color: colors.text,
    textAlign: 'center',
    marginBottom: spacing.lg,
  },
});
