import React from 'react';
import { View, Text, StyleSheet, TouchableOpacity } from 'react-native';
import { useAuth } from '../../../contexts/AuthContext';
import { colors, spacing } from '../../../theme';

interface DebugInfoProps {
  visible?: boolean;
}

export default function DebugInfo({ visible = false }: DebugInfoProps) {
  const { user, isAuthenticated, isDevMode, sessionToken } = useAuth();

  if (!visible) return null;

  return (
    <View style={styles.container}>
      <Text style={styles.title}>Debug Info</Text>
      <Text style={styles.info}>Authenticated: {isAuthenticated ? 'Yes' : 'No'}</Text>
      <Text style={styles.info}>Dev Mode: {isDevMode ? 'Yes' : 'No'}</Text>
      <Text style={styles.info}>User ID: {user?.id || 'None'}</Text>
      <Text style={styles.info}>Username: {user?.username || user?.display_name || 'None'}</Text>
      <Text style={styles.info}>Email: {user?.email || 'None'}</Text>
      <Text style={styles.info}>Session Token: {sessionToken ? 'Present' : 'None'}</Text>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    backgroundColor: colors.background,
    padding: spacing.sm,
    margin: spacing.sm,
    borderRadius: 8,
    borderWidth: 1,
    borderColor: colors.border,
  },
  title: {
    fontSize: 14,
    fontWeight: 'bold',
    color: colors.text,
    marginBottom: spacing.xs,
  },
  info: {
    fontSize: 12,
    color: colors.textSecondary,
    marginBottom: 2,
  },
});
