import React from 'react';
import { View, StyleSheet } from 'react-native';
import { colors, spacing } from '../../theme';
import { PageTitle } from '../../shared/components';
import HealthCheck from './HealthCheck';

export default function MapScreen() {
  return (
    <View style={styles.container}>
      <PageTitle title="Map" />
      <HealthCheck />
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: colors.background,
    padding: spacing.md,
  },
});
