import { computed, ref } from 'vue'
import {
  getBatchResults,
  getBatchStatus,
  getStoredActiveBatchId,
  setStoredActiveBatchId,
} from '../services/api'
import type {
  APIError,
  BatchFailureViewModel,
  BatchResultsData,
  BatchStatusData,
  LiveBatchViewModel,
  LiveCandidateCard,
  QueueItemViewModel,
} from '../types/analysis'

const activeBatchId = ref<string | null>(getStoredActiveBatchId())
const batchMeta = ref<BatchStatusData | null>(null)
const completedCandidates = ref<LiveCandidateCard[]>([])
const failedItems = ref<BatchFailureViewModel[]>([])
const batchError = ref<APIError | null>(null)
const selectedCandidateId = ref<string | null>(null)
const lastSeenCompletedCount = ref(0)
const isRestoredFromStorage = ref(false)
const isReconnecting = ref(false)

let pollTimer: number | null = null
let retryDelayMs = 1800

const clearPollTimer = () => {
  if (pollTimer !== null) {
    window.clearTimeout(pollTimer)
    pollTimer = null
  }
}

const formatUpdatedLabel = (status: string) => {
  if (status === 'processing') return 'Analyzing now'
  if (status === 'completed') return 'Ready to review'
  if (status === 'failed') return 'Needs attention'
  return 'Waiting in queue'
}

const buildQueueItems = (statusData: BatchStatusData, failedMap: Map<string, BatchFailureViewModel>): QueueItemViewModel[] => {
  return statusData.items
    .map((item) => ({
      ...item,
      error_code: item.error_code || failedMap.get(item.request_id)?.error_code,
      error_message: item.error_message || failedMap.get(item.request_id)?.error_message,
      updated_label: formatUpdatedLabel(item.status),
    }))
    .sort((left, right) => {
      const order = { processing: 0, queued: 1, failed: 2, completed: 3 } as Record<string, number>
      return (order[left.status] ?? 99) - (order[right.status] ?? 99)
    })
}

const mergeCandidates = (results: BatchResultsData) => {
  const existing = new Map(completedCandidates.value.map((candidate) => [candidate.processing_metadata.request_id, candidate]))
  results.candidates.forEach((candidate) => {
    const requestId = candidate.processing_metadata.request_id
    const existingCandidate = existing.get(requestId)
    const matchingItem = batchMeta.value?.items.find((item) => item.request_id === requestId)
    existing.set(requestId, {
      ...existingCandidate,
      ...candidate,
      is_cached: matchingItem?.cached ?? existingCandidate?.is_cached ?? false,
    })
  })
  completedCandidates.value = Array.from(existing.values()).sort((left, right) => right.matching_score.overall - left.matching_score.overall)

  const failures = new Map(failedItems.value.map((failure) => [failure.request_id, failure]))
  results.failed.forEach((failure) => {
    failures.set(failure.request_id, {
      ...failure,
      updated_label: 'Needs attention',
    })
  })
  failedItems.value = Array.from(failures.values())
}

const scheduleNextPoll = (runner: () => Promise<void>, delayMs: number) => {
  clearPollTimer()
  pollTimer = window.setTimeout(() => {
    void runner()
  }, delayMs)
}

export const useBatchLiveSession = () => {
  const hydrateResults = async (batchId: string) => {
    const response = await getBatchResults(batchId)
    if (response.status === 'error') {
      batchError.value = response.error
      return null
    }
    mergeCandidates(response.data)
    return response.data
  }

  const refreshResultsIfProgressAdvanced = async (batchId: string, currentCompleted: number) => {
    if (currentCompleted > lastSeenCompletedCount.value || completedCandidates.value.length === 0) {
      const results = await hydrateResults(batchId)
      if (results) {
        lastSeenCompletedCount.value = currentCompleted
      }
      return results
    }
    return null
  }

  const stop = () => {
    clearPollTimer()
  }

  const clearStoredBatch = () => {
    setStoredActiveBatchId(null)
    activeBatchId.value = null
  }

  const reset = () => {
    stop()
    clearStoredBatch()
    batchMeta.value = null
    completedCandidates.value = []
    failedItems.value = []
    batchError.value = null
    selectedCandidateId.value = null
    lastSeenCompletedCount.value = 0
    isRestoredFromStorage.value = false
    isReconnecting.value = false
    retryDelayMs = 1800
  }

  const selectCandidate = (requestId: string | null) => {
    selectedCandidateId.value = requestId
  }

  const start = async (
    batchId: string,
    options?: {
      restored?: boolean
      onFirstCandidate?: (candidate: LiveCandidateCard) => void | Promise<void>
      onCompleted?: (viewModel: LiveBatchViewModel) => void | Promise<void>
    },
  ) => {
    stop()
    activeBatchId.value = batchId
    setStoredActiveBatchId(batchId)
    batchError.value = null
    batchMeta.value = null
    completedCandidates.value = []
    failedItems.value = []
    selectedCandidateId.value = null
    lastSeenCompletedCount.value = 0
    isRestoredFromStorage.value = !!options?.restored
    isReconnecting.value = false
    retryDelayMs = 1800

    let firstCandidateAnnounced = false

    const poll = async () => {
      if (!activeBatchId.value) {
        return
      }

      const statusResponse = await getBatchStatus(batchId)
      if (statusResponse.status === 'error') {
        batchError.value = statusResponse.error
        isReconnecting.value = true
        retryDelayMs = Math.min(retryDelayMs + 1200, 6000)
        scheduleNextPoll(poll, retryDelayMs)
        return
      }

      batchError.value = null
      isReconnecting.value = false
      retryDelayMs = 1800
      batchMeta.value = statusResponse.data

      const previousCompletedCount = completedCandidates.value.length
      await refreshResultsIfProgressAdvanced(batchId, statusResponse.data.completed)

      if (!firstCandidateAnnounced && completedCandidates.value.length > 0) {
        firstCandidateAnnounced = true
        const firstCandidate = completedCandidates.value[0]
        selectCandidate(firstCandidate.processing_metadata.request_id)
        if (options?.onFirstCandidate) {
          await options.onFirstCandidate(firstCandidate)
        }
      } else if (completedCandidates.value.length > previousCompletedCount && !selectedCandidateId.value) {
        selectCandidate(completedCandidates.value[0]?.processing_metadata.request_id ?? null)
      }

      if (statusResponse.data.completed >= statusResponse.data.total && statusResponse.data.total > 0) {
        await hydrateResults(batchId)
        clearStoredBatch()
        if (options?.onCompleted) {
          await options.onCompleted(viewModel.value)
        }
        return
      }

      scheduleNextPoll(poll, 1800)
    }

    await poll()
    return viewModel.value
  }

  const resume = async (options?: Parameters<typeof start>[1]) => {
    const storedBatchId = getStoredActiveBatchId()
    if (!storedBatchId) {
      return false
    }
    await start(storedBatchId, { ...options, restored: true })
    return true
  }

  const processingItems = computed(() => {
    if (!batchMeta.value) {
      return [] as QueueItemViewModel[]
    }
    const failedMap = new Map(failedItems.value.map((failure) => [failure.request_id, failure]))
    return buildQueueItems(batchMeta.value, failedMap)
  })

  const viewModel = computed<LiveBatchViewModel>(() => ({
    batch_id: batchMeta.value?.batch_id || activeBatchId.value || '',
    status: batchMeta.value?.status || 'queued',
    job_id: batchMeta.value?.job_id,
    total: batchMeta.value?.total || 0,
    completed: batchMeta.value?.completed || 0,
    successful: batchMeta.value?.successful || 0,
    failed: batchMeta.value?.failed || 0,
    completed_candidates: completedCandidates.value,
    processing_items: processingItems.value,
    failed_items: failedItems.value,
    last_seen_completed_count: lastSeenCompletedCount.value,
    is_restored_from_storage: isRestoredFromStorage.value,
    is_reconnecting: isReconnecting.value,
  }))

  return {
    activeBatchId: computed(() => activeBatchId.value),
    batchMeta: computed(() => batchMeta.value),
    completedCandidates: computed(() => completedCandidates.value),
    failedItems: computed(() => failedItems.value),
    processingItems,
    batchError: computed(() => batchError.value),
    selectedCandidateId: computed(() => selectedCandidateId.value),
    isRestoredFromStorage: computed(() => isRestoredFromStorage.value),
    isReconnecting: computed(() => isReconnecting.value),
    viewModel,
    start,
    resume,
    stop,
    reset,
    selectCandidate,
    refreshResultsIfProgressAdvanced,
  }
}
