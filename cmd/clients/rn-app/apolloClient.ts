import { ApolloClient, InMemoryCache, createHttpLink } from '@apollo/client';
import { setContext } from '@apollo/client/link/context';
import AsyncStorage from '@react-native-async-storage/async-storage';

const uri =
  process.env.EXPO_PUBLIC_GRAPHQL_URI || 'http://localhost:8080/graphql';

const httpLink = createHttpLink({
  uri,
});

const authLink = setContext(async (_, { headers }) => {
  const token = await AsyncStorage.getItem('@rizon:auth_token');
  return {
    headers: {
      ...headers,
      authorization: token ? `Bearer ${token}` : '',
    },
  };
});

export const apolloClient = new ApolloClient({
  link: authLink.concat(httpLink),
  cache: new InMemoryCache(),
});
