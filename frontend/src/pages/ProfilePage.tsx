// src/pages/ProfilePage.tsx
import { useContext, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import axios from 'axios';
import { useQuery } from '@tanstack/react-query';
import { AuthContext } from '../context/AuthContext';

type Profile = {
  id: string;
  email: string;
  role: string;
};

export default function ProfilePage() {
  const { token, setToken } = useContext(AuthContext);
  const navigate = useNavigate();

  useEffect(() => {
    if (!token) {
      navigate('/login', { replace: true });
    }
  }, [token, navigate]);

   const { data, isLoading, error } = useQuery<Profile>({
       queryKey: ['profile'],
       queryFn: () =>
         axios
           .get<Profile>(`${import.meta.env.VITE_API_BASE_URL}/profile/me`, {
             headers: { Authorization: `Bearer ${token}` },
           })
           .then(res => res.data),
       enabled: !!token,
       retry: false,
     });

  if (isLoading) {
    return (
      <div className="flex justify-center items-center h-full">
        <p className="text-gray-500">Loading profile…</p>
      </div>
    );
  }

  if (error || !data) {
    return (
      <div className="max-w-md mx-auto mt-20 p-4 bg-red-100 text-red-700 rounded">
        <p>Failed to load profile. Please try logging in again.</p>
        <button
          onClick={() => {
            setToken('');
            navigate('/login', { replace: true });
          }}
          className="mt-4 px-4 py-2 bg-red-600 text-white rounded hover:bg-red-700"
        >
          Go to Login
        </button>
      </div>
    );
  }

  const handleLogout = () => {
    setToken('');
    navigate('/login', { replace: true });
  };

  return (
    <div className="max-w-lg mx-auto mt-20 p-6 bg-white shadow rounded">
      <h1 className="text-2xl font-semibold mb-4">My Profile</h1>
      <div className="space-y-2">
        <p>
          <span className="font-medium">ID:</span> {data.id}
        </p>
        <p>
          <span className="font-medium">Email:</span> {data.email}
        </p>
        <p>
          <span className="font-medium">Role:</span> {data.role}
        </p>
      </div>
      <button
        onClick={handleLogout}
        className="mt-6 px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700"
      >
        Log Out
      </button>
    </div>
  );
}