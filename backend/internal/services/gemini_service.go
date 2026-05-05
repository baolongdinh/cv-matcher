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
   - **Semantic Decomposition**: Break the JD into "Non-Negotiable Core", "Key Performance Indicators", and "Preferred Nice-to-haves".
   - **Evidence Audit**: Scrutinize the CV for explicit, quantified evidence. A technology listed without a project/outcome context is rated as "Basic Familiarity" only.
   - **Contrastive Analysis**: Explicitly compare the candidate's historical project complexity and scale against the JD's requirements.

### **SCORING RUBRIC (STRICT 0.0 - 10.0 SCALE)**
You must use the full 0-10 distribution to ensure **Differentiation**:
- **9.0 - 10.0 (Elite)**: Perfect alignment. Candidate has successfully solved the specific technical challenges mentioned in the JD at a higher or equal scale. Designs architecture, not just implements. 
- **7.5 - 8.9 (Premium)**: Hits all core requirements with strong, quantified evidence. Minor gaps in non-core areas.
- **5.5 - 7.4 (Functional)**: General competency but lacks the "Wow" factor or specific scale. Relies on "Participated in" rather than "Owned/Designed".
- **4.0 - 5.4 (Underqualified)**: Missing several core technical pillars. Experience is shallow or mismatched.
- **0.0 - 3.9 (Reject)**: Total mismatch or suspicious "keyword-stuffing" without any supporting project context.

### **DIMENSION DEFINITIONS & EXPLICABILITY**
For each dimension, the "evidence" field MUST start with a **Reasoning of Score** sentence that explains the delta from a 10.0 score.

**Matching Dimensions**:
- **Skills Alignment (25%%)**: Direct tech stack match. Note specific missing high-priority tools.
- **Experience Relevance (25%%)**: Focus on the *nature* of the work. Is it the same industry? Same scale? Same complexity?
- **Seniority Fit (15%%)**: Does the leadership level match? Look for decision-making and mentoring vs task execution.
- **Domain Knowledge (15%%)**: Understanding of the specific industry vertical or specialized field.
- **Soft Skills & Culture (10%%)**: Leadership, communication style evidenced by career progression and team roles.
- **Education (10%%)**: Degree relevance and specialized certifications.

**CV Quality Dimensions**:
- **Credibility**: Is the career path logical? Are the claims believable given the company and role?
- **Impact Quantification**: Do they use numbers (%%, $, time saved)?
- **Structure**: Is the hierarchy clear for a technical reviewer?

### **INPUT DATA**
**Job Description**:
%s

**Candidate CV**:
%s

### **OUTPUT INSTRUCTIONS**
- **Executive Summary**: A 2-3 sentence high-level synthesis explaining exactly why the candidate is a match or not. Must be personalized (e.g., "While the candidate has strong Python skills, they lack the specific high-scale distributed systems experience required for this Lead role.").
- **Technical Skills Alignment**: Extract the top 3-5 specific technical skills or domain competencies required in the JD and score the candidate for each.
- Return ONLY a valid JSON object.
- NO markdown formatting (no backticks).
- All "score" fields MUST be floats (e.g., 7.8).
- "gaps" must clearly state **Severity** (CRITICAL, MODERATE, MINOR).
- "recommendations" must include specific technical questions for the interviewer to probe identified gaps.

### **TARGET JSON STRUCTURE**
{
  "executive_summary": "...",
  "seniority_verdict": "...",
  "onboarding_estimate": "...",
  "matching_score": {
    "overall": 0.0,
    "skills_alignment": {"score": 0.0, "evidence": ["Reasoning: ...", "Detail 1...", "Detail 2..."], "gaps_identified": ["Gap String 1", "Gap String 2"], "confidence_level": 1.0},
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
