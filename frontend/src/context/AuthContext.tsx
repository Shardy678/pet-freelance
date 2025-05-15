import { createContext, useState, useEffect, ReactNode } from 'react'
import { useNavigate } from 'react-router-dom'

interface AuthContextType {
  token: string | null
  setToken: (t: string | null) => void
  logout: () => void
}

export const AuthContext = createContext<AuthContextType>({
  token: null,
  setToken: () => { },
  logout: () => { },
})

export function AuthProvider({ children }: { children: ReactNode }) {
  const [token, setTokenState] = useState<string | null>(() => {
    return localStorage.getItem('token')
  })
  const navigate = useNavigate()

  const setToken = (t: string | null) => {
    if (t) localStorage.setItem('token', t)
    else localStorage.removeItem('token')
    setTokenState(t)
  }

  const logout = () => {
    setToken(null)
    navigate('/login')
  }

  return (
    <AuthContext.Provider value={{ token, setToken, logout }}>
      {children}
    </AuthContext.Provider>
  )
}
