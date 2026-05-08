# CV-JD Matching Tool - API Contracts (v2.0)
 
## Endpoints
 
### POST /api/analyze/bulk
Upload multiple CVs for background processing.
 
**Request:** `multipart/form-data`
- `job_description`: string
- `cv_files`: file[] (multiple)
- `api_key`: string
 
**Response:** `202 Accepted`
```json
{
  "status": "success",
  "data": {
    "batch_id": "batch_987654321",
    "total_files": 15,
    "message": "Files queued for processing"
  }
}
```
 
### GET /api/jobs/:batch_id
Check status of a bulk processing job.
 
**Response:** `200 OK`
```json
{
  "status": "success",
  "data": {
    "batch_id": "batch_987654321",
    "status": "processing",
    "progress": {
      "total": 15,
      "completed": 7,
      "failed": 0,
      "percentage": 46.6
    },
    "results": [
      { "id": "cv_1", "status": "completed", "score": 8.5 },
      { "id": "cv_2", "status": "pending" }
    ]
  }
}
```
 
### GET /api/history
Retrieve analyzed candidates from Redis.
 
**Query Parameters:**
- `job_id`: string (optional filter)
- `sort_by`: `matching_score` | `quality_score` | `date`
- `order`: `asc` | `desc`
 
**Response:** `200 OK`
```json
{
  "status": "success",
  "data": [
    {
      "candidate_id": "req_123",
      "candidate_name": "Nguyen Van A",
      "job_title": "Senior Frontend Developer",
      "matching_score": 9.2,
      "quality_score": 8.5,
      "timestamp": "2024-05-05T10:30:00Z"
    }
  ]
}
```
 
### GET /api/history/:candidate_id
Get full analysis details for a specific candidate from Redis history.
 
**Response:** `200 OK` (Same format as single `/api/analyze`)
