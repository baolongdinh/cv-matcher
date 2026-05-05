# CV-JD Matching Tool - Technical Specification

## Feature Overview

A comprehensive CV-JD matching tool that enables recruiters to evaluate candidates against job descriptions through AI-powered analysis, providing detailed scoring, visual insights, and actionable recommendations.

## User Stories

### US-01: CV-JD Analysis Input
**As a recruiter**, I want to paste/upload a JD and upload a CV PDF, so that I can get a comprehensive evaluation of the candidate.

### US-02: Matching Score Visualization  
**As a recruiter**, I want to see a Matching Score /10 with detailed breakdown, so that I know exactly which areas the candidate fits or lacks.

### US-03: CV Quality Assessment
**As a recruiter**, I want to see a CV Quality Score /10 measuring credibility & professionalism, so that I can assess the candidate independently of the JD.

### US-04: Rich Output Analytics
**As a recruiter**, I want rich output data (charts, tables, written analysis), so that I can make a confident hiring decision.

### US-05: API Key Management
**As a user**, I want to input my Gemini API key in the UI, so that I don't need any backend setup.

## Functional Requirements

### FR-01: Input Processing
- The system shall accept job description text via paste or text input
- The system shall accept CV files in PDF format via file upload
- The system shall validate that uploaded files are readable PDF documents
- The system shall provide clear error messages for invalid file formats
- The system shall store user-provided API keys in browser memory only

### FR-02: Document Analysis
- The system shall convert PDF CV files to text format for processing
- The system shall extract and preserve document structure during conversion
- The system shall handle multi-page PDF documents completely
- The system shall provide fallback options when PDF conversion fails

### FR-03: AI-Powered Scoring
- The system shall analyze CV-JD compatibility using AI assessment
- The system shall generate a matching score out of 10 across 6 dimensions:
  - Skills Alignment (25% weight)
  - Experience Relevance (25% weight) 
  - Seniority Fit (15% weight)
  - Domain Knowledge (15% weight)
  - Soft Skills & Culture Fit (10% weight)
  - Education & Certifications (10% weight)
- The system shall generate a CV quality score out of 10 across 5 dimensions:
  - Credibility & Verifiability (25% weight)
  - Impact Quantification (25% weight)
  - Career Progression (20% weight)
  - CV Structure & Clarity (15% weight)
  - Completeness (15% weight)

### FR-04: Results Visualization
- The system shall display radar charts for multi-dimensional scoring breakdowns
- The system shall provide bar charts for individual dimension comparisons
- The system shall present executive summary with overall scores and recommendations
- The system shall list candidate strengths with supporting evidence
- The system shall identify gaps and concerns with specific examples
- The system shall generate interview recommendations based on weak areas
- The system shall flag potential red flags (job hopping, inflated claims, inconsistencies)

### FR-05: Error Handling & Edge Cases
- The system shall handle scanned/image-based PDFs with appropriate error messages
- The system shall validate API key validity before processing
- The system shall warn users for job descriptions under 50 words
- The system shall handle fresher CVs with no experience appropriately
- The system shall provide clear feedback for API quota exceeded scenarios

## User Scenarios & Testing

### Primary Flow: Complete CV-JD Analysis
1. User opens the application interface
2. User inputs their Gemini API key
3. User pastes job description text into the JD input field
4. User uploads CV PDF file
5. System validates inputs and converts PDF to text
6. System processes analysis through AI service
7. System displays comprehensive results with scores and visualizations
8. User reviews analysis and can download or share results

### Edge Case Flow: PDF Processing Failure
1. User uploads PDF file that cannot be processed
2. System detects conversion failure
3. System displays specific error message with suggested alternatives
4. User can retry with different file or paste CV text manually

### Error Flow: Invalid API Key
1. User submits analysis with invalid API key
2. System detects authentication failure
3. System displays clear error message about API key issues
4. User can update API key and retry analysis

## Success Criteria

### Performance Metrics
- Users can complete full CV-JD analysis in under 2 minutes
- System processes average PDF files (5 pages) in under 30 seconds
- 95% of analyses complete without technical errors
- System supports concurrent analysis of 10 users

### Quality Metrics  
- 90% of users report analysis accuracy meets or exceeds expectations
- 85% of users find visualizations helpful for decision making
- 80% of users can complete analysis without requiring support documentation
- User satisfaction score above 4.0/5.0 for overall experience

### Business Metrics
- Analysis completion rate above 75% for initiated sessions
- Average session duration under 5 minutes including review time
- 70% of users return for additional analyses within 30 days

## Key Entities

### Job Description
- **Content**: Raw text input by user
- **Validation**: Minimum length requirement (50 words)
- **Processing**: Used as reference for CV matching analysis

### CV Document  
- **Content**: PDF file uploaded by user
- **Conversion**: Processed to text format for analysis
- **Validation**: File format and readability checks

### Analysis Result
- **Matching Score**: Overall score /10 with dimension breakdowns
- **Quality Score**: CV assessment /10 with dimension breakdowns
- **Evidence**: Specific text excerpts supporting scoring decisions
- **Recommendations**: Interview questions and action items

### API Configuration
- **Key**: User-provided Gemini API key
- **Storage**: Browser memory only (session-based)
- **Validation**: Key format and access verification

## Assumptions

### Technical Assumptions
- Users have valid Gemini API keys with sufficient quota
- PDF files are text-based (not scanned images) for reliable processing
- Users have modern web browsers with JavaScript enabled
- Network connectivity is sufficient for API calls and file uploads
- Go 1.21+ runtime environment for backend services

### Business Assumptions  
- Job descriptions are provided in English language
- CVs are primarily professional documents with standard formatting
- Users are recruiters or hiring professionals familiar with candidate evaluation
- Analysis results are used for screening purposes, not final hiring decisions

### User Assumptions
- Users understand basic file upload and copy-paste operations
- Users can obtain API keys from external service providers
- Users have moderate technical literacy for web-based tools
- Users respect data privacy and don't upload sensitive personal information

## Dependencies

### External Services
- **Gemini API**: AI analysis engine for scoring and insights
- **UniDoc Library**: PDF to text conversion service (unidoc/unipdf)
- **Web Browser**: Runtime environment for user interface

### Data Requirements
- **Job Description Text**: Minimum 50 words for meaningful analysis
- **CV PDF Files**: Readable text-based PDF documents
- **API Keys**: Valid authentication tokens for external services

## Scope Boundaries

### In Scope
- CV-JD compatibility analysis with detailed scoring
- PDF document processing and text extraction
- Interactive data visualization with charts and graphs
- User interface for input management and results display
- Error handling for common failure scenarios

### Out of Scope
- Resume/CV creation or editing capabilities
- Database storage of analysis history
- Multi-user collaboration features
- Integration with applicant tracking systems
- Automated candidate sourcing or matching
- Compliance checking for employment regulations

## Constraints

### Technical Constraints
- Must function within web browser environment
- Cannot store sensitive data persistently without user consent
- Limited by external API rate limits and quotas
- Must handle varying PDF file sizes and qualities

### Business Constraints
- Must maintain data privacy and user confidentiality
- Cannot provide legal advice or employment recommendations
- Must clearly indicate AI-generated analysis limitations
- Should not replace human judgment in hiring decisions

### User Experience Constraints
- Interface must be intuitive for non-technical users
- Results must be easily understandable without technical expertise
- Processing time must be reasonable for user attention spans
- Error messages must be actionable and helpful

## Risk Mitigation

### Technical Risks
- **PDF Processing Failures**: Provide manual text input fallback
- **API Service Outages**: Clear error messaging and retry mechanisms  
- **Large File Handling**: File size limits and progress indicators
- **Browser Compatibility**: Test across major browsers and versions

### Business Risks
- **Analysis Accuracy**: Clear disclaimers about AI limitations
- **Data Privacy**: No persistent storage of sensitive information
- **User Expectations**: Transparent communication of capabilities
- **Dependency Risks**: Multiple fallback options for critical services

### User Experience Risks
- **Complexity**: Simplified interface with guided workflows
- **Error Recovery**: Clear paths to recover from failure states
- **Performance**: Optimized processing and responsive interface
- **Accessibility**: Compliance with accessibility standards
