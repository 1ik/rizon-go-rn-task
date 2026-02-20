import React from 'react';
import { StyleSheet, View } from 'react-native';
import { Button, Text } from 'react-native-paper';

export default function ReviewPromptStep() {
  return (
    <View style={styles.container}>
      <View style={styles.content}>
        <View style={styles.iconContainer}>
          <View style={styles.iconBackground}>
            <View style={styles.starburst}>
              <View style={styles.checkmark} />
            </View>
          </View>
        </View>
        <Text variant="headlineSmall" style={styles.title}>
          Got a minute to help us grow?
        </Text>
        <Text variant="bodyMedium" style={styles.subtitle}>
          It takes less than a minute and helps us a lot.
        </Text>
        <Button
          mode="contained"
          onPress={() => {}}
          style={styles.button}
          contentStyle={styles.buttonContent}
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
  iconBackground: {
    width: 80,
    height: 80,
    borderRadius: 20,
    backgroundColor: '#4285F4',
    alignItems: 'center',
    justifyContent: 'center',
  },
  starburst: {
    width: 50,
    height: 50,
    backgroundColor: '#C0C0C0',
    borderRadius: 25,
    alignItems: 'center',
    justifyContent: 'center',
    transform: [{ rotate: '45deg' }],
  },
  checkmark: {
    width: 20,
    height: 20,
    borderLeftWidth: 3,
    borderBottomWidth: 3,
    borderColor: '#FFFFFF',
    transform: [{ rotate: '-45deg' }, { translateX: -2 }, { translateY: -4 }],
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
    borderRadius: 8,
    backgroundColor: '#000000',
    width: '100%',
  },
  buttonContent: {
    paddingVertical: 8,
  },
});
