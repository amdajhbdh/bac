import { useState } from 'react'
import { useTranslation } from 'react-i18next'
import { Link, useNavigate } from 'react-router-dom'
import { useMutation } from '@tanstack/react-query'
import { Loader2 } from 'lucide-react'
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { authApi } from '@/lib/api'

export default function Login() {
  const { t } = useTranslation()
  const navigate = useNavigate()
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')

  const loginMutation = useMutation({
    mutationFn: () => authApi.login({ email, password }),
    onSuccess: () => {
      navigate('/')
    },
  })

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    loginMutation.mutate()
  }

  return (
    <div className="max-w-md mx-auto">
      <Card>
        <CardHeader className="text-center">
          <CardTitle className="text-2xl">{t('loginTitle')}</CardTitle>
          <CardDescription>
            {t('hasAccount')} <Link to="/register" className="text-primary-600 hover:underline">{t('register')}</Link>
          </CardDescription>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit} className="space-y-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">{t('email')}</label>
              <Input
                type="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                placeholder="email@example.com"
                required
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">{t('password')}</label>
              <Input
                type="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                placeholder="********"
                required
              />
            </div>
            <Button type="submit" className="w-full" disabled={loginMutation.isPending}>
              {loginMutation.isPending ? (
                <>
                  <Loader2 className="w-5 h-5 animate-spin ml-2" />
                  {t('loading')}
                </>
              ) : (
                t('login')
              )}
            </Button>
          </form>
        </CardContent>
      </Card>
    </div>
  )
}
