import React from 'react';
import { StyleSheet, View } from 'react-native';
import { Avatar, Text } from 'react-native-paper';
import { useAuth } from '../context/AuthContext';

export default function Greeting() {
  const { user } = useAuth();

  if (!user?.email) return null;

  // Get initials from email for avatar
  const getInitials = (email: string) => {
    const parts = email.split('@')[0];
    if (parts.length >= 2) {
      return parts.substring(0, 2).toUpperCase();
    }
    return email.substring(0, 1).toUpperCase();
  };

  return (
    <View style={styles.container}>
      <View style={styles.content}>
        <Avatar.Text
          size={56}
          label={getInitials(user.email)}
          style={styles.avatar}
          labelStyle={styles.avatarLabel}
        />
        <Text variant="titleMedium" style={styles.email}>
          {user.email}
        </Text>
      </View>
      <View style={styles.accentLine} />
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#ffffff',
    paddingHorizontal: 24,
    paddingTop: 32,
  },
  content: {
    alignItems: 'center',
    gap: 16,
  },
  avatar: {
    backgroundColor: '#f0f0f0',
  },
  avatarLabel: {
    color: '#666',
    fontWeight: '500',
    fontSize: 20,
  },
  email: {
    color: '#1a1a1a',
    fontWeight: '400',
    letterSpacing: 0.2,
  },
  accentLine: {
    width: 40,
    height: 3,
    backgroundColor: '#6200ee',
    borderRadius: 2,
    alignSelf: 'center',
    marginTop: 32,
  },
});
