import { Routes, Route, Navigate } from 'react-router-dom'
import Layout from './components/Layout'
import Home from './pages/Home'
import Submit from './pages/Submit'
import Solve from './pages/Solve'
import Predictions from './pages/Predictions'
import Leaderboard from './pages/Leaderboard'
import Profile from './pages/Profile'
import Login from './pages/Login'
import Register from './pages/Register'
import { useGlobalHotkeys } from './hooks/useGlobalHotkeys'
import { ErrorBoundary } from './components/ErrorBoundary'

function AppContent() {
  useGlobalHotkeys()

  return (
    <ErrorBoundary>
      <Routes>
        <Route path="/" element={<Layout />}>
          <Route index element={<Home />} />
          <Route path="submit" element={<Submit />} />
          <Route path="solve" element={<Solve />} />
          <Route path="predictions" element={<Predictions />} />
          <Route path="leaderboard" element={<Leaderboard />} />
          <Route path="profile" element={<Profile />} />
        </Route>
        <Route path="/login" element={<Login />} />
        <Route path="/register" element={<Register />} />
        <Route path="*" element={<Navigate to="/" replace />} />
      </Routes>
    </ErrorBoundary>
  )
}

export default function App() {
  return <AppContent />
}
