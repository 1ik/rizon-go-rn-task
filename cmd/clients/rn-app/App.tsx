import { ApolloProvider } from '@apollo/client';
import { StatusBar } from 'expo-status-bar';
import { StyleSheet, View } from 'react-native';
import { GestureHandlerRootView } from 'react-native-gesture-handler';
import { PaperProvider } from 'react-native-paper';
import { SafeAreaProvider } from 'react-native-safe-area-context';
import { apolloClient } from './apolloClient';
import Scaffold from './components/Scaffold';
import AuthenticatedView from './components/AuthenticatedView';
import { AuthProvider } from './context/AuthContext';

export default function App() {
  return (
    <GestureHandlerRootView style={styles.gestureRoot}>
      <SafeAreaProvider>
        <ApolloProvider client={apolloClient}>
          <PaperProvider>
            <AuthProvider>
              <View style={styles.container}>
                <StatusBar style="auto" />
                <Scaffold>
                  <AuthenticatedView />
                </Scaffold>
              </View>
            </AuthProvider>
          </PaperProvider>
        </ApolloProvider>
      </SafeAreaProvider>
    </GestureHandlerRootView>
  );
}

const styles = StyleSheet.create({
  gestureRoot: {
    flex: 1,
  },
  container: {
    flex: 1,
    backgroundColor: '#fff',
  },
});
