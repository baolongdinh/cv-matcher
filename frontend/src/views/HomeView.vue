<template>
  <div class="min-h-screen bg-[#f8fafc]">
    <aside class="fixed left-0 top-0 h-full w-72 bg-white border-r border-navy-50 p-8 hidden lg:block z-20">
      <div class="flex flex-col h-full">
        <div class="mb-12 px-2">
          <h1 class="text-2xl font-black text-navy-900 tracking-tighter leading-none">Recruiter Command</h1>
          <p class="text-[10px] text-navy-400 font-black uppercase tracking-[0.2em] mt-2">Data-Driven Hiring</p>
        </div>

        <nav class="space-y-4 flex-1 flex flex-col overflow-hidden">
          <div class="px-2">
            <a
              href="#"
              @click.prevent="openNewEvaluation"
              :class="view === 'input' ? 'bg-navy-900 text-white shadow-xl shadow-navy-900/10' : 'text-navy-400 hover:bg-navy-50 hover:text-navy-900'"
              class="flex items-center space-x-4 px-5 py-3.5 rounded-xl text-[13px] font-black transition-all duration-300 group"
            >
              <LayoutIcon class="w-5 h-5 transition-transform group-hover:scale-110" />
              <span class="uppercase tracking-widest">New Evaluation</span>
            </a>
          </div>

          <div class="px-2">
            <a
              href="#"
              @click.prevent="goToLiveBatch"
              :class="view === 'live-batch' ? 'bg-navy-900 text-white shadow-xl shadow-navy-900/10' : 'text-navy-400 hover:bg-navy-50 hover:text-navy-900'"
              class="flex items-center space-x-4 px-5 py-3.5 rounded-xl text-[13px] font-black transition-all duration-300 group"
            >
              <RefreshCcwIcon class="w-5 h-5 transition-transform group-hover:scale-110" />
              <span class="uppercase tracking-widest">Live Batch</span>
            </a>
          </div>

          <div class="px-2">
            <a
              href="#"
              @click.prevent="handleViewHistory"
              :class="view === 'history' ? 'bg-navy-900 text-white shadow-xl shadow-navy-900/10' : 'text-navy-400 hover:bg-navy-50 hover:text-navy-900'"
              class="flex items-center space-x-4 px-5 py-3.5 rounded-xl text-[13px] font-black transition-all duration-300 group"
            >
              <HistoryIcon class="w-5 h-5 transition-transform group-hover:scale-110" />
              <span class="uppercase tracking-widest">Candidate History</span>
            </a>
          </div>

          <div class="mt-6 flex-1 overflow-y-auto pr-2 px-2 pb-4">
            <div class="flex items-center justify-between px-3 mb-4">
              <div class="flex items-center space-x-2 text-navy-300 uppercase tracking-[0.1em] text-[10px] font-black">
                <HistoryIcon class="w-3.5 h-3.5" />
                <span>Past Evaluations</span>
              </div>
              <span v-if="historyItems.length > 0" class="text-[10px] font-black uppercase tracking-widest text-navy-300">
                {{ historyItems.length }}
              </span>
            </div>

            <div class="space-y-2">
              <div
                v-for="item in historyItems"
                :key="item.processing_metadata.request_id"
                @click="loadPastAnalysis(item)"
                class="group relative flex items-center justify-between p-3 rounded-xl border transition-all cursor-pointer"
                :class="currentResult?.processing_metadata.request_id === item.processing_metadata.request_id && view === 'candidate-detail'
                  ? 'bg-blue-50 border-blue-200 shadow-sm'
                  : 'bg-white border-navy-50 hover:border-navy-200 hover:shadow-sm'"
              >
                <div class="flex-1 min-w-0 pr-2">
                  <p class="text-xs font-bold text-navy-900 truncate" :title="item.processing_metadata.cv_file_name || 'candidate_cv.pdf'">
                    {{ item.processing_metadata.cv_file_name || 'candidate_cv.pdf' }}
                  </p>
                  <p class="text-[10px] text-navy-400 mt-1 uppercase tracking-wider font-bold">
                    Score: <span class="text-navy-800">{{ item.matching_score.overall.toFixed(1) }}/10</span>
                  </p>
                </div>
                <button
                  @click="(event) => handleDeleteHistory(item.processing_metadata.request_id, event)"
                  class="opacity-0 group-hover:opacity-100 p-2 text-red-400 hover:text-red-600 hover:bg-red-50 rounded-lg transition-all"
                  title="Delete record"
                >
                  <TrashIcon class="w-4 h-4" />
                </button>
              </div>

              <div v-if="historyItems.length === 0" class="text-center py-8">
                <p class="text-xs text-navy-300 font-bold uppercase tracking-wider">No history found</p>
              </div>
            </div>
          </div>
        </nav>
      </div>
    </aside>

    <main class="lg:ml-72 p-6 sm:p-8 lg:p-10 min-h-screen">
      <header class="flex flex-col gap-4 sm:flex-row sm:items-start sm:justify-between mb-10">
        <div>
          <h2 class="text-3xl font-black text-navy-900 tracking-tight">{{ pageTitle }}</h2>
          <div class="flex flex-wrap items-center gap-2 mt-1">
            <span class="text-[11px] font-bold text-navy-300 uppercase tracking-widest">Recruiter ID: #RE-2024-X</span>
            <span v-if="currentResult && view === 'candidate-detail'" class="text-[11px] font-bold text-navy-300 uppercase tracking-widest border-l border-navy-100 pl-2">
              Candidate: {{ currentResult.processing_metadata.request_id }}
            </span>
            <span v-if="liveBatch.batch_id" class="text-[11px] font-bold text-navy-300 uppercase tracking-widest border-l border-navy-100 pl-2">
              Batch: {{ liveBatch.batch_id.slice(-8) }}
            </span>
          </div>
        </div>

        <div class="flex items-center gap-3">
          <button
            v-if="hasActiveOrRecentBatch && view !== 'live-batch'"
            class="px-5 py-2.5 border border-navy-200 bg-white text-navy-900 rounded-xl font-black text-xs uppercase tracking-widest shadow-sm hover:bg-navy-50 transition-all"
            @click="goToLiveBatch"
          >
            Resume Live Batch
          </button>
          <button
            class="px-6 py-2.5 bg-navy-900 text-white rounded-xl font-black text-xs uppercase tracking-widest shadow-xl shadow-navy-900/10 hover:scale-[1.02] active:scale-95 transition-all"
            @click="openNewEvaluation"
          >
            New Evaluation
          </button>
        </div>
      </header>

      <div
        v-if="showPinnedLiveBanner"
        class="mb-8 p-5 bg-navy-900 rounded-3xl border border-navy-700 shadow-2xl flex flex-col gap-4 lg:flex-row lg:items-center lg:justify-between"
      >
        <div class="flex items-center gap-5">
          <div class="w-12 h-12 bg-primary-500/20 rounded-2xl flex items-center justify-center relative">
            <RefreshCcwIcon class="w-6 h-6 text-primary-400" :class="liveBatch.is_reconnecting ? '' : 'animate-spin'" />
          </div>
          <div>
            <h4 class="text-white font-black text-sm uppercase tracking-widest">
              {{ liveBatch.is_reconnecting ? 'Reconnecting Live Batch' : 'Live Batch Still Running' }}
            </h4>
            <p class="text-navy-400 text-xs font-bold mt-1 uppercase tracking-wider">
              {{ liveBatch.completed }} / {{ liveBatch.total }} processed • {{ liveBatch.successful }} ready • {{ liveBatch.failed }} failed
            </p>
          </div>
        </div>
        <button
          class="px-5 py-3 bg-white text-navy-950 rounded-xl font-black text-[11px] uppercase tracking-widest hover:scale-[1.02] active:scale-95 transition-all"
          @click="goToLiveBatch"
        >
          Open Live Dashboard
        </button>
      </div>

      <div v-if="view === 'input'" class="max-w-4xl mx-auto">
        <InputPanel :loading="loading" @start="handleStart" />
      </div>

      <div v-else-if="view === 'live-batch'">
        <LiveBatchPanel
          :live-batch="liveBatch"
          :selected-candidate-id="selectedCandidateId"
          @view="handleViewResult"
        />
      </div>

      <div v-else-if="view === 'candidate-detail' && currentResult">
        <ResultsPanel :result="currentResult" />
      </div>

      <div v-else-if="view === 'history'">
        <HistoryPanel :history="historyItems" @view="handleViewResult" />
      </div>

      <div v-if="error" class="mt-8 p-4 bg-red-50 border border-red-100 rounded-xl flex items-center space-x-3 text-red-700">
        <AlertCircleIcon class="w-5 h-5" />
        <span>{{ error }}</span>
        <button @click="error = null" class="ml-auto text-sm font-bold uppercase tracking-widest">Dismiss</button>
      </div>

      <div
        v-if="notificationBanner"
        class="mt-8 p-4 rounded-xl border flex items-center space-x-3"
        :class="notificationBanner.type === 'success'
          ? 'bg-emerald-50 border-emerald-100 text-emerald-700'
          : notificationBanner.type === 'warning'
            ? 'bg-amber-50 border-amber-100 text-amber-700'
            : 'bg-blue-50 border-blue-100 text-blue-700'"
      >
        <RefreshCcwIcon class="w-5 h-5" :class="notificationBanner.type === 'info' ? 'animate-spin' : ''" />
        <span>{{ notificationBanner.message }}</span>
        <button @click="notificationBanner = null" class="ml-auto text-sm font-bold uppercase tracking-widest">Dismiss</button>
      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref, watch } from 'vue'
import {
  Layout as LayoutIcon,
  AlertCircle as AlertCircleIcon,
  Trash2 as TrashIcon,
  History as HistoryIcon,
  RefreshCcw as RefreshCcwIcon,
} from 'lucide-vue-next'
import InputPanel from '../components/InputPanel.vue'
import ResultsPanel from '../components/ResultsPanel.vue'
import HistoryPanel from '../components/HistoryPanel.vue'
import LiveBatchPanel from '../components/LiveBatchPanel.vue'
import { analyzeBulk, analyzeCV } from '../services/api'
import { useAnalysisResult } from '../composables/useAnalysisResult'
import { useBatchLiveSession } from '../composables/useBatchLiveSession'
import { useSessionHistory } from '../composables/useSessionHistory'
import type { AnalysisResult, LiveCandidateCard } from '../types/analysis'

type ViewMode = 'input' | 'live-batch' | 'candidate-detail' | 'history'
type BannerState = {
  type: 'success' | 'warning' | 'info'
  message: string
}

const view = ref<ViewMode>('input')
const loading = ref(false)
const error = ref<string | null>(null)
const notificationBanner = ref<BannerState | null>(null)

const { currentResult, setResult } = useAnalysisResult()
const { historyItems, loadHistory, removeHistoryItem, historyError } = useSessionHistory()
const {
  activeBatchId,
  batchError,
  selectedCandidateId,
  viewModel,
  start,
  resume,
  stop,
  reset,
  selectCandidate,
} = useBatchLiveSession()

const liveBatch = computed(() => viewModel.value)
const hasActiveOrRecentBatch = computed(() => !!liveBatch.value.batch_id)
const showPinnedLiveBanner = computed(() => !!activeBatchId.value && view.value !== 'live-batch')

const pageTitle = computed(() => {
  if (view.value === 'input') return 'Evaluation Dashboard'
  if (view.value === 'live-batch') return 'Live Batch Dashboard'
  if (view.value === 'history') return 'Analysis History'
  return 'Candidate Report'
})

const announceFirstCandidate = async (candidate: LiveCandidateCard) => {
  notificationBanner.value = {
    type: 'success',
    message: `${candidate.processing_metadata.cv_file_name || 'A candidate'} is ready for review while the rest of the batch keeps processing.`,
  }
}

const onBatchCompleted = async () => {
  notificationBanner.value = {
    type: liveBatch.value.failed > 0 ? 'warning' : 'success',
    message: liveBatch.value.failed > 0
      ? `Batch completed with ${liveBatch.value.successful} successful and ${liveBatch.value.failed} failed results.`
      : `Batch completed successfully with ${liveBatch.value.successful} candidates ready for review.`,
  }
  await loadHistory()
}

const startLiveBatch = async (batchId: string, restored = false) => {
  view.value = 'live-batch'
  await start(batchId, {
    restored,
    onFirstCandidate: announceFirstCandidate,
    onCompleted: onBatchCompleted,
  })
}

const openNewEvaluation = () => {
  setResult(null)
  if (activeBatchId.value) {
    notificationBanner.value = {
      type: 'info',
      message: 'The live batch is still running in the background. You can return to it at any time from the banner or sidebar.',
    }
  }
  view.value = 'input'
}

const goToLiveBatch = () => {
  if (!hasActiveOrRecentBatch.value) {
    notificationBanner.value = {
      type: 'info',
      message: 'No active batch is available right now.',
    }
    return
  }
  view.value = 'live-batch'
}

const handleViewHistory = () => {
  void loadHistory()
  view.value = 'history'
}

const handleViewResult = (result: AnalysisResult) => {
  selectCandidate(result.processing_metadata.request_id)
  setResult(result)
  view.value = 'candidate-detail'
}

onMounted(() => {
  void loadHistory()
  void resume({
    onFirstCandidate: announceFirstCandidate,
    onCompleted: onBatchCompleted,
    restored: true,
  }).then((restored) => {
    if (restored) {
      view.value = 'live-batch'
      notificationBanner.value = {
        type: 'info',
        message: 'Restored an in-progress batch for this session.',
      }
    }
  })
})

onUnmounted(() => {
  stop()
})

const handleStart = async (data: { apiKey: string, jd: string, cvFiles: File[], jobId?: string }) => {
  loading.value = true
  error.value = null
  try {
    if (data.cvFiles.length > 1) {
      const response = await analyzeBulk(data.jd, data.cvFiles, data.apiKey, data.jobId)
      if (response.status === 'error') {
        throw new Error(response.error.message || 'Failed to start batch analysis')
      }
      notificationBanner.value = {
        type: 'info',
        message: `Bulk batch queued. ${response.data.total_files} CVs are being processed live.`,
      }
      await startLiveBatch(response.data.batch_id)
      return
    }

    const response = await analyzeCV(data.jd, data.cvFiles[0], data.apiKey, data.jobId)
    if (response.status === 'error') {
      throw new Error(response.error.message || 'Failed to analyze CV')
    }
    setResult(response.data)
    view.value = 'candidate-detail'
    await loadHistory()
  } catch (err: any) {
    error.value = err.response?.data?.error?.message || err.message || 'An unexpected error occurred'
  } finally {
    loading.value = false
  }
}

const loadPastAnalysis = (item: AnalysisResult) => {
  setResult(item)
  view.value = 'candidate-detail'
}

const handleDeleteHistory = async (id: string, event: Event) => {
  event.stopPropagation()
  try {
    const deleted = await removeHistoryItem(id)
    if (!deleted) {
      return
    }
    if (currentResult.value?.processing_metadata.request_id === id) {
      setResult(null)
      view.value = hasActiveOrRecentBatch.value ? 'live-batch' : 'history'
    }
  } catch (err) {
    console.error('Failed to delete history item', err)
  }
}

watch([historyError, batchError], () => {
  if (historyError.value) {
    error.value = historyError.value.message
  } else if (batchError.value) {
    error.value = batchError.value.message
  }
})

watch(() => liveBatch.value.is_reconnecting, (reconnecting) => {
  if (reconnecting) {
    notificationBanner.value = {
      type: 'warning',
      message: 'Connection is unstable. The live dashboard is retrying automatically without dropping current results.',
    }
  }
})

watch(() => activeBatchId.value, (nextBatchId, previousBatchId) => {
  if (!nextBatchId && previousBatchId && view.value === 'input' && hasActiveOrRecentBatch.value) {
    notificationBanner.value = {
      type: 'success',
      message: 'The background batch finished. You can review the full ranked list from the live dashboard.',
    }
  }
})

const clearLiveBatch = () => {
  reset()
}

defineExpose({ clearLiveBatch })
</script>
