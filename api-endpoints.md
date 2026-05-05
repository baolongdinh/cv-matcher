# CV-JD Matching Tool - API Contracts
 
## Overview
 
RESTful API contract for CV-JD analysis service. All endpoints accept and return JSON unless specified otherwise.
 
## Base URL
```
Development: http://localhost:8000
Production: https://api.cv-jd-matcher.com
```
 
## Authentication
No authentication required. API keys are provided per-request and not stored.
 
## Common Headers
```
Content-Type: application/json
Accept: application/json
```
 
## Common Response Format
 
### Success Response
```json
{
  "status": "success",
  "data": { /* response-specific data */ },
  "timestamp": "2025-06-17T10:30:00Z"
}
```
 
### Error Response
```json
{
  "status": "error", 
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Human-readable error description",
    "details": { /* additional error context */ }
  },
  "timestamp": "2025-06-17T10:30:00Z"
}
```
 
## Endpoints
 
### POST /api/analyze
Analyze CV against job description.
 
**Request:** `multipart/form-data`
```
job_description: string (required)
  - Job description text
  - Minimum 50 words, maximum 10,000 words
  - Plain text format
 
cv_file: file (required)  
  - PDF file containing CV
  - Maximum size: 10MB
  - Content-Type: application/pdf
 
api_key: string (required)
  - Gemini API key
  - Format: AIza[A-Za-z0-9_-]{35}
```
 
**Response:** `200 OK`
```json
{
  "status": "success",
  "data": {
    "request_id": "req_1234567890",
    "processing_status": "completed",
    "result": {
      "matching_score": {
        "overall": 7.8,
        "skills_alignment": {
          "score": 8.5,
          "evidence": [
            "5+ years of Python development experience",
            "Proficient in React and Vue.js frameworks"
          ],
          "gaps_identified": [
            "No experience with cloud platforms"
          ],
          "confidence_level": 0.85
        },
        "experience_relevance": {
          "score": 7.2,
          "evidence": [
            "3 years at fintech startup",
            "Led team of 5 developers"
          ],
          "gaps_identified": [
            "Limited enterprise experience"
          ],
          "confidence_level": 0.78
        },
        "seniority_fit": {
          "score": 8.0,
          "evidence": [
            "Senior Software Engineer title",
            "5 years total experience"
          ],
          "gaps_identified": [],
          "confidence_level": 0.82
        },
        "domain_knowledge": {
          "score": 7.5,
          "evidence": [
            "Fintech industry experience",
            "Understanding of payment systems"
          ],
          "gaps_identified": [
            "Limited regulatory knowledge"
          ],
          "confidence_level": 0.75
        },
        "soft_skills_culture": {
          "score": 8.2,
          "evidence": [
            "Led cross-functional teams",
            "Mentored junior developers"
          ],
          "gaps_identified": [],
          "confidence_level": 0.80
        },
        "education_certifications": {
          "score": 7.0,
          "evidence": [
            "BS Computer Science",
            "AWS Certified Developer"
          ],
          "gaps_identified": [
            "No advanced degree"
          ],
          "confidence_level": 0.70
        }
      },
      "quality_score": {
        "overall": 8.1,
        "credibility_verifiability": {
          "score": 8.5,
          "evidence": [
            "Specific company names and dates",
            "Quantifiable achievements with metrics"
          ],
          "confidence_level": 0.85
        },
        "impact_quantification": {
          "score": 7.8,
          "evidence": [
            "Improved performance by 40%",
            "Reduced costs by $200K annually"
          ],
          "confidence_level": 0.78
        },
        "career_progression": {
          "score": 8.2,
          "evidence": [
            "Clear advancement from Junior to Senior",
            "Increasing responsibility in each role"
          ],
          "confidence_level": 0.82
        },
        "cv_structure_clarity": {
          "score": 8.0,
          "evidence": [
            "Professional formatting",
            "Clear section organization"
          ],
          "confidence_level": 0.80
        },
        "completeness": {
          "score": 8.0,
          "evidence": [
            "All standard sections present",
            "Comprehensive contact information"
          ],
          "confidence_level": 0.80
        }
      },
      "strengths": [
        {
          "title": "Technical Leadership",
          "description": "Strong experience leading development teams and technical projects",
          "evidence": [
            "Led team of 5 developers on major product launch",
            "Architected microservices solution handling 1M+ requests"
          ],
          "relevance_score": 0.92
        },
        {
          "title": "Full-Stack Development", 
          "description": "Comprehensive experience across frontend and backend technologies",
          "evidence": [
            "Proficient in React, Vue.js, Node.js, Python",
            "Built complete web applications from scratch"
          ],
          "relevance_score": 0.88
        }
      ],
      "gaps": [
        {
          "title": "Cloud Platform Experience",
          "description": "Limited experience with major cloud platforms beyond basic usage",
          "severity": "medium",
          "suggestions": [
            "Consider AWS or Azure certification",
            "Highlight any cloud-related projects"
          ]
        }
      ],
      "recommendations": [
        {
          "type": "question",
          "title": "Cloud Architecture Design",
          "description": "Ask candidate to design a scalable cloud architecture for a web application",
          "target_area": "Cloud Platform Experience",
          "priority": "high"
        }
      ],
      "red_flags": [],
      "processing_metadata": {
        "processing_time": 23.5,
        "cv_pages_processed": 2,
        "jd_word_count": 156,
        "model_used": "gemini-2.0-flash",
        "timestamp": "2025-06-17T10:30:00Z",
        "request_id": "req_1234567890"
      }
    }
  }
}
```
 
**Error Responses:**
 
`400 Bad Request` - Validation errors
```json
{
  "status": "error",
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Invalid input parameters",
    "details": {
      "job_description": "Job description must be at least 50 words",
      "cv_file": "File size exceeds 10MB limit"
    }
  }
}
```
 
`401 Unauthorized` - Invalid API key
```json
{
  "status": "error",
  "error": {
    "code": "INVALID_API_KEY",
    "message": "Provided Gemini API key is invalid or has insufficient quota",
    "details": {
      "api_key_format": "Invalid API key format"
    }
  }
}
```
 
`413 Payload Too Large` - File too large
```json
{
  "status": "error",
  "error": {
    "code": "FILE_TOO_LARGE",
    "message": "Uploaded file exceeds maximum size limit",
    "details": {
      "max_size": "10MB",
      "received_size": "15.2MB"
    }
  }
}
```
 
`429 Too Many Requests` - Rate limit exceeded
```json
{
  "status": "error",
  "error": {
    "code": "RATE_LIMIT_EXCEEDED",
    "message": "Too many requests from this IP address",
    "details": {
      "retry_after": 60,
      "limit": "10 requests per minute"
    }
  }
}
```
 
`500 Internal Server Error` - Processing errors
```json
{
  "status": "error",
  "error": {
    "code": "PROCESSING_ERROR",
    "message": "Failed to process PDF file",
    "details": {
      "error_type": "PDF_PARSE_ERROR",
      "suggestion": "Try uploading a different PDF file or paste CV text manually"
    }
  }
}
```
 
`503 Service Unavailable` - External service issues
```json
{
  "status": "error",
  "error": {
    "code": "SERVICE_UNAVAILABLE",
    "message": "Gemini API is temporarily unavailable",
    "details": {
      "retry_after": 30,
      "service": "Google Gemini API"
    }
  }
}
```
 
### GET /api/health
Health check endpoint for monitoring service status.
 
**Response:** `200 OK`
```json
{
  "status": "success",
  "data": {
    "service": "cv-jd-matcher",
    "version": "1.0.0",
    "uptime": 86400,
    "dependencies": {
      "gemini_api": {
        "status": "available",
        "response_time_ms": 245
      },
      "unidoc": {
        "status": "available",
        "version": "v3.0.0"
      }
    },
    "system": {
      "cpu_usage": 0.25,
      "memory_usage": 0.45,
      "disk_usage": 0.60
    }
  }
}
```
 
**Error Response:** `503 Service Unavailable`
```json
{
  "status": "error",
  "error": {
    "code": "SERVICE_UNHEALTHY",
    "message": "One or more dependencies are unavailable",
    "details": {
      "gemini_api": {
        "status": "unavailable",
        "error": "Authentication failed"
      }
    }
  }
}
```
 
### GET /api/docs
API documentation endpoint (Swagger/OpenAPI).
 
**Response:** `200 OK`
- Returns OpenAPI JSON specification
- Interactive API documentation available at `/docs`
 
## Rate Limiting
 
### Per-Client Limits
- **Analysis endpoint**: 10 requests per minute per IP
- **Health endpoint**: 60 requests per minute per IP
- **Documentation endpoint**: 30 requests per minute per IP
 
### Response Headers
```
X-RateLimit-Limit: 10
X-RateLimit-Remaining: 7
X-RateLimit-Reset: 1640995200
```
 
## File Upload Constraints
 
### Supported Formats
- **PDF**: application/pdf (primary)
- **Size**: Maximum 10MB
- **Pages**: Maximum 50 pages
- **Content**: Text-based PDFs only (no scanned images)
 
### File Validation
- Magic number verification for PDF format
- Content scan for malicious content
- Structural validation for readability
 
## Error Codes Reference
 
| Code | HTTP Status | Description |
|------|-------------|-------------|
| VALIDATION_ERROR | 400 | Input validation failed |
| INVALID_API_KEY | 401 | Gemini API key invalid |
| FILE_TOO_LARGE | 413 | Uploaded file exceeds size limit |
| UNSUPPORTED_FORMAT | 415 | File format not supported |
| RATE_LIMIT_EXCEEDED | 429 | Too many requests |
| PROCESSING_ERROR | 500 | File processing failed |
| SERVICE_UNAVAILABLE | 503 | External service unavailable |
| SERVICE_UNHEALTHY | 503 | Health check failed |
| TIMEOUT_ERROR | 504 | Processing timeout |
 
## Response Time Expectations
 
### Analysis Endpoint
- **Average**: 15-30 seconds
- **95th percentile**: 45 seconds
- **Maximum**: 60 seconds (timeout)
 
### Health Endpoint
- **Average**: 50-100ms
- **95th percentile**: 200ms
- **Maximum**: 500ms
 
## SDK Integration Examples
 
### JavaScript/TypeScript
```typescript
interface AnalyzeRequest {
  jobDescription: string;
  cvFile: File;
  apiKey: string;
}
 
interface AnalyzeResponse {
  matchingScore: MatchingScore;
  qualityScore: QualityScore;
  strengths: Strength[];
  gaps: Gap[];
  recommendations: Recommendation[];
  redFlags: RedFlag[];
}
 
async function analyzeCV(request: AnalyzeRequest): Promise<AnalyzeResponse> {
  const formData = new FormData();
  formData.append('job_description', request.jobDescription);
  formData.append('cv_file', request.cvFile);
  formData.append('api_key', request.apiKey);
 
  const response = await fetch('/api/analyze', {
    method: 'POST',
    body: formData,
    headers: {
      'Accept': 'application/json'
    }
  });
 
  if (!response.ok) {
    const error = await response.json();
    throw new Error(error.error.message);
  }
 
  return response.json();
}
```
 
### Python
```python
import requests
from typing import Dict, Any
 
def analyze_cv(job_description: str, cv_file_path: str, api_key: str) -> Dict[str, Any]:
    url = "http://localhost:8000/api/analyze"
 
    with open(cv_file_path, 'rb') as cv_file:
        files = {'cv_file': cv_file}
        data = {
            'job_description': job_description,
            'api_key': api_key
        }
 
        response = requests.post(url, files=files, data=data)
        response.raise_for_status()
 
        return response.json()['data']
```
 
## Testing
 
### Request Validation
```bash
# Valid request
curl -X POST http://localhost:8000/api/analyze \
  -F "job_description=Senior Software Engineer with 5+ years experience..." \
  -F "cv_file=@resume.pdf" \
  -F "api_key=AIzaSyABC123..."
 
# Invalid job description (too short)
curl -X POST http://localhost:8000/api/analyze \
  -F "job_description=Engineer" \
  -F "cv_file=@resume.pdf" \
  -F "api_key=AIzaSyABC123..."
```
 
### Health Check
```bash
curl http://localhost:8000/api/health
```
 
## Versioning
 
API versioning follows semantic versioning:
- **Major**: Breaking changes
- **Minor**: New features, backward compatible
- **Patch**: Bug fixes, backward compatible
 
Current version: `v1.0.0`
 
Version specified in URL path for future compatibility:
- `/v1/api/analyze` (future)
- `/api/analyze` (current, defaults to v1)
 