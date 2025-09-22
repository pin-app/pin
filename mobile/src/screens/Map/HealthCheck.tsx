import React, { useState, useEffect } from 'react';
import { View, Text, StyleSheet, Button } from 'react-native';
import { colors, typography, spacing } from '../../theme';
import { apiService, HealthResponse } from '../../services/api';
import { useAuth } from '../../contexts/AuthContext';

// TODO: this is barely a health check anymore, this page in general is kinda useless and misleading rn
export default function HealthCheck() {
  const { isDevMode, user } = useAuth();

  return (
    <View style={styles.healthSection}>
      {isDevMode && (
        <View style={styles.devModeInfo}>
          <Text style={styles.devModeText}>🔧 Dev Mode Active</Text>
          <Text style={styles.devUserText}>User: {user?.username || user?.email || 'Unknown'}</Text>
        </View>
      )}
    </View>
  );
}

const styles = StyleSheet.create({
  healthSection: {
    flex: 1,
    alignItems: 'center',
    justifyContent: 'center',
  },
  sectionTitle: {
    fontSize: typography.fontSize.lg,
    fontWeight: typography.fontWeight.semibold,
    color: colors.text,
    marginBottom: spacing.md,
  },
  statusText: {
    fontSize: typography.fontSize.base,
    color: colors.text,
    marginBottom: spacing.sm,
  },
  errorText: {
    fontSize: typography.fontSize.base,
    color: colors.text,
    marginBottom: spacing.sm,
  },
  healthData: {
    alignItems: 'center',
    marginBottom: spacing.md,
  },
  healthText: {
    fontSize: typography.fontSize.base,
    color: colors.text,
    marginBottom: spacing.xs,
  },
  buttonContainer: {
    marginTop: spacing.md,
    width: 200,
  },
  devModeInfo: {
    backgroundColor: colors.textSecondary,
    padding: spacing.sm,
    borderRadius: 8,
    marginBottom: spacing.md,
    alignItems: 'center',
  },
  devModeText: {
    fontSize: typography.fontSize.sm,
    fontWeight: typography.fontWeight.semibold,
    color: colors.background,
    marginBottom: spacing.xs,
  },
  devUserText: {
    fontSize: typography.fontSize.xs,
    color: colors.background,
  },
});
