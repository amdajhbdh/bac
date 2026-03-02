import { Link } from 'react-router-dom'
import { useQuery } from '@tanstack/react-query'
import { useTranslation } from 'react-i18next'
import { 
  Upload, 
  Calculator, 
  TrendingUp, 
  Trophy, 
  ArrowLeft,
  Loader2
} from 'lucide-react'
import { Card, CardContent } from '@/components/ui/card'
import { healthApi, predictionsApi } from '@/lib/api'

const features = [
  { name: 'features.submitImage', description: 'features.submitDesc', href: '/submit', icon: Upload, color: 'bg-blue-500' },
  { name: 'features.solveProblem', description: 'features.solveDesc', href: '/solve', icon: Calculator, color: 'bg-green-500' },
  { name: 'features.viewPredictions', description: 'features.predictionsDesc', href: '/predictions', icon: TrendingUp, color: 'bg-purple-500' },
  { name: 'features.viewLeaderboard', description: 'features.leaderboardDesc', href: '/leaderboard', icon: Trophy, color: 'bg-yellow-500' },
]

export default function Home() {
  const { t } = useTranslation()

  const { data: health } = useQuery({
    queryKey: ['health'],
    queryFn: () => healthApi.check(),
    retry: false,
  })

  const { data: predictions, isLoading } = useQuery({
    queryKey: ['predictions', 'latest'],
    queryFn: () => predictionsApi.getLatest(),
  })

  return (
    <div className="space-y-8">
      {/* Hero Section */}
      <div className="text-center py-8">
        <h1 className="text-4xl font-bold text-gray-900 mb-4">
          {t('welcome')}
        </h1>
        <p className="text-xl text-gray-600 max-w-2xl mx-auto">
          {t('subtitle')}
        </p>
      </div>

      {/* Status Badge */}
      <div className="flex justify-center">
        <div className={`inline-flex items-center px-4 py-2 rounded-full text-sm ${
          health?.data?.status === 'ok' ? 'bg-green-100 text-green-800' : 'bg-yellow-100 text-yellow-800'
        }`}>
          <span className={`w-2 h-2 rounded-full ml-2 ${
            health?.data?.status === 'ok' ? 'bg-green-500' : 'bg-yellow-500'
          }`}></span>
          {t('serverStatus')}: {health?.data?.status === 'ok' ? t('working') : t('offline')}
        </div>
      </div>

      {/* Features Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        {features.map((feature) => (
          <Link
            key={feature.name}
            to={feature.href}
          >
            <Card className="hover:shadow-lg transition-shadow group cursor-pointer h-full">
              <CardContent className="flex items-start space-x-4 p-6">
                <div className={`${feature.color} p-3 rounded-lg`}>
                  <feature.icon className="w-6 h-6 text-white" />
                </div>
                <div className="flex-1">
                  <h3 className="text-lg font-semibold text-gray-900 group-hover:text-primary-600">
                    {t(feature.name)}
                  </h3>
                  <p className="text-gray-600 mt-1">{t(feature.description)}</p>
                </div>
                <ArrowLeft className="w-5 h-5 text-gray-400 group-hover:text-primary-600" />
              </CardContent>
            </Card>
          </Link>
        ))}
      </div>

      {/* Latest Predictions */}
      {isLoading ? (
        <div className="flex justify-center py-8">
          <Loader2 className="w-8 h-8 animate-spin text-primary-600" />
        </div>
      ) : predictions?.data && Array.isArray(predictions.data) && predictions.data.length > 0 ? (
        <Card>
          <CardContent className="p-6">
            <div className="flex items-center space-x-2 mb-4">
              <TrendingUp className="w-5 h-5 text-primary-600" />
              <h2 className="text-xl font-semibold">{t('latestPredictions')}</h2>
            </div>
            <div className="space-y-3">
              {predictions.data.slice(0, 3).map((pred: unknown) => (
                <Link
                  key={(pred as any)?.id || Math.random()}
                  to={`/predictions/${(pred as any)?.id}`}
                  className="block p-3 rounded-lg bg-gray-50 hover:bg-gray-100"
                >
                  <div className="flex justify-between items-center">
                    <span className="font-medium">{(pred as any)?.subject || 'N/A'}</span>
                    <span className="text-sm text-gray-500">{(pred as any)?.year}</span>
                  </div>
                </Link>
              ))}
            </div>
            <Link
              to="/predictions"
              className="block text-center text-primary-600 hover:text-primary-700 mt-4"
            >
              {t('viewAll')}
            </Link>
          </CardContent>
        </Card>
      ) : null}
    </div>
  )
}
