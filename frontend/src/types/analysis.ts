export interface DimensionScore {
  score: number
  evidence: string[]
  gaps_identified: string[]
  confidence_level: number
}

export interface MatchingScore {
  overall: number
  skills_alignment: DimensionScore
  experience_relevance: DimensionScore
  seniority_fit: DimensionScore
  domain_knowledge: DimensionScore
  soft_skills_culture: DimensionScore
  education_certifications: DimensionScore
}

export interface QualityScore {
  overall: number
  credibility_verifiability: DimensionScore
  impact_quantification: DimensionScore
  career_progression: DimensionScore
  cv_structure_clarity: DimensionScore
  completeness: DimensionScore
}

export interface Strength {
  title: string
  description: string
  evidence: string[]
  relevance_score: number
}

export interface Gap {
  title: string
  description: string
  severity: string
  suggestions: string[]
}

export interface Recommendation {
  type: string
  title: string
  description: string
  target_area: string
  priority: string
}

export interface RedFlag {
  title: string
  description: string
  severity: string
}

export interface ProcessingMetadata {
  processing_time: number
  cv_pages_processed: number
  jd_word_count: number
  model_used: string
  timestamp: string
  request_id: string
  cv_file_name?: string
  job_id?: string
  status?: string
  file_hash?: string
}

export interface SkillScore {
  name: string
  score: number
}

export interface AnalysisResult {
  executive_summary: string
  file_hash: string
  matching_score: MatchingScore
  quality_score: QualityScore
  technical_skills: SkillScore[]
  strengths: Strength[]
  gaps: Gap[]
  recommendations: Recommendation[]
  red_flags: RedFlag[]
  processing_metadata: ProcessingMetadata
}

export interface APIError {
  code: string
  message: string
  details?: unknown
}

export interface SuccessResponse<T> {
  status: 'success'
  data: T
}

export interface ErrorResponse {
  status: 'error'
  error: APIError
}

export type APIResponse<T> = SuccessResponse<T> | ErrorResponse

export interface BatchAnalyzeData {
  batch_id: string
  job_id?: string
  total_files: number
  message: string
  session_scoped: boolean
}

export interface BatchItemStatus {
  request_id: string
  cv_file_name: string
  status: 'queued' | 'processing' | 'completed' | 'failed'
  cached: boolean
  error_code?: string
  error_message?: string
}

export interface BatchStatusData {
  batch_id: string
  status: string
  job_id?: string
  total: number
  completed: number
  successful: number
  failed: number
  items: BatchItemStatus[]
}

export interface BatchNotificationData {
  batch_id: string
  status: string
  complete: boolean
  total: number
  completed: number
  successful: number
  failed: number
  message: string
  type: 'success' | 'warning' | 'info'
}

export interface BatchFailure {
  request_id: string
  cv_file_name: string
  error_code: string
  error_message: string
  details?: string
}

export interface BatchResultsData {
  batch_id: string
  candidates: AnalysisResult[]
  failed: BatchFailure[]
  total: number
  successful: number
  failed_count: number
}

export interface LiveCandidateCard extends AnalysisResult {
  is_cached: boolean
}

export interface QueueItemViewModel extends BatchItemStatus {
  updated_label: string
}

export interface BatchFailureViewModel extends BatchFailure {
  updated_label: string
}

export interface LiveBatchViewModel {
  batch_id: string
  status: string
  job_id?: string
  total: number
  completed: number
  successful: number
  failed: number
  completed_candidates: LiveCandidateCard[]
  processing_items: QueueItemViewModel[]
  failed_items: BatchFailureViewModel[]
  last_seen_completed_count: number
  is_restored_from_storage: boolean
  is_reconnecting: boolean
}
