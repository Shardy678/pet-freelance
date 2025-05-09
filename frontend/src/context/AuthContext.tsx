import { createContext, useState, useEffect } from 'react';
import type { ReactNode } from 'react';

interface AuthCtx { token: string; setToken: (t: string) => void; }
export const AuthContext = createContext<AuthCtx>({ token: '', setToken: () => {} });

export function AuthProvider({ children }: { children: ReactNode }) {
  const [token, setToken] = useState(() => localStorage.getItem('token') || '');
  useEffect(() => { localStorage.setItem('token', token); }, [token]);
  return (
    <AuthContext.Provider value={{ token, setToken }}>
      {children}
    </AuthContext.Provider>
  );
}
