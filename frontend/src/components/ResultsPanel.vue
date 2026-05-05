<template>
  <div class="space-y-8 animate-in fade-in slide-in-from-bottom-4 duration-700">
    <!-- Analysis Summary Header -->
    <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
      <div class="lg:col-span-2 bg-white p-8 rounded-3xl border border-navy-100 shadow-[0_8px_30px_rgb(0,0,0,0.04)]">
        <div class="flex items-center space-x-3 mb-6">
          <span class="px-3 py-1 rounded-full bg-emerald-100 text-emerald-700 text-[10px] font-black uppercase tracking-[0.2em]">Analysis Complete</span>
          <span class="text-[10px] font-bold text-navy-300 uppercase tracking-widest">ID: {{ result.processing_metadata.request_id }}</span>
        </div>
        
        <h2 class="text-4xl font-black text-navy-900 mb-2 tracking-tight">Overall Matching Score: {{ result.matching_score.overall }}/10</h2>
        <div class="flex items-center space-x-3 mb-8">
          <div class="w-3 h-3 rounded-full shadow-sm" :class="getScoreColor(result.matching_score.overall)"></div>
          <span class="text-sm font-bold text-navy-700">Hire Recommendation: <span :class="getScoreTextColor(result.matching_score.overall)" class="uppercase tracking-widest">{{ getRecommendationText(result.matching_score.overall) }}</span></span>
        </div>
        
        <div class="relative group">
          <div class="absolute -left-8 top-0 bottom-0 w-1.5 bg-navy-900 rounded-full opacity-100 group-hover:scale-y-105 transition-transform"></div>
          <p class="text-xl font-medium text-navy-700 leading-relaxed italic pr-4">
            "{{ getExecutiveSummary() }}"
          </p>
        </div>


      </div>

      <div class="bg-white p-8 rounded-3xl border border-navy-100 shadow-[0_8px_30px_rgb(0,0,0,0.04)] flex flex-col items-center">
        <div class="flex items-center justify-between w-full mb-8">
          <h3 class="text-sm font-black text-navy-900 uppercase tracking-widest">Competency Matrix</h3>
          <button class="p-1.5 text-navy-300 hover:text-navy-500 transition-colors">
            <InfoIcon class="w-4 h-4" />
          </button>
        </div>
        
        <div class="flex-1 w-full flex items-center justify-center min-h-[240px]">
          <RadarChart :score="result.matching_score" />
        </div>


      </div>
    </div>

    <!-- Details Grid -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-8">
      <!-- Top Strengths -->
      <div class="bg-white p-8 rounded-3xl border border-navy-100 shadow-[0_8px_30px_rgb(0,0,0,0.04)]">
        <h3 class="text-sm font-black text-navy-900 uppercase tracking-widest mb-8 flex items-center">
          <div class="w-8 h-8 bg-emerald-50 rounded-lg flex items-center justify-center mr-3">
            <CheckCircle2Icon class="w-4 h-4 text-emerald-500" />
          </div>
          Top Strengths
        </h3>
        
        <div class="space-y-4">
          <div v-for="(strength, idx) in result.strengths.slice(0, 3)" :key="idx" class="group flex items-start p-5 bg-navy-50/50 hover:bg-white border border-transparent hover:border-navy-100 hover:shadow-xl hover:shadow-navy-900/5 rounded-2xl transition-all duration-300">
            <div class="w-10 h-10 bg-white rounded-xl border border-navy-50 flex items-center justify-center mr-4 group-hover:scale-110 transition-transform">
              <SparklesIcon class="w-5 h-5 text-amber-400" />
            </div>
            <div class="flex-1">
              <h4 class="font-black text-navy-900 text-[15px] mb-1">{{ strength.title }}</h4>
              <p class="text-sm text-navy-500 leading-relaxed font-medium">{{ strength.description }}</p>
            </div>
          </div>
        </div>
      </div>

      <!-- Key Technical Alignment (Skill-by-Skill) -->
      <div class="bg-navy-900 p-8 rounded-3xl shadow-2xl relative overflow-hidden group">
        <div class="absolute top-0 right-0 p-8 opacity-[0.03] group-hover:opacity-[0.06] transition-opacity">
          <LayoutIcon class="w-64 h-64 text-white" />
        </div>
        
        <div class="relative z-10 h-full flex flex-col">
          <div class="flex items-center justify-between mb-8">
            <h3 class="text-sm font-black text-white uppercase tracking-widest">Key Technical Alignment</h3>
          </div>
          
          <div class="space-y-8 flex-1">
            <div v-for="(skill, idx) in result.technical_skills" :key="idx">
              <div class="flex justify-between items-end mb-3">
                <span class="text-[13px] font-bold text-white">{{ skill.name }}</span>
                <span class="text-[11px] font-black text-navy-400 uppercase tracking-widest">{{ skill.score.toFixed(1) }} / 10</span>
              </div>
              <div class="w-full bg-navy-800 h-2.5 rounded-full overflow-hidden border border-navy-700/50">
                <div 
                  class="bg-primary-500 h-full rounded-full transition-all duration-1000 shadow-[0_0_15px_rgba(14,165,233,0.3)]"
                  :style="{ width: (skill.score * 10) + '%' }"
                ></div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Final Verdict Grid -->
    <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
      <!-- Interview Recommendations -->
      <div class="lg:col-span-2 bg-white p-8 rounded-3xl border border-navy-100 shadow-[0_8px_30px_rgb(0,0,0,0.04)]">
        <div class="flex items-center justify-between mb-8">
          <h3 class="text-sm font-black text-navy-900 uppercase tracking-widest flex items-center">
            <div class="w-8 h-8 bg-indigo-50 rounded-lg flex items-center justify-center mr-3">
              <MessageSquareIcon class="w-4 h-4 text-indigo-500" />
            </div>
            Interview Guide & Risk Assessment
          </h3>

        </div>

        <div class="space-y-8">
          <div v-for="(rec, idx) in result.recommendations.slice(0, 3)" :key="idx" class="relative pl-8 group">
            <div class="absolute left-0 top-0 bottom-0 w-1 bg-navy-100 rounded-full group-hover:bg-navy-900 transition-colors"></div>
            <span class="text-[9px] font-black text-navy-300 uppercase tracking-[0.2em] block mb-2">Target Area: {{ rec.target_area }}</span>
            <p class="text-[17px] font-bold text-navy-800 leading-relaxed italic group-hover:text-navy-950 transition-colors">"{{ rec.description }}"</p>

          </div>
        </div>
      </div>

      <!-- Red Flags -->
      <div class="bg-white p-8 rounded-3xl border border-navy-100 shadow-[0_8px_30px_rgb(0,0,0,0.04)]">
        <h3 class="text-sm font-black text-navy-900 uppercase tracking-widest mb-8 flex items-center">
          <div class="w-8 h-8 bg-rose-50 rounded-lg flex items-center justify-center mr-3">
            <AlertCircleIcon class="w-4 h-4 text-rose-500" />
          </div>
          Red Flags
        </h3>
        
        <div class="space-y-4">
          <div v-for="(flag, idx) in result.red_flags" :key="idx" 
            class="p-5 rounded-2xl border-l-4 transition-all hover:shadow-lg"
            :class="flag.severity === 'high' ? 'bg-rose-50/50 border-rose-500' : 'bg-amber-50/50 border-amber-500'"
          >
            <div class="flex items-center justify-between mb-3">
              <span class="text-[9px] font-black uppercase tracking-widest" :class="flag.severity === 'high' ? 'text-rose-600' : 'text-amber-600'">{{ flag.severity }}</span>
              <AlertTriangleIcon class="w-4 h-4" :class="flag.severity === 'high' ? 'text-rose-500' : 'text-amber-500'" />
            </div>
            <h4 class="font-black text-navy-900 text-sm mb-1">{{ flag.title }}</h4>
            <p class="text-xs font-medium text-navy-500 leading-relaxed">{{ flag.description }}</p>
          </div>

          <div v-if="result.red_flags.length === 0" class="flex flex-col items-center justify-center py-16 text-center">
            <div class="w-16 h-16 bg-emerald-50 rounded-full flex items-center justify-center mb-6">
              <ShieldCheckIcon class="w-8 h-8 text-emerald-500" />
            </div>
            <h4 class="text-sm font-black text-navy-900 uppercase tracking-widest">Safe Profile</h4>
            <p class="text-xs font-medium text-navy-400 mt-2">No red flags detected by analysis engine.</p>
          </div>

          <div class="mt-10 p-6 bg-primary-50 rounded-2xl border border-primary-100 text-center relative overflow-hidden group">
            <div class="absolute inset-0 bg-primary-100 opacity-0 group-hover:opacity-20 transition-opacity"></div>
            <span class="text-[10px] font-black text-primary-600 uppercase tracking-[0.2em] relative z-10">Risk Score</span>
            <p class="text-2xl font-black text-navy-900 mt-1 uppercase tracking-tight relative z-10">{{ getRiskLevel(result.red_flags.length) }}</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { 
  Info as InfoIcon,
  CheckCircle2 as CheckCircle2Icon,
  Sparkles as SparklesIcon,
  MessageSquare as MessageSquareIcon,
  AlertTriangle as AlertTriangleIcon,
  AlertCircle as AlertCircleIcon,
  ShieldCheck as ShieldCheckIcon,
  Layout as LayoutIcon
} from 'lucide-vue-next'
import type { AnalysisResult } from '../types/analysis'
import RadarChart from './RadarChart.vue'

const props = defineProps<{
  result: AnalysisResult
}>()

const getScoreColor = (score: number) => {
  if (score >= 8) return 'bg-emerald-500 shadow-emerald-500/20'
  if (score >= 6) return 'bg-amber-500 shadow-amber-500/20'
  return 'bg-rose-500 shadow-rose-500/20'
}

const getScoreTextColor = (score: number) => {
  if (score >= 8) return 'text-emerald-600'
  if (score >= 6) return 'text-amber-600'
  return 'text-rose-600'
}

const getRecommendationText = (score: number) => {
  if (score >= 8.5) return 'Strong Yes'
  if (score >= 7) return 'Recommended'
  if (score >= 5) return 'Needs Review'
  return 'Not Recommended'
}

const getRiskLevel = (flagsCount: number) => {
  if (flagsCount === 0) return 'Low'
  if (flagsCount <= 2) return 'Moderate'
  return 'High Risk'
}

const getExecutiveSummary = () => {
  return props.result.executive_summary || 'Analysis summary is being generated...'
}
</script>

<style scoped>
@keyframes spin-slow {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}
.animate-spin-slow {
  animation: spin-slow 8s linear infinite;
}
</style>
