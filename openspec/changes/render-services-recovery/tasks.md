# Render Services Recovery - Tasks

## Phase 1: Diagnosis & Assessment (15 Tasks)

### 1.1 Service Status Check
- [ ] **1.1.1** Check Render dashboard for bac-api service status
- [ ] **1.1.2** Check Render dashboard for bac-agent service status
- [ ] **1.1.3** Verify GitHub repo https://github.com/amdajhbdh/bac exists and is accessible

### 1.2 Build Verification
- [ ] **1.2.1** Run `cd src/api && go build -o bac-api .` locally - verify it succeeds
- [ ] **1.2.2** Run `cd src/agent && go build -o agent ./cmd/main.go` locally - verify it succeeds
- [ ] **1.2.3** Check for any compile errors in Go code

### 1.3 Environment Variables Audit
- [ ] **1.3.1** Document current NEON_DB_URL (or mark as missing)
- [ ] **1.3.2** Document current TURSO_DB_URL (or mark as missing)
- [ ] **1.3.3** Document current REDIS_URL (or mark as missing)
- [ ] **1.3.4** Document current JWT_SECRET (or mark as missing)
- [ ] **1.3.5** Document current OLLAMA_URL (or mark as missing)

### 1.4 External Service Status
- [ ] **1.4.1** Verify Neon DB bac-project is still active
- [ ] **1.4.2** Verify Turso database is still active
- [ ] **1.4.3** Check if Redis instance is available

---

## Phase 2: Render Reconfiguration (45 Tasks)

### 2.1 Project Setup
- [ ] **2.1.1** Login to Render dashboard
- [ ] **2.1.2** Navigate to existing project or create new
- [ ] **2.1.3** Connect GitHub repository https://github.com/amdajhbdh/bac
- [ ] **2.1.4** Configure auto-deploy on push to first_bm branch
- [ ] **2.1.5** Set region to Oregon (us-west)

### 2.2 bac-api Service
- [ ] **2.2.1** Create web service: name=bac-api, type=web
- [ ] **2.2.2** Set root directory: src/api
- [ ] **2.2.3** Set build command: cd src/api && go build -o bac-api .
- [ ] **2.2.4** Set start command: ./bac-api
- [ ] **2.2.5** Set PORT environment variable: 8080
- [ ] **2.2.6** Add environment variable: NEON_DB_URL (from Neon console)
- [ ] **2.2.7** Add environment variable: TURSO_DB_URL (from Turso)
- [ ] **2.2.8** Add environment variable: REDIS_URL
- [ ] **2.2.9** Add environment variable: JWT_SECRET (generate with openssl rand -hex 32)
- [ ] **2.2.10** Add environment variable: OLLAMA_URL (https://bac-api.amdajhbdh.workers.dev)
- [ ] **2.2.11** Configure health check: GET /health, interval 1min
- [ ] **2.2.12** Deploy and check build logs for success

### 2.3 bac-agent Service
- [ ] **2.3.1** Create web service: name=bac-agent, type=web
- [ ] **2.3.2** Set root directory: src/agent
- [ ] **2.3.3** Set build command: cd src/agent && go build -o agent ./cmd/main.go
- [ ] **2.3.4** Set start command: ./agent -server
- [ ] **2.3.5** Set PORT environment variable: 8081
- [ ] **2.3.6** Add environment variable: NEON_DB_URL
- [ ] **2.3.7** Add environment variable: TURSO_DB_URL
- [ ] **2.3.8** Add environment variable: REDIS_URL
- [ ] **2.3.9** Add environment variable: OLLAMA_URL
- [ ] **2.3.10** Configure health check: GET /health, interval 1min
- [ ] **2.3.11** Deploy and check build logs for success

### 2.4 Keepalive Configuration
- [ ] **2.4.1** Decide: External pinger OR paid tier OR Cloudflare primary
- [ ] **2.4.2** If external pinger: Create simple ping script
- [ ] **2.4.3** If external pinger: Set up cron job to ping every 10 minutes
- [ ] **2.4.4** If paid tier: Note cost in project documentation

---

## Phase 3: Environment Variables Setup (30 Tasks)

### 3.1 Database Configuration
- [ ] **3.1.1** Get Neon connection string from console
- [ ] **3.1.2** Test Neon connection from local machine
- [ ] **3.1.3** Get Turso connection string from dashboard
- [ ] **3.1.4** Test Turso connection from local machine
- [ ] **3.1.5** Configure connection pooling settings in Render
- [ ] **3.1.6** Set connection timeout to 30 seconds
- [ ] **3.1.7** Configure max connections based on Render tier

### 3.2 Security Configuration
- [ ] **3.2.1** Generate JWT_SECRET: openssl rand -hex 32
- [ ] **3.2.2** Add JWT_SECRET to Render environment variables
- [ ] **3.2.3** Configure JWT expiration: 24 hours access, 7 days refresh
- [ ] **3.2.4** Configure CORS to allow appropriate origins
- [ ] **3.2.5** Review and remove any hardcoded secrets in code

### 3.3 AI Service Configuration
- [ ] **3.3.1** Set OLLAMA_URL to Cloudflare Worker fallback
- [ ] **3.3.2** Configure AI model priority: Ollama → Cloudflare → DeepSeek → OpenAI
- [ ] **3.3.3** Set AI call timeout to 30 seconds
- [ ] **3.3.4** Configure retry logic: 3 retries with exponential backoff

### 3.4 External Services
- [ ] **3.4.1** Set up Elasticsearch endpoint (or mark as future)
- [ ] **3.4.2** Configure Garage/MinIO storage endpoint
- [ ] **3.4.3** Set storage credentials securely

---

## Phase 4: Database Setup & Migration (35 Tasks)

### 4.1 Schema Management
- [ ] **4.1.1** Verify Neon schema tables exist
- [ ] **4.1.2** Create migrations folder if not exists: sql/migrations/
- [ ] **4.1.3** Write users, questions, predictions tables initial migration for
- [ ] **4.1.4** Enable vector extension: CREATE EXTENSION vector
- [ ] **4.1.5** Enable timescaledb extension: CREATE EXTENSION timescaledb
- [ ] **4.1.6** Create user_progress table
- [ ] **4.1.7** Create badges table
- [ ] **4.1.8** Create sessions table
- [ ] **4.1.9** Run migrations on Neon
- [ ] **4.1.10** Verify schema version in migration history

### 4.2 Data Seeding
- [ ] **4.2.1** Add 50 sample math questions
- [ ] **4.2.2** Add 30 sample physics questions
- [ ] **4.2.3** Add 20 sample chemistry questions
- [ ] **4.2.4** Add step-by-step solutions to all questions
- [ ] **4.2.5** Add difficulty ratings (1-3 scale)
- [ ] **4.2.6** Add subject tags: math, physics, chemistry
- [ ] **4.2.7** Add chapter tags for curriculum mapping
- [ ] **4.2.8** Generate embeddings for vector search
- [ ] **4.2.9** Import historical exam questions from db/pdfs/
- [ ] **4.2.10** Create admin user for testing
- [ ] **4.2.11** Create 5-10 test users
- [ ] **4.2.12** Verify data integrity - check for duplicates

### 4.3 Vector Search Setup
- [ ] **4.3.1** Generate question embeddings using bge-small model
- [ ] **4.3.2** Store embeddings in pgvector (384 dimensions)
- [ ] **4.3.3** Create vector index (HNSW or IVFFlat)
- [ ] **4.3.4** Test similarity search
- [ ] **4.3.5** Add hybrid search (vector + full-text)
- [ ] **4.3.6** Optimize index for recall/performance

### 4.4 Turso Cache Setup
- [ ] **4.4.1** Verify Turso connection from bac-api
- [ ] **4.4.2** Create cache schema
- [ ] **4.4.3** Configure TTL settings
- [ ] **4.4.4** Set up cache invalidation on data changes
- [ ] **4.4.5** Test cache hit rate

---

## Phase 5: Keepalive & Reliability (25 Tasks)

### 5.1 Health Checks
- [ ] **5.1.1** Add GET /health endpoint to bac-api
- [ ] **5.1.2** Add GET /health endpoint to bac-agent
- [ ] **5.1.3** Add GET /health/db endpoint to check DB connection
- [ ] **5.1.4** Add GET /health/redis endpoint to check Redis
- [ ] **5.1.5** Add GET /health/ai endpoint to check AI service
- [ ] **5.1.6** Configure Render health checks on dashboard
- [ ] **5.1.7** Set health check interval to 1 minute
- [ ] **5.1.8** Test health check behavior

### 5.2 Auto-Wake Solutions
- [ ] **5.2.1** Implement external ping script
- [ ] **5.2.2** Set up cron job: */10 * * * * curl https://bac-api.onrender.com/health
- [ ] **5.2.3** Or: Upgrade to Render paid tier (always-on)
- [ ] **5.2.4** Or: Use Cloudflare as primary, disable Render auto-sleep concern
- [ ] **5.2.5** Set up UptimeRobot for monitoring
- [ ] **5.2.6** Configure Slack alerts on service down
- [ ] **5.2.7** Set up email notifications
- [ ] **5.2.8** Document keepalive solution in README

### 5.3 Rate Limiting
- [ ] **5.3.1** Implement per-user rate limit: 100 req/min
- [ ] **5.3.2** Implement per-IP rate limit: 50 req/min
- [ ] **5.3.3** Add rate limit headers: X-RateLimit-Limit, X-RateLimit-Remaining
- [ ] **5.3.4** Configure rate limit storage in Redis
- [ ] **5.3.5** Add 429 response with retry-after header

---

## Phase 6: API Enhancements (40 Tasks)

### 6.1 Core Endpoints Verification
- [ ] **6.1.1** Test POST /solve with math problem
- [ ] **6.1.2** Test GET /questions returns data
- [ ] **6.1.3** Test POST /questions adds new question
- [ ] **6.1.4** Test GET /leaderboard
- [ ] **6.1.5** Test POST /auth/register
- [ ] **6.1.6** Test POST /auth/login
- [ ] **6.1.7** Test POST /practice/answer

### 6.2 Additional Endpoints
- [ ] **6.2.1** Add GET /questions/:id endpoint
- [ ] **6.2.2** Add POST /questions/search endpoint (full-text)
- [ ] **6.2.3** Add GET /questions/random endpoint
- [ ] **6.2.4** Add GET /subjects endpoint
- [ ] **6.2.5** Add GET /chapters endpoint
- [ ] **6.2.6** Add GET /users/profile endpoint
- [ ] **6.2.7** Add GET /users/progress endpoint
- [ ] **6.2.8** Add GET /predictions endpoint
- [ ] **6.2.9** Add GET /badges endpoint

### 6.3 File Handling
- [ ] **6.3.1** Add POST /upload/image endpoint
- [ ] **6.3.2** Add POST /upload/pdf endpoint
- [ ] **6.3.3** Add POST /ocr endpoint (image to text)
- [ ] **6.3.4** Configure file size limit: 10MB
- [ ] **6.3.5** Add file type validation

### 6.4 AI Enhancements
- [ ] **6.4.1** Add POST /solve/stream endpoint (streaming response)
- [ ] **6.4.2** Add POST /solve/steps endpoint (step-by-step)
- [ ] **6.4.3** Add POST /explain endpoint (explain concept)
- [ ] **6.4.4** Add POST /hint endpoint (give hints)
- [ ] **6.4.5** Add POST /similar endpoint (find similar problems)
- [ ] **6.4.6** Add POST /practice/generate endpoint
- [ ] **6.4.7** Add POST /quiz/generate endpoint
- [ ] **6.4.8** Add AI response caching
- [ ] **6.4.9** Implement fallback chain: Ollama → Cloudflare → DeepSeek
- [ ] **6.4.10** Add model selection parameter
- [ ] **6.4.11** Add temperature control

---

## Phase 7: Cloudflare Integration (30 Tasks)

### 7.1 Architecture Decision
- [ ] **7.1.1** Decide: Cloudflare-first OR Render-first
- [ ] **7.1.2** Document architecture decision
- [ ] **7.1.3** Set up unified monitoring dashboard

### 7.2 Cloudflare-Primary Implementation
- [ ] **7.2.1** Keep Cloudflare Worker as main API
- [ ] **7.2.2** Add D1 schema for complete questions table
- [ ] **7.2.3** Seed D1 with questions from Neon
- [ ] **7.2.4** Use Workers AI as primary AI
- [ ] **7.2.5** Use Vectorize for search
- [ ] **7.2.6** Enable R2 storage in Cloudflare dashboard
- [ ] **7.2.7** Add R2 file upload endpoint
- [ ] **7.2.8** Keep Render as cold backup
- [ ] **7.2.9** Add periodic sync: Neon → D1 daily

### 7.3 Render Backup Implementation
- [ ] **7.3.1** Keep bac-api on Render for redundancy
- [ ] **7.3.2** Add failover logic: Try Render, then Cloudflare
- [ ] **7.3.3** Configure Cloudflare Workers as origin
- [ ] **7.3.4** Test failover manually

---

## Phase 8: Frontend Development (45 Tasks)

### 8.1 Cloudflare Pages Deployment
- [ ] **8.1.1** Deploy src/cloudflare-pages/ to Cloudflare Pages
- [ ] **8.1.2** Configure custom domain (if available)
- [ ] **8.1.3** Enable HTTPS
- [ ] **8.1.4** Configure caching for static assets
- [ ] **8.1.5** Add 404 error page

### 8.2 Frontend Features
- [ ] **8.2.1** Update API_URL in frontend to point to Cloudflare Worker
- [ ] **8.2.2** Test user registration UI
- [ ] **8.2.3** Test login UI with JWT storage
- [ ] **8.2.4** Test problem submission (text input)
- [ ] **8.2.5** Test solution display
- [ ] **8.2.6** Test practice mode
- [ ] **8.2.7** Test progress tracking
- [ ] **8.2.8** Test leaderboard display

### 8.3 User Experience
- [ ] **8.3.1** Add loading states (spinners)
- [ ] **8.3.2** Add error handling with user-friendly messages
- [ ] **8.3.3** Add dark mode support
- [ ] **8.3.4** Add offline support with service worker
- [ ] **8.3.5** Add PWA manifest for installability
- [ ] **8.3.6** Add French translations
- [ ] **8.3.7** Add Arabic translations (RTL support)
- [ ] **8.3.8** Test responsive design on mobile
- [ ] **8.3.9** Test keyboard navigation
- [ ] **8.3.10** Add screen reader support (ARIA labels)

---

## Phase 9: Testing (35 Tasks)

### 9.1 Unit Tests
- [ ] **9.1.1** Write tests for solver package
- [ ] **9.1.2** Write tests for OCR package
- [ ] **9.1.3** Write tests for memory package
- [ ] **9.1.4** Write tests for auth package
- [ ] **9.1.5** Write tests for API handlers
- [ ] **9.1.6** Run test suite: go test ./...
- [ ] **9.1.7** Achieve >80% test coverage

### 9.2 Integration Tests
- [ ] **9.2.1** Test API ↔ Database CRUD operations
- [ ] **9.2.2** Test API ↔ AI solve endpoint
- [ ] **9.2.3** Test API ↔ Storage file operations
- [ ] **9.2.4** Test full auth flow: register → login → protected
- [ ] **9.2.5** Test end-to-end solve flow
- [ ] **9.2.6** Set up automated integration tests in CI

### 9.3 Manual Testing
- [ ] **9.3.1** Test user registration with valid/invalid inputs
- [ ] **9.3.2** Test login with correct/wrong passwords
- [ ] **9.3.3** Test problem solving with various inputs
- [ ] **9.3.4** Test file upload (small, large, invalid)
- [ ] **9.3.5** Test search functionality
- [ ] **9.3.6** Test leaderboard points calculation

### 9.4 Load Testing
- [ ] **9.4.1** Set up k6 or wrk for load testing
- [ ] **9.4.2** Test baseline: single user response time
- [ ] **9.4.3** Test 10 concurrent users
- [ ] **9.4.4** Test 50 concurrent users
- [ ] **9.4.5** Identify and fix bottlenecks

---

## Phase 10: Monitoring & Observability (30 Tasks)

### 10.1 Metrics
- [ ] **10.1.1** Set up Prometheus for metrics collection
- [ ] **10.1.2** Create Grafana dashboard for visualization
- [ ] **10.1.3** Add request count, latency, error rate metrics
- [ ] **10.1.4** Add AI call metrics (latency, cost)
- [ ] **10.1.5** Add database metrics (connections, query time)
- [ ] **10.1.6** Add cache hit rate metrics
- [ ] **10.1.7** Set up alerts for >5% error rate
- [ ] **10.1.8** Set up alerts for >5s latency

### 10.2 Logging
- [ ] **10.2.1** Configure structured JSON logging
- [ ] **10.2.2** Set log levels: Debug, Info, Warn, Error
- [ ] **10.2.3** Add correlation IDs for request tracing
- [ ] **10.2.4** Set up log aggregation
- [ ] **10.2.5** Create log dashboards for errors

### 10.3 Uptime Monitoring
- [ ] **10.3.1** Set up UptimeRobot (free tier)
- [ ] **10.3.2** Monitor bac-api /health endpoint
- [ ] **10.3.3** Monitor bac-agent /health endpoint
- [ ] **10.3.4** Monitor Cloudflare Worker /health endpoint
- [ ] **10.3.5** Configure Slack notifications on downtime

---

## Phase 11: Security (25 Tasks)

### 11.1 Authentication
- [ ] **11.1.1** Verify JWT implementation
- [ ] **11.1.2** Verify password hashing (bcrypt)
- [ ] **11.1.3** Add TOTP two-factor authentication
- [ ] **11.1.4** Implement session management
- [ ] **11.1.5** Add rate limiting for auth endpoints

### 11.2 Authorization
- [ ] **11.2.1** Implement RBAC: admin, teacher, student roles
- [ ] **11.2.2** Add role-based endpoint protection
- [ ] **11.2.3** Add resource ownership checks
- [ ] **11.2.4** Add audit logging for sensitive operations

### 11.3 Data Protection
- [ ] **11.3.1** Verify TLS on all connections
- [ ] **11.3.2** Configure Cloudflare WAF rules
- [ ] **11.3.3** Add input sanitization (prevent injection)
- [ ] **11.3.4** Add output encoding (prevent XSS)
- [ ] **11.3.5** Configure CORS properly
- [ ] **11.3.6** Add security headers (HSTS, CSP)

---

## Phase 12: Documentation (20 Tasks)

### 12.1 API Documentation
- [ ] **12.1.1** Verify Swagger/OpenAPI documentation
- [ ] **12.1.2** Add request/response examples
- [ ] **12.1.3** Document error codes
- [ ] **12.1.4** Document rate limits

### 12.2 Developer Documentation
- [ ] **12.2.1** Update README with current architecture
- [ ] **12.2.2** Add local development setup guide
- [ ] **12.2.3** Add deployment instructions
- [ ] **12.2.4** Add troubleshooting section
- [ ] **12.2.5** Document environment variables

### 12.3 User Documentation
- [ ] **12.3.1** Create user guide
- [ ] **12.3.2** Add FAQ section
- [ ] **12.3.3** Add video tutorials (optional)

---

## Phase 13: Content Pipeline (30 Tasks)

### 13.1 Question Bank Expansion
- [ ] **13.1.1** Add 100 more math questions
- [ ] **13.1.2** Add 50 more physics questions
- [ ] **13.1.3** Add 50 more chemistry questions
- [ ] **13.1.4** Add solutions to new questions
- [ ] **13.1.5** Tag questions by exam year
- [ ] **13.1.6** Add images to visual questions

### 13.2 Content Aggregation
- [ ] **13.2.1** Set up YouTube content scraper
- [ ] **13.2.2** Set up web scraper for education sites
- [ ] **13.2.3** Download and index PDFs
- [ ] **13.2.4** Integrate OpenStax textbooks

### 13.3 Media Processing
- [ ] **13.3.1** Process uploaded images (resize, optimize)
- [ ] **13.3.2** Extract text from PDFs
- [ ] **13.3.3** Generate thumbnails
- [ ] **13.3.4** Add virus scanning for uploads

---

## Phase 14: Advanced Features (35 Tasks)

### 14.1 Animation System
- [ ] **14.1.1** Deploy Noon animation engine
- [ ] **14.1.2** Add POST /animate endpoint
- [ ] **14.1.3** Create math animation templates
- [ ] **14.1.4** Create physics animation templates
- [ ] **14.1.5** Add AI-to-Noon code generator
- [ ] **14.1.6** Add video export (MP4, GIF)
- [ ] **14.1.7** Cache generated animations

### 14.2 Prediction Engine
- [ ] **14.2.1** Analyze question frequency patterns
- [ ] **14.2.2** Identify exam trends
- [ ] **14.2.3** Build prediction model
- [ ] **14.2.4** Add GET /predictions/:subject endpoint
- [ ] **14.2.5** Add confidence scores
- [ ] **14.2.6** Track prediction accuracy over time

### 14.3 Gamification
- [ ] **14.3.1** Define points system
- [ ] **14.3.2** Create badge definitions
- [ ] **14.3.3** Implement levels (Bronze, Silver, Gold)
- [ ] **14.3.4** Add daily challenges
- [ ] **14.3.5** Add streak tracking
- [ ] **14.3.6** Create leaderboard by time period

---

## Phase 15: DevOps & Automation (25 Tasks)

### 15.1 CI/CD
- [ ] **15.1.1** Set up GitHub Actions workflow
- [ ] **15.1.2** Add go fmt check to CI
- [ ] **15.1.3** Add go vet check to CI
- [ ] **15.1.4** Add go test to CI
- [ ] **15.1.5** Add go build to CI
- [ ] **15.1.6** Configure auto-deploy on merge to main

### 15.2 Infrastructure as Code
- [ ] **15.2.1** Update render.yaml with all services
- [ ] **15.2.2** Add environment groups
- [ ] **15.2.3** Version control render.yaml

### 15.3 Backup & Recovery
- [ ] **15.3.1** Set up automated database backups
- [ ] **15.3.2** Test backup restoration
- [ ] **15.3.3** Document disaster recovery plan

---

## Phase 16: Maintenance (20 Tasks)

### 16.1 Performance Optimization
- [ ] **16.1.1** Profile API endpoints
- [ ] **16.1.2** Add database indexes for slow queries
- [ ] **16.1.3** Optimize AI response caching
- [ ] **16.1.4** Add CDN for static assets
- [ ] **16.1.5** Optimize frontend bundle size

### 16.2 Tech Debt
- [ ] **16.2.1** Update Go dependencies
- [ ] **16.2.2** Fix all lint warnings
- [ ] **16.2.3** Increase test coverage to 80%+
- [ ] **16.2.4** Remove dead code
- [ ] **16.2.5** Refactor large functions
- [ ] **16.2.6** Add type annotations where missing

---

## Summary

| Phase | Name | Tasks |
|-------|------|-------|
| 1 | Diagnosis & Assessment | 15 |
| 2 | Render Reconfiguration | 45 |
| 3 | Environment Variables | 30 |
| 4 | Database Setup | 35 |
| 5 | Keepalive & Reliability | 25 |
| 6 | API Enhancements | 40 |
| 7 | Cloudflare Integration | 30 |
| 8 | Frontend Development | 45 |
| 9 | Testing | 35 |
| 10 | Monitoring | 30 |
| 11 | Security | 25 |
| 12 | Documentation | 20 |
| 13 | Content Pipeline | 30 |
| 14 | Advanced Features | 35 |
| 15 | DevOps & Automation | 25 |
| 16 | Maintenance | 20 |
| **TOTAL** | | **455 TASKS** |
