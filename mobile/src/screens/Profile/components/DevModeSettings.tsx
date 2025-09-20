import React, { useState } from 'react';
import {
  View,
  StyleSheet,
  Text,
  TouchableOpacity,
  TextInput,
  Alert,
  Modal,
} from 'react-native';
import { useAuth } from '../../../contexts/AuthContext';
import { colors, spacing } from '../../../theme';

interface DevModeSettingsProps {
  visible: boolean;
  onClose: () => void;
}

export default function DevModeSettings({ visible, onClose }: DevModeSettingsProps) {
  const { isDevMode, setDevMode, setDevUser, logout } = useAuth();
  const [devUserId, setDevUserId] = useState('');

  const handleDevModeToggle = async (enabled: boolean) => {
    try {
      await setDevMode(enabled);
      if (!enabled) {
        setDevUserId('');
      }
    } catch (error) {
      Alert.alert('Error', 'Failed to toggle dev mode');
    }
  };

  const handleDevUserSubmit = async () => {
    if (!devUserId.trim()) {
      Alert.alert('Error', 'Please enter a user ID');
      return;
    }

    try {
      await setDevUser(devUserId.trim());
      Alert.alert('Success', 'Dev user set successfully');
      onClose();
    } catch (error) {
      Alert.alert('Error', 'Failed to set dev user');
    }
  };

  const handleLogout = () => {
    Alert.alert(
      'Logout',
      'Are you sure you want to logout?',
      [
        { text: 'Cancel', style: 'cancel' },
        { text: 'Logout', style: 'destructive', onPress: logout },
      ]
    );
  };

  return (
    <Modal
      visible={visible}
      animationType="slide"
      presentationStyle="pageSheet"
      onRequestClose={onClose}
    >
      <View style={styles.container}>
        <View style={styles.header}>
          <Text style={styles.title}>Development Settings</Text>
          <TouchableOpacity onPress={onClose} style={styles.closeButton}>
            <Text style={styles.closeButtonText}>Done</Text>
          </TouchableOpacity>
        </View>

        <View style={styles.content}>
          <View style={styles.section}>
            <View style={styles.sectionHeader}>
              <Text style={styles.sectionTitle}>Dev Mode</Text>
              <TouchableOpacity
                style={[styles.toggle, isDevMode && styles.toggleActive]}
                onPress={() => handleDevModeToggle(!isDevMode)}
              >
                <Text style={[styles.toggleText, isDevMode && styles.toggleTextActive]}>
                  {isDevMode ? 'ON' : 'OFF'}
                </Text>
              </TouchableOpacity>
            </View>
            <Text style={styles.sectionDescription}>
              Enable dev mode to bypass authentication and use mock data
            </Text>
          </View>

          {isDevMode && (
            <View style={styles.section}>
              <Text style={styles.sectionTitle}>Dev User ID</Text>
              <Text style={styles.sectionDescription}>
                Enter a user ID to use for API requests in dev mode
              </Text>
              <TextInput
                style={styles.input}
                placeholder="Enter user ID (e.g., 123e4567-e89b-12d3-a456-426614174000)"
                value={devUserId}
                onChangeText={setDevUserId}
                autoCapitalize="none"
                autoCorrect={false}
                placeholderTextColor={colors.textSecondary}
              />
              <TouchableOpacity
                style={[styles.button, !devUserId.trim() && styles.buttonDisabled]}
                onPress={handleDevUserSubmit}
                disabled={!devUserId.trim()}
              >
                <Text style={[styles.buttonText, !devUserId.trim() && styles.buttonTextDisabled]}>
                  Set Dev User
                </Text>
              </TouchableOpacity>
            </View>
          )}

          <View style={styles.section}>
            <Text style={styles.sectionTitle}>API Base URL</Text>
            <Text style={styles.sectionDescription}>
              {process.env.EXPO_PUBLIC_API_BASE || 'http://localhost:8080'}
            </Text>
          </View>

          <View style={styles.section}>
            <TouchableOpacity
              style={styles.logoutButton}
              onPress={handleLogout}
            >
              <Text style={styles.logoutButtonText}>Logout</Text>
            </TouchableOpacity>
          </View>
        </View>
      </View>
    </Modal>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: colors.background,
  },
  header: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    padding: spacing.md,
    borderBottomWidth: 1,
    borderBottomColor: colors.border,
  },
  title: {
    fontSize: 18,
    fontWeight: '600',
    color: colors.text,
  },
  closeButton: {
    padding: spacing.sm,
  },
  closeButtonText: {
    fontSize: 16,
    color: colors.textSecondary,
    fontWeight: '500',
  },
  content: {
    flex: 1,
    padding: spacing.md,
  },
  section: {
    marginBottom: spacing.lg,
  },
  sectionHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: spacing.sm,
  },
  sectionTitle: {
    fontSize: 16,
    fontWeight: '600',
    color: colors.text,
  },
  sectionDescription: {
    fontSize: 14,
    color: colors.textSecondary,
    marginBottom: spacing.sm,
  },
  toggle: {
    paddingHorizontal: spacing.sm,
    paddingVertical: spacing.xs,
    borderRadius: 16,
    backgroundColor: colors.border,
  },
  toggleActive: {
    backgroundColor: colors.textSecondary,
  },
  toggleText: {
    fontSize: 12,
    fontWeight: '600',
    color: colors.textSecondary,
  },
  toggleTextActive: {
    color: colors.background,
  },
  input: {
    borderWidth: 1,
    borderColor: colors.border,
    borderRadius: 8,
    padding: spacing.sm,
    fontSize: 14,
    color: colors.text,
    backgroundColor: colors.background,
    marginBottom: spacing.sm,
  },
  button: {
    backgroundColor: colors.textSecondary,
    paddingHorizontal: spacing.md,
    paddingVertical: spacing.sm,
    borderRadius: 8,
    alignSelf: 'flex-start',
  },
  buttonDisabled: {
    backgroundColor: colors.border,
  },
  buttonText: {
    color: colors.background,
    fontSize: 14,
    fontWeight: '500',
  },
  buttonTextDisabled: {
    color: colors.textSecondary,
  },
  logoutButton: {
    backgroundColor: '#dc2626',
    paddingHorizontal: spacing.md,
    paddingVertical: spacing.sm,
    borderRadius: 8,
    alignItems: 'center',
  },
  logoutButtonText: {
    color: 'white',
    fontSize: 16,
    fontWeight: '600',
  },
});
