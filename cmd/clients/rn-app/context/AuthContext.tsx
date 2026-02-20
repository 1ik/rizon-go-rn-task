import AsyncStorage from '@react-native-async-storage/async-storage';
import * as Linking from 'expo-linking';
import React, {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useState,
} from 'react';
import {
  useLoginWithEmailAndSecretMutation,
  useMeLazyQuery,
} from '../graphql/generated/graphql';

export type User = { email: string };

const TOKEN_STORAGE_KEY = '@rizon:auth_token';

type AuthContextValue = {
  user: User | null;
  token: string | null;
  isLoading: boolean;
  authError: string | null;
  login: (userData: User) => void;
  logout: () => void;
  loginWithEmailAndSecret: (email: string, secret: string) => void;
  loadStoredToken: () => Promise<void>;
  handleDeepLink: (url: string | null) => void;
  clearAuthError: () => void;
};

const AuthContext = createContext<AuthContextValue | null>(null);

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [token, setToken] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [authError, setAuthError] = useState<string | null>(null);

  const [getMe] = useMeLazyQuery({
    fetchPolicy: 'network-only',
  });

  const [loginWithEmailAndSecretMutation] = useLoginWithEmailAndSecretMutation();

  const login = (userData: User) => {
    setUser(userData);
    setAuthError(null);
  };

  const logout = async () => {
    setUser(null);
    setToken(null);
    setAuthError(null);
    await AsyncStorage.removeItem(TOKEN_STORAGE_KEY);
  };

  const clearAuthError = useCallback(() => {
    setAuthError(null);
    setIsLoading(false);
  }, []);

  const fetchAndSetUser = useCallback(async (): Promise<boolean> => {
    try {
      const meResult = await getMe();
      const userData = meResult.data?.me;
      
      if (userData) {
        setUser({ email: userData.email });
        return true;
      }
      return false;
    } catch (error) {
      console.error('Failed to fetch user data:', error);
      return false;
    }
  }, [getMe]);

  const loginWithEmailAndSecret = useCallback(
    (email: string, secret: string): void => {
      setIsLoading(true);
      setAuthError(null);

      // Fire and forget - handle async internally
      (async () => {
        try {
          // Call LoginWithEmailAndSecret mutation
          const result = await loginWithEmailAndSecretMutation({
            variables: { email, secret },
          });
          debugger;

          const authToken = result.data?.loginWithEmailAndSecret;
          if (!authToken) {
            throw new Error('No token received from authentication');
          }

          // Store token
          await AsyncStorage.setItem(TOKEN_STORAGE_KEY, authToken);
          setToken(authToken);

          // Fetch and set user data
          const success = await fetchAndSetUser();
          if (!success) {
            throw new Error('Failed to fetch user data');
          }
        } catch (error: any) {
          debugger;
          // Clear token on error
          await AsyncStorage.removeItem(TOKEN_STORAGE_KEY);
          setToken(null);
          
          // Set error message
          const errorMessage = error?.message || 'Authentication failed';
          setAuthError(errorMessage);
        } finally {
          debugger;
          setIsLoading(false);
        }
      })();
    },
    [fetchAndSetUser, loginWithEmailAndSecretMutation]
  );

  const handleDeepLink = useCallback(
    (url: string | null): void => {
      if (!url) return;

      try {
        const parsed = Linking.parse(url);
        // Check if it's our deep link scheme
        // URL format: rizon://email-auth?email=...&secret=...
        // The 'email-auth' can be parsed as hostname or path depending on platform
        const isEmailAuthLink =
          parsed.scheme === 'rizon' &&
          (parsed.hostname === 'email-auth' ||
            parsed.path === 'email-auth' ||
            parsed.path === '/email-auth' ||
            url.includes('email-auth'));
        
        if (isEmailAuthLink) {
          const email = parsed.queryParams?.email as string | undefined;
          const secret = parsed.queryParams?.secret as string | undefined;

          if (!email || !secret) {
            console.error('Invalid deep link: missing email or secret');
            return;
          }

          loginWithEmailAndSecret(email, secret);
        }
      } catch (error) {
        console.error('Error parsing deep link:', error);
        setIsLoading(false);
      }
    },
    [loginWithEmailAndSecret]
  );

  const loadStoredToken = useCallback(async (): Promise<void> => {
    try {
      setIsLoading(true);
      const storedToken = await AsyncStorage.getItem(TOKEN_STORAGE_KEY);

      if (!storedToken) {
        setIsLoading(false);
        return;
      }

      setToken(storedToken);

      // Try to fetch user data to validate token
      const success = await fetchAndSetUser();
      if (!success) {
        // Token is invalid or expired, clear it
        await AsyncStorage.removeItem(TOKEN_STORAGE_KEY);
        setToken(null);
      }
    } catch (error) {
      // Error reading from storage, clear token
      await AsyncStorage.removeItem(TOKEN_STORAGE_KEY);
      setToken(null);
    } finally {
      setIsLoading(false);
    }
  }, [fetchAndSetUser]);

  // Load stored token on mount
  useEffect(() => {
    loadStoredToken();
  }, [loadStoredToken]);

  return (
    <AuthContext.Provider
      value={{
        user,
        token,
        isLoading,
        authError,
        login,
        logout,
        loginWithEmailAndSecret,
        loadStoredToken,
        handleDeepLink,
        clearAuthError,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth(): AuthContextValue {
  const ctx = useContext(AuthContext);
  if (!ctx) throw new Error('useAuth must be used within AuthProvider');
  return ctx;
}
