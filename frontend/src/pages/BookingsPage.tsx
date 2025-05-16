import { useContext } from 'react'
import { useQuery } from '@tanstack/react-query'
import { AuthContext } from '@/context/AuthContext'
import axios from 'axios'
import { Link } from 'react-router-dom'
import { format, parseISO } from 'date-fns'
import {
  Card,
  CardHeader,
  CardTitle,
  CardDescription,
  CardContent,
} from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'

type Booking = {
  id: string
  offerId: string
  slotId: string
  ownerId: string
  status: string
  createdAt: string
  updatedAt: string
}

export default function BookingsPage() {
  const { token } = useContext(AuthContext)

  const { data: bookings, isLoading, isError } = useQuery<Booking[]>({
    queryKey: ['bookings'],
    queryFn: () =>
      axios
        .get<Booking[]>(
          `${import.meta.env.VITE_API_BASE_URL}/bookings`,
          { headers: { Authorization: `Bearer ${token}` } }
        )
        .then(r => r.data),
    enabled: !!token,
  })

  if (isLoading) {
    return (
      <div className="container mx-auto px-4 py-8 text-center">
        Loading your bookings…
      </div>
    )
  }

  if (isError) {
    return (
      <div className="container mx-auto px-4 py-8 text-center text-red-600">
        Failed to load bookings.
      </div>
    )
  }

  if (!bookings || bookings.length === 0) {
    return (
      <div className="container mx-auto px-4 py-8 text-center text-gray-600">
        You have no bookings yet.
      </div>
    )
  }

  return (
    <div className="container mx-auto px-4 py-8 space-y-6">
      <h1 className="text-3xl font-bold">My Bookings</h1>

      <div className="space-y-4">
        {bookings.map(booking => (
          <Card key={booking.id}>
            <CardHeader>
              <CardTitle>Booking #{booking.id}</CardTitle>
              <CardDescription>
                Created on {format(parseISO(booking.createdAt), 'PPP p')}
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="flex flex-col sm:flex-row sm:justify-between sm:items-center space-y-4 sm:space-y-0">
                <div className="space-y-1">
                  <p>
                    <span className="font-medium">Status:</span>{' '}
                    <Badge
                      variant={
                        booking.status === 'pending'
                          ? 'outline'
                          : booking.status === 'confirmed'
                          ? 'secondary'
                          : 'destructive'
                      }
                      className="capitalize"
                    >
                      {booking.status}
                    </Badge>
                  </p>
                </div>
                <Link to={`/bookings/${booking.id}`}>
                  <Button variant="outline">View Details</Button>
                </Link>
              </div>
            </CardContent>
          </Card>
        ))}
      </div>
    </div>
  )
}
