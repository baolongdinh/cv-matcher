# CV-JD Matching Tool - Quick Start Guide
 
## Overview
 
Get your CV-JD matching tool running in under 30 minutes with this comprehensive setup guide.
 
## Prerequisites
 
### Required Software
- **Go 1.21+** - Backend development
- **Node.js 18+** - Frontend development  
- **Git** - Version control
- **Modern web browser** - Frontend testing
 
### Required Accounts
- **Google AI Studio** - For Gemini API key
  - Visit https://aistudio.google.com/app/apikey
  - Create new API key
  - Copy key for setup
 
## Project Setup
 
### 1. Clone Repository
```bash
git clone <repository-url>
cd cv-jd-matcher
```
 
### 2. Backend Setup
 
#### Initialize Go Module
```bash
# Initialize Go module
go mod init cv-jd-matcher

# Create go.mod with dependencies
cat > go.mod << 'EOF'
module cv-jd-matcher

go 1.21

require (
    github.com/gin-gonic/gin v1.9.1
    github.com/gin-contrib/cors v1.4.0
    github.com/joho/godotenv v1.4.0
    github.com/unidoc/unipdf/v3 v3.0.0
    google.golang.org/api v0.149.0
    google.golang.org/genai v0.1.0
)
EOF

# Download dependencies
go mod tidy
```
 
#### Environment Configuration
```bash
# Create .env file
cat > .env << EOF
# Gemini API Configuration
GEMINI_API_KEY=your_api_key_here
 
# Server Configuration  
HOST=0.0.0.0
PORT=8000
DEBUG=true
 
# CORS Configuration
CORS_ORIGINS=http://localhost:3000,http://127.0.0.1:3000
 
# File Upload Limits
MAX_FILE_SIZE=10485760  # 10MB
MAX_PAGES=50
 
# Rate Limiting
RATE_LIMIT_PER_MINUTE=10
EOF
 
# Replace with your actual API key
sed -i 's/your_api_key_here/ACTUAL_API_KEY/' .env
```
 
#### Create Backend Structure
```bash
mkdir -p cmd/server internal/{handlers,services,models,utils}
```
 
#### Create Main Application
```bash
# Create cmd/server/main.go
cat > cmd/server/main.go << 'EOF'
package main

import (
    "log"
    "os"
    "github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"
    "github.com/joho/godotenv"
)

func main() {
    // Load environment variables
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found")
    }

    // Initialize Gin router
    r := gin.Default()

    // CORS middleware
    config := cors.DefaultConfig()
    config.AllowOrigins = []string{os.Getenv("CORS_ORIGINS", "http://localhost:3000")}
    config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
    config.AllowHeaders = []string{"*"}
    r.Use(cors.New(config))

    // Basic routes
    r.GET("/", func(c *gin.Context) {
        c.JSON(200, gin.H{"message": "CV-JD Matcher API"})
    })

    // TODO: Add API routes
    // r.POST("/api/analyze", handlers.Analyze)
    // r.GET("/api/health", handlers.Health)

    port := os.Getenv("PORT", "8000")
    r.Run(":" + port)
}
EOF
```
 
#### Start Backend Server
```bash
go run cmd/server/main.go
```

**Expected Output:**
```
[GIN-debug] [WARNING] Creating an Engine instance with the Logger and Recovery middleware already attached.

[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:   export GIN_MODE=release
 - using code:  gin.SetMode(gin.ReleaseMode)

[GIN-debug] GET    /                         --> main.main.func1 (3 handlers)
[GIN-debug] Listening and serving HTTP on :8000
```
 
### 3. Frontend Setup
 
#### Create Vue Project
```bash
# In new terminal
cd cv-jd-matcher
 
# Create Vue project
npm create vue@latest frontend
cd frontend
 
# Select options:
# ✅ TypeScript
# ✅ Router  
# ❌ Pinia (we'll use composition API)
# ✅ ESLint
# ❌ Prettier
```
 
#### Install Dependencies
```bash
npm install
 
# Additional dependencies
npm install axios chart.js vue-chartjs
npm install @types/node -D
npm install tailwindcss postcss autoprefixer -D
npm install @headlessui/vue @heroicons/vue
```
 
#### Configure Tailwind CSS
```bash
# Initialize Tailwind
npx tailwindcss init -p
 
# Configure tailwind.config.js
cat > tailwind.config.js << 'EOF'
/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{vue,js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {},
  },
  plugins: [],
}
EOF
 
# Update CSS
cat > src/assets/main.css << 'EOF'
@tailwind base;
@tailwind components;
@tailwind utilities;
EOF
```
 
#### Create Project Structure
```bash
mkdir -p src/{components,views,services,types,composables,utils}
```
 
#### Create Main Application
```bash
# Create App.vue
cat > src/App.vue << 'EOF'
<template>
  <div id="app">
    <header class="bg-white shadow">
      <div class="max-w-7xl mx-auto py-6 px-4">
        <h1 class="text-3xl font-bold text-gray-900">CV-JD Matcher</h1>
      </div>
    </header>
 
    <main class="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
      <HomeView />
    </main>
  </div>
</template>
 
<script setup lang="ts">
import HomeView from './views/HomeView.vue'
</script>
 
<style>
#app {
  min-height: 100vh;
  background-color: #f9fafb;
}
</style>
EOF
 
# Create HomeView
cat > src/views/HomeView.vue << 'EOF'
<template>
  <div class="px-4 py-6 sm:px-0">
    <div class="border-4 border-dashed border-gray-200 rounded-lg p-8">
      <h2 class="text-2xl font-semibold mb-4">CV-JD Analysis Tool</h2>
      <p class="text-gray-600 mb-6">
        Upload your CV and paste a job description to get a comprehensive analysis.
      </p>
 
      <!-- Placeholder for components -->
      <div class="bg-blue-50 p-4 rounded-lg">
        <p class="text-blue-800">🚧 Application setup complete! Components coming next...</p>
      </div>
    </div>
  </div>
</template>
 
<script setup lang="ts">
// Component implementation will go here
</script>
EOF
```
 
#### Start Frontend Server
```bash
npm run dev
```
 
**Expected Output:**
```
  VITE v5.0.0  ready in 523 ms
 
  ➜  Local:   http://localhost:3000/
  ➜  Network: http://192.168.1.100:3000/
  ➜  press h to show help
```
 
### 4. Verify Setup
 
#### Test Backend
```bash
# Test health endpoint
curl http://localhost:8000/
 
# Expected: {"message": "CV-JD Matcher API"}
```
 
#### Test Frontend
- Open browser to http://localhost:3000
- Should see "CV-JD Analysis Tool" heading
 
## Development Workflow
 
### 1. Component Development
```bash
# Create input component
cat > frontend/src/components/InputPanel.vue << 'EOF'
<template>
  <div class="bg-white p-6 rounded-lg shadow">
    <h3 class="text-lg font-semibold mb-4">Input Panel</h3>
    <!-- Component implementation -->
  </div>
</template>
 
<script setup lang="ts">
// Component logic
</script>
EOF
```
 
### 2. API Integration
```bash
# Create API service
cat > frontend/src/services/api.ts << 'EOF'
import axios from 'axios'
 
const API_BASE_URL = 'http://localhost:8000/api'
 
export const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
})
 
export default api
EOF
```
 
### 3. Type Definitions
```bash
# Create types
cat > frontend/src/types/analysis.ts << 'EOF'
export interface AnalysisRequest {
  jobDescription: string
  cvFile: File
  apiKey: string
}
 
export interface AnalysisResult {
  matchingScore: MatchingScore
  qualityScore: QualityScore
  // ... other interfaces
}
 
interface MatchingScore {
  overall: number
  skillsAlignment: DimensionScore
  // ... other dimensions
}
 
interface DimensionScore {
  score: number
  evidence: string[]
  gapsIdentified: string[]
  confidenceLevel: number
}
EOF
```
 
## Testing
 
### Backend Testing
```bash
# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific test
go test ./internal/handlers -v
```
 
### Frontend Testing
```bash
# Install test dependencies
npm install -D vitest @vue/test-utils jsdom
 
# Create test file
cat > tests/components.test.ts << 'EOF'
import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import InputPanel from '../src/components/InputPanel.vue'
 
describe('InputPanel', () => {
  it('renders properly', () => {
    const wrapper = mount(InputPanel)
    expect(wrapper.text()).toContain('Input Panel')
  })
})
EOF
 
# Run tests
npm run test
```
 
## Common Issues & Solutions
 
### Backend Issues
 
**Issue:** ModuleNotFoundError for gin
```bash
# Solution: Ensure Go module is initialized
go mod init cv-jd-matcher
 
# Reinstall if needed
go mod tidy
```
 
**Issue:** Port already in use
```bash
# Solution: Change port or kill existing process
# Change port in .env:
PORT=8001
 
# Or kill existing process:
lsof -ti:8000 | xargs kill -9
```
 
**Issue:** Gemini API key errors
```bash
# Solution: Verify API key format and quota
# Key should start with "AIza"
curl -H "Content-Type: application/json" \
     -d '{"contents":[{"parts":[{"text":"Hello"}]}]}' \
     "https://generativelanguage.googleapis.com/v1beta/models/gemini-2.0-flash:generateContent?key=YOUR_API_KEY"
```
 
### Frontend Issues
 
**Issue:** npm create vue fails
```bash
# Solution: Use alternative method
npm init vue@latest frontend
# Or
yarn create vue frontend
```
 
**Issue:** Tailwind CSS not working
```bash
# Solution: Verify PostCSS configuration
cat > postcss.config.js << 'EOF'
export default {
  plugins: {
    tailwindcss: {},
    autoprefixer: {},
  },
}
EOF
```
 
**Issue:** CORS errors
```bash
# Solution: Verify CORS origins in backend .env
# Ensure frontend URL is listed
CORS_ORIGINS=http://localhost:3000,http://127.0.0.1:3000
```
 
## Production Deployment
 
### Backend Deployment
```bash
# Build binary
go build -o bin/server cmd/server/main.go

# Run production binary
./bin/server

# Or build for specific platform
go build -o bin/server-linux cmd/server/main.go
GOOS=windows GOARCH=amd64 go build -o bin/server.exe cmd/server/main.go
```
 
### Frontend Deployment
```bash
# Build for production
npm run build
 
# Output in dist/ directory
# Serve with any static web server
python -m http.server 3000 --directory dist
```
 
## Next Steps
 
1. **Implement Core Components**
   - InputPanel.vue for file upload and text input
   - ScoreCard.vue for displaying results
   - Chart components for visualizations
 
2. **Add API Endpoints**
   - POST /api/analyze for CV-JD analysis
   - GET /api/health for service monitoring
 
3. **Integration Testing**
   - End-to-end workflow testing
   - Error handling validation
   - Performance optimization
 
4. **Polish & Refinement**
   - UI/UX improvements
   - Error message refinement
   - Accessibility compliance
 
## Support Resources
 
- **Backend Documentation**: https://gin-gonic.com/docs/
- **Frontend Documentation**: https://vuejs.org/guide/
- **Gemini API**: https://ai.google.dev/docs
- **UniDoc**: https://unidoc.io/
- **Chart.js**: https://www.chartjs.org/docs/
 
## Troubleshooting Checklist
 
- [ ] Go 1.21+ installed and working
- [ ] All backend dependencies downloaded without errors
- [ ] Backend server starts on port 8000
- [ ] Node.js 18+ installed