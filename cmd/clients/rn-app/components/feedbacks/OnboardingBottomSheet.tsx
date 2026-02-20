import BottomSheet, { BottomSheetBackdrop, BottomSheetView } from '@gorhom/bottom-sheet';
import React, { useCallback, useMemo, useRef } from 'react';
import { StyleSheet, View } from 'react-native';
import { Button, Text } from 'react-native-paper';

/**
 * Example onboarding / feedback bottom sheet component.
 * Uses @gorhom/bottom-sheet for smooth gesture interactions.
 */
export default function OnboardingBottomSheet() {
  const bottomSheetRef = useRef<BottomSheet>(null);

  // Define snap points (25%, 50%, 90% of screen height)
  const snapPoints = useMemo(() => ['25%', '50%', '90%'], []);

  // Handle sheet changes
  const handleSheetChanges = useCallback((index: number) => {
    console.log('Bottom sheet index changed:', index);
  }, []);

  // Handle close button press
  const handleClose = useCallback(() => {
    bottomSheetRef.current?.close();
  }, []);

  // Handle expand button press
  const handleExpand = useCallback(() => {
    bottomSheetRef.current?.expand();
  }, []);

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
      index={0} // Start at first snap point (25%)
      snapPoints={snapPoints}
      onChange={handleSheetChanges}
      enablePanDownToClose
      backdropComponent={renderBackdrop}
      backgroundStyle={styles.bottomSheetBackground}
      handleIndicatorStyle={styles.handleIndicator}
    >
      <BottomSheetView style={styles.contentContainer}>
        <View style={styles.header}>
          <Text variant="headlineSmall" style={styles.title}>
            Welcome! 👋
          </Text>
          <Text variant="bodyMedium" style={styles.subtitle}>
            This is an example bottom sheet component
          </Text>
        </View>

        <View style={styles.body}>
          <Text variant="bodyLarge" style={styles.description}>
            You can drag the handle to resize the sheet, or use the buttons below to control it programmatically.
          </Text>

          <View style={styles.features}>
            <Text variant="bodyMedium" style={styles.featureText}>
              ✓ Smooth gesture interactions
            </Text>
            <Text variant="bodyMedium" style={styles.featureText}>
              ✓ Multiple snap points
            </Text>
            <Text variant="bodyMedium" style={styles.featureText}>
              ✓ Keyboard handling support
            </Text>
            <Text variant="bodyMedium" style={styles.featureText}>
              ✓ Scrollable content support
            </Text>
          </View>
        </View>

        <View style={styles.actions}>
          <Button
            mode="contained"
            onPress={handleExpand}
            style={styles.button}
            contentStyle={styles.buttonContent}
          >
            Expand Sheet
          </Button>
          <Button
            mode="outlined"
            onPress={handleClose}
            style={styles.button}
            contentStyle={styles.buttonContent}
          >
            Close
          </Button>
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
    paddingHorizontal: 24,
    paddingBottom: 32,
  },
  header: {
    marginBottom: 24,
    alignItems: 'center',
  },
  title: {
    fontWeight: '600',
    color: '#1a1a1a',
    marginBottom: 8,
  },
  subtitle: {
    color: '#666',
    textAlign: 'center',
  },
  body: {
    flex: 1,
    marginBottom: 24,
  },
  description: {
    color: '#333',
    marginBottom: 20,
    lineHeight: 22,
  },
  features: {
    gap: 12,
  },
  featureText: {
    color: '#555',
    lineHeight: 20,
  },
  actions: {
    gap: 12,
  },
  button: {
    borderRadius: 8,
  },
  buttonContent: {
    paddingVertical: 6,
  },
});
