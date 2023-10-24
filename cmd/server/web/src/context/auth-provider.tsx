import { createContext, PropsWithChildren, useState, ReactNode } from "react";

interface AuthContextProps {
  auth?: Auth;
  setAuth?: React.Dispatch<React.SetStateAction<Auth | undefined>>;
}

const AuthContext = createContext<AuthContextProps>({});

type User = {
  id: string;
  username: string;
  full_name: string;
  role: string;
  password_changed_at: string;
  created_at: string;
}

type Auth = {
  user: User;
  access_token: string;
  refresh_token: string;
  access_token_expires_at: string;
  refresh_token_expires_at: string;
}

export const AuthProvider = ({ children }: PropsWithChildren<{}>): ReactNode => {
  const [auth, setAuth] = useState<Auth>();

  return (
    <AuthContext.Provider value={{ auth, setAuth }}>
      {children}
    </AuthContext.Provider>
  );
};

export default AuthContext;
