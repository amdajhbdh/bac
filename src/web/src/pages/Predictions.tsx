import { useQuery } from '@tanstack/react-query'
import { useTranslation } from 'react-i18next'
import { TrendingUp, BookOpen, Calendar, Loader2 } from 'lucide-react'
import { Card, CardContent } from '@/components/ui/card'
import { predictionsApi } from '@/lib/api'

export default function Predictions() {
  const { t } = useTranslation()

  const { data, isLoading } = useQuery({
    queryKey: ['predictions'],
    queryFn: () => predictionsApi.list(),
  })

  const subjects = [
    { key: 'math', label: 'الرياضيات', color: 'bg-blue-100 text-blue-700' },
    { key: 'pc', label: 'الفيزياء', color: 'bg-green-100 text-green-700' },
    { key: 'svt', label: 'العلوم', color: 'bg-purple-100 text-purple-700' },
    { key: 'philosophie', label: 'الفلسفة', color: 'bg-yellow-100 text-yellow-700' },
  ]

  if (isLoading) {
    return (
      <div className="flex justify-center py-12">
        <Loader2 className="w-8 h-8 animate-spin text-primary-600" />
      </div>
    )
  }

  const predictionsData = Array.isArray(data?.data) ? data.data : []

  return (
    <div className="space-y-6">
      <div className="text-center">
        <h1 className="text-3xl font-bold text-gray-900">{t('predictionsTitle')}</h1>
        <p className="text-gray-600 mt-2">{t('predictionsDesc')}</p>
      </div>

      {/* Subject Filters */}
      <div className="flex flex-wrap justify-center gap-2">
        {subjects.map((subject) => (
          <button
            key={subject.key}
            className={`px-4 py-2 rounded-full text-sm font-medium ${subject.color}`}
          >
            {subject.label}
          </button>
        ))}
      </div>

      {/* Predictions Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        {predictionsData.map((pred: unknown) => (
          <Card key={(pred as any)?.id || Math.random()} className="hover:shadow-lg transition-shadow cursor-pointer">
            <CardContent className="p-6">
              <div className="flex items-start justify-between mb-3">
                <div className="flex items-center space-x-2">
                  <div className="p-2 bg-primary-100 rounded-lg">
                    <BookOpen className="w-5 h-5 text-primary-600" />
                  </div>
                  <div>
                    <h3 className="font-semibold">{(pred as any)?.subject || 'N/A'}</h3>
                    <p className="text-sm text-gray-500">{(pred as any)?.chapter}</p>
                  </div>
                </div>
                <div className="flex items-center space-x-1 text-green-600">
                  <TrendingUp className="w-4 h-4" />
                  <span className="text-sm font-medium">{(pred as any)?.probability}%</span>
                </div>
              </div>

              <p className="text-gray-600 text-sm mb-3">{(pred as any)?.description}</p>

              <div className="flex items-center justify-between text-sm text-gray-500">
                <div className="flex items-center space-x-1">
                  <Calendar className="w-4 h-4" />
                  <span>{t('year')}: {(pred as any)?.year}</span>
                </div>
                {(pred as any)?.topics && Array.isArray((pred as any).topics) && (
                  <div className="flex gap-1">
                    {(pred as any).topics.slice(0, 2).map((topic: string, i: number) => (
                      <span key={i} className="px-2 py-0.5 bg-gray-100 rounded text-xs">
                        {topic}
                      </span>
                    ))}
                  </div>
                )}
              </div>
            </CardContent>
          </Card>
        ))}

        {predictionsData.length === 0 && (
          <div className="col-span-2 text-center py-12 text-gray-500">
            <TrendingUp className="w-12 h-12 mx-auto mb-4 text-gray-300" />
            <p>{t('noPredictions')}</p>
          </div>
        )}
      </div>
    </div>
  )
}
