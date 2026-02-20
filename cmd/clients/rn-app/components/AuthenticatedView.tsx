import React from 'react';
import { View } from 'react-native';
import Greeting from './Greeting';
import Header from './Header';
import OnboardingBottomSheet from './feedbacks/OnboardingBottomSheet';

export default function AuthenticatedView() {
  return (
    <View style={{ flex: 1 }}>
      <Header />
        <Greeting />
        <OnboardingBottomSheet />
    </View>
  );
}
