import * as Linking from 'expo-linking';
import React from 'react';
import { Image, Platform, StyleSheet, View } from 'react-native';
import { Button, Text } from 'react-native-paper';

const APP_STORE_URLS: Record<string, string> = {
  ios: 'https://apps.apple.com/us/app/rizon-stablecoin-finance/id6745082515',
  android: 'https://play.google.com/store/apps/details?id=com.rizon.app',
};

interface ReviewPromptStepProps {
  onLeaveReviewPressed?: () => void;
}

export default function ReviewPromptStep({ onLeaveReviewPressed }: ReviewPromptStepProps) {
  const step3Image = require('../../assets/step3.png');

  const handleLeaveReview = async () => {
    const url = APP_STORE_URLS[Platform.OS] ?? APP_STORE_URLS.android;
    try {
      await Linking.openURL(url);
      onLeaveReviewPressed?.();
    } catch (err) {
      console.error('Failed to open app store:', err);
    }
  };

  return (
    <View style={styles.container}>
      <View style={styles.content}>
        <View style={styles.iconContainer}>
          <Image source={step3Image} style={styles.image} resizeMode="contain" />
        </View>
        <Text variant="headlineSmall" style={styles.title}>
          Got a minute to help us grow?
        </Text>
        <Text variant="bodyMedium" style={styles.subtitle}>
          It takes less than a minute and helps us a lot.
        </Text>
        <Button
          mode="contained"
          onPress={handleLeaveReview}
          style={styles.button}
        >
          Leave a review
        </Button>
      </View>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
  },
  content: {
    width: '100%',
    maxWidth: 360,
    alignItems: 'center',
  },
  iconContainer: {
    marginBottom: 24,
    alignItems: 'center',
    justifyContent: 'center',
  },
  image: {
    width: 80,
    height: 80,
  },
  title: {
    fontWeight: '600',
    color: '#1a1a1a',
    marginBottom: 8,
    textAlign: 'center',
  },
  subtitle: {
    color: '#666',
    textAlign: 'center',
    marginBottom: 32,
    opacity: 0.8,
  },
  button: {
    marginTop: 8,
  },
});
