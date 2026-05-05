# CV-JD Matching Tool - Technical Implementation Plan
 
## Technical Context
 
### Architecture Overview
- **Architecture Pattern**: Client-Server with REST API
- **Backend**: Go + Gin (HTTP framework) 
- **Frontend**: Vue 3 + TypeScript + Vite
- **AI Engine**: Google Gemini API (gemini-2.0-flash)
- **Document Processing**: UniDoc (unidoc/unipdf)
- **Deployment**: Single-page application with stateless backend
 
### Technology Stack Decisions
 
#### Backend Components
- **Gin**: HTTP web framework for high-performance API endpoints
- **UniDoc**: Go PDF processing library for text extraction
- **google.golang.org/genai**: Official Go Gemini SDK for structured output
- **Go structs**: Data validation and serialization
- **multipart/form-data**: File upload handling via Gin
 
#### Frontend Components  
- **Vue 3**: Composition API for reactive state management
- **TypeScript**: Type safety and better development experience
- **Vite**: Fast build tool and development server
- **Tailwind CSS**: Utility-first styling without design system complexity
- **Chart.js + vue-chartjs**: Radar and bar chart visualizations
- **Axios**: HTTP client with multipart upload support
 
### Integration Patterns
- **PDF Processing**: Direct Go integration via UniDoc
- **AI Service**: Structured JSON output from Gemini API
- **File Upload**: Multipart form data to backend
- **Real-time Updates**: Server-sent events for processing status
 
## Constitution Check
 
### Development Principles
- **Code Quality**: Clean, maintainable Go and TypeScript code
- **Testing**: Unit tests for core logic, integration tests for API endpoints
- **Error Handling**: Comprehensive error handling with user-friendly messages
- **Security**: API keys only in browser memory, no persistent storage
- **Performance**: Concurrent processing, optimized PDF handling
 
### Architectural Constraints
- **Stateless Backend**: No session storage, all state in frontend
- **Browser Compatibility**: Modern browsers with ES6+ support
- **API Rate Limits**: Handle Gemini API quota gracefully
- **File Size Limits**: Reasonable PDF upload limits (e.g., 10MB)
 
## Phase 0: Research & Architecture Decisions
 
### Research Findings
 
#### UniDoc Integration
**Decision**: Use UniDoc Go library for PDF processing
**Rationale**: 
- Native Go library with active maintenance
- Better performance than Python alternatives
- Direct Go integration avoids external dependencies
- Supports various PDF formats and structures
- Excellent text extraction capabilities

**Alternatives Considered**:
- Subprocess with external PDF tools (more complex)
- Go bindings for Python libraries (less reliable)
- OCR services (overkill for text-based PDFs)
 
#### Gemini API Structured Output
**Decision**: Use gemini-2.0-flash with JSON response format
**Rationale**:
- Faster processing than gemini-1.5-pro
- Native structured output support
- Better cost/performance ratio
- Reliable JSON parsing with official SDK
 
**Alternatives Considered**:
- gemini-1.5-pro (slower, more expensive)
- Text parsing with regex (less reliable)
- Custom prompt engineering (more maintenance)
 
#### Frontend State Management
**Decision**: Vue 3 Composition API with reactive state
**Rationale**:
- Built-in reactivity without external state management
- TypeScript integration for type safety
- Smaller bundle size than Vuex/Pinia
- Adequate for single-page application scope
 
**Alternatives Considered**:
- Pinia (overkill for this scope)
- Redux (complex, not Vue-native)
- Local state only (limiting for component communication)
 
## Phase 1: System Design
 
### Data Model
 
#### Core Entities
 
```go
// Analysis Request
type AnalysisRequest struct {
    JobDescription string `json:"job_description"`
    CVFile         []byte  `json:"cv_file"`
    APIKey         string  `json:"api_key"`
}

// Analysis Result  
type AnalysisResult struct {
    MatchingScore  MatchingScore  `json:"matching_score"`
    QualityScore   QualityScore   `json:"quality_score"`
    Strengths      []Strength     `json:"strengths"`
    Gaps          []Gap          `json:"gaps"`
    Recommendations []Recommendation `json:"recommendations"`
    RedFlags      []RedFlag      `json:"red_flags"`
}

type MatchingScore struct {
    Overall          float64 `json:"overall"`
    SkillsAlignment  DimensionScore `json:"skills_alignment"`
    ExperienceRelevance DimensionScore `json:"experience_relevance"`
    SeniorityFit     DimensionScore `json:"seniority_fit"`
    DomainKnowledge  DimensionScore `json:"domain_knowledge"`
    SoftSkillsCulture DimensionScore `json:"soft_skills_culture"`
    EducationCertifications DimensionScore `json:"education_certifications"`
}

type QualityScore struct {
    Overall          float64 `json:"overall"`
    CredibilityVerifiability DimensionScore `json:"credibility_verifiability"`
    ImpactQuantification DimensionScore `json:"impact_quantification"`
    CareerProgression DimensionScore `json:"career_progression"`
    CVStructureClarity DimensionScore `json:"cv_structure_clarity"`
    Completeness     DimensionScore `json:"completeness"`
}
```
 
#### API Contracts
 
##### POST /api/analyze
```json
{
  "request": {
    "job_description": "string",
    "cv_file": "string",  // Base64 encoded PDF
    "api_key": "string"
  }
}
 
{
  "response": {
    "status": "processing|completed|error",
    "result": AnalysisResult,
    "error_message": "string"
  }
}
```
 
##### GET /api/health
```json
{
  "status": "healthy",
  "version": "string",
  "dependencies": {
    "gemini_api": "available|unavailable",
    "unidoc": "available"
  }
}
```
 
### Backend Architecture
 
#### Service Layer
```
backend/
├── cmd/
│   └── server/
│       └── main.go        # Gin app setup and routing
├── internal/
│   ├── config/
│   │   └── config.go     # Settings, CORS, environment
│   ├── handlers/
│   │   ├── analyze.go     # POST /api/analyze endpoint
│   │   └── health.go      # GET /api/health endpoint
│   ├── services/
│   │   ├── pdf_service.go # PDF → Text via UniDoc
│   │   ├── gemini_service.go # Gemini API integration
│   │   └── prompt_builder.go # System prompt + rubric builder
│   ├── models/
│   │   ├── analysis.go    # Go struct request/response models
│   │   └── scoring.go     # Score calculation models
│   └── utils/
│       ├── error_handlers.go # Custom error handling
│       └── validators.go     # Input validation
├── go.mod                  # Go module definition
└── go.sum                  # Dependency checksums
```
 
#### Key Service Implementations
 
**PDF Service**:
```go
type PDFService struct {
    // UniDoc configuration
}

func (s *PDFService) ConvertToText(pdfContent []byte) (string, error) {
    // UniDoc native conversion
    // Error handling for corrupted files
    // Size limit validation
}
```

**Gemini Service**:
```go
type GeminiService struct {
    client *genai.GenerativeModel
}

func (s *GeminiService) AnalyzeCVJD(jdText, cvText string) (*AnalysisResult, error) {
    // Structured JSON output
    // Error handling for API limits
    // Retry logic for transient failures
}
```
 
### Frontend Architecture
 
#### Component Structure
```
frontend/src/
├── components/
│   ├── InputPanel.vue        # JD paste + CV upload
│   ├── ScoreCard.vue         # Overall scores display
│   ├── RadarChart.vue        # 6-dimension matching visualization
│   ├── BarChart.vue          # Quality score breakdown
│   ├── DimensionDetail.vue   # Individual dimension analysis
│   ├── StrengthsGaps.vue     # Strengths and gaps lists
│   ├── InterviewQuestions.vue # Recommended questions
│   └── RedFlags.vue          # Warning indicators
├── views/
│   └── HomeView.vue          # Main application layout
├── services/
│   └── api.ts                # HTTP client with Axios
├── types/
│   └── analysis.ts           # TypeScript interfaces
├── composables/
│   ├── useAnalysis.ts        # Analysis state management
│   └── useFileUpload.ts      # File upload handling
└── utils/
    ├── validation.ts         # Form validation
    └── chartHelpers.ts       # Chart configuration
```
 
#### State Management Pattern
```typescript
// useAnalysis.ts
export function useAnalysis() {
  const state = reactive({
    jobDescription: '',
    cvFile: null,
    apiKey: '',
    result: null,
    loading: false,
    error: null
  })
 
  const actions = {
    async analyze() {
      // API call with progress tracking
      // Error handling and retry logic
    },
 
    reset() {
      // Clear all state
    }
  }
 
  return { state, actions }
}
```
 
### Implementation Phases
 
#### Phase 1.1: Backend Foundation (Week 1)
1. **Gin Setup**
   - Project structure creation
   - CORS configuration
   - Environment variable setup
   - Basic health endpoint

2. **PDF Processing Service**
   - UniDoc integration
   - File upload handling
   - Error handling for invalid PDFs
   - Size limit validation

3. **Gemini Service Integration**
   - API key validation
   - Basic analysis prompt
   - Structured output parsing
   - Error handling for API limits
 
#### Phase 1.2: Frontend Foundation (Week 1)
1. **Vue 3 + Vite Setup**
   - Project initialization
   - TypeScript configuration
   - Tailwind CSS setup
   - Basic routing
 
2. **Core Components**
   - Input panel with file upload
   - Basic API service layer
   - Loading and error states
   - Responsive layout
 
#### Phase 1.3: Analysis Engine (Week 2)
1. **Scoring System**
   - Prompt engineering for 6+5 dimensions
   - Structured JSON response format
   - Score calculation logic
   - Evidence extraction
 
2. **API Integration**
   - Full analyze endpoint
   - Multipart file upload
   - Progress tracking
   - Comprehensive error handling
 
#### Phase 1.4: Visualization Layer (Week 2)
1. **Chart Components**
   - Radar chart for matching dimensions
   - Bar chart for quality scores
   - Responsive design
   - Interactive tooltips
 
2. **Results Display**
   - Score cards with animations
   - Strengths/gaps presentation
   - Interview recommendations
   - Red flag indicators
 
#### Phase 1.5: Polish & Testing (Week 3)
1. **User Experience**
   - Loading states and progress indicators
   - Error message refinement
   - Input validation improvements
   - Accessibility compliance
 
2. **Testing & Quality**
   - Unit tests for services
   - Integration tests for API
   - Frontend component tests
   - End-to-end testing
 
## Quick Start Guide
 
### Development Setup

#### Backend
```bash
# Initialize Go module
go mod init cv-jd-matcher

# Install dependencies
go get github.com/gin-gonic/gin
go get github.com/gin-contrib/cors
go get github.com/joho/godotenv
go get github.com/unidoc/unipdf/v3
go get google.golang.org/api

# Run development server
go run cmd/server/main.go
```
 
#### Frontend
```bash
# Install Node.js dependencies
npm create vue@latest frontend
cd frontend
npm install
 
# Install additional dependencies
npm install axios chart.js vue-chartjs
npm install @types/node tailwindcss
 
# Run development server
npm run dev
```
 
### Environment Configuration
```bash
# .env (backend)
GEMINI_API_KEY=your_api_key_here
CORS_ORIGINS=http://localhost:3000
MAX_FILE_SIZE=10485760  # 10MB
```
 
### API Usage Example
```typescript
// Frontend API call
const formData = new FormData();
formData.append('job_description', jdText);
formData.append('cv_file', cvFile);
formData.append('api_key', apiKey);
 
const response = await axios.post('/api/analyze', formData, {
  headers: { 'Content-Type': 'multipart/form-data' }
});
```
 
## Risk Mitigation Strategies
 
### Technical Risks
- **PDF Processing Failures**: Manual text input fallback, clear error messages
- **API Rate Limits**: Exponential backoff, user notification of limits
- **Large File Handling**: Progressive upload, client-side validation
- **Browser Compatibility**: Polyfills for older browsers, graceful degradation
 
### Development Risks
- **Integration Complexity**: Modular service design, comprehensive testing
- **Performance Optimization**: Async processing, lazy loading, caching
- **Security**: API key validation, input sanitization, CORS configuration
 
## Success Metrics
 
### Technical Metrics
- API response time < 30 seconds for average files
- 99% uptime for backend services
- < 5% error rate for valid inputs
- Support for concurrent 10+ users
 
### User Experience Metrics
- Analysis completion rate > 80%
- Average session duration < 5 minutes
- User satisfaction > 4.0/5.0
- Zero critical security vulnerabilities
 
## Next Steps
 
1. **Immediate**: Set up development environments and project structure
2. **Week 1**: Implement backend foundation and basic frontend
3. **Week 2**: Complete analysis engine and visualization layer  
4. **Week 3**: Testing, polish, and deployment preparation
5. **Post-Launch**: Monitor performance, gather user feedback, iterate
 
This plan provides a comprehensive technical roadmap for implementing the CV-JD matching tool following the specification requirements while leveraging the chosen technology stack.
 