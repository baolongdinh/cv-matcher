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

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL ||
  (window.location.hostname === 'localhost' || window.location.hostname === '127.0.0.1'
    ? 'http://localhost:8000/api'
    : `http://${window.location.hostname}:8000/api`)

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

const getFileHash = async (file: File): Promise<string> => {
  try {
    const buffer = await file.arrayBuffer()

    // Check if crypto.subtle is available (only in Secure Contexts: HTTPS or localhost)
    if (window.crypto && window.crypto.subtle) {
      const hashBuffer = await crypto.subtle.digest('SHA-256', buffer)
      const hashArray = Array.from(new Uint8Array(hashBuffer))
      return hashArray.map(b => b.toString(16).padStart(2, '0')).join('')
    }

    // Fallback for non-secure contexts (e.g., accessing via network IP without HTTPS)
    console.warn('[API] crypto.subtle not available. Using fallback hash.')
    let hash = 0
    const view = new Uint8Array(buffer)
    for (let i = 0; i < view.length; i++) {
      hash = ((hash << 5) - hash) + view[i]
      hash |= 0 // Convert to 32bit integer
    }
    return `fallback-${file.name}-${file.size}-${hash}`
  } catch (e) {
    console.error('[API] Error generating file hash:', e)
    return `error-${file.name}-${file.size}-${Date.now()}`
  }
}

export const analyzeCV = async (
  jd: string,
  cvFile: File,
  apiKey: string
): Promise<AnalysisResponse> => {
  const fileHash = await getFileHash(cvFile)
  const formData = new FormData()
  formData.append('job_description', jd)
  formData.append('cv_file', cvFile)
  formData.append('api_key', apiKey)

  const response = await api.post<AnalysisResponse>('/analyze', formData, {
    headers: {
      'Content-Type': 'multipart/form-data',
      'X-File-Hash': fileHash
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
