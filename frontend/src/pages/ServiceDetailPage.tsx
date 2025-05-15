// src/pages/ServiceDetailPage.tsx
import { useParams, Link } from 'react-router-dom';
import { useQuery, useQueryClient } from '@tanstack/react-query';
import axios from 'axios';
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from '@/components/ui/card';

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
    const qc = useQueryClient();

    const {
        data: service,
        isLoading: svcLoading,
        isError: svcError,
    } = useQuery<Service>({
        queryKey: ['service', id],
        queryFn: () =>
            axios
                .get<Service>(`${import.meta.env.VITE_API_BASE_URL}/services/${id}`)
                .then((r) => r.data),
        initialData: () =>
            (qc.getQueryData<Service[]>(['services']) || []).find((s) => s.id === id),
        enabled: !!id,
        staleTime: Infinity,         
        refetchOnMount: false,        
    });


    // 2) Fetch the offers
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
                .then((r) => r.data),
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
          From ${service.basePrice.toFixed(2)} &mdash; {service.defaultDurationMin} min
        </p>
      </header>

      {/* Offers Section */}
      <section>
        <h2 className="text-2xl font-semibold mb-4">Offers</h2>

        {offersLoading && <p>Loading offers…</p>}
        {offersError && <p className="text-red-600">Failed to load offers.</p>}
        {!offersLoading && offers && offers.length === 0 && (
          <p className="text-gray-500">No offers available yet.</p>
        )}

        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
          {offers?.map((offer) => (
            <Card key={offer.id} className="hover:shadow-lg transition">
              <Link to={`/offers/${offer.id}`} className="block">
                <CardHeader>
                  <CardTitle>{offer.title}</CardTitle>
                  <CardDescription className="line-clamp-3">
                    {offer.description}
                  </CardDescription>
                </CardHeader>
                <CardContent>
                  <div className="flex items-center justify-between text-sm text-gray-500">
                    <span>
                      {offer.currency} {offer.price.toFixed(2)}
                      {offer.priceType === 'hourly' ? '/hr' : ''}
                    </span>
                    <span>≈ {offer.durationEstimateMin} min</span>
                  </div>
                </CardContent>
              </Link>
            </Card>
          ))}
        </div>
      </section>
    </div>
    );
}
