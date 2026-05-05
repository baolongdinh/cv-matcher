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
	CVFileName       string  `json:"cv_file_name"`
}

type SkillScore struct {
	Name  string  `json:"name"`
	Score float64 `json:"score"`
}

type AnalysisResult struct {
	ExecutiveSummary   string             `json:"executive_summary"`
	MatchingScore      MatchingScore      `json:"matching_score"`
	QualityScore       QualityScore       `json:"quality_score"`
	TechnicalSkills    []SkillScore       `json:"technical_skills"`
	Strengths          []Strength         `json:"strengths"`
	Gaps               []Gap              `json:"gaps"`
	Recommendations    []Recommendation   `json:"recommendations"`
	RedFlags           []RedFlag          `json:"red_flags"`
	ProcessingMetadata ProcessingMetadata `json:"processing_metadata"`
}
