import { useState, useContext } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import axios from 'axios';
import { AuthContext } from '../context/AuthContext';

export default function LoginPage() {
  const [email, setEmail] = useState('');
  const [pw, setPw]     = useState('');
  const { setToken }    = useContext(AuthContext);
  const nav             = useNavigate();

  async function submit(e: React.FormEvent) {
    e.preventDefault();
    const res = await axios.post(`${import.meta.env.VITE_API_BASE_URL}/auth/login`, {
      email, password: pw
    });
    setToken(res.data.token);
    nav('/profile');
  }

  return (
    <form onSubmit={submit} className="max-w-md mx-auto mt-20 p-4">
      <h1 className="text-2xl mb-4">Login</h1>
      <input
        className="w-full p-2 border mb-2"
        type="email" placeholder="Email"
        value={email} onChange={e => setEmail(e.target.value)}
      />
      <input
        className="w-full p-2 border mb-4"
        type="password" placeholder="Password"
        value={pw} onChange={e => setPw(e.target.value)}
      />
      <button className="px-4 py-2 bg-blue-600 text-white rounded">Log In</button>
      <p className="mt-4 text-center text-sm text-gray-600">
        Don’t have an account?{' '}
        <Link to="/register" className="text-blue-600 hover:underline">
          Register here
        </Link>
      </p>
    </form>
  );
}
