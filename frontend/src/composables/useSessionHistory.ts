import { computed, ref } from 'vue'
import { deleteHistoryItem, fetchHistory } from '../services/api'
import type { AnalysisResult, APIError } from '../types/analysis'

const historyItems = ref<AnalysisResult[]>([])
const loadingHistory = ref(false)
const historyError = ref<APIError | null>(null)

const sortHistory = (items: AnalysisResult[]) => {
  return [...items].sort((a, b) =>
    new Date(b.processing_metadata.timestamp).getTime() - new Date(a.processing_metadata.timestamp).getTime(),
  )
}

export const useSessionHistory = () => {
  const loadHistory = async (jobID?: string) => {
    loadingHistory.value = true
    historyError.value = null
    try {
      const response = await fetchHistory(jobID)
      if (response.status === 'error') {
        historyError.value = response.error
        return
      }
      historyItems.value = sortHistory(response.data)
    } finally {
      loadingHistory.value = false
    }
  }

  const removeHistoryItem = async (id: string) => {
    historyError.value = null
    const response = await deleteHistoryItem(id)
    if (response.status === 'error') {
      historyError.value = response.error
      return false
    }
    historyItems.value = historyItems.value.filter((item) => item.processing_metadata.request_id !== id)
    return true
  }

  return {
    historyItems: computed(() => historyItems.value),
    loadingHistory: computed(() => loadingHistory.value),
    historyError: computed(() => historyError.value),
    loadHistory,
    removeHistoryItem,
  }
}
