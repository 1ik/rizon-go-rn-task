import AsyncStorage from '@react-native-async-storage/async-storage';
import React, {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useState,
} from 'react';
import { useSubmitFeedbackMutation } from '../graphql/generated/graphql';
import { useAuth } from './AuthContext';

const getReviewStorageKey = (email: string) => `@rizon:has_left_review:${email}`;
const getOnboardingSeenKey = (email: string) => `rizon:has_seen_onboarding:${email}`;

type OnboardingContextValue = {
  isLoading: boolean;
  error: string | null;
  isSubmitting: boolean;
  submissionError: string | null;
  hasSeenOnboarding: boolean;
  markOnboardingSeen: () => Promise<void>;
  submitFeedback: (content: string) => Promise<void>;
  submitReview: (email: string) => Promise<void>;
  clearSubmissionError: () => void;
};

const OnboardingContext = createContext<OnboardingContextValue | null>(null);

export function OnboardingProvider({ children }: { children: React.ReactNode }) {
  const { user } = useAuth();
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [submissionError, setSubmissionError] = useState<string | null>(null);
  const [hasSeenOnboarding, setHasSeenOnboarding] = useState<boolean>(false);

  // Load hasSeenOnboarding from AsyncStorage for current user
  useEffect(() => {
    const loadOnboardingSeen = async () => {
      if (!user?.email) {
        setIsLoading(false);
        return;
      }
      try {
        const key = getOnboardingSeenKey(user.email);
        const value = await AsyncStorage.getItem(key);
        console.log(key, value);
        setHasSeenOnboarding(value === 'true');
      } catch (err) {
        console.error('Failed to load onboarding seen status:', err);
        setHasSeenOnboarding(false);
      } finally {
        setIsLoading(false);
      }
    };
    loadOnboardingSeen();
  }, [user?.email]);

  const markOnboardingSeen = useCallback(async () => {
    if (!user?.email) return;
    try {
      const key = getOnboardingSeenKey(user.email);
      await AsyncStorage.setItem(key, 'true');
      console.log('set onboarding to seen', key);
    } catch (err) {
      console.error('Failed to save onboarding seen status:', err);
    }
  }, [user, getOnboardingSeenKey]);

  const [submitFeedbackMutation, { loading: mutationLoading }] = useSubmitFeedbackMutation();

  const submitFeedback = useCallback(
    async (content: string): Promise<void> => {
     

      if (!content || content.trim() === '') {
        setSubmissionError('Please enter your feedback');
        throw new Error('Feedback content cannot be empty');
      }

      setSubmissionError(null);

      try {
        const result = await submitFeedbackMutation({
          variables: {
            deviceId: '1234',
            content: content.trim(),
          },
        });

        if (!result.data?.submitFeedback) {
          throw new Error('Failed to submit feedback');
        }

        // After successful submission, refetch the feedback
        setSubmissionError(null);
      } catch (err: any) {
        const errorMessage = err?.message || 'Failed to submit feedback';
        setSubmissionError(errorMessage);
        throw err;
      }
    },
    [submitFeedbackMutation]
  );

  const clearSubmissionError = useCallback(() => {
    setSubmissionError(null);
  }, []);

  const submitReview = useCallback(async (email: string): Promise<void> => {
    try {
      // Store review flag for this specific user's email
      const reviewKey = getReviewStorageKey(email);
      await AsyncStorage.setItem(reviewKey, 'true');

    } catch (err) {
      console.error('Failed to save review status:', err);
      throw new Error('Failed to save review status');
    }
  }, [user?.email]);

  const isLoadingState = isLoading || mutationLoading;

  const value = useMemo(
    () => ({
      isLoading: isLoadingState,
      error,
      isSubmitting: mutationLoading,
      submissionError,
      hasSeenOnboarding,
      markOnboardingSeen,
      submitFeedback,
      submitReview,
      clearSubmissionError,
    }),
    [
      isLoadingState,
      error,
      mutationLoading,
      submissionError,
      hasSeenOnboarding,
      markOnboardingSeen,
      submitFeedback,
      submitReview,
      clearSubmissionError,
    ]
  );

  return <OnboardingContext.Provider value={value}>{children}</OnboardingContext.Provider>;
}

export function useOnboarding(): OnboardingContextValue {
  const ctx = useContext(OnboardingContext);
  if (!ctx) throw new Error('useOnboarding must be used within OnboardingProvider');
  return ctx;
}
