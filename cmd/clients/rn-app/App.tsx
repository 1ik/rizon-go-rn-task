import { ApolloProvider } from '@apollo/client';
import { StatusBar } from 'expo-status-bar';
import { StyleSheet, View } from 'react-native';
import { PaperProvider } from 'react-native-paper';
import { SafeAreaProvider } from 'react-native-safe-area-context';
import { apolloClient } from './apolloClient';
import { AuthProvider } from './context/AuthContext';
import Scaffold from './components/Scaffold';
import AuthenticatedView from './components/AuthenticatedView';

export default function App() {
  return (
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
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#fff',
  },
});
