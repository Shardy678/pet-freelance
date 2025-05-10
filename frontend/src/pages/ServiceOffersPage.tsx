import { useQuery } from '@tanstack/react-query';
import axios from 'axios';
import { Link } from 'react-router-dom';

type ServiceOffer = {
  id: string;
  title: string;
  description: string;
  price: number;
  currency: string;
  priceType: 'hourly' | 'fixed';
  durationEstimateMin: number;
  isActive: boolean;
};

export default function ServiceOffersPage() {
  const { data, isLoading, isError, error } = useQuery<ServiceOffer[]>({
    queryKey: ['offers'],
    queryFn: () =>
      axios
        .get<ServiceOffer[]>(`${import.meta.env.VITE_API_BASE_URL}/offers`)
        .then((res) => res.data),
    staleTime: 1000 * 60 * 5,
  });

  if (isLoading) {
    return (
      <div className="flex justify-center mt-20 text-gray-500">
        Loading service offers…
      </div>
    );
  }

  if (isError) {
    return (
      <div className="max-w-md mx-auto mt-20 p-4 bg-red-100 text-red-700 rounded">
        <p>Failed to load offers: {(error as any).message}</p>
      </div>
    );
  }

  return (
    <div className="max-w-5xl mx-auto px-4 mt-10">
      <h1 className="text-3xl font-bold mb-6">Service Offers</h1>
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
        {data!.map((offer) => (
          <Link
            key={offer.id}
            to={`/offers/${offer.id}`}
            className="block p-4 border rounded hover:shadow-lg transition"
          >
            <h2 className="text-xl font-semibold mb-2">{offer.title}</h2>
            <p className="text-gray-600 text-sm mb-4 line-clamp-3">
              {offer.description}
            </p>
            <div className="text-sm text-gray-500">
              {offer.currency} {offer.price.toFixed(2)}{' '}
              {offer.priceType === 'hourly' ? '/hr' : ''}
            </div>
            <div className="text-sm text-gray-500">
              ≈ {offer.durationEstimateMin} min
            </div>
          </Link>
        ))}
      </div>
    </div>
  );
}
