import {createContext, Dispatch, PropsWithChildren, SetStateAction, useMemo, useState} from "react";

interface AuthContextProps {
  auth?: Auth | null;
  setAuth?: Dispatch<SetStateAction<Auth | null>>;
}

const AuthContext = createContext<AuthContextProps>({});

export type User = {
  id: string;
  username: string;
  full_name: string;
  role: string;
  password_changed_at: string;
  created_at: string;
}

export type Auth = {
  user: User;
  session_id: string;
  access_token: string;
  refresh_token: string;
  access_token_expires_at: string;
  refresh_token_expires_at: string;
}

export const AuthProvider = ({ children }: PropsWithChildren<{}>) => {
  // Loading auth data from localStorage
  const savedAuthData = localStorage.getItem('authData');
  const initialAuth: Auth | null = savedAuthData ? JSON.parse(savedAuthData) : null;

  // Allow the state to be Auth or null
  const [auth, setAuth] = useState<Auth | null>(initialAuth);

  const value = useMemo(() => ({ auth, setAuth }), [auth, setAuth]);

  return (
    <AuthContext.Provider value={value}>
      {children}
    </AuthContext.Provider>
  );
};

export default AuthContext;
