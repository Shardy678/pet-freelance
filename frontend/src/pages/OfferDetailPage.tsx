import { useParams } from 'react-router-dom';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import axios from 'axios';
import {
    format,
    addDays,
    startOfDay,
    isSameDay,
    parseISO,
} from 'date-fns'
import { useCallback, useMemo, useState } from 'react';

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
    const today = startOfDay(new Date())
    const [selectedDate, setSelectedDate] = useState<Date>(today)
    const [viewMode, setViewMode] = useState<'week' | 'month'>('week')
    const [selectedSlot, setSelectedSlot] = useState<Slot | null>(null)
    const [isModalOpen, setModalOpen] = useState(false)

    const nextSevenDays = useMemo(
        () => Array.from({ length: 7 }).map((_, i) => addDays(today, i)),
        [today]
    )



    // date-picker change for month view
    function onDateChange(e: React.ChangeEvent<HTMLInputElement>) {
        const d = e.target.value
        if (d) setSelectedDate(parseISO(d + 'T00:00:00'))
    }

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

    const getSlotsForDate = useCallback(
        (date: Date) => {
            const ds = format(date, 'yyyy-MM-dd')
            return (
                slots?.filter((s) =>
                    format(parseISO(s.startTime), 'yyyy-MM-dd') === ds && !s.isBooked
                ) || []
            )
        },
        [slots]
    )

    if (loadingOffer) return <p className="p-4 text-center">Loading offer…</p>;
    if (errorOffer || !offer)
        return <p className="p-4 text-center text-red-600">Offer not found.</p>;

    return (
        <>
            <div className="max-w-3xl mx-auto px-4 py-8">
                {/* Offer Header */}
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

                {/* Availability Header */}
                <div className="mb-4 flex justify-between items-center">
                    <h2 className="text-2xl font-semibold">Availability</h2>
                    <div className="space-x-2">
                        <button
                            onClick={() => setViewMode('week')}
                            className={
                                viewMode === 'week'
                                    ? 'px-3 py-1 bg-blue-600 text-white rounded'
                                    : 'px-3 py-1 border rounded'
                            }
                        >
                            Week
                        </button>
                        <button
                            onClick={() => setViewMode('month')}
                            className={
                                viewMode === 'month'
                                    ? 'px-3 py-1 bg-blue-600 text-white rounded'
                                    : 'px-3 py-1 border rounded'
                            }
                        >
                            Month
                        </button>
                    </div>
                </div>

                {/* Week View */}
                {viewMode === 'week' ? (
                    <>
                        <div className="grid grid-cols-7 gap-2 mb-6">
                            {nextSevenDays.map((date) => {
                                const daySlots = getSlotsForDate(date)
                                const isSel = isSameDay(date, selectedDate)
                                return (
                                    <div
                                        key={format(date, 'yyyy-MM-dd')}
                                        onClick={() => setSelectedDate(date)}
                                        className={
                                            'flex flex-col items-center p-2 rounded cursor-pointer ' +
                                            (isSel
                                                ? 'border border-blue-600 bg-blue-100'
                                                : daySlots.length > 0
                                                    ? 'border border-gray-300 hover:border-blue-400'
                                                    : 'border border-gray-300 bg-gray-100 text-gray-400')
                                        }
                                    >
                                        <div className="text-sm font-medium">
                                            {format(date, 'EEE')}
                                        </div>
                                        <div className="text-lg">{format(date, 'd')}</div>
                                        <div className="text-xs mt-1">
                                            {daySlots.length > 0
                                                ? `${daySlots.length} slot${daySlots.length > 1 ? 's' : ''
                                                }`
                                                : 'No slots'}
                                        </div>
                                    </div>
                                )
                            })}
                        </div>
                    </>
                ) : (
                    /* Month View */
                    <div className="mb-6">
                        <input
                            type="date"
                            value={format(selectedDate, 'yyyy-MM-dd')}
                            onChange={(e) => {
                                const d = e.target.value
                                if (d) setSelectedDate(parseISO(d + 'T00:00:00'))
                            }}
                            className="border px-2 py-1 rounded"
                            min={format(today, 'yyyy-MM-dd')}
                        />
                    </div>
                )}

                {/* Slot List for Selected Date */}
                <div>
                    <h3 className="text-lg font-medium mb-4">
                        Available Slots for {format(selectedDate, 'MMMM d, yyyy')}
                    </h3>

                    {loadingSlots && <p>Loading slots…</p>}
                    {errorSlots && (
                        <p className="text-red-600">Failed to load slots.</p>
                    )}
                    {!loadingSlots && getSlotsForDate(selectedDate).length === 0 && (
                        <p className="text-gray-500">
                            No available slots for this day.
                        </p>
                    )}

                    <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 gap-2">
                        {getSlotsForDate(selectedDate).map((slot) => (
                            <button
                                key={slot.id}
                                onClick={() => {
                                    setSelectedSlot(slot)
                                    setModalOpen(true)
                                }}
                                className="border py-4 rounded text-center hover:bg-blue-50 disabled:opacity-50"
                            >
                                {format(parseISO(slot.startTime), 'HH:mm')} –{' '}
                                {format(parseISO(slot.endTime), 'HH:mm')}
                            </button>
                        ))}
                    </div>
                </div>
            </div>
            {/* Booking Confirmation Modal */}
            {isModalOpen && selectedSlot && (
                <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50">
                    <div className="bg-white rounded-lg shadow-lg max-w-md w-full">
                        <div className="px-6 py-4">
                            <h2 className="text-xl font-semibold mb-4">Confirm Booking</h2>
                            <p className="mb-2">
                                <strong>Date:</strong>{' '}
                                {format(parseISO(selectedSlot.startTime), 'MMMM d, yyyy')}
                            </p>
                            <p className="mb-2">
                                <strong>Time:</strong>{' '}
                                {format(parseISO(selectedSlot.startTime), 'HH:mm')} –{' '}
                                {format(parseISO(selectedSlot.endTime), 'HH:mm')}
                            </p>
                            <p className="mb-4">
                                <strong>Price:</strong> {offer.currency}{' '}
                                {offer.price.toFixed(2)}
                            </p>
                            <div className="flex justify-end space-x-2">
                                <button
                                    onClick={() => setModalOpen(false)}
                                    className="px-4 py-2 bg-gray-200 rounded hover:bg-gray-300"
                                >
                                    Cancel
                                </button>
                                <button
                                    onClick={() => {
                                        bookSlot.mutate(selectedSlot.id)
                                        setModalOpen(false)
                                    }}
                                    className="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700 disabled:opacity-50"
                                >
                                </button>
                            </div>
                        </div>
                    </div>
                </div>
            )}
        </>
    );
}
