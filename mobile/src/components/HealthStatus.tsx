import React from 'react';
import { View, Text, StyleSheet, ActivityIndicator } from 'react-native';
import { HealthResponse } from '../services/api';

interface HealthStatusProps {
  healthData: HealthResponse | null;
  loading: boolean;
  error: string | null;
}

export const HealthStatus: React.FC<HealthStatusProps> = ({ 
  healthData, 
  loading, 
  error 
}) => {
  if (loading) {
    return (
      <View style={styles.container}>
        <ActivityIndicator size="large" color="#007AFF" />
        <Text style={styles.loadingText}>Checking backend health...</Text>
      </View>
    );
  }

  if (error) {
    return (
      <View style={styles.container}>
        <Text style={styles.errorText}>❌ Backend Error</Text>
        <Text style={styles.errorMessage}>{error}</Text>
      </View>
    );
  }

  if (healthData) {
    return (
      <View style={styles.container}>
        <Text style={styles.successText}>✅ Backend Connected</Text>
        <Text style={styles.statusText}>Status: {healthData.status}</Text>
        <Text style={styles.serviceText}>Service: {healthData.service}</Text>
        <Text style={styles.timestampText}>
          Last checked: {new Date(healthData.timestamp).toLocaleTimeString()}
        </Text>
      </View>
    );
  }

  return null;
};

const styles = StyleSheet.create({
  container: {
    alignItems: 'center',
    padding: 20,
    backgroundColor: '#f8f9fa',
    borderRadius: 10,
    margin: 20,
    shadowColor: '#000',
    shadowOffset: {
      width: 0,
      height: 2,
    },
    shadowOpacity: 0.1,
    shadowRadius: 3.84,
    elevation: 5,
  },
  loadingText: {
    marginTop: 10,
    fontSize: 16,
    color: '#666',
  },
  successText: {
    fontSize: 20,
    fontWeight: 'bold',
    color: '#28a745',
    marginBottom: 10,
  },
  errorText: {
    fontSize: 20,
    fontWeight: 'bold',
    color: '#dc3545',
    marginBottom: 10,
  },
  errorMessage: {
    fontSize: 14,
    color: '#666',
    textAlign: 'center',
  },
  statusText: {
    fontSize: 16,
    color: '#333',
    marginBottom: 5,
  },
  serviceText: {
    fontSize: 16,
    color: '#333',
    marginBottom: 5,
  },
  timestampText: {
    fontSize: 14,
    color: '#666',
    fontStyle: 'italic',
  },
});
