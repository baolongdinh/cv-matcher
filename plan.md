# CV-JD Matching Tool - Technical Implementation Plan (v2.0)
 
## New Architectural Decisions (Bulk & Persistence)
 
### Queueing Strategy
**Decision**: Use Redis Lists for a lightweight Task Queue.
**Rationale**: 
- Simplifies the "Internal Tool" nature (no need for complex brokers like RabbitMQ).
- Native Go support via `go-redis`.
- Enables sequential processing to avoid Gemini API 429 (Rate Limit) errors.
 
### Persistence Strategy
**Decision**: Redis Hashes + Sorted Sets.
**Rationale**:
- **Hashes (`cv:result:{id}`)**: Efficiently store full JSON analysis objects.
- **Sorted Sets (`cv:ranking:{job_id}`)**: Store candidate IDs with their `matching_score` as the score. This allows Redis to handle the sorting natively: `ZREVRANGE cv:ranking:dev_job 0 -1 WITHSCORES`.
 
## Implementation Phases (Updated)
 
### Phase 2.1: Redis Integration (Backend)
1. **Redis Client Setup**
   - Configure connection pool in `internal/utils/redis.go`.
   - Add `REDIS_URL` to `.env`.

2. **History Service**
   - Create `internal/services/history_service.go`.
   - Implement `SaveResult(result)`, `GetHistory()`, and `GetRankedCandidates(jobID)`.

3. **Background Worker**
   - Implement a background loop in Go that pops CVs from the `cv_queue` and calls the `GeminiService`.
   - Update `AnalysisResult` with `batch_id` and `status` (pending/processing/done).
 
### Phase 2.2: Bulk API & Polling (Backend & Frontend)
1. **New Endpoints**
   - `POST /api/analyze/bulk`: Queues multiple files.
   - `GET /api/jobs/:id`: Returns current status of a batch.
   - `GET /api/history`: Returns persistent results from Redis.

2. **Frontend Polling**
   - Implement a `useQueue` composable to poll the status of a bulk upload until completion.
 
### Phase 2.3: HR Management View (Frontend)
1. **History Table Component**
   - List all candidates with their scores.
   - Implement sorting by clicking column headers (calling Redis Sorted Sets).
   - Add filtering by Job Title.

## Data Flow (Bulk Analysis)
1. **User**: Uploads 10 CVs.
2. **Backend**: 
   - Generates a `batch_id`.
   - Pushes 10 tasks to Redis `cv_queue`.
   - Returns `batch_id` to Frontend.
3. **Worker**: 
   - Pops task, extracts text, calls Gemini.
   - Saves JSON result to Redis Hash.
   - Adds Score to Redis Sorted Set for that Job.
4. **Frontend**: Polls `/api/jobs/{batch_id}` and updates progress bar.
5. **User**: Views "History" and sorts by "Matching Score".
