<template>
  <div class="space-y-8 animate-in fade-in slide-in-from-bottom-4 duration-700">
    <!-- Batch Summary Header -->
    <div class="bg-white p-8 rounded-3xl border border-navy-100 shadow-[0_8px_30px_rgb(0,0,0,0.04)]">
      <div class="flex items-center justify-between mb-6">
        <div class="flex items-center space-x-3">
          <span class="px-3 py-1 rounded-full bg-emerald-100 text-emerald-700 text-[10px] font-black uppercase tracking-[0.2em]">Batch Complete</span>
          <span class="text-[10px] font-bold text-navy-300 uppercase tracking-widest">Batch: {{ batchId.slice(-8) }}</span>
        </div>
        <button 
          @click="handleBack"
          class="px-4 py-2 bg-navy-900 hover:bg-navy-800 text-white text-[11px] font-black uppercase tracking-widest rounded-xl transition-all hover:scale-105 active:scale-95 shadow-lg shadow-navy-900/10"
        >
          Back to History
        </button>
      </div>
      
      <h2 class="text-3xl font-black text-navy-900 mb-4">Batch Analysis Results</h2>
      <div class="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
        <div class="text-center p-4 bg-emerald-50 rounded-2xl">
          <div class="text-2xl font-black text-emerald-600">{{ candidates.length }}</div>
          <div class="text-xs font-bold text-emerald-700 uppercase tracking-widest">Successful</div>
        </div>
        <div class="text-center p-4 bg-amber-50 rounded-2xl">
          <div class="text-2xl font-black text-amber-600">{{ failedCount }}</div>
          <div class="text-xs font-bold text-amber-700 uppercase tracking-widest">Failed</div>
        </div>
        <div class="text-center p-4 bg-navy-50 rounded-2xl">
          <div class="text-2xl font-black text-navy-600">{{ candidates.length + failedCount }}</div>
          <div class="text-xs font-bold text-navy-700 uppercase tracking-widest">Total Processed</div>
        </div>
      </div>
    </div>

    <!-- Candidates Grid -->
    <div class="grid grid-cols-1 gap-6">
      <div v-if="candidates.length === 0" class="bg-white rounded-3xl border border-navy-100 p-20 text-center flex flex-col items-center">
        <div class="w-20 h-20 bg-navy-50 rounded-2xl flex items-center justify-center mb-6">
          <HistoryIcon class="w-8 h-8 text-navy-200" />
        </div>
        <h3 class="text-lg font-black text-navy-900 uppercase tracking-widest">No Successful Candidates</h3>
        <p class="text-sm text-navy-400 max-w-xs mt-2 leading-relaxed">All candidates failed analysis or no results were generated.</p>
      </div>

      <div 
        v-for="(candidate, idx) in sortedCandidates" 
        :key="candidate.processing_metadata.request_id"
        class="group bg-white rounded-2xl border border-navy-100 p-6 hover:border-primary-500 hover:shadow-2xl hover:shadow-primary-500/5 transition-all duration-300"
      >
        <!-- Rank and Score -->
        <div class="flex items-start justify-between mb-6">
          <div class="flex items-center space-x-4">
            <div class="w-12 text-2xl font-black text-navy-100 group-hover:text-primary-100 transition-colors">
              {{ (idx + 1).toString().padStart(2, '0') }}
            </div>
            <div>
              <h4 class="font-black text-navy-900 text-lg tracking-tight group-hover:text-primary-600 transition-colors">
                {{ candidate.processing_metadata.cv_file_name }}
              </h4>
              <div class="flex items-center space-x-3 mt-1 text-[10px] font-bold text-navy-300 uppercase tracking-widest">
                <span class="flex items-center"><ClockIcon class="w-3 h-3 mr-1" /> {{ formatDate(candidate.processing_metadata.timestamp) }}</span>
                <span class="border-l border-navy-50 pl-3">ID: {{ candidate.processing_metadata.request_id.slice(-8) }}</span>
              </div>
            </div>
          </div>
          
          <div class="flex items-center space-x-4">
            <div class="text-center">
              <p class="text-[10px] font-black text-navy-300 uppercase tracking-widest mb-1">Match Score</p>
              <div class="flex items-end justify-center">
                <span class="text-2xl font-black text-navy-900 leading-none">{{ candidate.matching_score.overall }}</span>
                <span class="text-xs font-bold text-navy-300 ml-1">/10</span>
              </div>
            </div>
            <button 
              @click="handleViewCandidate(candidate)"
              class="px-6 py-2.5 bg-navy-900 text-white rounded-xl font-black text-[10px] uppercase tracking-widest hover:scale-[1.05] active:scale-95 transition-all shadow-xl shadow-navy-900/10"
            >
              View Report
            </button>
          </div>
        </div>

        <!-- Executive Summary -->
        <div class="bg-navy-50 rounded-xl p-4 mb-4">
          <p class="text-sm font-medium text-navy-700 italic">
            "{{ candidate.executive_summary }}"
          </p>
        </div>

        <!-- Key Skills Preview -->
        <div class="flex items-center space-x-2">
          <span class="text-[10px] font-bold text-navy-300 uppercase tracking-widest">Top Skills:</span>
          <div class="flex flex-wrap gap-2">
            <span 
              v-for="skill in candidate.technical_skills.slice(0, 3)" 
              :key="skill.name"
              class="px-2 py-1 bg-primary-100 text-primary-700 text-[9px] font-bold rounded-full"
            >
              {{ skill.name }} ({{ skill.score.toFixed(1) }})
            </span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { 
  History as HistoryIcon,
  Clock as ClockIcon
} from 'lucide-vue-next'
import type { AnalysisResult } from '../types/analysis'

const props = defineProps<{
  candidates: AnalysisResult[]
  batchId: string
  failedCount: number
}>()

const emit = defineEmits(['view', 'back'])

const sortedCandidates = computed(() => {
  return [...props.candidates].sort((a, b) => b.matching_score.overall - a.matching_score.overall)
})

const formatDate = (ts: string) => {
  return new Date(ts).toLocaleDateString('en-US', {
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

const handleViewCandidate = (candidate: AnalysisResult) => {
  emit('view', candidate)
}

const handleBack = () => {
  emit('back')
}
</script>
