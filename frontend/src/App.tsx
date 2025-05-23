import { Navigate, Route, Routes } from 'react-router-dom'
import LoginPage from './pages/LoginPage';
import RegisterPage from './pages/RegisterPage';
import ProfilePage from './pages/ProfilePage';
import ServicesPage from './pages/ServicesPage';
import ServiceDetailPage from './pages/ServiceDetailPage';
import OfferDetailPage from './pages/OfferDetailPage';
import Header from './components/Header';
import BookingsPage from './pages/BookingsPage';
import Dashboard from './pages/Dashboard';

function App() {
  return (
    <>
      <Header />
      <Routes>
        <Route path="/" element={<Navigate to="/services" replace />}></Route>
        <Route path="/login" element={<LoginPage />} />
        <Route path="/register" element={<RegisterPage />} />
        <Route path="/profile" element={<ProfilePage />} />
        <Route path="/services" element={<ServicesPage />} />
        <Route path="/services/:id" element={<ServiceDetailPage />} />
        <Route path="/offers/:id" element={<OfferDetailPage />} />
        <Route path="/bookings" element={<BookingsPage />} />
        <Route path="/dashboard" element={<Dashboard />} />
      </Routes>
    </>

  )
}

export default App
