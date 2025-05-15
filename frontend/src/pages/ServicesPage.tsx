import { useContext } from 'react';
import { Link } from 'react-router-dom';
import { useQuery } from '@tanstack/react-query';
import axios from 'axios';
import { AuthContext } from '../context/AuthContext';
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from '@/components/ui/card';

type Service = {
  id: string;
  name: string;
  description?: string;
  basePrice: number;
  defaultDurationMin: number;
};

export default function ServicesPage() {
  const { token } = useContext(AuthContext);

  const { data, isLoading, isError, error } = useQuery<Service[]>({
    queryKey: ['services'],
    queryFn: () =>
      axios
        .get<Service[]>(`${import.meta.env.VITE_API_BASE_URL}/services`, {
          headers: token ? { Authorization: `Bearer ${token}` } : undefined,
        })
        .then(res => res.data),
    staleTime: 1000 * 60 * 5, // 5m
  });

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-full py-20">
        <p className="text-gray-500">Loading services…</p>
      </div>
    );
  }

  if (isError) {
    return (
      <div className="max-w-md mx-auto mt-20 p-4 bg-red-100 text-red-700 rounded">
        <p>Error loading services: {(error as any).message}</p>
      </div>
    );
  }

  return (
    <div className="max-w-4xl mx-auto mt-10 px-4">
      <h1 className="text-3xl font-bold mb-6">Available Services</h1>
      <div className="grid gap-6 grid-cols-1 sm:grid-cols-2 lg:grid-cols-3">
        {data!.map((svc) => (
          <Card key={svc.id} className="hover:shadow-lg transition">
            <Link to={`/services/${svc.id}`} className="block">
              <CardHeader>
                <CardTitle>{svc.name}</CardTitle>
                {svc.description && (
                  <CardDescription className="line-clamp-3">
                    {svc.description}
                  </CardDescription>
                )}
              </CardHeader>
              <CardContent>
                <p className="text-sm text-gray-500">
                  From ${svc.basePrice.toFixed(2)} &mdash; {svc.defaultDurationMin} min
                </p>
              </CardContent>
            </Link>
          </Card>
        ))}
      </div>
    </div>
  );
}
