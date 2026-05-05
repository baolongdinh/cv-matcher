# CV-JD Matcher Tool

A professional internal tool for recruiters to analyze CVs against Job Descriptions using Google Gemini AI.
<img width="3126" height="3818" alt="10 0 0 125_5173_" src="https://github.com/user-attachments/assets/d1359288-9e0b-4e6d-acaa-11a0954f1eb2" />

## Features
- **AI-Powered Matching**: Detailed scoring across 6 dimensions.
- **CV Quality Analysis**: Independent assessment of resume professionalism.
- **Visual Insights**: Radar charts and bar charts for quick evaluation.
- **Interview Guide**: AI-generated questions based on identified gaps.
- **Privacy First**: API keys and documents are processed and kept in session/memory.

## Project Structure
- `/backend`: Go + Gin API server.
- `/frontend`: Vue 3 + Tailwind CSS dashboard.

## Quick Start

### 1. Prerequisites
- Go 1.21+
- Node.js 18+
- Google Gemini API Key ([Get it here](https://aistudio.google.com/app/apikey))

### 2. Backend Setup
```bash
cd backend
cp .env.example .env
# No API key needed in .env - users provide their own in the UI
go run cmd/server/main.go
```
The API will be available at `http://localhost:8000`.

### 3. Frontend Setup
```bash
cd frontend
npm install
npm run dev
```
The UI will be available at `http://localhost:5173`.

## UI Orientation
The UI is designed as a professional "Recruiter Command" dashboard:
- **Workspace**: Input your API key, JD, and upload the CV PDF.
- **Match Summary**: View overall scores, competency matrix, strengths, and gaps.
- **Analysis Results**: Deep dive into AI-generated recommendations and red flags.

## Tech Stack
- **Backend**: Go, Gin, UniDoc (PDF processing), Gemini SDK.
- **Frontend**: Vue 3, TypeScript, Tailwind CSS, Chart.js, Lucide Icons.
