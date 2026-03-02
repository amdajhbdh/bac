import { useState, useEffect } from 'react'
import { useTranslation } from 'react-i18next'
import { useMutation } from '@tanstack/react-query'
import { Calculator, Loader2, Copy, Check } from 'lucide-react'
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { solverApi } from '@/lib/api'

const examples = [
  '2x + 5 = 15',
  'x^2 - 4 = 0',
  '∫2x dx',
]

export default function Solve() {
  const { t } = useTranslation()
  const [problem, setProblem] = useState('')
  const [result, setResult] = useState<any>(null)
  const [copied, setCopied] = useState(false)

  const solveMutation = useMutation({
    mutationFn: (problem: string) => solverApi.solve(problem),
    onSuccess: (data) => {
      if (data?.data) {
        setResult(data.data)
      }
    },
  })

  useEffect(() => {
    const handleSubmit = () => {
      if (problem.trim()) {
        solveMutation.mutate(problem)
      }
    }

    window.addEventListener('hotkey-submit', handleSubmit)
    return () => window.removeEventListener('hotkey-submit', handleSubmit)
  }, [problem])

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    if (problem.trim()) {
      solveMutation.mutate(problem)
    }
  }

  const handleCopy = () => {
    if (result?.solution) {
      navigator.clipboard.writeText(result.solution)
      setCopied(true)
      setTimeout(() => setCopied(false), 2000)
    }
  }

  return (
    <div className="max-w-3xl mx-auto space-y-6">
      <Card>
        <CardHeader className="text-center">
          <CardTitle className="text-2xl">{t('solveTitle')}</CardTitle>
          <CardDescription>{t('solveDesc')}</CardDescription>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit} className="space-y-4">
            <textarea
              value={problem}
              onChange={(e) => setProblem(e.target.value)}
              placeholder={t('enterProblem')}
              className="input h-32 resize-none"
              dir="ltr"
            />
            <Button
              type="submit"
              disabled={!problem.trim() || solveMutation.isPending}
              className="w-full"
            >
              {solveMutation.isPending ? (
                <>
                  <Loader2 className="w-5 h-5 animate-spin ml-2" />
                  {t('solving')}
                </>
              ) : (
                <>
                  <Calculator className="w-5 h-5 ml-2" />
                  {t('solve')}
                </>
              )}
            </Button>
          </form>
        </CardContent>
      </Card>

      {/* Result */}
      {result && (
        <Card>
          <CardHeader>
            <div className="flex justify-between items-center">
              <CardTitle>{t('solution')}</CardTitle>
              <Button variant="outline" size="sm" onClick={handleCopy}>
                {copied ? <Check className="w-4 h-4 ml-1" /> : <Copy className="w-4 h-4 ml-1" />}
                {copied ? t('copied') : t('copy')}
              </Button>
            </div>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="bg-gray-50 rounded-lg p-4">
              <pre className="whitespace-pre-wrap text-sm" dir="ltr">
                {result.solution}
              </pre>
            </div>

            <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
              <div className="bg-gray-50 rounded-lg p-3 text-center">
                <div className="text-lg font-bold text-primary-600">
                  {result.confidence ? `${(result.confidence * 100).toFixed(0)}%` : '-'}
                </div>
                <div className="text-xs text-gray-500">{t('confidence')}</div>
              </div>
              <div className="bg-gray-50 rounded-lg p-3 text-center">
                <div className="text-lg font-bold text-green-600">{result.steps || '-'}</div>
                <div className="text-xs text-gray-500">{t('steps')}</div>
              </div>
              <div className="bg-gray-50 rounded-lg p-3 text-center">
                <div className="text-lg font-bold text-purple-600">{result.subject || '-'}</div>
                <div className="text-xs text-gray-500">{t('subject')}</div>
              </div>
              <div className="bg-gray-50 rounded-lg p-3 text-center">
                <div className="text-lg font-bold text-yellow-600 truncate">{result.model || '-'}</div>
                <div className="text-xs text-gray-500">{t('model')}</div>
              </div>
            </div>

            {result.concepts && Array.isArray(result.concepts) && result.concepts.length > 0 && (
              <div>
                <h3 className="text-sm font-medium text-gray-700 mb-2">{t('concepts')}:</h3>
                <div className="flex flex-wrap gap-2">
                  {result.concepts.map((concept: string, i: number) => (
                    <span
                      key={i}
                      className="px-3 py-1 bg-primary-50 text-primary-700 rounded-full text-sm"
                    >
                      {concept}
                    </span>
                  ))}
                </div>
              </div>
            )}

            {result.similar_found > 0 && (
              <div className="text-sm text-gray-500">
                {t('similarFound', { count: result.similar_found })}
              </div>
            )}
          </CardContent>
        </Card>
      )}

      {/* Error */}
      {solveMutation.isError && (
        <Card className="bg-red-50 border-red-200">
          <CardContent className="p-6">
            <p className="text-red-800">{t('error')}</p>
          </CardContent>
        </Card>
      )}

      {/* Examples */}
      <Card>
        <CardHeader>
          <CardTitle>{t('examples')}</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="space-y-2">
            {examples.map((ex) => (
              <button
                key={ex}
                onClick={() => setProblem(ex)}
                className="block w-full text-left px-3 py-2 rounded-lg bg-gray-50 hover:bg-gray-100 text-sm"
                dir="ltr"
              >
                {ex}
              </button>
            ))}
          </div>
        </CardContent>
      </Card>
    </div>
  )
}
