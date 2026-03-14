/// <reference types="vite/client" />

import axios from 'axios'

const API_BASE_URL = import.meta.env.VITE_API_URL || '/api/v1'

export const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
})

// Auth
export const authApi = {
  register: (data: { email: string; password: string; name: string }) =>
    api.post('/auth/register', data),
  login: (data: { email: string; password: string }) =>
    api.post('/auth/login', data),
  refresh: (token: string) =>
    api.post('/auth/refresh', { token }),
}

// Questions
export const questionsApi = {
  list: (params?: { subject?: string; chapter?: number }) =>
    api.get('/questions', { params }),
  get: (id: string) => api.get(`/questions/${id}`),
  create: (data: unknown) => api.post('/questions', data),
  update: (id: string, data: unknown) => api.put(`/questions/${id}`, data),
  delete: (id: string) => api.delete(`/questions/${id}`),
}

// Submission (OCR)
export const submissionApi = {
  submitImage: (file: File) => {
    const formData = new FormData()
    formData.append('file', file)
    return api.post('/submit/image', formData, {
      headers: { 'Content-Type': 'multipart/form-data' },
    })
  },
  submitPDF: (file: File) => {
    const formData = new FormData()
    formData.append('file', file)
    return api.post('/submit/pdf', formData, {
      headers: { 'Content-Type': 'multipart/form-data' },
    })
  },
  submitURL: (url: string) => api.post('/submit/url', { url }),
  getStatus: (jobId: string) => api.get(`/submit/status/${jobId}`),
}

// Solver
export const solverApi = {
  solve: (problem: string) => api.post('/solve', { problem }),
  getSteps: (id: string) => api.get(`/solve/${id}/steps`),
  getAnimation: (id: string) => api.get(`/solve/${id}/animation`),
}

// Predictions
export const predictionsApi = {
  list: () => api.get('/predictions'),
  get: (id: string) => api.get(`/predictions/${id}`),
  getBySubject: (subject: string) => api.get(`/predictions/subject/${subject}`),
  getLatest: () => api.get('/predictions/latest'),
}

// User
export const userApi = {
  getProfile: () => api.get('/user/me'),
  updateProfile: (data: unknown) => api.put('/user/me', data),
  getProgress: () => api.get('/user/progress'),
  getStats: () => api.get('/user/stats'),
  getBadges: () => api.get('/user/badges'),
}

// Practice
export const practiceApi = {
  start: (subject: string) => api.post('/practice/start', { subject }),
  answer: (sessionId: string, questionId: string, answer: string) =>
    api.post('/practice/answer', { session_id: sessionId, question_id: questionId, answer }),
  end: (sessionId: string) => api.post(`/practice/${sessionId}/end`),
}

// Leaderboard
export const leaderboardApi = {
  get: () => api.get('/leaderboard'),
}

// Subjects
export const subjectsApi = {
  list: () => api.get('/subjects'),
}

// Health
export const healthApi = {
  check: () => api.get('/health'),
}
