import { useQuery } from '@tanstack/react-query'
import { useTranslation } from 'react-i18next'
import { Trophy, Medal, Loader2 } from 'lucide-react'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { leaderboardApi } from '@/lib/api'

export default function Leaderboard() {
  const { t } = useTranslation()

  const { data, isLoading } = useQuery({
    queryKey: ['leaderboard'],
    queryFn: () => leaderboardApi.get(),
  })

  if (isLoading) {
    return (
      <div className="flex justify-center py-12">
        <Loader2 className="w-8 h-8 animate-spin text-primary-600" />
      </div>
    )
  }

  const getRankIcon = (rank: number) => {
    if (rank === 1) return <Medal className="w-6 h-6 text-yellow-500" />
    if (rank === 2) return <Medal className="w-6 h-6 text-gray-400" />
    if (rank === 3) return <Medal className="w-6 h-6 text-amber-600" />
    return <span className="text-lg font-bold text-gray-500">#{rank}</span>
  }

  const getRankBg = (rank: number) => {
    if (rank === 1) return 'bg-yellow-50 border-yellow-200'
    if (rank === 2) return 'bg-gray-50 border-gray-200'
    if (rank === 3) return 'bg-amber-50 border-amber-200'
    return ''
  }

  const leaderboardData = Array.isArray(data?.data) ? data.data : []

  return (
    <div className="max-w-2xl mx-auto space-y-6">
      <Card>
        <CardHeader className="text-center">
          <CardTitle className="text-2xl flex items-center justify-center space-x-2">
            <Trophy className="w-6 h-6 text-yellow-500" />
            <span>{t('leaderboardTitle')}</span>
          </CardTitle>
        </CardHeader>
        <CardContent>
          {leaderboardData.length > 0 ? (
            <div className="space-y-3">
              {leaderboardData.map((user: unknown, index: number) => (
                <div
                  key={(user as any)?.id || index}
                  className={`flex items-center justify-between p-4 rounded-lg border ${getRankBg(index + 1)}`}
                >
                  <div className="flex items-center space-x-4">
                    <div className="w-10 flex justify-center">
                      {getRankIcon(index + 1)}
                    </div>
                    <div className="w-10 h-10 bg-gray-200 rounded-full flex items-center justify-center">
                      <span className="font-medium">{(user as any)?.name?.[0] || '?'}</span>
                    </div>
                    <div>
                      <div className="font-medium">{(user as any)?.name || 'Anonymous'}</div>
                      {(user as any)?.school && (
                        <div className="text-sm text-gray-500">{(user as any)?.school}</div>
                      )}
                    </div>
                  </div>
                  <div className="text-left">
                    <div className="text-lg font-bold text-primary-600">{(user as any)?.points}</div>
                    <div className="text-xs text-gray-500">{t('points')}</div>
                  </div>
                </div>
              ))}
            </div>
          ) : (
            <div className="text-center py-8 text-gray-500">
              <Trophy className="w-12 h-12 mx-auto mb-4 text-gray-300" />
              <p>{t('noLeaderboard')}</p>
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  )
}
