import React, { useState, useEffect } from 'react';
import { StyleSheet, Text, View, Button } from 'react-native';
import { apiService, HealthResponse } from './src/services/api';

export default function App() {
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
      <Text style={styles.title}>Pin</Text>
      <Text>Backend Health</Text>

      {loading && <Text>Checking...</Text>}
      {error && <Text>Error: {error}</Text>}
      {!loading && !error && healthData && (
        <>
          <Text>Status: {healthData.status}</Text>
          <Text>Service: {healthData.service}</Text>
          <Text>Time: {new Date(healthData.timestamp).toLocaleTimeString()}</Text>
        </>
      )}

      <View style={styles.buttonContainer}>
        <Button title="Check Health" onPress={checkHealth} />
      </View>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    alignItems: 'center',
    justifyContent: 'center',
    padding: 16,
  },
  title: {
    fontSize: 20,
    marginBottom: 8,
  },
  buttonContainer: {
    marginTop: 16,
    width: 200,
  },
});
