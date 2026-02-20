import BottomSheet, {
  BottomSheetBackdrop,
  BottomSheetScrollView,
  BottomSheetView,
} from '@gorhom/bottom-sheet';
import React, { useCallback, useEffect, useMemo, useRef, useState } from 'react';
import { Animated, StyleSheet, useWindowDimensions, View } from 'react-native';
import FeedbackFormStep from './FeedbackFormStep';
import InitialQuestionStep from './InitialQuestionStep';
import ReviewPromptStep from './ReviewPromptStep';

const SLIDE_DURATION = 300;
const HORIZONTAL_PADDING = 24;

type Step = 0 | 1 | 2; // 0 = initial, 1 = feedback, 2 = review

/**
 * Multi-step onboarding / feedback bottom sheet component.
 * Uses @gorhom/bottom-sheet with animated slide transitions between steps.
 */
export default function OnboardingBottomSheet() {
  const bottomSheetRef = useRef<BottomSheet>(null);
  const { width } = useWindowDimensions();
  const contentWidth = width - HORIZONTAL_PADDING * 2;
  const [step, setStep] = useState<Step>(0);
  const slideProgress = useRef(new Animated.Value(0)).current;

  // Define snap points: Index 0 = 40%, Index 1 = 50%, Index 2 = 90%
  const snapPoints = useMemo(() => ['40%', '90%'], []);

  // Animate slide progress when step changes
  useEffect(() => {
    Animated.timing(slideProgress, {
      toValue: step,
      duration: SLIDE_DURATION,
      useNativeDriver: true,
    }).start();
  }, [step, slideProgress]);

  // Handler to expand sheet to 90% (index 2)
  const handleExpandSheet = useCallback(() => {
    bottomSheetRef.current?.snapToIndex(2);
  }, []);

  // Navigation handlers
  const handleYesLovingIt = useCallback(() => {
  
    setStep(1);
  }, [handleExpandSheet]);

  const handleNotYet = useCallback(() => {
    setStep(2);
  }, []);

  // Interpolation for sliding the entire row
  const sliderTranslateX = slideProgress.interpolate({
    inputRange: [0, 1, 2],
    outputRange: [0, -contentWidth, -contentWidth * 2],
  });

  // Render backdrop with opacity
  const renderBackdrop = useCallback(
    (props: any) => (
      <BottomSheetBackdrop
        {...props}
        disappearsOnIndex={-1}
        appearsOnIndex={0}
        opacity={0.5}
      />
    ),
    []
  );

  return (
    <BottomSheet
      ref={bottomSheetRef}
      index={0}
      snapPoints={snapPoints}
      animateOnMount={true}
      enablePanDownToClose
      backdropComponent={renderBackdrop}
      backgroundStyle={styles.bottomSheetBackground}
      handleIndicatorStyle={styles.handleIndicator}
      keyboardBehavior="extend"
      keyboardBlurBehavior="restore"
      android_keyboardInputMode="adjustResize"
    >
      <BottomSheetView style={styles.contentContainer}>
        <View style={styles.sliderWrapper}>
          <Animated.View
            style={[
              styles.sliderRow,
              { width: contentWidth * 3 },
              { transform: [{ translateX: sliderTranslateX }] },
            ]}
          >
            <View style={[styles.panel, { width: contentWidth }]}>
              <InitialQuestionStep
                onYesLovingIt={handleYesLovingIt}
                onNotYet={handleNotYet}
              />
            </View>
            <BottomSheetScrollView
              style={[styles.panel, { width: contentWidth }]}
              contentContainerStyle={styles.scrollContent}
              keyboardShouldPersistTaps="handled"
            >
              <FeedbackFormStep onInputFocus={handleExpandSheet} />
            </BottomSheetScrollView>
            <View style={[styles.panel, { width: contentWidth }]}>
              <ReviewPromptStep />
            </View>
          </Animated.View>
        </View>
      </BottomSheetView>
    </BottomSheet>
  );
}

const styles = StyleSheet.create({
  bottomSheetBackground: {
    backgroundColor: '#ffffff',
    borderTopLeftRadius: 20,
    borderTopRightRadius: 20,
  },
  handleIndicator: {
    backgroundColor: '#ccc',
    width: 40,
    height: 4,
  },
  contentContainer: {
    flex: 1,
    paddingHorizontal: HORIZONTAL_PADDING,
    paddingBottom: 32,
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
  },
  scrollContent: {
    flexGrow: 1,
    paddingBottom: 20,
  },
});
