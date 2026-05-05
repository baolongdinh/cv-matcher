import axios from 'axios'
import type { AnalysisResponse, AnalysisResult } from '../types/analysis'

export interface HistoryResponse {
  status: 'success' | 'error'
  data?: AnalysisResult[]
  error?: {
    code: string
    message: string
  }
}

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8000/api'

let sessionId = localStorage.getItem('cv_matcher_session_id')
if (!sessionId) {
  sessionId = crypto.randomUUID ? crypto.randomUUID() : `sess_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`
  localStorage.setItem('cv_matcher_session_id', sessionId)
}

const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'X-Session-ID': sessionId
  }
})

export const analyzeCV = async (
  jd: string,
  cvFile: File,
  apiKey: string
): Promise<AnalysisResponse> => {
  const formData = new FormData()
  formData.append('job_description', jd)
  formData.append('cv_file', cvFile)
  formData.append('api_key', apiKey)

  const response = await api.post<AnalysisResponse>('/analyze', formData, {
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  })

  return response.data
}

export const fetchHistory = async (): Promise<HistoryResponse> => {
  const response = await api.get<HistoryResponse>('/history')
  return response.data
}

export const deleteHistoryItem = async (id: string) => {
  const response = await api.delete(`/history/${id}`)
  return response.data
}

export const checkHealth = async () => {
  const response = await api.get('/health')
  return response.data
}

export default api
