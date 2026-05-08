<template>
  <div class="space-y-8 animate-in fade-in slide-in-from-bottom-4 duration-700">
    <!-- Header with Sort Options -->
    <div class="flex items-center justify-between">
      <div>
        <h2 class="text-3xl font-black text-navy-900 tracking-tight">Candidate Ranking</h2>
        <p class="text-sm text-navy-400 font-medium mt-1 uppercase tracking-widest">Sorted by Matching Score</p>
      </div>
      
      <div class="flex items-center space-x-4">
        <div class="bg-white border border-navy-100 rounded-xl p-1 flex items-center shadow-sm">
          <button 
            @click="sortBy = 'score'"
            :class="sortBy === 'score' ? 'bg-navy-900 text-white shadow-lg' : 'text-navy-400 hover:bg-navy-50'"
            class="px-4 py-2 rounded-lg text-[11px] font-black uppercase tracking-widest transition-all"
          >
            By Score
          </button>
          <button 
            @click="sortBy = 'date'"
            :class="sortBy === 'date' ? 'bg-navy-900 text-white shadow-lg' : 'text-navy-400 hover:bg-navy-50'"
            class="px-4 py-2 rounded-lg text-[11px] font-black uppercase tracking-widest transition-all"
          >
            By Date
          </button>
        </div>
      </div>
    </div>

    <!-- Candidate List -->
    <div class="grid grid-cols-1 gap-4">
      <div v-if="history.length === 0" class="bg-white rounded-3xl border border-navy-100 p-20 text-center flex flex-col items-center">
        <div class="w-20 h-20 bg-navy-50 rounded-2xl flex items-center justify-center mb-6">
          <HistoryIcon class="w-8 h-8 text-navy-200" />
        </div>
        <h3 class="text-lg font-black text-navy-900 uppercase tracking-widest">No candidates found</h3>
        <p class="text-sm text-navy-400 max-w-xs mt-2 leading-relaxed">Start an evaluation to see candidates ranked by their potential.</p>
      </div>

      <div 
        v-for="(candidate, idx) in sortedHistory" 
        :key="candidate.processing_metadata.request_id"
        class="group bg-white rounded-2xl border border-navy-100 p-6 hover:border-primary-500 hover:shadow-2xl hover:shadow-primary-500/5 transition-all duration-300 flex items-center"
      >
        <!-- Rank Number -->
        <div class="w-12 text-2xl font-black text-navy-100 group-hover:text-primary-100 transition-colors">
          {{ (idx + 1).toString().padStart(2, '0') }}
        </div>

        <!-- Candidate Info -->
        <div class="flex-1 flex items-center space-x-6">
          <div class="relative">
            <div class="absolute -inset-1 bg-gradient-to-r from-primary-400 to-indigo-500 rounded-full opacity-0 group-hover:opacity-20 transition-opacity"></div>
            <img :src="`https://ui-avatars.com/api/?name=${candidate.processing_metadata.cv_file_name || 'Candidate'}&background=f8fafc&color=0f172a&bold=true`" class="w-12 h-12 rounded-full relative border-2 border-white shadow-sm" />
          </div>
          <div>
            <h4 class="font-black text-navy-900 text-lg tracking-tight group-hover:text-primary-600 transition-colors">{{ candidate.processing_metadata.cv_file_name }}</h4>
            <div class="flex items-center space-x-3 mt-1 text-[10px] font-bold text-navy-300 uppercase tracking-widest">
              <span class="flex items-center"><ClockIcon class="w-3 h-3 mr-1" /> {{ formatDate(candidate.processing_metadata.timestamp) }}</span>
              <span class="border-l border-navy-50 pl-3">ID: {{ candidate.processing_metadata.request_id.slice(-8) }}</span>
            </div>
          </div>
        </div>

        <!-- Scores -->
        <div class="flex items-center space-x-12 px-8 border-x border-navy-50">
          <div class="text-center">
            <p class="text-[10px] font-black text-navy-300 uppercase tracking-widest mb-1">Match Score</p>
            <div class="flex items-end justify-center">
              <span class="text-2xl font-black text-navy-900 leading-none">{{ candidate.matching_score.overall }}</span>
              <span class="text-xs font-bold text-navy-300 ml-1">/10</span>
            </div>
          </div>
          <div class="text-center">
            <p class="text-[10px] font-black text-navy-300 uppercase tracking-widest mb-1">CV Quality</p>
            <div class="flex items-end justify-center">
              <span class="text-2xl font-black text-navy-700 leading-none">{{ candidate.quality_score.overall }}</span>
              <span class="text-xs font-bold text-navy-300 ml-1">/10</span>
            </div>
          </div>
        </div>

        <!-- Actions -->
        <div class="pl-8 flex items-center space-x-3">
          <button 
            @click="$emit('view', candidate)"
            class="px-6 py-2.5 bg-navy-900 text-white rounded-xl font-black text-[10px] uppercase tracking-widest hover:scale-[1.05] active:scale-95 transition-all shadow-xl shadow-navy-900/10"
          >
            View Report
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { 
  History as HistoryIcon,
  Clock as ClockIcon
} from 'lucide-vue-next'
import type { AnalysisResult } from '../types/analysis'

const props = defineProps<{
  history: AnalysisResult[]
}>()

const emit = defineEmits(['view'])

const sortBy = ref<'score' | 'date'>('score')

const sortedHistory = computed(() => {
  const items = [...props.history]
  if (sortBy.value === 'date') {
    return items.sort((a, b) =>
      new Date(b.processing_metadata.timestamp).getTime() - new Date(a.processing_metadata.timestamp).getTime(),
    )
  }
  return items.sort((a, b) => b.matching_score.overall - a.matching_score.overall)
})

const formatDate = (ts: string) => {
  return new Date(ts).toLocaleDateString('en-US', {
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}
</script>
