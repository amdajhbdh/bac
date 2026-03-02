import { useEffect } from 'react'
import { useNavigate } from 'react-router-dom'

export const useGlobalHotkeys = () => {
  const navigate = useNavigate()

  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      // Don't trigger if user is typing in an input
      if (
        e.target instanceof HTMLInputElement ||
        e.target instanceof HTMLTextAreaElement
      ) {
        return
      }

      // h - Home
      if (e.key === 'h' && !e.ctrlKey && !e.metaKey) {
        navigate('/')
      }
      // s - Submit
      else if (e.key === 's' && !e.ctrlKey && !e.metaKey) {
        navigate('/submit')
      }
      // k - Solve
      else if (e.key === 'k' && !e.ctrlKey && !e.metaKey) {
        navigate('/solve')
      }
      // p - Predictions
      else if (e.key === 'p' && !e.ctrlKey && !e.metaKey) {
        navigate('/predictions')
      }
      // l - Leaderboard
      else if (e.key === 'l' && !e.ctrlKey && !e.metaKey) {
        navigate('/leaderboard')
      }
      // u - Profile
      else if (e.key === 'u' && !e.ctrlKey && !e.metaKey) {
        navigate('/profile')
      }
    }

    window.addEventListener('keydown', handleKeyDown)
    return () => window.removeEventListener('keydown', handleKeyDown)
  }, [navigate])
}
