import React from 'react';
import { View } from 'react-native';
import { useFeedback } from '../context/FeedbackContext';
import Greeting from './Greeting';
import Header from './Header';
import OnboardingBottomSheet from './feedbacks/OnboardingBottomSheet';

export default function AuthenticatedView() {
  const { feedback, isOnBoardingComplete } = useFeedback();

  return (
    <View style={{ flex: 1 }}>
      <Header />
      <Greeting />
      {!feedback && !isOnBoardingComplete && <OnboardingBottomSheet />}
    </View>
  );
}
