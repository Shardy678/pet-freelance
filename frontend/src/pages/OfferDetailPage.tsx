import { useParams } from 'react-router-dom';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import axios from 'axios';
import { format, addDays } from 'date-fns';

type Offer = {
    id: string;
    title: string;
    description: string;
    price: number;
    currency: string;
    priceType: 'hourly' | 'fixed';
    durationEstimateMin: number;
};

type Slot = {
    id: string;
    startTime: string;
    endTime: string;
    isBooked: boolean;
};

export default function OfferDetailPage() {
    const { id } = useParams<{ id: string }>();
    const qc = useQueryClient();

    // 1) Fetch offer details
    const {
        data: offer,
        isLoading: loadingOffer,
        isError: errorOffer,
    } = useQuery<Offer>({
        queryKey: ['offer', id],
        queryFn: () =>
            axios
                .get<Offer>(`${import.meta.env.VITE_API_BASE_URL}/offers/${id}`)
                .then(r => r.data),
        enabled: !!id,
    });

    // compute date range: now → +7 days
    const now = new Date();
    const stripMs = (iso: string) => iso.replace(/\.\d+Z$/, 'Z');

    const fromISO = stripMs(now.toISOString());
    const toISO = stripMs(addDays(now, 7).toISOString());


    // 2) Fetch available slots
    const {
        data: slots,
        isLoading: loadingSlots,
        isError: errorSlots,
    } = useQuery<Slot[]>({
        queryKey: ['slots', id, fromISO, toISO],
        queryFn: () =>
            axios
                .get<Slot[]>(
                    `${import.meta.env.VITE_API_BASE_URL}/offers/${id}/slots`,
                    {
                        params: { available: true, from: fromISO, to: toISO },
                    }
                )
                .then(r => r.data),
        enabled: !!id,
    });

    // 3) Prepare booking mutation
    const bookSlot = useMutation<void, any, string>({
        mutationFn: slotId =>
            axios.post(
                `${import.meta.env.VITE_API_BASE_URL}/bookings`,
                { offer_id: id, slot_id: slotId },
            ),
        onSuccess: () => {
            // refetch slots so that booked ones disappear
            qc.invalidateQueries({
                queryKey: ['slots', id, fromISO, toISO],
                exact: true,
            });
        },
    });

    if (loadingOffer) return <p className="p-4 text-center">Loading offer…</p>;
    if (errorOffer || !offer)
        return <p className="p-4 text-center text-red-600">Offer not found.</p>;

    return (
        <div className="max-w-3xl mx-auto px-4 py-8">
            {/* Offer Info */}
            <div className="mb-8">
                <h1 className="text-3xl font-bold mb-2">{offer.title}</h1>
                <p className="text-gray-700 mb-4">{offer.description}</p>
                <div className="text-lg font-semibold">
                    {offer.currency} {offer.price.toFixed(2)}{' '}
                    {offer.priceType === 'hourly' ? '/hr' : ''}
                </div>
                <p className="text-sm text-gray-500">
                    ≈ {offer.durationEstimateMin} min
                </p>
            </div>

            {/* Availability */}
            <div>
                <h2 className="text-2xl font-semibold mb-4">Available Slots (Next 7 days)</h2>

                {loadingSlots && <p>Loading slots…</p>}
                {errorSlots && <p className="text-red-600">Failed to load slots.</p>}
                {!loadingSlots && slots && slots.length === 0 && (
                    <p className="text-gray-500">No slots available.</p>
                )}

                <ul className="space-y-3">
                    {slots?.map(slot => (
                        <li
                            key={slot.id}
                            className="flex justify-between items-center p-3 border rounded"
                        >
                            <span>
                                {format(new Date(slot.startTime), 'PPpp')} –{' '}
                                {format(new Date(slot.endTime), 'pp')}
                            </span>
                            <button
                                onClick={() => bookSlot.mutate(slot.id)}
                                disabled={bookSlot.status === "pending"}
                                className="px-3 py-1 bg-blue-600 text-white rounded hover:bg-blue-700 disabled:opacity-50"
                            >
                                {bookSlot.status === "pending" ? 'Booking…' : 'Book'}
                            </button>
                        </li>
                    ))}
                </ul>
            </div>
        </div>
    );
}
