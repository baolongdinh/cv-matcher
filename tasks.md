# CV-JD Matching Tool - Implementation Tasks (v2.0)

## Phase 2: Redis & Bulk Processing

### Redis Infrastructure
- **TASK-039**: Setup Redis connection utility
  - Description: Initialize `go-redis` client with connection pooling
  - Files: `internal/utils/redis.go`

- **TASK-040**: Implement Redis Result Storage
  - Description: Logic to save/retrieve full analysis JSON in Redis Hashes
  - Files: `internal/services/history_service.go`

- **TASK-041**: Implement Redis Sorting Logic
  - Description: Use Redis Sorted Sets to store and retrieve candidate rankings
  - Files: `internal/services/history_service.go`

### Queue & Background Worker
- **TASK-042**: Create Task Queue Producer
  - Description: Endpoint to accept multiple files and push tasks to Redis list
  - Files: `internal/handlers/bulk_analyze.go`

- **TASK-043**: Implement Background Worker
  - Description: Loop to pop tasks, process with Gemini, and update status in Redis
  - Files: `internal/workers/analysis_worker.go`

- **TASK-044**: Implement Job Status Tracking
  - Description: API to return progress of a specific `batch_id`
  - Files: `internal/handlers/job_status.go`

### Frontend Enhancements
- **TASK-045**: Create Bulk Upload Component
  - Description: Support multi-file selection and batch upload UI
  - Files: `frontend/src/components/BulkUpload.vue`

- **TASK-046**: Implement Progress Tracking View
  - Description: Real-time progress bar for bulk processing jobs
  - Files: `frontend/src/views/JobProgressView.vue`

- **TASK-047**: Create HR History Dashboard
  - Description: Table view with sorting and filtering for all candidates
  - Files: `frontend/src/views/HistoryView.vue`

- **TASK-048**: Integrate Sorting with Redis
  - Description: Call history API with sort parameters
  - Files: `frontend/src/services/api.ts`
