import { Navigate, Route, Routes } from 'react-router-dom'
import { useContext } from 'react';
import type { JSX } from 'react';
import { AuthContext } from './context/AuthContext';
import LoginPage from './pages/LoginPage';
import RegisterPage from './pages/RegisterPage';
import ProfilePage from './pages/ProfilePage';
import ServicesPage from './pages/ServicesPage';
import ServiceOffersPage from './pages/ServiceOffersPage';
import ServiceDetailPage from './pages/ServiceDetailPage';

function RequireAuth({ children }: { children: JSX.Element }) {
  const { token } = useContext(AuthContext);
  return token ? children : <Navigate to="/login" replace />;
}


function App() {
  return (
    <Routes>
      <Route path="/" element={<Navigate to="/services" replace />}></Route>
      <Route path="/login" element={<LoginPage />} />
      <Route path="/register" element={<RegisterPage />} />
      <Route path="/profile" element={<ProfilePage />} />
      <Route path="/services" element={<ServicesPage />} />
      <Route path="/services/:id" element={<ServiceDetailPage />} />


    </Routes>
  )
}

export default App
