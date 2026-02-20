import React, { useState } from 'react';
import { StyleSheet, View } from 'react-native';
import { Button, Text, TextInput } from 'react-native-paper';
import { useFeedback } from '../../context/FeedbackContext';

interface FeedbackFormStepProps {
  onInputFocus?: () => void;
}

export default function FeedbackFormStep({ onInputFocus }: FeedbackFormStepProps) {
  const [feedback, setFeedback] = useState('');
  const { submitFeedback, isSubmitting, submissionError, clearSubmissionError } = useFeedback();

  const handleInputFocus = () => {
    onInputFocus?.();
  };

  return (
    <View style={styles.container}>
      <View style={styles.content}>
        <Text variant="headlineSmall" style={styles.title}>
          Help us improve Rizon
        </Text>
        <Text variant="bodyMedium" style={styles.subtitle}>
          Tell us what didn't feel right, we read every message.
        </Text>
        <TextInput
          mode="outlined"
          multiline
          numberOfLines={6}
          placeholder="Type your feedback here..."
          value={feedback}
          onChangeText={(text) => {
            setFeedback(text);
            // Clear error when user starts typing
            if (submissionError) {
              clearSubmissionError();
            }
          }}
          onFocus={handleInputFocus}
          style={styles.input}
          contentStyle={styles.inputContent}
          error={!!submissionError}
          disabled={isSubmitting}
        />
        {submissionError && (
          <Text variant="bodySmall" style={styles.errorText}>
            {submissionError}
          </Text>
        )}
        <Button
          mode="contained"
          onPress={() => {
            submitFeedback(feedback);
            // Clear input optimistically
            setFeedback('');
          }}
          style={styles.button}
          disabled={isSubmitting || !feedback.trim()}
          loading={isSubmitting}
        >
          Send feedback
        </Button>
      </View>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    justifyContent: 'flex-start',
    alignItems: 'center',
    paddingTop: 20,
  },
  content: {
    width: '100%',
    maxWidth: 360,
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
    marginBottom: 24,
    opacity: 0.8,
  },
  input: {
    marginBottom: 24,
    minHeight: 120,
  },
  inputContent: {
    paddingVertical: 12,
    textAlignVertical: 'top',
  },
  errorText: {
    color: '#d32f2f',
    marginTop: -20,
    marginBottom: 16,
    paddingHorizontal: 4,
  },
  button: {
    marginTop: 8,
  },
});
