package models

type AnalysisRequest struct {
	JobDescription string `json:"job_description"`
	APIKey         string `json:"api_key"`
}

type DimensionScore struct {
	Score           float64  `json:"score"`
	Evidence        []string `json:"evidence"`
	GapsIdentified  []string `json:"gaps_identified"`
	ConfidenceLevel float64  `json:"confidence_level"`
}

type MatchingScore struct {
	Overall                 float64        `json:"overall"`
	SkillsAlignment         DimensionScore `json:"skills_alignment"`
	ExperienceRelevance     DimensionScore `json:"experience_relevance"`
	SeniorityFit            DimensionScore `json:"seniority_fit"`
	DomainKnowledge         DimensionScore `json:"domain_knowledge"`
	SoftSkillsCulture       DimensionScore `json:"soft_skills_culture"`
	EducationCertifications DimensionScore `json:"education_certifications"`
}

type QualityScore struct {
	Overall                  float64        `json:"overall"`
	CredibilityVerifiability DimensionScore `json:"credibility_verifiability"`
	ImpactQuantification     DimensionScore `json:"impact_quantification"`
	CareerProgression        DimensionScore `json:"career_progression"`
	CVStructureClarity       DimensionScore `json:"cv_structure_clarity"`
	Completeness             DimensionScore `json:"completeness"`
}

type Strength struct {
	Title          string   `json:"title"`
	Description    string   `json:"description"`
	Evidence       []string `json:"evidence"`
	RelevanceScore float64  `json:"relevance_score"`
}

type Gap struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Severity    string   `json:"severity"`
	Suggestions []string `json:"suggestions"`
}

type Recommendation struct {
	Type        string `json:"type"`
	Title       string `json:"title"`
	Description string `json:"description"`
	TargetArea  string `json:"target_area"`
	Priority    string `json:"priority"`
}

type RedFlag struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Severity    string `json:"severity"`
}

type ProcessingMetadata struct {
	ProcessingTime   float64 `json:"processing_time"`
	CVPagesProcessed int     `json:"cv_pages_processed"`
	JDWordCount      int     `json:"jd_word_count"`
	ModelUsed        string  `json:"model_used"`
	Timestamp        string  `json:"timestamp"`
	RequestID        string  `json:"request_id"`
	CVFileName       string  `json:"cv_file_name,omitempty"`
	JobID            string  `json:"job_id,omitempty"`
	Status           string  `json:"status,omitempty"` // pending, processing, completed, failed
	FileHash         string  `json:"file_hash,omitempty"`
}

type SkillScore struct {
	Name  string  `json:"name"`
	Score float64 `json:"score"`
}

type AnalysisResult struct {
	ExecutiveSummary   string             `json:"executive_summary"`
	FileHash           string             `json:"file_hash"`
	MatchingScore      MatchingScore      `json:"matching_score"`
	QualityScore       QualityScore       `json:"quality_score"`
	TechnicalSkills    []SkillScore       `json:"technical_skills"`
	Strengths          []Strength         `json:"strengths"`
	Gaps               []Gap              `json:"gaps"`
	Recommendations    []Recommendation   `json:"recommendations"`
	RedFlags           []RedFlag          `json:"red_flags"`
	ProcessingMetadata ProcessingMetadata `json:"processing_metadata"`
}

type APIError struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

type ErrorResponse struct {
	Status string   `json:"status"`
	Error  APIError `json:"error"`
}

type SuccessResponse[T any] struct {
	Status string `json:"status"`
	Data   T      `json:"data"`
}

type BatchItemStatus struct {
	RequestID    string `json:"request_id"`
	CVFileName   string `json:"cv_file_name"`
	Status       string `json:"status"`
	Cached       bool   `json:"cached"`
	ErrorCode    string `json:"error_code,omitempty"`
	ErrorMessage string `json:"error_message,omitempty"`
}

type BatchStatusData struct {
	BatchID    string            `json:"batch_id"`
	Status     string            `json:"status"`
	JobID      string            `json:"job_id,omitempty"`
	Total      int               `json:"total"`
	Completed  int               `json:"completed"`
	Successful int               `json:"successful"`
	Failed     int               `json:"failed"`
	Items      []BatchItemStatus `json:"items"`
}

type BatchAnalyzeData struct {
	BatchID       string `json:"batch_id"`
	JobID         string `json:"job_id,omitempty"`
	TotalFiles    int    `json:"total_files"`
	Message       string `json:"message"`
	SessionScoped bool   `json:"session_scoped"`
}

type BatchNotificationData struct {
	BatchID    string `json:"batch_id"`
	Status     string `json:"status"`
	Complete   bool   `json:"complete"`
	Total      int    `json:"total"`
	Completed  int    `json:"completed"`
	Successful int    `json:"successful"`
	Failed     int    `json:"failed"`
	Message    string `json:"message"`
	Type       string `json:"type"`
}

type BatchFailure struct {
	RequestID    string `json:"request_id"`
	CVFileName   string `json:"cv_file_name"`
	ErrorCode    string `json:"error_code"`
	ErrorMessage string `json:"error_message"`
	Details      string `json:"details,omitempty"`
}

type BatchResultsData struct {
	BatchID     string           `json:"batch_id"`
	Candidates  []AnalysisResult `json:"candidates"`
	Failed      []BatchFailure   `json:"failed"`
	Total       int              `json:"total"`
	Successful  int              `json:"successful"`
	FailedCount int              `json:"failed_count"`
}
