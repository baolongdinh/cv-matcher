import axios from 'axios'
import type {
  APIResponse,
  AnalysisResult,
  BatchAnalyzeData,
  BatchNotificationData,
  BatchResultsData,
  BatchStatusData,
} from '../types/analysis'

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL ||
  (window.location.hostname === 'localhost' || window.location.hostname === '127.0.0.1'
    ? 'http://localhost:8000/api'
    : `http://${window.location.hostname}:8000/api`)

const SESSION_STORAGE_KEY = 'cv_matcher_session_id'
const ACTIVE_BATCH_STORAGE_KEY = 'cv_matcher_active_batch'
const ACTIVE_BATCH_VIEW_STORAGE_KEY = 'cv_matcher_active_batch_view'

let sessionId = localStorage.getItem(SESSION_STORAGE_KEY)
if (!sessionId) {
  sessionId = crypto.randomUUID ? crypto.randomUUID() : `sess_${Date.now()}_${Math.random().toString(36).slice(2, 11)}`
  localStorage.setItem(SESSION_STORAGE_KEY, sessionId)
}

const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'X-Session-ID': sessionId,
  },
})

const getFileHash = async (file: File): Promise<string> => {
  try {
    const buffer = await file.arrayBuffer()
    if (window.crypto && window.crypto.subtle) {
      const hashBuffer = await crypto.subtle.digest('SHA-256', buffer)
      const hashArray = Array.from(new Uint8Array(hashBuffer))
      return hashArray.map((byte) => byte.toString(16).padStart(2, '0')).join('')
    }

    let hash = 0
    const view = new Uint8Array(buffer)
    for (let i = 0; i < view.length; i += 1) {
      hash = ((hash << 5) - hash) + view[i]
      hash |= 0
    }
    return `fallback-${file.name}-${file.size}-${hash}`
  } catch (error) {
    console.error('[API] Error generating file hash:', error)
    return `error-${file.name}-${file.size}-${Date.now()}`
  }
}

export const getSessionId = () => sessionId as string

export const getActiveBatchStorageKey = () => `${ACTIVE_BATCH_STORAGE_KEY}:${getSessionId()}`

export const getStoredActiveBatchId = (): string | null => localStorage.getItem(getActiveBatchStorageKey())

export const setStoredActiveBatchId = (batchId: string | null) => {
  const key = getActiveBatchStorageKey()
  if (!batchId) {
    localStorage.removeItem(key)
    return
  }
  localStorage.setItem(key, batchId)
}

export const getActiveBatchViewStorageKey = () => `${ACTIVE_BATCH_VIEW_STORAGE_KEY}:${getSessionId()}`

export const getStoredActiveBatchView = (): string | null => localStorage.getItem(getActiveBatchViewStorageKey())

export const setStoredActiveBatchView = (view: string | null) => {
  const key = getActiveBatchViewStorageKey()
  if (!view) {
    localStorage.removeItem(key)
    return
  }
  localStorage.setItem(key, view)
}

export const analyzeCV = async (
  jd: string,
  cvFile: File,
  apiKey: string,
  jobID?: string,
): Promise<APIResponse<AnalysisResult>> => {
  const fileHash = await getFileHash(cvFile)
  const formData = new FormData()
  formData.append('job_description', jd)
  formData.append('cv_file', cvFile)
  formData.append('api_key', apiKey)
  if (jobID) formData.append('job_id', jobID)

  const response = await api.post<APIResponse<AnalysisResult>>('/analyze', formData, {
    headers: {
      'Content-Type': 'multipart/form-data',
      'X-File-Hash': fileHash,
    },
  })

  return response.data
}

export const fetchHistory = async (jobID?: string, limit = 50): Promise<APIResponse<AnalysisResult[]>> => {
  const response = await api.get<APIResponse<AnalysisResult[]>>('/history', {
    params: { job_id: jobID, limit },
  })
  return response.data
}

export const deleteHistoryItem = async (id: string): Promise<APIResponse<{ message: string }>> => {
  const response = await api.delete<APIResponse<{ message: string }>>(`/history/${id}`)
  return response.data
}

export const checkHealth = async () => {
  const response = await api.get('/health')
  return response.data
}

export const getBatchStatus = async (batchId: string): Promise<APIResponse<BatchStatusData>> => {
  const response = await api.get<APIResponse<BatchStatusData>>(`/jobs/${batchId}`)
  return response.data
}

export const getBatchNotification = async (batchId: string): Promise<APIResponse<BatchNotificationData>> => {
  const response = await api.get<APIResponse<BatchNotificationData>>(`/jobs/${batchId}/notification`)
  return response.data
}

export const getBatchResults = async (batchId: string): Promise<APIResponse<BatchResultsData>> => {
  const response = await api.get<APIResponse<BatchResultsData>>(`/jobs/${batchId}/results`)
  return response.data
}

export const analyzeBulk = async (
  jd: string,
  cvFiles: File[],
  apiKey: string,
  jobID?: string,
): Promise<APIResponse<BatchAnalyzeData>> => {
  const formData = new FormData()
  formData.append('job_description', jd)
  formData.append('api_key', apiKey)
  if (jobID) formData.append('job_id', jobID)

  cvFiles.forEach((file) => {
    formData.append('cv_files', file)
  })

  const response = await api.post<APIResponse<BatchAnalyzeData>>('/analyze/bulk', formData, {
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  })

  return response.data
}

export default api
