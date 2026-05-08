# CV-JD Matching Tool - Technical Specification (v2.0 - Bulk & Redis)

## Feature Overview

A comprehensive CV-JD matching tool that enables recruiters to evaluate multiple candidates simultaneously through AI-powered bulk analysis. The system uses a queue-based architecture for reliability and Redis for persistent history and high-performance sorting.

## User Stories (New & Updated)

### US-06: Bulk CV Analysis
**As a recruiter**, I want to upload multiple CVs at once against a single JD, so that I can process a whole batch of candidates efficiently.

### US-07: Analysis Queue Tracking
**As a recruiter**, I want to see the progress of my uploaded CVs in real-time, so that I know how many are processed and how many are pending.

### US-08: Persistent Analysis History
**As an HR manager**, I want my previous analysis results to be saved, so that I can revisit them later without re-uploading or re-processing.

### US-09: Candidate Ranking & Sorting
**As a recruiter**, I want to sort my candidates by their matching scores, so that I can quickly identify the top-tier talent from a large pool.

## Functional Requirements

### FR-01: Input Processing (Updated)
- The system shall accept multiple CV files (up to 20 per batch) in PDF format.
- The system shall validate each file in the batch before queueing.

### FR-06: Bulk Queueing System
- The system shall use a Redis-based task queue to handle background processing.
- The system shall assign a unique `batch_id` to each bulk upload session.
- The system shall process CVs sequentially or in parallel depending on API quota limits.

### FR-07: Redis Persistence
- The system shall store every successful analysis result in Redis.
- The system shall index results by `overall_score` to enable high-speed sorting.
- The system shall store metadata such as `upload_date`, `candidate_name`, and `job_title`.

### FR-08: Candidate Management Dashboard
- The system shall provide a "History" view displaying all analyzed candidates.
- The system shall allow sorting by:
  - Overall Matching Score (Highest to Lowest)
  - CV Quality Score
  - Upload Date (Latest First)
- The system shall support filtering by Job Title or Score Range.

## Technical Architecture (v2.0)

### Infrastructure
- **Redis**: Primary data store for task queueing (Asynq/Custom List) and result persistence (Hashes/Sorted Sets).
- **Background Worker**: Dedicated Go routine/service to consume CVs from Redis and call Gemini API.

### Data Storage Strategy
- **Result Hashes**: `cv:result:{id}` stores the full JSON analysis.
- **Score Sorted Set**: `cv:scores:{job_id}` stores candidate IDs ordered by score for O(log N) sorting.
- **Batch Metadata**: `batch:{id}` tracks the progress (total, completed, failed).

## Success Criteria (Updated)

### Performance Metrics
- Sorting through 1,000+ candidate results should take under 100ms.
- Queue processing should gracefully handle Gemini API rate limits (429 errors) with automatic retries.

### Quality Metrics
- Zero data loss for analysis results once stored in Redis.
- Real-time UI updates should reflect accurate queue status.
