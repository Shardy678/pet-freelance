import { useContext } from 'react'
import { useQuery } from '@tanstack/react-query'
import { AuthContext } from '@/context/AuthContext'
import axios from 'axios'
import { formatDistanceToNow, parseISO } from 'date-fns'

import {
  Card,
  CardHeader,
  CardTitle,
  CardContent,
} from '@/components/ui/card'
import { Calendar } from 'lucide-react'

type Activity = {
  id: string
  title: string
  message: string
  type: string
  createdAt: string
}

export default function RecentActivity() {
  const { token } = useContext(AuthContext)

  const { data: activities, isLoading } = useQuery<Activity[]>({
    queryKey: ['activities'],
    queryFn: () =>
      axios
        .get<Activity[]>(
          `${import.meta.env.VITE_API_BASE_URL}/activities?limit=5`,
          { headers: { Authorization: `Bearer ${token}` } }
        )
        .then((res) => res.data),
    enabled: !!token,
  })

  if (isLoading || !activities) return null

  return (
    <Card className="max-w-md shadow-none">
      <CardHeader>
        <CardTitle>Recent activity</CardTitle>
      </CardHeader>
      <CardContent className="space-y-4">
        {activities.map((act) => (
          <div
            key={act.id}
            className="flex items-center space-x-3 border-b pb-3 last:border-0 last:pb-0"
          >
            {act.type === 'appointment' && (
              <div className="bg-gray-100 rounded-full p-3 flex-shrink-0">
                <Calendar className="w-5 h-5 text-red-500" />
              </div>
            )}
            <div className="space-y-1">
              <p className="font-medium">{act.title}</p>
              <p className="text-sm text-gray-600">{act.message}</p>
              <p className="text-xs text-gray-400">
                {formatDistanceToNow(parseISO(act.createdAt), { addSuffix: true })}
              </p>
            </div>
          </div>
        ))}
      </CardContent>

    </Card>
  )
}
