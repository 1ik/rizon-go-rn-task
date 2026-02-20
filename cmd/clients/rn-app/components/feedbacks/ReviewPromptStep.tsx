import * as Linking from 'expo-linking';
import React from 'react';
import { Image, Platform, StyleSheet, View } from 'react-native';
import { Button, Text } from 'react-native-paper';
import { useAuth } from '../../context/AuthContext';
import { useFeedback } from '../../context/FeedbackContext';

const APP_STORE_URLS = {
  ios: {
    app: 'itms-apps://itunes.apple.com/app/id6745082515',
    web: 'https://apps.apple.com/us/app/rizon-stablecoin-finance/id6745082515',
  },
  android: {
    app: 'market://details?id=com.rizon.app',
    web: 'https://play.google.com/store/apps/details?id=com.rizon.app',
  },
};

export default function ReviewPromptStep() {
  const step3Image = require('../../assets/step3.png');
  const { user } = useAuth();
  const { submitReview } = useFeedback();

  const handleLeaveReview = async () => {
    const urls =
      Platform.OS === 'ios' ? APP_STORE_URLS.ios : APP_STORE_URLS.android;
    
    try {
      // Try to open in the app store app first
      const canOpen = await Linking.canOpenURL(urls.app);
      if (canOpen) {
        await Linking.openURL(urls.app);
      } else {
        // Fallback to web URL
        await Linking.openURL(urls.web);
      }

      // Store review flag with user email
      if (user?.email) {
        await submitReview(user.email);
      }
    } catch (error) {
      // If app URL fails, try web URL as fallback
      try {
        await Linking.openURL(urls.web);
        // Store review flag even if we opened web URL
        if (user?.email) {
          await submitReview(user.email);
        }
      } catch (webError) {
        console.error('Failed to open app store:', webError);
      }
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
