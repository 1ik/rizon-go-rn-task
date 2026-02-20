import React from 'react';
import { Image, StyleSheet, View } from 'react-native';
import { Button, Text } from 'react-native-paper';

interface InitialQuestionStepProps {
  onYesLovingIt: () => void;
  onNotYet: () => void;
}

export default function InitialQuestionStep({
  onYesLovingIt,
  onNotYet,
}: InitialQuestionStepProps) {
  const logo = require('../../assets/rizon-logo.webp');

  return (
    <View style={styles.container}>
      <View style={styles.content}>
        <Image source={logo} style={styles.logo} resizeMode="contain" />
        <Text variant="headlineSmall" style={styles.title}>
          Enjoying Rizon so far?
        </Text>
        <Text variant="bodyMedium" style={styles.subtitle}>
          Your feedback helps us build a better money experience.
        </Text>
        <View style={styles.buttonRow}>
          <Button
            mode="outlined"
            onPress={onNotYet}
            style={styles.button}
            contentStyle={styles.buttonContent}
          >
            Not yet
          </Button>
          <Button
            mode="contained"
            onPress={onYesLovingIt}
            style={[styles.button, styles.buttonPrimary]}
            contentStyle={styles.buttonContent}
          >
            Yes, loving it
          </Button>
        </View>
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
  logo: {
    width: 80,
    height: 80,
    marginBottom: 24,
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
  buttonRow: {
    flexDirection: 'row',
    gap: 12,
    width: '100%',
  },
  button: {
    flex: 1,
    borderRadius: 8,
  },
  buttonPrimary: {
    backgroundColor: '#000000',
  },
  buttonContent: {
    paddingVertical: 8,
  },
});
