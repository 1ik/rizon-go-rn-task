import React from 'react';
import { StyleSheet, View } from 'react-native';
import { ActivityIndicator, Button, Text } from 'react-native-paper';
import { useAuth } from '../../context/AuthContext';

export default function AuthenticatingScreen() {
  const { authError, clearAuthError } = useAuth();

  if (authError) {
    return (
      <View style={styles.container}>
        <Text variant="titleLarge" style={styles.errorTitle}>
          Authentication Failed
        </Text>
        <Text style={styles.errorText}>{authError}</Text>
        <Button
          mode="contained"
          onPress={clearAuthError}
          style={styles.button}
        >
          Back to Login
        </Button>
      </View>
    );
  }

  return (
    <View style={styles.container}>
      <ActivityIndicator size="large" style={styles.spinner} />
      <Text variant="titleLarge" style={styles.text}>
        Authenticating...
      </Text>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    padding: 24,
  },
  spinner: {
    marginBottom: 16,
  },
  text: {
    opacity: 0.8,
  },
  errorTitle: {
    marginBottom: 16,
    color: '#B00020',
  },
  errorText: {
    marginBottom: 24,
    textAlign: 'center',
    opacity: 0.8,
    color: '#B00020',
  },
  button: {
    marginTop: 8,
  },
});
