import { gql } from '@apollo/client';
import * as Apollo from '@apollo/client';
export type Maybe<T> = T | null;
export type InputMaybe<T> = Maybe<T>;
export type Exact<T extends { [key: string]: unknown }> = { [K in keyof T]: T[K] };
export type MakeOptional<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]?: Maybe<T[SubKey]> };
export type MakeMaybe<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]: Maybe<T[SubKey]> };
export type MakeEmpty<T extends { [key: string]: unknown }, K extends keyof T> = { [_ in K]?: never };
export type Incremental<T> = T | { [P in keyof T]?: P extends ' $fragmentName' | '__typename' ? T[P] : never };
const defaultOptions = {} as const;
/** All built-in and custom scalars, mapped to their actual values */
export type Scalars = {
  ID: { input: string; output: string; }
  String: { input: string; output: string; }
  Boolean: { input: boolean; output: boolean; }
  Int: { input: number; output: number; }
  Float: { input: number; output: number; }
  Time: { input: any; output: any; }
};

export type Feedback = {
  __typename?: 'Feedback';
  content: Scalars['String']['output'];
  createdAt: Scalars['Time']['output'];
  deviceId: Scalars['String']['output'];
  id: Scalars['ID']['output'];
  updatedAt: Scalars['Time']['output'];
  userId: Scalars['ID']['output'];
};

export type Mutation = {
  __typename?: 'Mutation';
  loginWithEmailAndSecret: Scalars['String']['output'];
  requestEmailAuthLink: Scalars['Boolean']['output'];
  submitFeedback: Scalars['Boolean']['output'];
};


export type MutationLoginWithEmailAndSecretArgs = {
  email: Scalars['String']['input'];
  secret: Scalars['String']['input'];
};


export type MutationRequestEmailAuthLinkArgs = {
  email: Scalars['String']['input'];
};


export type MutationSubmitFeedbackArgs = {
  content: Scalars['String']['input'];
  deviceId: Scalars['String']['input'];
};

export type Query = {
  __typename?: 'Query';
  getUserFeedbackOnDevice?: Maybe<Feedback>;
  hello: Scalars['String']['output'];
  me: User;
};


export type QueryGetUserFeedbackOnDeviceArgs = {
  deviceId: Scalars['String']['input'];
};

export type User = {
  __typename?: 'User';
  createdAt: Scalars['Time']['output'];
  email: Scalars['String']['output'];
  id: Scalars['ID']['output'];
  updatedAt: Scalars['Time']['output'];
};

export type GetUserFeedbackOnDeviceQueryVariables = Exact<{
  deviceId: Scalars['String']['input'];
}>;


export type GetUserFeedbackOnDeviceQuery = { __typename?: 'Query', getUserFeedbackOnDevice?: { __typename?: 'Feedback', id: string, userId: string, deviceId: string, content: string, createdAt: any, updatedAt: any } | null };

export type HelloQueryVariables = Exact<{ [key: string]: never; }>;


export type HelloQuery = { __typename?: 'Query', hello: string };

export type LoginWithEmailAndSecretMutationVariables = Exact<{
  email: Scalars['String']['input'];
  secret: Scalars['String']['input'];
}>;


export type LoginWithEmailAndSecretMutation = { __typename?: 'Mutation', loginWithEmailAndSecret: string };

export type MeQueryVariables = Exact<{ [key: string]: never; }>;


export type MeQuery = { __typename?: 'Query', me: { __typename?: 'User', id: string, email: string, createdAt: any, updatedAt: any } };

export type RequestEmailAuthLinkMutationVariables = Exact<{
  email: Scalars['String']['input'];
}>;


export type RequestEmailAuthLinkMutation = { __typename?: 'Mutation', requestEmailAuthLink: boolean };

export type SubmitFeedbackMutationVariables = Exact<{
  deviceId: Scalars['String']['input'];
  content: Scalars['String']['input'];
}>;


export type SubmitFeedbackMutation = { __typename?: 'Mutation', submitFeedback: boolean };


export const GetUserFeedbackOnDeviceDocument = gql`
    query GetUserFeedbackOnDevice($deviceId: String!) {
  getUserFeedbackOnDevice(deviceId: $deviceId) {
    id
    userId
    deviceId
    content
    createdAt
    updatedAt
  }
}
    `;

/**
 * __useGetUserFeedbackOnDeviceQuery__
 *
 * To run a query within a React component, call `useGetUserFeedbackOnDeviceQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetUserFeedbackOnDeviceQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetUserFeedbackOnDeviceQuery({
 *   variables: {
 *      deviceId: // value for 'deviceId'
 *   },
 * });
 */
export function useGetUserFeedbackOnDeviceQuery(baseOptions: Apollo.QueryHookOptions<GetUserFeedbackOnDeviceQuery, GetUserFeedbackOnDeviceQueryVariables> & ({ variables: GetUserFeedbackOnDeviceQueryVariables; skip?: boolean; } | { skip: boolean; }) ) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<GetUserFeedbackOnDeviceQuery, GetUserFeedbackOnDeviceQueryVariables>(GetUserFeedbackOnDeviceDocument, options);
      }
export function useGetUserFeedbackOnDeviceLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<GetUserFeedbackOnDeviceQuery, GetUserFeedbackOnDeviceQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<GetUserFeedbackOnDeviceQuery, GetUserFeedbackOnDeviceQueryVariables>(GetUserFeedbackOnDeviceDocument, options);
        }
// @ts-ignore
export function useGetUserFeedbackOnDeviceSuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<GetUserFeedbackOnDeviceQuery, GetUserFeedbackOnDeviceQueryVariables>): Apollo.UseSuspenseQueryResult<GetUserFeedbackOnDeviceQuery, GetUserFeedbackOnDeviceQueryVariables>;
export function useGetUserFeedbackOnDeviceSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetUserFeedbackOnDeviceQuery, GetUserFeedbackOnDeviceQueryVariables>): Apollo.UseSuspenseQueryResult<GetUserFeedbackOnDeviceQuery | undefined, GetUserFeedbackOnDeviceQueryVariables>;
export function useGetUserFeedbackOnDeviceSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<GetUserFeedbackOnDeviceQuery, GetUserFeedbackOnDeviceQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<GetUserFeedbackOnDeviceQuery, GetUserFeedbackOnDeviceQueryVariables>(GetUserFeedbackOnDeviceDocument, options);
        }
export type GetUserFeedbackOnDeviceQueryHookResult = ReturnType<typeof useGetUserFeedbackOnDeviceQuery>;
export type GetUserFeedbackOnDeviceLazyQueryHookResult = ReturnType<typeof useGetUserFeedbackOnDeviceLazyQuery>;
export type GetUserFeedbackOnDeviceSuspenseQueryHookResult = ReturnType<typeof useGetUserFeedbackOnDeviceSuspenseQuery>;
export type GetUserFeedbackOnDeviceQueryResult = Apollo.QueryResult<GetUserFeedbackOnDeviceQuery, GetUserFeedbackOnDeviceQueryVariables>;
export const HelloDocument = gql`
    query Hello {
  hello
}
    `;

/**
 * __useHelloQuery__
 *
 * To run a query within a React component, call `useHelloQuery` and pass it any options that fit your needs.
 * When your component renders, `useHelloQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useHelloQuery({
 *   variables: {
 *   },
 * });
 */
export function useHelloQuery(baseOptions?: Apollo.QueryHookOptions<HelloQuery, HelloQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<HelloQuery, HelloQueryVariables>(HelloDocument, options);
      }
export function useHelloLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<HelloQuery, HelloQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<HelloQuery, HelloQueryVariables>(HelloDocument, options);
        }
// @ts-ignore
export function useHelloSuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<HelloQuery, HelloQueryVariables>): Apollo.UseSuspenseQueryResult<HelloQuery, HelloQueryVariables>;
export function useHelloSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<HelloQuery, HelloQueryVariables>): Apollo.UseSuspenseQueryResult<HelloQuery | undefined, HelloQueryVariables>;
export function useHelloSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<HelloQuery, HelloQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<HelloQuery, HelloQueryVariables>(HelloDocument, options);
        }
export type HelloQueryHookResult = ReturnType<typeof useHelloQuery>;
export type HelloLazyQueryHookResult = ReturnType<typeof useHelloLazyQuery>;
export type HelloSuspenseQueryHookResult = ReturnType<typeof useHelloSuspenseQuery>;
export type HelloQueryResult = Apollo.QueryResult<HelloQuery, HelloQueryVariables>;
export const LoginWithEmailAndSecretDocument = gql`
    mutation LoginWithEmailAndSecret($email: String!, $secret: String!) {
  loginWithEmailAndSecret(email: $email, secret: $secret)
}
    `;
export type LoginWithEmailAndSecretMutationFn = Apollo.MutationFunction<LoginWithEmailAndSecretMutation, LoginWithEmailAndSecretMutationVariables>;

/**
 * __useLoginWithEmailAndSecretMutation__
 *
 * To run a mutation, you first call `useLoginWithEmailAndSecretMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useLoginWithEmailAndSecretMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [loginWithEmailAndSecretMutation, { data, loading, error }] = useLoginWithEmailAndSecretMutation({
 *   variables: {
 *      email: // value for 'email'
 *      secret: // value for 'secret'
 *   },
 * });
 */
export function useLoginWithEmailAndSecretMutation(baseOptions?: Apollo.MutationHookOptions<LoginWithEmailAndSecretMutation, LoginWithEmailAndSecretMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<LoginWithEmailAndSecretMutation, LoginWithEmailAndSecretMutationVariables>(LoginWithEmailAndSecretDocument, options);
      }
export type LoginWithEmailAndSecretMutationHookResult = ReturnType<typeof useLoginWithEmailAndSecretMutation>;
export type LoginWithEmailAndSecretMutationResult = Apollo.MutationResult<LoginWithEmailAndSecretMutation>;
export type LoginWithEmailAndSecretMutationOptions = Apollo.BaseMutationOptions<LoginWithEmailAndSecretMutation, LoginWithEmailAndSecretMutationVariables>;
export const MeDocument = gql`
    query Me {
  me {
    id
    email
    createdAt
    updatedAt
  }
}
    `;

/**
 * __useMeQuery__
 *
 * To run a query within a React component, call `useMeQuery` and pass it any options that fit your needs.
 * When your component renders, `useMeQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useMeQuery({
 *   variables: {
 *   },
 * });
 */
export function useMeQuery(baseOptions?: Apollo.QueryHookOptions<MeQuery, MeQueryVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useQuery<MeQuery, MeQueryVariables>(MeDocument, options);
      }
export function useMeLazyQuery(baseOptions?: Apollo.LazyQueryHookOptions<MeQuery, MeQueryVariables>) {
          const options = {...defaultOptions, ...baseOptions}
          return Apollo.useLazyQuery<MeQuery, MeQueryVariables>(MeDocument, options);
        }
// @ts-ignore
export function useMeSuspenseQuery(baseOptions?: Apollo.SuspenseQueryHookOptions<MeQuery, MeQueryVariables>): Apollo.UseSuspenseQueryResult<MeQuery, MeQueryVariables>;
export function useMeSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<MeQuery, MeQueryVariables>): Apollo.UseSuspenseQueryResult<MeQuery | undefined, MeQueryVariables>;
export function useMeSuspenseQuery(baseOptions?: Apollo.SkipToken | Apollo.SuspenseQueryHookOptions<MeQuery, MeQueryVariables>) {
          const options = baseOptions === Apollo.skipToken ? baseOptions : {...defaultOptions, ...baseOptions}
          return Apollo.useSuspenseQuery<MeQuery, MeQueryVariables>(MeDocument, options);
        }
export type MeQueryHookResult = ReturnType<typeof useMeQuery>;
export type MeLazyQueryHookResult = ReturnType<typeof useMeLazyQuery>;
export type MeSuspenseQueryHookResult = ReturnType<typeof useMeSuspenseQuery>;
export type MeQueryResult = Apollo.QueryResult<MeQuery, MeQueryVariables>;
export const RequestEmailAuthLinkDocument = gql`
    mutation RequestEmailAuthLink($email: String!) {
  requestEmailAuthLink(email: $email)
}
    `;
export type RequestEmailAuthLinkMutationFn = Apollo.MutationFunction<RequestEmailAuthLinkMutation, RequestEmailAuthLinkMutationVariables>;

/**
 * __useRequestEmailAuthLinkMutation__
 *
 * To run a mutation, you first call `useRequestEmailAuthLinkMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useRequestEmailAuthLinkMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [requestEmailAuthLinkMutation, { data, loading, error }] = useRequestEmailAuthLinkMutation({
 *   variables: {
 *      email: // value for 'email'
 *   },
 * });
 */
export function useRequestEmailAuthLinkMutation(baseOptions?: Apollo.MutationHookOptions<RequestEmailAuthLinkMutation, RequestEmailAuthLinkMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<RequestEmailAuthLinkMutation, RequestEmailAuthLinkMutationVariables>(RequestEmailAuthLinkDocument, options);
      }
export type RequestEmailAuthLinkMutationHookResult = ReturnType<typeof useRequestEmailAuthLinkMutation>;
export type RequestEmailAuthLinkMutationResult = Apollo.MutationResult<RequestEmailAuthLinkMutation>;
export type RequestEmailAuthLinkMutationOptions = Apollo.BaseMutationOptions<RequestEmailAuthLinkMutation, RequestEmailAuthLinkMutationVariables>;
export const SubmitFeedbackDocument = gql`
    mutation SubmitFeedback($deviceId: String!, $content: String!) {
  submitFeedback(deviceId: $deviceId, content: $content)
}
    `;
export type SubmitFeedbackMutationFn = Apollo.MutationFunction<SubmitFeedbackMutation, SubmitFeedbackMutationVariables>;

/**
 * __useSubmitFeedbackMutation__
 *
 * To run a mutation, you first call `useSubmitFeedbackMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useSubmitFeedbackMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [submitFeedbackMutation, { data, loading, error }] = useSubmitFeedbackMutation({
 *   variables: {
 *      deviceId: // value for 'deviceId'
 *      content: // value for 'content'
 *   },
 * });
 */
export function useSubmitFeedbackMutation(baseOptions?: Apollo.MutationHookOptions<SubmitFeedbackMutation, SubmitFeedbackMutationVariables>) {
        const options = {...defaultOptions, ...baseOptions}
        return Apollo.useMutation<SubmitFeedbackMutation, SubmitFeedbackMutationVariables>(SubmitFeedbackDocument, options);
      }
export type SubmitFeedbackMutationHookResult = ReturnType<typeof useSubmitFeedbackMutation>;
export type SubmitFeedbackMutationResult = Apollo.MutationResult<SubmitFeedbackMutation>;
export type SubmitFeedbackMutationOptions = Apollo.BaseMutationOptions<SubmitFeedbackMutation, SubmitFeedbackMutationVariables>;