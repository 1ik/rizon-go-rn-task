import AsyncStorage from '@react-native-async-storage/async-storage';
import React, {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useState,
} from 'react';
import { Platform } from 'react-native';
import {
  Feedback,
  useGetUserFeedbackOnDeviceLazyQuery,
  useSubmitFeedbackMutation,
} from '../graphql/generated/graphql';
import { useAuth } from './AuthContext';

const getReviewStorageKey = (email: string) => `@rizon:has_left_review:${email}`;

// Dynamically import expo-application for platform-specific device IDs
async function getAndroidId(): Promise<string | null> {
  if (Platform.OS !== 'android') {
    return null;
  }
  
  try {
    const Application = await import('expo-application');
    // expo-application exports getAndroidId as a named export
    if (Application && typeof Application.getAndroidId === 'function') {
      return await Application.getAndroidId();
    }
    return null;
  } catch (err) {
    console.error('Failed to load expo-application:', err);
    return null;
  }
}

// Get iOS device identifier using Identifier for Vendor (IDFV)
// This is a real device identifier tied to the vendor
async function getIOSDeviceId(): Promise<string | null> {
  if (Platform.OS !== 'ios') {
    return null;
  }
  
  try {
    const Application = await import('expo-application');
    // expo-application exports getIosIdForVendorAsync for iOS
    if (Application && typeof Application.getIosIdForVendorAsync === 'function') {
      const idfv = await Application.getIosIdForVendorAsync();
      return idfv || null;
    }
    return null;
  } catch (err) {
    console.error('Failed to get iOS device ID:', err);
    return null;
  }
}

type FeedbackContextValue = {
  feedback: Feedback | null;
  isLoading: boolean;
  error: string | null;
  isSubmitting: boolean;
  submissionError: string | null;
  hasLeftReview: boolean;
  isOnBoardingComplete: boolean;
  submitFeedback: (content: string) => Promise<void>;
  submitReview: (email: string) => Promise<void>;
  refetchFeedback: () => Promise<void>;
  clearSubmissionError: () => void;
};

const FeedbackContext = createContext<FeedbackContextValue | null>(null);

export function FeedbackProvider({ children }: { children: React.ReactNode }) {
  const { user } = useAuth();
  const [feedback, setFeedback] = useState<Feedback | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [deviceId, setDeviceId] = useState<string | null>(null);
  const [submissionError, setSubmissionError] = useState<string | null>(null);
  const [hasLeftReview, setHasLeftReview] = useState<boolean>(false);
  const [isOnBoardingComplete, setIsOnBoardingComplete] = useState<boolean>(false);

  useEffect(() => {
    setIsOnBoardingComplete(hasLeftReview || feedback !== null);
  }, [hasLeftReview, feedback]);

  // Load review status from AsyncStorage for current user
  useEffect(() => {
    const loadReviewStatus = async () => {
      if (!user?.email) {
        setHasLeftReview(false);
        return;
      }

      try {
        const reviewKey = getReviewStorageKey(user.email);
        const reviewData = await AsyncStorage.getItem(reviewKey);
        if (reviewData === 'true') {
          setHasLeftReview(true);
        } else {
          setHasLeftReview(false);
        }
      } catch (err) {
        console.error('Failed to load review status:', err);
        setHasLeftReview(false);
      }
    };
    loadReviewStatus();
  }, [user?.email]);

  // Get real device ID on mount
  useEffect(() => {
    const getDeviceId = async () => {
      try {
        let realDeviceId: string | null = null;

        if (Platform.OS === 'android') {
          // Android: Use getAndroidId() for real device ID
          realDeviceId = await getAndroidId();
        } else if (Platform.OS === 'ios') {
          // iOS: Use getIosIdForVendorAsync() for Identifier for Vendor (IDFV)
          realDeviceId = await getIOSDeviceId();
        }

        if (!realDeviceId) {
          throw new Error('Unable to get device ID');
        }

        setDeviceId(realDeviceId);
      } catch (err) {
        console.error('Failed to get device ID:', err);
        setError('Failed to get device ID');
        setIsLoading(false);
      }
    };

    getDeviceId();
  }, []);

  const [getUserFeedback, { loading: queryLoading, data: queryData, error: queryError }] =
    useGetUserFeedbackOnDeviceLazyQuery({
      fetchPolicy: 'network-only',
    });

  const [submitFeedbackMutation, { loading: mutationLoading }] = useSubmitFeedbackMutation();

  // Handle query data updates
  useEffect(() => {
    if (queryData) {
      setFeedback(queryData.getUserFeedbackOnDevice || null);
      setError(null);
      setIsLoading(false);
    }
  }, [queryData]);

  // Handle query error updates
  useEffect(() => {
    if (queryError) {
      setError(queryError.message);
      setFeedback(null);
      setIsLoading(false);
    }
  }, [queryError]);

  const fetchFeedback = useCallback(async () => {
    if (!deviceId) return;
    
    setIsLoading(true);
    setError(null);
    try {
      await getUserFeedback({ variables: { deviceId } });
    } catch (err: any) {
      setError(err.message || 'Failed to fetch feedback');
      setIsLoading(false);
    }
  }, [deviceId, getUserFeedback]);

  // Fetch feedback when device ID is available
  useEffect(() => {
    if (deviceId) {
      fetchFeedback();
    }
  }, [deviceId, fetchFeedback]);

  const submitFeedback = useCallback(
    async (content: string): Promise<void> => {
      if (!deviceId) {
        throw new Error('Device ID not available');
      }

      if (!content || content.trim() === '') {
        setSubmissionError('Please enter your feedback');
        throw new Error('Feedback content cannot be empty');
      }

      setSubmissionError(null);

      try {
        const result = await submitFeedbackMutation({
          variables: {
            deviceId,
            content: content.trim(),
          },
        });

        if (!result.data?.submitFeedback) {
          throw new Error('Failed to submit feedback');
        }

        // After successful submission, refetch the feedback
        await fetchFeedback();
        setSubmissionError(null);
      } catch (err: any) {
        const errorMessage = err?.message || 'Failed to submit feedback';
        setSubmissionError(errorMessage);
        throw err;
      }
    },
    [deviceId, submitFeedbackMutation, fetchFeedback]
  );

  const refetchFeedback = useCallback(async () => {
    await fetchFeedback();
  }, [fetchFeedback]);

  const clearSubmissionError = useCallback(() => {
    setSubmissionError(null);
  }, []);

  const submitReview = useCallback(async (email: string): Promise<void> => {
    try {
      // Store review flag for this specific user's email
      const reviewKey = getReviewStorageKey(email);
      await AsyncStorage.setItem(reviewKey, 'true');

      setHasLeftReview(true);
    } catch (err) {
      console.error('Failed to save review status:', err);
      throw new Error('Failed to save review status');
    }
  }, [user?.email]);

  const isLoadingState = isLoading || queryLoading || mutationLoading;

  const value = useMemo(
    () => ({
      feedback,
      isLoading: isLoadingState,
      error,
      isSubmitting: mutationLoading,
      submissionError,
      hasLeftReview,
      isOnBoardingComplete,
      submitFeedback,
      submitReview,
      refetchFeedback,
      clearSubmissionError,
    }),
    [
      feedback,
      isLoadingState,
      error,
      mutationLoading,
      submissionError,
      hasLeftReview,
      isOnBoardingComplete,
      submitFeedback,
      submitReview,
      refetchFeedback,
      clearSubmissionError,
    ]
  );

  return <FeedbackContext.Provider value={value}>{children}</FeedbackContext.Provider>;
}

export function useFeedback(): FeedbackContextValue {
  const ctx = useContext(FeedbackContext);
  if (!ctx) throw new Error('useFeedback must be used within FeedbackProvider');
  return ctx;
}
