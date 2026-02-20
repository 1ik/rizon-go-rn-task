import React, { useCallback, useEffect, useMemo, useRef, useState } from 'react';
import {
  Animated,
  KeyboardAvoidingView,
  Platform,
  StyleSheet,
  useWindowDimensions,
  View,
} from 'react-native';
import { Button, Text, TextInput } from 'react-native-paper';
import { useRequestEmailAuthLinkMutation } from '../../graphql/generated/graphql';

type Screen = 'form' | 'checkEmail';

const SLIDE_DURATION = 300;

const HORIZONTAL_PADDING = 24;

export default function LoginForm() {
  const { width } = useWindowDimensions();
  const contentWidth = width - HORIZONTAL_PADDING * 2;
  const [email, setEmail] = useState('');
  const [screen, setScreen] = useState<Screen>('form');
  const [isSubmitting, setIsSubmitting] = useState(false);
  const slideProgress = useRef(new Animated.Value(0)).current;

  const [requestEmailAuthLink, { loading, error }] =
    useRequestEmailAuthLinkMutation();

  const isLoading = loading || isSubmitting;

  const formTranslateX = slideProgress.interpolate({
    inputRange: [0, 1],
    outputRange: [0, -contentWidth],
  });
  const checkEmailTranslateX = slideProgress.interpolate({
    inputRange: [0, 1],
    outputRange: [0, -contentWidth],
  });

  useEffect(() => {
    if (screen === 'checkEmail') {
      Animated.timing(slideProgress, {
        toValue: 1,
        duration: SLIDE_DURATION,
        useNativeDriver: true,
      }).start();
    } else {
      Animated.timing(slideProgress, {
        toValue: 0,
        duration: SLIDE_DURATION,
        useNativeDriver: true,
      }).start();
    }
  }, [screen, slideProgress]);

  const handleSubmit = useCallback(async () => {
    const trimmed = email.trim();
    if (!trimmed || isLoading) return;
    
    setIsSubmitting(true);
    try {
      // Add 2 second delay before making the API call so that the loading is visible, only for demo purpose.
      await new Promise(resolve => setTimeout(resolve, 2000));
      const result = await requestEmailAuthLink({ variables: { email: trimmed } });
      if (result.data?.requestEmailAuthLink) {
        setScreen('checkEmail');
      }
    } catch {
      // Error is available in `error` from the hook
      // Don't slide on error - stay on form view
    } finally {
      setIsSubmitting(false);
    }
  }, [email, isLoading, requestEmailAuthLink]);

  const handleBack = useCallback(() => {
    setScreen('form');
  }, []);

  const Form = useMemo(
    () => (
      <View style={styles.form}>
        <Text variant="titleLarge" style={styles.title}>
          Sign in using email
        </Text>
        <Text style={styles.hint}>
          A link will be sent to your email. Enter your email below to log in.
        </Text>
        <TextInput
          label="Email"
          mode="outlined"
          value={email}
          onChangeText={setEmail}
          keyboardType="email-address"
          autoCapitalize="none"
          autoCorrect={false}
          style={styles.input}
          editable={!isLoading}
        />
        {error && <Text style={styles.errorText}>{error.message}</Text>}
        <Button
          mode="contained"
          onPress={handleSubmit}
          style={styles.button}
          loading={isLoading}
          disabled={isLoading}
        >
          Send login link
        </Button>
      </View>
    ),
    [email, isLoading, error, handleSubmit]
  );

  const CheckEmail = useMemo(
    () => (
      <View style={styles.form}>
        <Text variant="titleLarge" style={styles.title}>
          Check your email
        </Text>
        <Text style={styles.hint}>
          Please look into your email for a link for you to sign in.
        </Text>
        <Button mode="outlined" onPress={handleBack} style={styles.button}>
          Back
        </Button>
      </View>
    ),
    [handleBack]
  );

  return (
    <KeyboardAvoidingView
      style={styles.container}
      behavior={Platform.OS === 'ios' ? 'padding' : undefined}
    >
      <View style={styles.sliderWrapper}>
        <View style={[styles.sliderRow, { width: contentWidth * 2 }]}>
          <Animated.View
            style={[
              styles.panel,
              { width: contentWidth },
              { transform: [{ translateX: formTranslateX }] },
            ]}
          >
            {Form}
          </Animated.View>
          <Animated.View
            style={[
              styles.panel,
              { width: contentWidth },
              { transform: [{ translateX: checkEmailTranslateX }] },
            ]}
          >
            {CheckEmail}
          </Animated.View>
        </View>
      </View>
    </KeyboardAvoidingView>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    justifyContent: 'center',
    padding: HORIZONTAL_PADDING,
  },
  sliderWrapper: {
    flex: 1,
    overflow: 'hidden',
  },
  sliderRow: {
    flexDirection: 'row',
    flex: 1,
  },
  panel: {
    flex: 1,
    justifyContent: 'center',
  },
  form: {
    width: '100%',
    maxWidth: 360,
    alignSelf: 'center',
  },
  title: {
    marginBottom: 8,
  },
  hint: {
    marginBottom: 20,
    opacity: 0.8,
  },
  input: {
    marginBottom: 16,
  },
  errorText: {
    color: '#B00020',
    marginBottom: 8,
  },
  button: {
    marginTop: 8,
  },
});
