<template>
  <div class="grid grid-cols-1 lg:grid-cols-2 gap-8 items-start">
    <!-- Left Column: JD & API Key -->
    <div class="space-y-6">
      <div class="bg-white rounded-2xl border border-navy-100 shadow-[0_8px_30px_rgb(0,0,0,0.04)] overflow-hidden">
        <div class="px-6 py-5 border-b border-navy-50 flex items-center justify-between bg-white">
          <h3 class="text-[17px] font-bold text-navy-900">Evaluation Configuration</h3>
          <div class="flex items-center space-x-2">
            <div class="w-2 h-2 rounded-full bg-emerald-500 animate-pulse"></div>
            <span class="text-[11px] font-bold text-emerald-600 uppercase tracking-widest">API Active</span>
          </div>
        </div>
        <div class="p-6 space-y-5">
          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <div class="flex items-center justify-between mb-2">
                <label class="text-[13px] font-bold text-navy-700 uppercase tracking-tight">Gemini API Key</label>
                <span class="text-[10px] font-medium text-navy-400 uppercase tracking-widest">Encrypted</span>
              </div>
              <div class="relative group">
                <input 
                  v-model="apiKey"
                  type="password" 
                  placeholder="••••••••••••••••••••••••••••••••••••"
                  class="block w-full px-4 py-3 bg-navy-50 border-none rounded-xl focus:ring-2 focus:ring-primary-500/20 text-navy-900 text-sm transition-all group-hover:bg-navy-100"
                />
                <KeyIcon class="w-4 h-4 absolute right-4 top-1/2 -translate-y-1/2 text-navy-300 group-hover:text-navy-400 transition-colors" />
              </div>
            </div>
            <div>
              <div class="flex items-center justify-between mb-2">
                <label class="text-[13px] font-bold text-navy-700 uppercase tracking-tight">Job Title / ID</label>
                <span class="text-[10px] font-medium text-navy-400 uppercase tracking-widest">For grouping</span>
              </div>
              <div class="relative group">
                <input 
                  v-model="jobId"
                  type="text" 
                  placeholder="e.g. Senior Frontend Dev"
                  class="block w-full px-4 py-3 bg-navy-50 border-none rounded-xl focus:ring-2 focus:ring-primary-500/20 text-navy-900 text-sm transition-all group-hover:bg-navy-100 font-bold"
                />
                <BriefcaseIcon class="w-4 h-4 absolute right-4 top-1/2 -translate-y-1/2 text-navy-300 group-hover:text-navy-400 transition-colors" />
              </div>
            </div>
          </div>
          
          <div>
            <label class="block text-[13px] font-bold text-navy-700 uppercase tracking-tight mb-2">Job Description & Requirements</label>
            <textarea
              v-model="jd"
              rows="12"
              placeholder="Paste the technical job description, key performance indicators, and required technology stack here for the AI to analyze candidate alignment..."
              class="block w-full px-4 py-4 bg-navy-50 border-none rounded-xl focus:ring-2 focus:ring-primary-500/20 text-navy-900 text-sm transition-all hover:bg-navy-100 resize-none"
            ></textarea>
            <div class="mt-3 flex items-start space-x-2 opacity-60">
              <InfoIcon class="w-3.5 h-3.5 mt-0.5 text-navy-400" />
              <p class="text-[11px] font-medium text-navy-500 leading-tight uppercase tracking-wider">Analysis engine: Gemini 1.5 Flash - Data-rich processing enabled</p>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Right Column: Upload & Process -->
    <div class="space-y-6">
      <!-- Upload Card -->
      <div 
        class="bg-white rounded-2xl border-2 border-dashed transition-all duration-300 min-h-[380px] flex flex-col items-center justify-center p-8 text-center"
        :class="isDragging ? 'border-primary-500 bg-primary-50/30' : 'border-navy-100 bg-white hover:border-navy-200'"
        @dragover.prevent="isDragging = true"
        @dragleave.prevent="isDragging = false"
        @drop.prevent="handleDrop"
        @click="fileInput?.click()"
      >
        <input 
          type="file" 
          ref="fileInput" 
          class="hidden" 
          accept=".pdf"
          multiple
          @change="handleFileChange"
        />
        
        <div class="w-20 h-20 bg-navy-50 rounded-2xl flex items-center justify-center mb-6 group-hover:scale-110 transition-transform">
          <FileUpIcon class="w-8 h-8 text-navy-400" />
        </div>
        
        <h3 class="text-lg font-bold text-navy-900 mb-2">Upload Candidate CVs</h3>
        <p class="text-sm text-navy-500 max-w-[280px] leading-relaxed mb-8">
          Drag and drop multiple candidate resumes in PDF format.
        </p>
        
        <button 
          type="button"
          class="px-8 py-2.5 bg-white border border-navy-200 rounded-lg text-xs font-black text-navy-900 hover:bg-navy-50 transition-all uppercase tracking-widest shadow-sm"
        >
          Select Files
        </button>
        
        <div v-if="cvFiles.length > 0" class="mt-8 flex flex-wrap justify-center gap-2">
          <div v-for="file in cvFiles" :key="file.name" class="flex items-center space-x-2 bg-emerald-50 px-3 py-1.5 rounded-full border border-emerald-100">
            <CheckCircle2Icon class="w-3 h-3 text-emerald-500" />
            <span class="text-[10px] font-bold text-emerald-700 truncate max-w-[120px]">{{ file.name }}</span>
          </div>
        </div>
      </div>

      <!-- Action Card -->
      <div class="bg-navy-900 rounded-2xl p-8 shadow-2xl relative overflow-hidden group">
        <div class="absolute top-0 right-0 p-8 opacity-10 group-hover:opacity-20 transition-opacity">
          <ZapIcon class="w-32 h-32 text-white" />
        </div>
        
        <div class="relative z-10">
          <div class="flex items-center space-x-4 mb-8">
            <div class="w-12 h-12 bg-primary-500 rounded-xl flex items-center justify-center shadow-lg shadow-primary-500/20">
              <ZapIcon class="w-6 h-6 text-white" />
            </div>
            <div>
              <h3 class="text-xl font-bold text-white tracking-tight">Process Evaluation</h3>
              <p class="text-sm text-navy-400">Cross-reference candidate data with requirements.</p>
            </div>
          </div>


          <button
            @click="handleStart"
            :disabled="loading || !canSubmit"
            class="w-full py-4 bg-white hover:bg-navy-50 disabled:opacity-30 disabled:cursor-not-allowed rounded-xl text-sm font-black text-navy-950 flex items-center justify-center space-x-3 transition-all active:scale-[0.98]"
          >
            <Loader2Icon v-if="loading" class="w-5 h-5 animate-spin" />
            <RocketIcon v-else class="w-5 h-5" />
            <span class="uppercase tracking-widest">{{ loading ? 'Analyzing...' : 'Start Analysis Command' }}</span>
          </button>
        </div>
      </div>

    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { 
  Info as InfoIcon, 
  Key as KeyIcon,
  Briefcase as BriefcaseIcon,
  FileUp as FileUpIcon,
  Zap as ZapIcon, 
  Rocket as RocketIcon,
  Loader2 as Loader2Icon,
  CheckCircle2 as CheckCircle2Icon
} from 'lucide-vue-next'

const props = defineProps<{
  loading: boolean
}>()

const emit = defineEmits(['start'])

const apiKey = ref(localStorage.getItem('gemini_api_key') || '')
const jobId = ref(localStorage.getItem('cv_matcher_job_id') || '')
const jd = ref(localStorage.getItem('cv_matcher_jd') || '')
const fileInput = ref<HTMLInputElement | null>(null)
const cvFiles = ref<File[]>([])
const isDragging = ref(false)

const canSubmit = computed(() => {
  return apiKey.value.length > 20 && jd.value.length > 50 && cvFiles.value.length > 0
})

const handleFileChange = (e: Event) => {
  const target = e.target as HTMLInputElement
  if (target.files) {
    cvFiles.value = Array.from(target.files)
  }
}

const handleDrop = (e: DragEvent) => {
  isDragging.value = false
  if (e.dataTransfer?.files) {
    cvFiles.value = Array.from(e.dataTransfer.files)
  }
}

const handleStart = () => {
  localStorage.setItem('gemini_api_key', apiKey.value)
  localStorage.setItem('cv_matcher_job_id', jobId.value)
  localStorage.setItem('cv_matcher_jd', jd.value)
  emit('start', {
    apiKey: apiKey.value,
    jobId: jobId.value,
    jd: jd.value,
    cvFiles: cvFiles.value
  })
}
</script>
