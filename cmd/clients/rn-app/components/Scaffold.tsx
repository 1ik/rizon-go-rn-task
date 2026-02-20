import React from 'react';
import { StyleSheet, View } from 'react-native';
import { useAuth } from '../context/AuthContext';
import LoginForm from './auth/LoginForm';

/** If user is not authenticated, shows the login form; otherwise renders children. */
export default function Scaffold({ children }: { children: React.ReactNode }) {
  const { user } = useAuth();

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
});
