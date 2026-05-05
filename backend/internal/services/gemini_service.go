package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"cv-jd-matcher/internal/models"

	"google.golang.org/genai"
)

type GeminiService struct {
	apiKey string
}

func NewGeminiService(apiKey string) (*GeminiService, error) {
	apiKey = strings.TrimSpace(apiKey)
	if apiKey == "" {
		return nil, fmt.Errorf("gemini api key is required from client")
	}

	log.Printf("[GeminiService] Initialized with API Key (length: %d)", len(apiKey))

	return &GeminiService{
		apiKey: apiKey,
	}, nil
}

func (s *GeminiService) AnalyzeCVJD(ctx context.Context, jdText, cvText string) (*models.AnalysisResult, error) {
	prompt := fmt.Sprintf(`
You are the **Lead Technical Talent Assessment Board** for a Tier-1 Global Engineering Organization. 
Your objective is to produce a definitive, hyper-accurate, and uninflated technical match analysis between a Candidate CV and a Job Description (JD).

### **EXECUTIVE ASSESSMENT PROTOCOL**
1. **Persona**: Act as a skeptical, highly-experienced Engineering Manager sitting on a hiring board with a Senior Architect and a Technical Recruiter.
2. **Methodology**: 
   - **Semantic Decomposition**: Break the JD into "Non-Negotiable Tech Pillars", "Business Logic/Scale Context", and "Leadership/Soft Skills".
   - **Dreyfus Model Audit**: Classify every key skill into one of 5 levels:
     - **Novice**: Basic understanding, needs constant guidance.
     - **Competent**: Can complete tasks decently but lacks deep optimization knowledge.
     - **Proficient**: Operates independently; understands "Why" things work.
     - **Expert**: Designs systems; anticipates failure modes; optimizes for scale/performance.
     - **Master/Outlier**: Redefines the field; has unique, rare, or exceptional contributions (e.g., patent-level work, solving impossible scale issues).
   - **Scale & Complexity Audit**: Scrutinize the CV for EXPLICIT evidence of handling specific scale (RPM, P99 latency, TB of data). If scale is not mentioned, assume "Standard Low-Scale".
   - **Exceptional Signal Search**: Look for "Outlier" signals: Open source leadership, unique architectural ownership, solving specific "hard" industry problems.

### **SCORING RUBRIC (HYPER-STRICT 0.0 - 10.0 SCALE)**
You must use the full 0-10 distribution to ensure **Differentiation**:
- **9.3 - 10.0 (Exceptional/Outlier)**: The 1%%. Candidate brings unique expertise that solves the JD's most difficult problems immediately. Clear evidence of "Master" level depth.
- **8.0 - 9.2 (Elite Senior)**: Strong alignment at "Expert" level. Has designed and owned similar systems at a similar or higher scale.
- **6.5 - 7.9 (Premium Functional)**: Hits all core requirements at "Proficient" level. Solid, dependable, but lacks the architectural "Wow" factor.
- **5.0 - 6.4 (Developing)**: Competent in tasks but needs mentorship for architectural decisions or scale handling.
- **0.0 - 4.9 (Gap/Mismatch)**: Missing core technical pillars or mismatched complexity levels.

### **DIMENSION DEFINITIONS**
For each dimension, the "evidence" field MUST start with a **Reasoning of Score** sentence explaining the technical depth gap from a 10.0.

**Matching Dimensions**:
- **Skills Alignment (25%%)**: Score based on the **Dreyfus Level** of the core tech stack.
- **Experience Relevance (25%%)**: Focus on **Complexity & Scale**. Is the previous scale comparable?
- **Seniority Fit (15%%)**: Look for **Decision Ownership**. Did they decide the tech stack or just use it?
- **Domain Knowledge (15%%)**: Specialized industry vertical deep-dives.
- **Soft Skills & Culture (10%%)**: Communication style evidenced by mentoring, leadership, or project outcomes.
- **Education (10%%)**: Degree relevance and specialized certifications.

### **INPUT DATA**
**Job Description**:
%s

**Candidate CV**:
%s

### **OUTPUT INSTRUCTIONS**
- **Executive Summary**: A brutal, honest 2-3 sentence synthesis. Start with specific fit (e.g., "The candidate is a Proficient SDE but lacks the Expert-level concurrency handling required for our Low-Latency requirements.")
- **Technical Skills Alignment**: Extract 3-5 key skills. Score them strictly using the Dreyfus context.
- Return ONLY a valid JSON object. No markdown. All scores are floats.

### **TARGET JSON STRUCTURE**
{
  "executive_summary": "...",
  "matching_score": {
    "overall": 0.0,
    "skills_alignment": {"score": 0.0, "evidence": ["Reasoning: ...", "Detail 1...", "Detail 2..."], "gaps_identified": ["Gap String 1"], "confidence_level": 1.0},
    "experience_relevance": {"score": 0.0, "evidence": [], "gaps_identified": [], "confidence_level": 1.0},
    "seniority_fit": {"score": 0.0, "evidence": [], "gaps_identified": [], "confidence_level": 1.0},
    "domain_knowledge": {"score": 0.0, "evidence": [], "gaps_identified": [], "confidence_level": 1.0},
    "soft_skills_culture": {"score": 0.0, "evidence": [], "gaps_identified": [], "confidence_level": 1.0},
    "education_certifications": {"score": 0.0, "evidence": [], "gaps_identified": [], "confidence_level": 1.0}
  },
  "quality_score": {
    "overall": 0.0,
    "credibility_verifiability": {"score": 0.0, "evidence": [], "gaps_identified": [], "confidence_level": 1.0},
    "impact_quantification": {"score": 0.0, "evidence": [], "gaps_identified": [], "confidence_level": 1.0},
    "career_progression": {"score": 0.0, "evidence": [], "gaps_identified": [], "confidence_level": 1.0},
    "cv_structure_clarity": {"score": 0.0, "evidence": [], "gaps_identified": [], "confidence_level": 1.0},
    "completeness": {"score": 0.0, "evidence": [], "gaps_identified": [], "confidence_level": 1.0}
  },
  "technical_skills": [
    {"name": "Skill Name", "score": 0.0}
  ],
  "strengths": [{"title": "", "description": "", "evidence": [], "relevance_score": 0.0}],
  "gaps": [{"title": "", "description": "", "severity": "CRITICAL|MODERATE|MINOR", "suggestions": []}],
  "recommendations": [{"type": "INTERVIEW_QUESTION", "title": "", "description": "", "target_area": "", "priority": "HIGH|MEDIUM|LOW"}],
  "red_flags": [{"title": "", "description": "", "severity": "HIGH|LOW"}]
}
`, jdText, cvText)

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  s.apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create genai client: %w", err)
	}

	model := "gemini-2.5-flash-lite"

	temperature := float32(0.3) // Lower temperature for more stable, objective analysis

	config := &genai.GenerateContentConfig{
		ResponseMIMEType: "application/json",
		Temperature:      &temperature,
		MaxOutputTokens:  int32(8192),
	}

	parts := []*genai.Part{
		{Text: prompt},
	}

	// Add timeout context
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 120*time.Second)
	defer cancel()

	result, err := client.Models.GenerateContent(ctxWithTimeout, model, []*genai.Content{{Parts: parts}}, config)
	if err != nil {
		return nil, fmt.Errorf("gemini api error: %w", err)
	}

	if len(result.Candidates) == 0 || result.Candidates[0].Content == nil || len(result.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("empty response from gemini")
	}

	jsonText := result.Candidates[0].Content.Parts[0].Text

	// Clean up potential markdown blocks if JSON mode fails to strip them
	jsonText = strings.TrimPrefix(jsonText, "```json")
	jsonText = strings.TrimPrefix(jsonText, "```")
	jsonText = strings.TrimSuffix(jsonText, "```")
	jsonText = strings.TrimSpace(jsonText)

	var analysisResult models.AnalysisResult
	if err := json.Unmarshal([]byte(jsonText), &analysisResult); err != nil {
		return nil, fmt.Errorf("failed to parse analysis JSON: %v\nJSON: %s", err, jsonText)
	}

	return &analysisResult, nil
}
