import * as Linking from 'expo-linking';
import React, { useEffect } from 'react';
import { ActivityIndicator, StyleSheet, View } from 'react-native';
import { useAuth } from '../context/AuthContext';
import AuthenticatingScreen from './auth/AuthenticatingScreen';
import LoginForm from './auth/LoginForm';

/** If user is not authenticated, shows the login form; otherwise renders children. */
export default function Scaffold({ children }: { children: React.ReactNode }) {
  const { user, isLoading, authError, handleDeepLink } = useAuth();

  useEffect(() => {
    // Handle initial URL (cold start)
    Linking.getInitialURL().then((url) => {
      handleDeepLink(url);
    });

    // Listen for deep link events (warm start)
    const subscription = Linking.addEventListener('url', (event) => {
      handleDeepLink(event.url);
    });

    return () => {
      subscription.remove();
    };
  }, [handleDeepLink]);

  // Show authenticating screen during loading or when there's an auth error
  if (isLoading || authError) {
    return (
      <View style={styles.container}>
        <AuthenticatingScreen />
      </View>
    );
  }

  if (!user) {
    return (
      <View style={styles.container}>
        <LoginForm />
      </View>
    );
  }

  return <View style={styles.container}>{children}</View>;
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#fff',
  },
  center: {
    justifyContent: 'center',
    alignItems: 'center',
  },
});
