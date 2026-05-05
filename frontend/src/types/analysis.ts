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

export interface AnalysisResponse {
  status: 'success' | 'error'
  data?: AnalysisResult
  error?: {
    code: string
    message: string
    details?: any
  }
}
