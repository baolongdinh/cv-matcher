import { ref } from 'vue'
import type { AnalysisResult } from '../types/analysis'

const currentResult = ref<AnalysisResult | null>(null)

export const useAnalysisResult = () => {
  const setResult = (result: AnalysisResult | null) => {
    currentResult.value = result
  }

  return {
    currentResult,
    setResult,
  }
}
