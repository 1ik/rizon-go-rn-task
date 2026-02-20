import React from 'react';
import { StyleSheet, Text, View } from 'react-native';
import { useAuth } from '../context/AuthContext';

export default function Greeting() {
  const { user } = useAuth();
  if (!user?.email) return null;
  return (
    <View style={styles.container}>
      <Text style={styles.text}>Hello, {user.email}</Text>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    justifyContent: 'center',
    padding: 24,
  },
  text: {
    fontSize: 24,
    fontWeight: '600',
    color: '#333',
  },
});
