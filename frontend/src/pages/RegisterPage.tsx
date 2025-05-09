import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import axios from 'axios';
import { useMutation } from '@tanstack/react-query';

type RegisterInput = {
    email: string;
    password: string;
    role: 'owner' | 'freelancer';
  };

export default function RegisterPage() {
  const [email, setEmail]     = useState('');
  const [password, setPassword] = useState('');
  const [role, setRole]       = useState<'owner' | 'freelancer'>('owner');
  const [error, setError]     = useState<string | null>(null);
  const navigate               = useNavigate();

  const registerMutation = useMutation<unknown, any, RegisterInput>({
    mutationFn: data =>
      axios.post(`${import.meta.env.VITE_API_BASE_URL}/auth/register`, data),
    onSuccess: () => navigate('/login'),
    onError: err =>
      setError(err.response?.data?.error || 'Registration failed'),
  });

  const onSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    registerMutation.mutate({ email, password, role });
  };

  return (
    <div className="max-w-md mx-auto mt-20 p-6 bg-white shadow rounded">
      <h1 className="text-2xl font-semibold mb-6">Register</h1>
      {error && <div className="mb-4 text-red-600">{error}</div>}
      <form onSubmit={onSubmit} className="space-y-4">
        <div>
          <label className="block text-sm font-medium mb-1">Email</label>
          <input
            type="email"
            required
            className="w-full border px-3 py-2 rounded"
            value={email}
            onChange={e => setEmail(e.target.value)}
          />
        </div>

        <div>
          <label className="block text-sm font-medium mb-1">Password</label>
          <input
            type="password"
            required
            minLength={8}
            className="w-full border px-3 py-2 rounded"
            value={password}
            onChange={e => setPassword(e.target.value)}
          />
        </div>

        <div>
          <label className="block text-sm font-medium mb-1">Role</label>
          <select
            className="w-full border px-3 py-2 rounded"
            value={role}
            onChange={e => setRole(e.target.value as any)}
          >
            <option value="owner">Pet Owner</option>
            <option value="freelancer">Freelancer</option>
          </select>
        </div>

        <button
          type="submit"
          className="w-full bg-blue-600 text-white py-2 rounded hover:bg-blue-700 transition"
        >
          Create Account
        </button>
      </form>

      <p className="mt-4 text-center text-sm text-gray-600">
        Already have an account?{' '}
        <button
          onClick={() => navigate('/login')}
          className="text-blue-600 hover:underline"
        >
          Log in
        </button>
      </p>
    </div>
  );
}
