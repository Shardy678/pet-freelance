import { useContext, useState } from 'react';
import { useParams } from 'react-router-dom';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import axios from 'axios';
import {
  format,
  addDays,
  startOfDay,
  isSameDay,
  parseISO,
} from 'date-fns';

import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Calendar } from '@/components/ui/calendar';
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogFooter,
} from '@/components/ui/dialog';
import { Tabs, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { Clock, DollarSign, User, Star, CheckCircle } from 'lucide-react';
import { AuthContext } from '@/context/AuthContext';

type Service = {
  id: string;
  name: string;
  category: string;
};

type Freelancer = {
  id: string;
  name: string;
  photo: string;
  badge: string;
  rating: number;
  reviewCount: number;
  bio: string;
  yearsExperience: number;
};

type Offer = {
  id: string;
  title: string;
  description: string;
  price: number;
  currency: string;
  priceType: 'hourly' | 'fixed';
  durationEstimateMin: number;
  service: Service;
  freelancer: Freelancer;
};

type Slot = {
  id: string;
  startTime: string; // full ISO
  endTime: string;   // full ISO
  isBooked: boolean;
};

export default function OfferDetailPage() {
  const { id } = useParams<{ id: string }>();
  const qc = useQueryClient();
  const { token, setToken } = useContext(AuthContext);

  // 1) Load the Offer
  const { data: offer, isLoading: loadingOffer, isError: errorOffer } = useQuery<Offer>({
    queryKey: ['offer', id],
    queryFn: () =>
      axios
        .get<Offer>(`${import.meta.env.VITE_API_BASE_URL}/offers/${id}`)
        .then((r) => r.data),
    enabled: !!id,
  });

  // 2) Build a 7-day window (stripping ms so Go parses it)
  const stripMs = (iso: string) => iso.replace(/\.\d+Z$/, 'Z');
  const today = startOfDay(new Date());
  const fromISO = stripMs(today.toISOString());
  const toISO = stripMs(addDays(today, 7).toISOString());

  

  // 3) Fetch available slots in that window
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
          { params: { available: true, from: fromISO, to: toISO } }
        )
        .then((r) => r.data),
    enabled: !!id,
  });

  // 4) Booking mutation
  const bookSlot = useMutation<void, any, string>({
    mutationFn: (slotId) =>
      axios.post(`${import.meta.env.VITE_API_BASE_URL}/bookings`, {
        offer_id: id,
        slot_id: slotId,
      },
      { headers: { Authorization: `Bearer ${token}` } }
    ),
    onSuccess: () => {
      qc.invalidateQueries({
        queryKey: ['slots', id, fromISO, toISO],
        exact: true,
      });
      setIsBookingModalOpen(false);
      setSelectedSlot(null);
    },
  });

  // 5) UI state
  const [selectedDate, setSelectedDate] = useState<Date>(today);
  const [selectedSlot, setSelectedSlot] = useState<Slot | null>(null);
  const [isBookingModalOpen, setIsBookingModalOpen] = useState(false);
  const [viewMode, setViewMode] = useState<'week' | 'month'>('week');

  const slotsForDate = (date: Date) =>
    slots?.filter(
      (s) =>
        format(parseISO(s.startTime), 'yyyy-MM-dd') ===
        format(date, 'yyyy-MM-dd') && !s.isBooked
    ) || [];

  const getNextSevenDays = () =>
    Array.from({ length: 7 }).map((_, i) => addDays(today, i));

  const formatDuration = (m: number) =>
    m < 60 ? `${m} min` : `${Math.floor(m / 60)} hr${m % 60 ? ` ${m % 60} min` : ''}`;

  if (loadingOffer) return <p className="p-4 text-center">Loading offer…</p>;
  if (errorOffer || !offer)
    return <p className="p-4 text-center text-red-600">Offer not found.</p>;

  

  return (
    <div className="container mx-auto px-4 py-8">
      {/* Offer Header */}
      <div className="mb-8">
        <div className="flex items-center gap-2 mb-2">
          <h1 className="text-3xl font-bold">{offer.title}</h1>
          {/* <Badge variant="secondary" className="capitalize">
            {offer.service.category}
          </Badge> */}
        </div>
        <p className="text-lg text-muted-foreground mb-6">
          {offer.description}
        </p>
        <div className="flex flex-wrap items-center gap-6">
          <div className="flex items-center">
            <DollarSign className="h-5 w-5 mr-1 text-muted-foreground" />
            <span className="text-xl font-medium">
              {offer.currency} {offer.price.toFixed(2)}
            </span>
            <span className="text-muted-foreground ml-1">
              {offer.priceType === 'hourly' ? 'per hour' : 'fixed price'}
            </span>
          </div>
          <div className="flex items-center">
            <Clock className="h-5 w-5 mr-1 text-muted-foreground" />
            <span className="text-xl font-medium">
              {formatDuration(offer.durationEstimateMin)}
            </span>
            <span className="text-muted-foreground ml-1">
              estimated duration
            </span>
          </div>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
        {/* Freelancer Info */}
        <Card className="lg:col-span-1">
          <CardHeader>
            <CardTitle>About the Freelancer</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="flex flex-col items-center text-center mb-6">
              {/* <Avatar className="h-24 w-24 mb-4">
                <AvatarImage
                  src={offer.freelancer.photo}
                  alt={offer.freelancer.name}
                />
                <AvatarFallback>
                  {offer.freelancer.name.charAt(0)}
                </AvatarFallback>
              </Avatar>
              <h3 className="text-xl font-semibold">
                {offer.freelancer.name}
              </h3>
              <div className="flex items-center mt-1 mb-2">
                <Star className="h-4 w-4 mr-1 text-yellow-500 fill-yellow-500" />
                <span className="font-medium">
                  {offer.freelancer.rating}
                </span>
                <span className="text-muted-foreground text-sm ml-1">
                  ({offer.freelancer.reviewCount} reviews)
                </span>
              </div>
              <Badge variant="outline" className="mb-4">
                {offer.freelancer.badge}
              </Badge>
              <p className="text-muted-foreground mb-4">
                {offer.freelancer.bio}
              </p>
              <div className="flex items-center justify-center mb-4">
                <User className="h-4 w-4 mr-2 text-muted-foreground" />
                <span>
                  {offer.freelancer.yearsExperience} years
                  experience
                </span>
              </div> */}
            </div>
          </CardContent>
        </Card>

        {/* Availability */}
        <Card className="lg:col-span-2">
          <CardHeader>
            <div className="flex justify-between items-center">
              <CardTitle>Availability</CardTitle>
              <Tabs>
                <TabsList>
                  <TabsTrigger
                    value="week"
                    onClick={() => setViewMode('week')}
                  >
                    Week
                  </TabsTrigger>
                  <TabsTrigger
                    value="month"
                    onClick={() => setViewMode('month')}
                  >
                    Month
                  </TabsTrigger>
                </TabsList>
              </Tabs>
            </div>
          </CardHeader>

          <CardContent>
            {viewMode === 'week' ? (
              <div className="space-y-6">
                <div className="grid grid-cols-7 gap-2">
                  {getNextSevenDays().map((date) => {
                    const daySlots = slotsForDate(date);
                    const isSel = isSameDay(date, selectedDate);
                    return (
                      <div
                        key={format(date, 'yyyy-MM-dd')}
                        className={`flex flex-col items-center p-2 rounded-md cursor-pointer border ${isSel
                            ? 'border-primary bg-primary/10'
                            : daySlots.length > 0
                              ? 'border-muted hover:border-primary/50'
                              : 'border-muted bg-muted/30'
                          }`}
                        onClick={() => setSelectedDate(date)}
                      >
                        <div className="text-sm font-medium">
                          {format(date, 'EEE')}
                        </div>
                        <div className="text-lg">
                          {format(date, 'd')}
                        </div>
                        <div className="text-xs text-muted-foreground mt-1">
                          {daySlots.length > 0
                            ? `${daySlots.length} slots`
                            : 'No slots'}
                        </div>
                      </div>
                    );
                  })}
                </div>

                <div className="mt-6">
                  <h3 className="text-lg font-medium mb-4">
                    Available Slots for{' '}
                    {format(selectedDate, 'MMMM d, yyyy')}
                  </h3>
                  <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 gap-2">
                    {slotsForDate(selectedDate).length > 0 ? (
                      slotsForDate(selectedDate).map((slot) => (
                        <Button
                          key={slot.id}
                          variant="outline"
                          className="text-center py-6"
                          onClick={() => {
                            setSelectedSlot(slot);
                            setIsBookingModalOpen(true);
                          }}
                        >
                          {format(parseISO(slot.startTime), 'HH:mm')} -{' '}
                          {format(parseISO(slot.endTime), 'HH:mm')}
                        </Button>
                      ))
                    ) : (
                      <p className="text-muted-foreground col-span-full">
                        No available slots for this day.
                      </p>
                    )}
                  </div>
                </div>
              </div>
            ) : (
              <Calendar
                mode="single"
                selected={selectedDate}
                onSelect={(d) => d && setSelectedDate(d)}
                disabled={(d) => {
                  const ds = format(d, 'yyyy-MM-dd');
                  return (
                    !slots?.some(
                      (s) =>
                        format(parseISO(s.startTime), 'yyyy-MM-dd') ===
                        ds && !s.isBooked
                    ) || d < today
                  );
                }}
                className="mx-auto"
              />
            )}
          </CardContent>
        </Card>
      </div>

      {/* Booking Modal */}
      <Dialog
        open={isBookingModalOpen}
        onOpenChange={setIsBookingModalOpen}
      >
        <DialogContent className="sm:max-w-[425px]">
          <DialogHeader>
            <DialogTitle>Confirm Booking</DialogTitle>
            <DialogDescription>
              Please review your booking details before confirming.
            </DialogDescription>
          </DialogHeader>

          {selectedSlot && (
            <div className="space-y-4 py-4">
              <div className="space-y-2">
                <h4 className="font-medium">Service</h4>
                <p>{offer.title}</p>
              </div>

              <div className="grid grid-cols-2 gap-4">
                <div className="space-y-2">
                  <h4 className="font-medium">Date</h4>
                  <p>
                    {format(
                      parseISO(selectedSlot.startTime),
                      'MMMM d, yyyy'
                    )}
                  </p>
                </div>

                <div className="space-y-2">
                  <h4 className="font-medium">Time</h4>
                  <p>
                    {format(parseISO(selectedSlot.startTime), 'HH:mm')} -{' '}
                    {format(parseISO(selectedSlot.endTime), 'HH:mm')}
                  </p>
                </div>
              </div>

              <div className="space-y-2">
                <h4 className="font-medium">Price</h4>
                <p className="text-lg font-semibold">
                  {offer.currency} {offer.price.toFixed(2)}
                </p>
              </div>
            </div>
          )}

          <DialogFooter>
            <Button
              variant="outline"
              onClick={() => setIsBookingModalOpen(false)}
            >
              Cancel
            </Button>
            <Button
              onClick={() => bookSlot.mutate(selectedSlot!.id)}
              className="gap-2"
            >
              <CheckCircle className="h-4 w-4" />
              Confirm booking
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  );
}
