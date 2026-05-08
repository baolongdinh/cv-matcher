package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"cv-jd-matcher/internal/models"
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
1. **Persona**: Act as a skeptical, highly-experienced Engineering Manager sitting on a hiring board.
2. **Methodology**: 
   - **Semantic Decomposition**: Break the JD into "Non-Negotiable Tech Pillars".
   - **Dreyfus Model Audit**: Classify every key skill into one of 5 levels: Novice, Competent, Proficient, Expert, Master.
   - **Scale & Complexity Audit**: Scrutinize the CV for EXPLICIT evidence of handling specific scale.

### **SCORING RUBRIC (HYPER-STRICT 0.0 - 10.0 SCALE)**
- **9.3 - 10.0 (Exceptional)**: The 1%%. Master level.
- **8.0 - 9.2 (Elite Senior)**: Expert level.
- **6.5 - 7.9 (Premium Functional)**: Proficient level.
- **5.0 - 6.4 (Developing)**: Competent level.
- **0.0 - 4.9 (Gap/Mismatch)**: Missing core technical pillars.

### **INPUT DATA**
**Job Description**:
%s

**Candidate CV**:
%s

### **OUTPUT INSTRUCTIONS**
- Return ONLY a valid JSON object. No markdown. All scores are floats.
- Ensure the JSON matches the target structure exactly.

### **TARGET JSON STRUCTURE**
{
  "executive_summary": "...",
  "matching_score": {
    "overall": 0.0,
    "skills_alignment": {"score": 0.0, "evidence": [], "gaps_identified": [], "confidence_level": 1.0},
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

	// API Endpoint for Gemini 2.5 Flash Lite
	url := "https://generativelanguage.googleapis.com/v1beta/models/gemini-2.5-flash-lite:generateContent"

	// Create request body
	requestBody := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"parts": []map[string]interface{}{
					{"text": prompt},
				},
			},
		},
		"generationConfig": map[string]interface{}{
			"responseMimeType": "application/json",
			"temperature":      0.3,
			"maxOutputTokens":  8192,
		},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-goog-api-key", s.apiKey)

	client := &http.Client{Timeout: 120 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("api request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("gemini api returned error (status %d): %s", resp.StatusCode, string(body))
	}

	var geminiResponse struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}

	if err := json.Unmarshal(body, &geminiResponse); err != nil {
		return nil, fmt.Errorf("failed to parse gemini response: %w", err)
	}

	if len(geminiResponse.Candidates) == 0 || len(geminiResponse.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("empty response from gemini")
	}

	jsonText := geminiResponse.Candidates[0].Content.Parts[0].Text

	// Pre-processing to ensure valid JSON
	// Sometimes AI returns text before or after the JSON block even with responseMimeType
	start := strings.Index(jsonText, "{")
	end := strings.LastIndex(jsonText, "}")
	if start != -1 && end != -1 && end > start {
		jsonText = jsonText[start : end+1]
	}

	var analysisResult models.AnalysisResult
	if err := json.Unmarshal([]byte(jsonText), &analysisResult); err != nil {
		return nil, fmt.Errorf("failed to parse analysis JSON: %v\nJSON: %s", err, jsonText)
	}

	return &analysisResult, nil
}
