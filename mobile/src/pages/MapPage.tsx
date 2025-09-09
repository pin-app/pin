import React, { useState, useEffect } from 'react';
import { View, Text, StyleSheet, Button } from 'react-native';
import { colors, typography, spacing } from '../theme';
import { apiService, HealthResponse } from '../services/api';

export default function MapPage() {
  const [healthData, setHealthData] = useState<HealthResponse | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const checkHealth = async () => {
    setLoading(true);
    setError(null);

    try {
      const data = await apiService.checkHealth();
      setHealthData(data);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Unknown error');
      setHealthData(null);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    checkHealth();
  }, []);

  return (
    <View style={styles.container}>
      <Text style={styles.title}>Map</Text>
      
      <View style={styles.healthSection}>
        <Text style={styles.sectionTitle}>Backend Health</Text>

        {loading && <Text style={styles.statusText}>Checking...</Text>}
        {error && <Text style={styles.errorText}>Error: {error}</Text>}
        {!loading && !error && healthData && (
          <View style={styles.healthData}>
            <Text style={styles.healthText}>Status: {healthData.status}</Text>
            <Text style={styles.healthText}>Service: {healthData.service}</Text>
            <Text style={styles.healthText}>Checks: {healthData.count}</Text>
            <Text style={styles.healthText}>Time: {new Date(healthData.timestamp).toLocaleTimeString()}</Text>
          </View>
        )}

        <View style={styles.buttonContainer}>
          <Button title="Check Health" onPress={checkHealth} />
        </View>
      </View>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: colors.background,
    padding: spacing.md,
  },
  title: {
    fontSize: typography.fontSize['2xl'],
    fontWeight: typography.fontWeight.bold,
    color: colors.text,
    textAlign: 'center',
    marginBottom: spacing.lg,
  },
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
});
