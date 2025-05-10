// src/pages/ServiceDetailPage.tsx
import { useParams, Link } from 'react-router-dom';
import { useQuery } from '@tanstack/react-query';
import axios from 'axios';

type Service = {
  id: string;
  name: string;
  description?: string;
  basePrice: number;
  defaultDurationMin: number;
};

type ServiceOffer = {
  id: string;
  title: string;
  description: string;
  price: number;
  currency: string;
  priceType: 'hourly' | 'fixed';
  durationEstimateMin: number;
};

export default function ServiceDetailPage() {
  const { id } = useParams<{ id: string }>();

  // 1) Fetch the service itself
  const {
    data: service,
    isLoading: svcLoading,
    isError: svcError,
  } = useQuery<Service>({
    queryKey: ['service', id],
    queryFn: () =>
      axios
        .get<Service>(`${import.meta.env.VITE_API_BASE_URL}/services/${id}`)
        .then((res) => res.data),
    enabled: !!id,
  });

  // 2) Fetch offers filtered by service_id
  const {
    data: offers,
    isLoading: offersLoading,
    isError: offersError,
  } = useQuery<ServiceOffer[]>({
    queryKey: ['offersByService', id],
    queryFn: () =>
      axios
        .get<ServiceOffer[]>(
          `${import.meta.env.VITE_API_BASE_URL}/offers`,
          { params: { service_id: id } }
        )
        .then((res) => res.data),
    enabled: !!id,
  });

  // 3) Loading / error / empty guards
  if (svcLoading) return <p className="p-4 text-center">Loading service…</p>;
  if (svcError || !service)
    return <p className="p-4 text-center text-red-600">Service not found.</p>;

  return (
    <div className="max-w-4xl mx-auto mt-10 px-4">
      {/* Service Header */}
      <header className="mb-8">
        <h1 className="text-3xl font-bold">{service.name}</h1>
        {service.description && (
          <p className="mt-2 text-gray-600">{service.description}</p>
        )}
        <p className="mt-4 text-sm text-gray-500">
          From ${service.basePrice.toFixed(2)} &mdash;{' '}
          {service.defaultDurationMin} min
        </p>
      </header>

      {/* Offers Section */}
      <section>
        <h2 className="text-2xl font-semibold mb-4">Offers</h2>

        {offersLoading && <p>Loading offers…</p>}
        {offersError && (
          <p className="text-red-600">Failed to load offers.</p>
        )}
        {!offersLoading && offers && offers.length === 0 && (
          <p className="text-gray-500">No offers available yet.</p>
        )}

        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
          {offers?.map((offer) => (
            <Link
              key={offer.id}
              to={`/offers/${offer.id}`}
              className="block p-4 border rounded hover:shadow-lg transition"
            >
              <h3 className="text-xl font-semibold mb-1">{offer.title}</h3>
              <p className="text-gray-600 text-sm mb-2 line-clamp-3">
                {offer.description}
              </p>
              <div className="text-sm text-gray-500">
                {offer.currency} {offer.price.toFixed(2)}{' '}
                {offer.priceType === 'hourly' ? '/hr' : ''}
              </div>
              <div className="text-sm text-gray-500">
                ≈ {offer.durationEstimateMin} min
              </div>
            </Link>
          ))}
        </div>
      </section>
    </div>
  );
}
