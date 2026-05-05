<template>
  <div class="h-64 w-full">
    <Radar :data="chartData" :options="chartOptions" />
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import {
  Chart as ChartJS,
  RadialLinearScale,
  PointElement,
  LineElement,
  Filler,
  Tooltip,
  Legend,
} from 'chart.js'
import { Radar } from 'vue-chartjs'
import type { MatchingScore } from '../types/analysis'

ChartJS.register(
  RadialLinearScale,
  PointElement,
  LineElement,
  Filler,
  Tooltip,
  Legend
)

const props = defineProps<{
  score: MatchingScore
}>()

const chartData = computed(() => ({
  labels: [
    'SKILLS',
    'EXP',
    'SENIORITY',
    'DOMAIN',
    'SOFT SKILLS',
    'EDUCATION',
  ],
  datasets: [
    {
      label: 'Candidate Score',
      backgroundColor: 'rgba(14, 165, 233, 0.1)',
      borderColor: '#0ea5e9',
      borderWidth: 3,
      pointBackgroundColor: '#0ea5e9',
      pointBorderColor: '#fff',
      pointBorderWidth: 2,
      pointRadius: 4,
      pointHoverRadius: 6,
      data: [
        props.score.skills_alignment.score,
        props.score.experience_relevance.score,
        props.score.seniority_fit.score,
        props.score.domain_knowledge.score,
        props.score.soft_skills_culture.score,
        props.score.education_certifications.score,
      ],
    },
  ],
}))

const chartOptions = {
  responsive: true,
  maintainAspectRatio: false,
  scales: {
    r: {
      grid: {
        color: 'rgba(15, 23, 42, 0.05)',
        lineWidth: 1,
      },
      angleLines: {
        color: 'rgba(15, 23, 42, 0.05)',
      },
      pointLabels: {
        font: {
          family: 'Inter',
          size: 9,
          weight: 900,
        },
        color: '#94a3b8',
      },
      ticks: {
        display: false,
        stepSize: 2,
      },
      suggestedMin: 0,
      suggestedMax: 10,
    },
  },
  plugins: {
    legend: {
      display: false,
    },
    tooltip: {
      enabled: true,
      backgroundColor: '#0f172a',
      titleFont: { size: 12, weight: 'bold' as const },
      bodyFont: { size: 12 },
      padding: 12,
      cornerRadius: 8,
      displayColors: false,
    }
  },
}
</script>
