import { computed, ref } from 'vue'
import {
  getBatchNotification,
  getBatchResults,
  getStoredActiveBatchId,
  setStoredActiveBatchId,
} from '../services/api'
import type {
  APIError,
  BatchNotificationData,
  BatchResultsData,
} from '../types/analysis'

const activeBatchId = ref<string | null>(getStoredActiveBatchId())
const batchNotification = ref<BatchNotificationData | null>(null)
const batchResults = ref<BatchResultsData | null>(null)
const batchError = ref<APIError | null>(null)
let pollTimer: number | null = null

const clearPollTimer = () => {
  if (pollTimer !== null) {
    window.clearInterval(pollTimer)
    pollTimer = null
  }
}

export const useBatchAnalysis = () => {
  const stopPolling = () => {
    clearPollTimer()
  }

  const resetBatchState = () => {
    stopPolling()
    activeBatchId.value = null
    batchNotification.value = null
    batchResults.value = null
    batchError.value = null
    setStoredActiveBatchId(null)
  }

  const loadBatchResults = async (batchId: string) => {
    const response = await getBatchResults(batchId)
    if (response.status === 'error') {
      batchError.value = response.error
      return null
    }
    batchResults.value = response.data
    return response.data
  }

  const refreshBatchNotification = async (batchId: string) => {
    const response = await getBatchNotification(batchId)
    if (response.status === 'error') {
      batchError.value = response.error
      return null
    }
    batchNotification.value = response.data
    return response.data
  }

  const startPolling = async (
    batchId: string,
    onCompleted?: (results: BatchResultsData, notification: BatchNotificationData) => Promise<void> | void,
  ) => {
    stopPolling()
    activeBatchId.value = batchId
    batchError.value = null
    batchResults.value = null
    setStoredActiveBatchId(batchId)

    const tick = async () => {
      const notification = await refreshBatchNotification(batchId)
      if (!notification) {
        stopPolling()
        return
      }

      if (notification.complete) {
        stopPolling()
        const results = await loadBatchResults(batchId)
        if (results && onCompleted) {
          await onCompleted(results, notification)
        }
        setStoredActiveBatchId(null)
        activeBatchId.value = null
      }
    }

    await tick()
    if (!activeBatchId.value) {
      return
    }

    pollTimer = window.setInterval(() => {
      void tick()
    }, 3000)
  }

  const resumePolling = async (
    onCompleted?: (results: BatchResultsData, notification: BatchNotificationData) => Promise<void> | void,
  ) => {
    const storedBatchId = getStoredActiveBatchId()
    if (!storedBatchId) {
      return false
    }
    await startPolling(storedBatchId, onCompleted)
    return true
  }

  return {
    activeBatchId: computed(() => activeBatchId.value),
    batchNotification: computed(() => batchNotification.value),
    batchResults: computed(() => batchResults.value),
    batchError: computed(() => batchError.value),
    startPolling,
    stopPolling,
    resumePolling,
    loadBatchResults,
    resetBatchState,
  }
}
