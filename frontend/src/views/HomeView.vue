<template>
  <div class="min-h-screen bg-[#f8fafc]">
    <!-- Sidebar -->
    <aside class="fixed left-0 top-0 h-full w-72 bg-white border-r border-navy-50 p-8 hidden lg:block z-20">
      <div class="flex flex-col h-full">
        <div class="mb-12 px-2">
          <h1 class="text-2xl font-black text-navy-900 tracking-tighter leading-none">Recruiter Command</h1>
          <p class="text-[10px] text-navy-400 font-black uppercase tracking-[0.2em] mt-2">Data-Driven Hiring</p>
        </div>

        <nav class="space-y-4 flex-1 flex flex-col overflow-hidden">
          <div class="px-2">
            <a href="#" @click.prevent="reset()" 
              :class="view === 'input' ? 'bg-navy-900 text-white shadow-xl shadow-navy-900/10' : 'text-navy-400 hover:bg-navy-50 hover:text-navy-900'" 
              class="flex items-center space-x-4 px-5 py-3.5 rounded-xl text-[13px] font-black transition-all duration-300 group"
            >
              <LayoutIcon class="w-5 h-5 transition-transform group-hover:scale-110" />
              <span class="uppercase tracking-widest">New Evaluation</span>
            </a>
          </div>

          <!-- History Section -->
          <div class="mt-6 flex-1 overflow-y-auto pr-2 px-2 pb-4">
            <div class="flex items-center space-x-2 text-navy-300 px-3 mb-4 uppercase tracking-[0.1em] text-[10px] font-black">
              <HistoryIcon class="w-3.5 h-3.5" />
              <span>Past Evaluations</span>
            </div>

            <div class="space-y-2">
              <div v-for="item in historyList" :key="item.processing_metadata.request_id"
                @click="loadPastAnalysis(item)"
                class="group relative flex items-center justify-between p-3 rounded-xl border transition-all cursor-pointer"
                :class="result?.processing_metadata.request_id === item.processing_metadata.request_id 
                  ? 'bg-blue-50 border-blue-200 shadow-sm' 
                  : 'bg-white border-navy-50 hover:border-navy-200 hover:shadow-sm'"
              >
                <div class="flex-1 min-w-0 pr-2">
                  <p class="text-xs font-bold text-navy-900 truncate" :title="item.processing_metadata.cv_file_name || 'candidate_cv.pdf'">
                    {{ item.processing_metadata.cv_file_name || 'candidate_cv.pdf' }}
                  </p>
                  <p class="text-[10px] text-navy-400 mt-1 uppercase tracking-wider font-bold">
                    Score: <span class="text-navy-800">{{ item.matching_score.overall.toFixed(1) }}/100</span>
                  </p>
                </div>
                <button @click="(e) => handleDeleteHistory(item.processing_metadata.request_id, e)" 
                  class="opacity-0 group-hover:opacity-100 p-2 text-red-400 hover:text-red-600 hover:bg-red-50 rounded-lg transition-all"
                  title="Delete record"
                >
                  <TrashIcon class="w-4 h-4" />
                </button>
              </div>
              
              <div v-if="historyList.length === 0" class="text-center py-8">
                <p class="text-xs text-navy-300 font-bold uppercase tracking-wider">No history found</p>
              </div>
            </div>
          </div>
        </nav>

      </div>
    </aside>

    <!-- Main Content -->
    <main class="lg:ml-72 p-10 min-h-screen">
      <!-- Header -->
      <header class="flex items-center justify-between mb-12">
        <div>
          <h2 class="text-3xl font-black text-navy-900 tracking-tight">{{ view === 'input' ? 'Evaluation Dashboard' : 'Analysis Results' }}</h2>
          <div class="flex items-center space-x-2 mt-1">
            <span class="text-[11px] font-bold text-navy-300 uppercase tracking-widest">Recruiter ID: #RE-2024-X</span>
            <span v-if="result" class="text-[11px] font-bold text-navy-300 uppercase tracking-widest border-l border-navy-100 pl-2">Candidate: {{ result.processing_metadata.request_id }}</span>
          </div>
        </div>

        <div class="flex items-center space-x-6">
          <button class="px-6 py-2.5 bg-navy-900 text-white rounded-xl font-black text-xs uppercase tracking-widest shadow-xl shadow-navy-900/10 hover:scale-[1.02] active:scale-95 transition-all" @click="reset">New Evaluation</button>
        </div>
      </header>

      <!-- Input Section -->
      <div v-if="view === 'input'" class="max-w-4xl mx-auto">
        <InputPanel :loading="loading" @start="handleStart" />
      </div>

      <!-- Results Section -->
      <div v-else-if="view === 'results' && result">
        <ResultsPanel :result="result" />
      </div>

      <!-- Error State -->
      <div v-if="error" class="mt-8 p-4 bg-red-50 border border-red-100 rounded-xl flex items-center space-x-3 text-red-700">
        <AlertCircleIcon class="w-5 h-5" />
        <span>{{ error }}</span>
        <button @click="error = null" class="ml-auto text-sm font-bold uppercase tracking-widest">Dismiss</button>
      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { 
  Layout as LayoutIcon,
  AlertCircle as AlertCircleIcon,
  Trash2 as TrashIcon,
  History as HistoryIcon
} from 'lucide-vue-next'
import InputPanel from '../components/InputPanel.vue'
import ResultsPanel from '../components/ResultsPanel.vue'
import { analyzeCV, fetchHistory, deleteHistoryItem } from '../services/api'
import type { AnalysisResult } from '../types/analysis'

const view = ref<'input' | 'results'>('input')
const loading = ref(false)
const result = ref<AnalysisResult | null>(null)
const error = ref<string | null>(null)
const historyList = ref<AnalysisResult[]>([])

const loadHistory = async () => {
  try {
    const res = await fetchHistory()
    if (res.status === 'success' && res.data) {
      historyList.value = res.data.sort((a, b) => 
        new Date(b.processing_metadata.timestamp).getTime() - new Date(a.processing_metadata.timestamp).getTime()
      )
    }
  } catch (err) {
    console.error('Failed to load history', err)
  }
}

onMounted(() => {
  loadHistory()
})

const handleStart = async (data: { apiKey: string, jd: string, cvFile: File }) => {
  loading.value = true
  error.value = null
  try {
    const response = await analyzeCV(data.jd, data.cvFile, data.apiKey)
    if (response.status === 'success' && response.data) {
      result.value = response.data
      view.value = 'results'
      await loadHistory() // refresh sidebar
    } else {
      error.value = response.error?.message || 'Failed to analyze CV'
    }
  } catch (err: any) {
    error.value = err.response?.data?.error?.message || 'An unexpected error occurred during analysis'
    console.error(err)
  } finally {
    loading.value = false
  }
}

const loadPastAnalysis = (item: AnalysisResult) => {
  result.value = item
  view.value = 'results'
}

const handleDeleteHistory = async (id: string, e: Event) => {
  e.stopPropagation()
  try {
    await deleteHistoryItem(id)
    if (result.value && result.value.processing_metadata.request_id === id) {
      reset()
    }
    await loadHistory()
  } catch (err) {
    console.error('Failed to delete history item', err)
  }
}

const reset = () => {
  view.value = 'input'
  result.value = null
  error.value = null
}

</script>
