<template>
  <div class="space-y-8 animate-in fade-in slide-in-from-bottom-4 duration-700">
    <section class="bg-navy-900 rounded-3xl border border-navy-700 shadow-2xl p-8 text-white">
      <div class="flex flex-col gap-6 lg:flex-row lg:items-start lg:justify-between">
        <div class="space-y-4">
          <div class="flex flex-wrap items-center gap-3">
            <span
              class="px-3 py-1 rounded-full text-[10px] font-black uppercase tracking-[0.2em]"
              :class="isComplete ? 'bg-emerald-400/20 text-emerald-300' : 'bg-primary-400/20 text-primary-300'"
            >
              {{ isComplete ? 'Batch Complete' : 'Batch Running' }}
            </span>
            <span class="text-[10px] font-bold text-navy-300 uppercase tracking-widest">
              Batch: {{ liveBatch.batch_id.slice(-8) }}
            </span>
            <span v-if="liveBatch.job_id" class="text-[10px] font-bold text-navy-300 uppercase tracking-widest">
              Job: {{ liveBatch.job_id }}
            </span>
            <span v-if="liveBatch.is_restored_from_storage" class="text-[10px] font-bold text-amber-300 uppercase tracking-widest">
              Restored Session
            </span>
          </div>

          <div>
            <h2 class="text-3xl font-black tracking-tight">Live Candidate Processing</h2>
            <p class="mt-2 text-sm text-navy-300 max-w-2xl">
              {{ summaryMessage }}
            </p>
          </div>
        </div>

        <div class="w-full lg:max-w-sm space-y-4">
          <div class="flex items-center justify-between text-[11px] font-black uppercase tracking-widest text-navy-300">
            <span>Progress</span>
            <span>{{ liveBatch.completed }} / {{ liveBatch.total }}</span>
          </div>
          <div class="w-full bg-navy-800 h-3 rounded-full overflow-hidden border border-navy-700">
            <div
              class="h-full bg-primary-500 transition-all duration-700 shadow-[0_0_20px_rgba(14,165,233,0.45)]"
              :style="{ width: progressPercentage + '%' }"
            ></div>
          </div>
          <div class="grid grid-cols-3 gap-3">
            <div class="bg-white/5 border border-white/10 rounded-2xl p-4">
              <p class="text-[10px] font-black uppercase tracking-widest text-navy-400">Ready</p>
              <p class="mt-2 text-2xl font-black text-white">{{ liveBatch.successful }}</p>
            </div>
            <div class="bg-white/5 border border-white/10 rounded-2xl p-4">
              <p class="text-[10px] font-black uppercase tracking-widest text-navy-400">Queue</p>
              <p class="mt-2 text-2xl font-black text-white">{{ queuedCount }}</p>
            </div>
            <div class="bg-white/5 border border-white/10 rounded-2xl p-4">
              <p class="text-[10px] font-black uppercase tracking-widest text-navy-400">Failed</p>
              <p class="mt-2 text-2xl font-black text-white">{{ liveBatch.failed }}</p>
            </div>
          </div>
          <div v-if="liveBatch.is_reconnecting" class="rounded-2xl border border-amber-300/20 bg-amber-300/10 px-4 py-3 text-sm text-amber-100">
            Connection is unstable. Retrying in the background without dropping current results.
          </div>
        </div>
      </div>
    </section>

    <section class="grid grid-cols-1 xl:grid-cols-[1.45fr_0.95fr] gap-8">
      <div class="space-y-6">
        <div class="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
          <div>
            <h3 class="text-3xl font-black text-navy-900 tracking-tight">Completed Candidates</h3>
            <p class="text-sm text-navy-400 font-medium mt-1 uppercase tracking-widest">
              {{ completedCandidatesLabel }}
            </p>
          </div>

          <div class="bg-white border border-navy-100 rounded-xl p-1 flex items-center shadow-sm w-full sm:w-auto">
            <button
              @click="sortBy = 'score'"
              :class="sortBy === 'score' ? 'bg-navy-900 text-white shadow-lg' : 'text-navy-400 hover:bg-navy-50'"
              class="flex-1 sm:flex-none px-4 py-2 rounded-lg text-[11px] font-black uppercase tracking-widest transition-all"
            >
              By Score
            </button>
            <button
              @click="sortBy = 'date'"
              :class="sortBy === 'date' ? 'bg-navy-900 text-white shadow-lg' : 'text-navy-400 hover:bg-navy-50'"
              class="flex-1 sm:flex-none px-4 py-2 rounded-lg text-[11px] font-black uppercase tracking-widest transition-all"
            >
              By Date
            </button>
          </div>
        </div>

        <div v-if="sortedCandidates.length === 0" class="bg-white rounded-3xl border border-navy-100 p-8 shadow-[0_8px_30px_rgb(0,0,0,0.04)]">
          <div class="flex items-start gap-5">
            <div class="w-14 h-14 bg-primary-50 rounded-2xl flex items-center justify-center">
              <RefreshCcwIcon class="w-6 h-6 text-primary-500 animate-spin" />
            </div>
            <div class="space-y-4 flex-1">
              <div>
                <h4 class="text-xl font-black text-navy-900">Batch is running</h4>
                <p class="text-sm text-navy-500 mt-2 max-w-xl">
                  The first completed candidate will appear here immediately. We keep the queue visible below so the user always sees progress instead of an empty dashboard.
                </p>
              </div>
              <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div
                  v-for="skeleton in 2"
                  :key="skeleton"
                  class="rounded-2xl border border-navy-100 bg-navy-50/70 p-5 animate-pulse"
                >
                  <div class="h-4 w-2/3 bg-navy-100 rounded"></div>
                  <div class="mt-4 h-8 w-1/4 bg-navy-100 rounded"></div>
                  <div class="mt-5 h-3 w-full bg-navy-100 rounded"></div>
                  <div class="mt-2 h-3 w-5/6 bg-navy-100 rounded"></div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <div class="grid grid-cols-1 gap-5">
          <article
            v-for="candidate in sortedCandidates"
            :key="candidate.processing_metadata.request_id"
            class="group bg-white rounded-3xl border p-6 transition-all duration-300 shadow-[0_8px_30px_rgb(0,0,0,0.04)]"
            :class="selectedCandidateId === candidate.processing_metadata.request_id
              ? 'border-primary-400 shadow-[0_16px_40px_rgba(14,165,233,0.12)]'
              : 'border-navy-100 hover:border-primary-400 hover:shadow-[0_16px_40px_rgba(14,165,233,0.08)]'"
          >
            <div class="flex flex-col gap-6 lg:flex-row lg:items-start lg:justify-between">
              <div class="space-y-4 flex-1">
                <div class="flex flex-wrap items-center gap-3">
                  <span class="px-3 py-1 rounded-full bg-primary-50 text-primary-700 text-[10px] font-black uppercase tracking-[0.2em]">
                    {{ candidate.processing_metadata.request_id.slice(-8) }}
                  </span>
                  <span v-if="candidate.is_cached" class="px-3 py-1 rounded-full bg-amber-50 text-amber-700 text-[10px] font-black uppercase tracking-[0.2em]">
                    Cached
                  </span>
                  <span class="text-[10px] font-black uppercase tracking-widest text-navy-300">
                    {{ formatDate(candidate.processing_metadata.timestamp) }}
                  </span>
                </div>
                <div>
                  <h4 class="text-xl font-black text-navy-900 tracking-tight">{{ candidate.processing_metadata.cv_file_name }}</h4>
                  <p class="mt-3 text-sm leading-relaxed text-navy-600">
                    "{{ candidate.executive_summary || 'No executive summary available yet.' }}"
                  </p>
                </div>
                <div class="flex flex-wrap gap-2">
                  <span
                    v-for="skill in candidate.technical_skills.slice(0, 4)"
                    :key="skill.name"
                    class="px-2.5 py-1 rounded-full bg-navy-50 text-navy-700 text-[10px] font-black uppercase tracking-widest"
                  >
                    {{ skill.name }} {{ skill.score.toFixed(1) }}
                  </span>
                </div>
              </div>

              <div class="flex flex-col items-start gap-4 lg:items-end lg:min-w-[210px]">
                <div class="grid grid-cols-2 gap-3 w-full lg:w-auto">
                  <div class="rounded-2xl bg-navy-50 px-4 py-3 text-center min-w-[92px]">
                    <p class="text-[10px] font-black uppercase tracking-widest text-navy-300">Match</p>
                    <p class="mt-1 text-2xl font-black text-navy-900">{{ candidate.matching_score.overall.toFixed(1) }}</p>
                  </div>
                  <div class="rounded-2xl bg-emerald-50 px-4 py-3 text-center min-w-[92px]">
                    <p class="text-[10px] font-black uppercase tracking-widest text-emerald-700/70">Quality</p>
                    <p class="mt-1 text-2xl font-black text-emerald-700">{{ candidate.quality_score.overall.toFixed(1) }}</p>
                  </div>
                </div>
                <button
                  @click="handleViewCandidate(candidate)"
                  class="w-full lg:w-auto px-6 py-3 bg-navy-900 text-white rounded-xl font-black text-[11px] uppercase tracking-widest hover:scale-[1.03] active:scale-95 transition-all shadow-xl shadow-navy-900/10"
                >
                  View Report
                </button>
              </div>
            </div>
          </article>
        </div>
      </div>

      <div class="space-y-6">
        <div class="bg-white rounded-3xl border border-navy-100 p-6 shadow-[0_8px_30px_rgb(0,0,0,0.04)]">
          <div class="flex items-center justify-between">
            <div>
              <h3 class="text-xl font-black text-navy-900 tracking-tight">Processing Queue</h3>
              <p class="text-sm text-navy-400 mt-1 uppercase tracking-widest">Queued, processing, and failed items</p>
            </div>
            <div class="lg:hidden bg-navy-50 rounded-xl p-1 flex items-center">
              <button
                v-for="tab in mobileTabs"
                :key="tab"
                @click="mobileTab = tab"
                :class="mobileTab === tab ? 'bg-navy-900 text-white shadow-lg' : 'text-navy-400'"
                class="px-3 py-2 rounded-lg text-[10px] font-black uppercase tracking-widest transition-all"
              >
                {{ tab }}
              </button>
            </div>
          </div>

          <div class="mt-6 space-y-3" :class="{ 'hidden lg:block': mobileTab !== 'Queue' }">
            <div
              v-for="item in visibleQueueItems"
              :key="item.request_id"
              class="rounded-2xl border p-4 transition-all"
              :class="item.status === 'processing'
                ? 'border-primary-200 bg-primary-50/60'
                : item.status === 'failed'
                  ? 'border-rose-200 bg-rose-50/60'
                  : item.status === 'completed'
                    ? 'border-emerald-200 bg-emerald-50/60'
                    : 'border-navy-100 bg-navy-50/70'"
            >
              <div class="flex items-start justify-between gap-4">
                <div class="min-w-0">
                  <p class="text-sm font-black text-navy-900 truncate">{{ item.cv_file_name }}</p>
                  <p class="mt-1 text-[11px] font-black uppercase tracking-widest"
                    :class="item.status === 'failed' ? 'text-rose-600' : item.status === 'completed' ? 'text-emerald-600' : 'text-navy-400'">
                    {{ item.updated_label }}
                  </p>
                  <p v-if="item.error_message" class="mt-2 text-xs text-rose-600 leading-relaxed">
                    {{ item.error_message }}
                  </p>
                </div>
                <span
                  class="shrink-0 px-2.5 py-1 rounded-full text-[10px] font-black uppercase tracking-[0.18em]"
                  :class="item.status === 'processing'
                    ? 'bg-primary-100 text-primary-700'
                    : item.status === 'failed'
                      ? 'bg-rose-100 text-rose-700'
                      : item.status === 'completed'
                        ? 'bg-emerald-100 text-emerald-700'
                        : 'bg-navy-100 text-navy-700'"
                >
                  {{ item.status }}
                </span>
              </div>
            </div>
            <div v-if="visibleQueueItems.length === 0" class="rounded-2xl border border-navy-100 bg-navy-50/70 p-4 text-sm text-navy-500">
              All items have been processed. Review completed candidates or failed items below.
            </div>
          </div>

          <div class="mt-6 space-y-3 lg:hidden" v-if="mobileTab === 'Failed'">
            <div
              v-for="failure in liveBatch.failed_items"
              :key="failure.request_id"
              class="rounded-2xl border border-rose-200 bg-rose-50/60 p-4"
            >
              <p class="text-sm font-black text-navy-900">{{ failure.cv_file_name }}</p>
              <p class="mt-2 text-[11px] font-black uppercase tracking-widest text-rose-600">{{ failure.error_code }}</p>
              <p class="mt-2 text-xs text-rose-700 leading-relaxed">{{ failure.error_message }}</p>
            </div>
            <div v-if="liveBatch.failed_items.length === 0" class="rounded-2xl border border-navy-100 bg-navy-50/70 p-4 text-sm text-navy-500">
              No failed candidates so far.
            </div>
          </div>

          <div class="mt-6 space-y-3 lg:hidden" v-if="mobileTab === 'Completed'">
            <div
              v-for="candidate in sortedCandidates.slice(0, 3)"
              :key="candidate.processing_metadata.request_id"
              class="rounded-2xl border border-emerald-200 bg-emerald-50/60 p-4"
            >
              <div class="flex items-center justify-between gap-3">
                <div class="min-w-0">
                  <p class="text-sm font-black text-navy-900 truncate">{{ candidate.processing_metadata.cv_file_name }}</p>
                  <p class="mt-1 text-[11px] font-black uppercase tracking-widest text-emerald-700">Ready to review</p>
                </div>
                <button
                  @click="handleViewCandidate(candidate)"
                  class="px-4 py-2 bg-navy-900 text-white rounded-xl text-[10px] font-black uppercase tracking-widest"
                >
                  Open
                </button>
              </div>
            </div>
            <div v-if="sortedCandidates.length === 0" class="rounded-2xl border border-navy-100 bg-navy-50/70 p-4 text-sm text-navy-500">
              Completed candidates will appear here as soon as the first CV finishes.
            </div>
          </div>
        </div>

        <div class="hidden lg:block bg-white rounded-3xl border border-navy-100 p-6 shadow-[0_8px_30px_rgb(0,0,0,0.04)]">
          <h3 class="text-xl font-black text-navy-900 tracking-tight">Failed Items</h3>
          <p class="text-sm text-navy-400 mt-1 uppercase tracking-widest">Immediate diagnostics for recruiter confidence</p>
          <div class="mt-6 space-y-3">
            <div
              v-for="failure in liveBatch.failed_items"
              :key="failure.request_id"
              class="rounded-2xl border border-rose-200 bg-rose-50/60 p-4"
            >
              <p class="text-sm font-black text-navy-900">{{ failure.cv_file_name }}</p>
              <p class="mt-2 text-[11px] font-black uppercase tracking-widest text-rose-600">{{ failure.error_code }}</p>
              <p class="mt-2 text-xs text-rose-700 leading-relaxed">{{ failure.error_message }}</p>
            </div>
            <div v-if="liveBatch.failed_items.length === 0" class="rounded-2xl border border-navy-100 bg-navy-50/70 p-4 text-sm text-navy-500">
              No failed candidates so far.
            </div>
          </div>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { RefreshCcw as RefreshCcwIcon } from 'lucide-vue-next'
import type { LiveBatchViewModel, LiveCandidateCard } from '../types/analysis'

const props = defineProps<{
  liveBatch: LiveBatchViewModel
  selectedCandidateId?: string | null
}>()

const emit = defineEmits<{
  view: [candidate: LiveCandidateCard]
}>()

const sortBy = ref<'score' | 'date'>('score')
const mobileTab = ref<'Completed' | 'Queue' | 'Failed'>('Queue')
const mobileTabs: Array<'Completed' | 'Queue' | 'Failed'> = ['Completed', 'Queue', 'Failed']

const sortedCandidates = computed(() => {
  const items = [...props.liveBatch.completed_candidates]
  if (sortBy.value === 'date') {
    return items.sort((left, right) =>
      new Date(right.processing_metadata.timestamp).getTime() - new Date(left.processing_metadata.timestamp).getTime(),
    )
  }
  return items.sort((left, right) => right.matching_score.overall - left.matching_score.overall)
})

const queuedCount = computed(() => props.liveBatch.processing_items.filter((item) => item.status !== 'completed' && item.status !== 'failed').length)

const progressPercentage = computed(() => {
  if (props.liveBatch.total === 0) return 0
  return Math.round((props.liveBatch.completed / props.liveBatch.total) * 100)
})

const isComplete = computed(() => props.liveBatch.total > 0 && props.liveBatch.completed >= props.liveBatch.total)

const summaryMessage = computed(() => {
  if (isComplete.value) {
    return `All ${props.liveBatch.total} CVs have been processed. Keep reviewing ranked candidates below without leaving this screen.`
  }
  if (props.liveBatch.completed_candidates.length > 0) {
    return `${props.liveBatch.completed_candidates.length} candidate${props.liveBatch.completed_candidates.length > 1 ? 's are' : ' is'} already ready for review while the rest of the queue keeps running.`
  }
  return 'We are processing the queue now. The first completed candidate will appear instantly in the ranked list below.'
})

const completedCandidatesLabel = computed(() => {
  if (sortedCandidates.value.length === 0) {
    return 'Waiting for the first completed candidate'
  }
  return `${sortedCandidates.value.length} candidate${sortedCandidates.value.length > 1 ? 's' : ''} ready for immediate review`
})

const visibleQueueItems = computed(() => props.liveBatch.processing_items.filter((item) => item.status !== 'completed'))

const handleViewCandidate = (candidate: LiveCandidateCard) => {
  emit('view', candidate)
}

const formatDate = (timestamp: string) => {
  return new Date(timestamp).toLocaleDateString('en-US', {
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  })
}
</script>
