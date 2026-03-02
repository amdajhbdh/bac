import { useTranslation } from 'react-i18next'
import { useQuery } from '@tanstack/react-query'
import { User, Award, TrendingUp, Loader2 } from 'lucide-react'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { userApi } from '@/lib/api'

export default function Profile() {
  const { t } = useTranslation()

  const { data: profile, isLoading: profileLoading } = useQuery({
    queryKey: ['profile'],
    queryFn: () => userApi.getProfile(),
  })

  const { data: stats, isLoading: statsLoading } = useQuery({
    queryKey: ['stats'],
    queryFn: () => userApi.getStats(),
  })

  const { data: badges, isLoading: badgesLoading } = useQuery({
    queryKey: ['badges'],
    queryFn: () => userApi.getBadges(),
  })

  const isLoading = profileLoading || statsLoading || badgesLoading

  if (isLoading) {
    return (
      <div className="flex justify-center py-12">
        <Loader2 className="w-8 h-8 animate-spin text-primary-600" />
      </div>
    )
  }

  return (
    <div className="max-w-4xl mx-auto space-y-6">
      {/* Profile Header */}
      <Card>
        <CardContent className="p-6">
          <div className="flex items-center space-x-6">
            <div className="w-20 h-20 bg-primary-100 rounded-full flex items-center justify-center">
              <User className="w-10 h-10 text-primary-600" />
            </div>
            <div className="flex-1">
              <h1 className="text-2xl font-bold">{profile?.data?.name || 'Student'}</h1>
              <p className="text-gray-500">{profile?.data?.email || 'user@example.com'}</p>
              <div className="flex items-center space-x-4 mt-2">
                <span className="px-3 py-1 bg-primary-100 text-primary-700 rounded-full text-sm">
                  {profile?.data?.level || 'Level 1'}
                </span>
                <span className="text-sm text-gray-500">
                  {stats?.data?.points || 0} {t('points')}
                </span>
              </div>
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Stats Grid */}
      <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
        <Card>
          <CardContent className="p-4 text-center">
            <TrendingUp className="w-6 h-6 mx-auto text-green-500 mb-2" />
            <div className="text-2xl font-bold">{stats?.data?.streak || 0}</div>
            <div className="text-sm text-gray-500">Day Streak</div>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="p-4 text-center">
            <Award className="w-6 h-6 mx-auto text-yellow-500 mb-2" />
            <div className="text-2xl font-bold">{stats?.data?.solved || 0}</div>
            <div className="text-sm text-gray-500">Problems Solved</div>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="p-4 text-center">
            <TrendingUp className="w-6 h-6 mx-auto text-blue-500 mb-2" />
            <div className="text-2xl font-bold">{stats?.data?.accuracy || 0}%</div>
            <div className="text-sm text-gray-500">Accuracy</div>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="p-4 text-center">
            <User className="w-6 h-6 mx-auto text-purple-500 mb-2" />
            <div className="text-2xl font-bold">{stats?.data?.rank || '-'}</div>
            <div className="text-sm text-gray-500">{t('rank')}</div>
          </CardContent>
        </Card>
      </div>

      {/* Badges */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center space-x-2">
            <Award className="w-5 h-5" />
            <span>{t('badges')}</span>
          </CardTitle>
        </CardHeader>
        <CardContent>
          {Array.isArray(badges?.data) && badges.data.length > 0 ? (
            <div className="grid grid-cols-3 md:grid-cols-6 gap-4">
              {badges.data.map((badge: unknown) => (
                <div key={(badge as any)?.id || Math.random()} className="text-center">
                  <div className="w-12 h-12 mx-auto bg-gray-100 rounded-full flex items-center justify-center">
                    <Award className="w-6 h-6 text-yellow-500" />
                  </div>
                  <div className="text-xs mt-1">{(badge as any)?.name}</div>
                </div>
              ))}
            </div>
          ) : (
            <p className="text-center text-gray-500 py-4">No badges yet</p>
          )}
        </CardContent>
      </Card>

      {/* Progress */}
      <Card>
        <CardHeader>
          <CardTitle>{t('progress')}</CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          {['math', 'pc', 'svt', 'philosophie'].map((subject) => (
            <div key={subject}>
              <div className="flex justify-between text-sm mb-1">
                <span className="capitalize">{subject}</span>
                <span className="text-gray-500">
                  {profile?.data?.progress?.[subject] || 0}%
                </span>
              </div>
              <div className="h-2 bg-gray-100 rounded-full overflow-hidden">
                <div
                  className="h-full bg-primary-500 rounded-full"
                  style={{ width: `${profile?.data?.progress?.[subject] || 0}%` }}
                />
              </div>
            </div>
          ))}
        </CardContent>
      </Card>
    </div>
  )
}
