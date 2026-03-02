import { useState } from 'react'
import { useTranslation } from 'react-i18next'
import { useMutation } from '@tanstack/react-query'
import { Upload, FileImage, Loader2, Check, AlertCircle } from 'lucide-react'
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { submissionApi } from '@/lib/api'

export default function Submit() {
  const { t } = useTranslation()
  const [file, setFile] = useState<File | null>(null)
  const [preview, setPreview] = useState<string | null>(null)
  const [result, setResult] = useState<any>(null)

  const submitMutation = useMutation({
    mutationFn: async () => {
      if (!file) return null
      const isPDF = file.type === 'application/pdf'
      const submitFn = isPDF ? submissionApi.submitPDF : submissionApi.submitImage
      return submitFn(file)
    },
    onSuccess: (data) => {
      if (data?.data) {
        setResult(data.data)
      }
    },
  })

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const selected = e.target.files?.[0]
    if (selected) {
      setFile(selected)
      setResult(null)
      if (selected.type.startsWith('image/')) {
        setPreview(URL.createObjectURL(selected))
      } else {
        setPreview(null)
      }
    }
  }

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    submitMutation.mutate()
  }

  return (
    <div className="max-w-2xl mx-auto space-y-6">
      <Card>
        <CardHeader className="text-center">
          <CardTitle className="text-2xl">{t('submitTitle')}</CardTitle>
          <CardDescription>{t('submitDesc')}</CardDescription>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit} className="space-y-6">
            {/* File Input */}
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                {t('chooseFile')}
              </label>
              <div className="border-2 border-dashed border-gray-300 rounded-lg p-8 text-center hover:border-primary-500 transition-colors">
                <input
                  type="file"
                  accept="image/*,.pdf"
                  onChange={handleFileChange}
                  className="hidden"
                  id="file-input"
                />
                <label htmlFor="file-input" className="cursor-pointer">
                  {preview ? (
                    <img src={preview} alt="Preview" className="max-h-64 mx-auto rounded-lg" />
                  ) : (
                    <div className="space-y-2">
                      <Upload className="w-12 h-12 mx-auto text-gray-400" />
                      <p className="text-gray-600">{t('clickToUpload')}</p>
                      <p className="text-sm text-gray-400">{t('supportedFormats')}</p>
                    </div>
                  )}
                </label>
              </div>
            </div>

            {file && (
              <div className="flex items-center space-x-2 text-sm text-gray-600">
                <FileImage className="w-4 h-4" />
                <span>{file.name}</span>
                <span className="text-gray-400">({(file.size / 1024).toFixed(1)} KB)</span>
              </div>
            )}

            {/* Submit Button */}
            <Button
              type="submit"
              disabled={!file || submitMutation.isPending}
              className="w-full"
            >
              {submitMutation.isPending ? (
                <>
                  <Loader2 className="w-5 h-5 animate-spin ml-2" />
                  {t('sending')}
                </>
              ) : (
                <>
                  <Upload className="w-5 h-5 ml-2" />
                  {t('send')}
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
            <CardTitle className="flex items-center space-x-2">
              <Check className="w-5 h-5 text-green-500" />
              <span>{t('result')}</span>
            </CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            {result.text && (
              <div className="bg-gray-50 rounded-lg p-4">
                <pre className="whitespace-pre-wrap text-sm">{result.text}</pre>
              </div>
            )}

            <div className="grid grid-cols-2 gap-4 text-sm">
              <div>
                <span className="text-gray-500">{t('ocrSource')}:</span>
                <span className="mr-2 font-medium">{result.source || 'N/A'}</span>
              </div>
              <div>
                <span className="text-gray-500">{t('confidence')}:</span>
                <span className="mr-2 font-medium">
                  {result.confidence ? `${(result.confidence * 100).toFixed(0)}%` : 'N/A'}
                </span>
              </div>
              {result.nlm_analysis && (
                <>
                  <div>
                    <span className="text-gray-500">{t('subject')}:</span>
                    <span className="mr-2 font-medium">{result.nlm_analysis.subject}</span>
                  </div>
                  <div>
                    <span className="text-gray-500">{t('concepts')}:</span>
                    <span className="mr-2 font-medium">
                      {result.nlm_analysis.concepts?.join(', ')}
                    </span>
                  </div>
                </>
              )}
            </div>
          </CardContent>
        </Card>
      )}

      {/* Error */}
      {submitMutation.isError && (
        <Card className="bg-red-50 border-red-200">
          <CardContent className="flex items-center space-x-2 text-red-800 p-6">
            <AlertCircle className="w-5 h-5" />
            <span>{t('error')}</span>
          </CardContent>
        </Card>
      )}
    </div>
  )
}
