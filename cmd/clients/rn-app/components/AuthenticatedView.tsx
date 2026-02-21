import React from 'react';
import { View } from 'react-native';
import { useOnboarding } from '../context/OnboardingContext';
import Greeting from './Greeting';
import Header from './Header';
import OnboardingBottomSheet from './feedbacks/OnboardingBottomSheet';

export default function AuthenticatedView() {
  const { hasSeenOnboarding } = useOnboarding();

  return (
    <View style={{ flex: 1 }}>
      <Header />
      <Greeting />
      {!hasSeenOnboarding && <OnboardingBottomSheet />}
    </View>
  );
}
