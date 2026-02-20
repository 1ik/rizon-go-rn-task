import { ApolloClient, InMemoryCache } from '@apollo/client';

const uri =
  process.env.EXPO_PUBLIC_GRAPHQL_URI || 'http://localhost:8080/graphql';

export const apolloClient = new ApolloClient({
  uri,
  cache: new InMemoryCache(),
});
