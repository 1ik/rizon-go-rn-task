# Execution Plan: Email Authentication & Onboarding Flow

## Overview
This document breaks down the development into two phases with granular subtasks. Each subtask should be completed and verified before moving to the next one.

---

## Phase 1: User Authentication (Up to User Logged In)

### Authentication Flow Diagram

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                        STEP 1: REQUEST EMAIL AUTH                           │
└─────────────────────────────────────────────────────────────────────────────┘

    User
     |
     | Enters email address
     v
┌─────────────────┐
│   Frontend App  │
└─────────────────┘
     |
     | Validate email format
     |
     | Request email auth
     v
┌─────────────────┐
│   Backend API   │
└─────────────────┘
     |
     | Check if email exists in in-memory map
     v
┌─────────────────┐
│  In-Memory Map  │ <-- Check: email exists?
│   (email=>hash) │
└─────────────────┘
     |
     | (Rate limiting check)
     |
     +---> Email exists in map -----> Rate limit: "Try after X time"
     |                                |
     |                                v
     |                           ┌─────────────────┐
     |                           │  Frontend App   │ <-- Error response
     |                           └─────────────────┘
     |
     +---> Email NOT in map -----> Continue
                              |
                              | Generate hash = SHA-256(email, salt)
                              |
                              | Store in in-memory map:
                              |   Key: email
                              |   Value: hash (with expiration)
                              v
                         ┌─────────────────┐
                         │  In-Memory Map  │ <-- Store: email => hash
                         │   (email=>hash) │     (Expires after given time)
                         └─────────────────┘
                              |
                              | Schedule cleanup task to delete after expiration
                              v
                         ┌─────────────────┐
                         │  Task Scheduler │ <-- Schedule: Delete(email) after expires_at
                         │   (Background)  │
                         └─────────────────┘
                              |
                              | Generate email auth: /email-auth?email={email}&secret={hash}
                              |
                              | Send email with email auth
                              v
                         ┌─────────────────┐
                         │  Email Service  │ <-- app://email-auth?email=user@example.com&secret=abc123
                         │     (Mock)      │     (Logs to console in mock)
                         └─────────────────┘
     |
     | Success response
     v
┌─────────────────┐
│   Frontend App  │ <-- Show "Check your email" message
└─────────────────┘
     |
     v
    User


┌─────────────────────────────────────────────────────────────────────────────┐
│                    STEP 2: USER CLICKS EMAIL AUTH LINK                           │
└─────────────────────────────────────────────────────────────────────────────┘

    User
     |
     | Clicks email auth in email
     v
┌─────────────────┐
│  Email Service  │
└─────────────────┘
     |
     | Opens deep link: app://email-auth?email=user@example.com&secret={hash}
     v
┌─────────────────┐
│  Deep Link      │
│    Handler      │
└─────────────────┘


┌─────────────────────────────────────────────────────────────────────────────┐
│                    STEP 3: DEEP LINK HANDLING                               │
└─────────────────────────────────────────────────────────────────────────────┘

    Deep Link Handler
     |
     | (Two scenarios)
     |
     +---> Cold Start (App Closed) -----> App launches from deep link
     |                                    |
     |                                    v
     |                              ┌─────────────────┐
     |                              │  Frontend App   │
     |                              └─────────────────┘
     |                                    |
     |                                    | Extract email and secret from URL
     |                                    |
     +---> Warm Start (App Running) -----> Deep link event received
                                          |
                                          v
                                    ┌─────────────────┐
                                    │  Frontend App   │
                                    └─────────────────┘
                                          |
                                          | Extract email and secret from URL


┌─────────────────────────────────────────────────────────────────────────────┐
│                  STEP 4: VERIFY EMAIL AUTH                            │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────┐
│  Frontend App   │
└─────────────────┘
     |
     | Verify email auth (send email and secret in request body)
     v
┌─────────────────┐
│   Backend API   │
└─────────────────┘
     |
     | Extract email from request body
     |
     | Extract secret from request body
     |
     | Check if email exists in in-memory map
     v
┌─────────────────┐
│  In-Memory Map  │ <-- Check: email exists?
│   (email=>hash) │
└─────────────────┘
     |
     | (Validation checks)
     |
     +---> Email NOT in map -----> Error: "Timed out" or "Does not exist"
     |     (Email either never existed or was auto-deleted after expiration)
     |
     +---> Email exists in map -----> Continue
                              |
                              | Get hash value for email from map
                              | Get(email) -> hash
                              |
                              | Compare secret from request body with hash from map
                              |
                              +---> Secret matches hash -----> Continue
                              |
                              +---> Secret does NOT match hash -----> Error: "Invalid secret"
                              |
                              | (Email remains in map - deletion handled by scheduler)
                              |
                              | Create session token
                              |
                              | Check if user exists by email
                              v
                         ┌─────────────────┐
                         │   Database      │ <-- Query user by email
                         └─────────────────┘
                              |
                              | (if user doesn't exist)
                              |
                              | Create user record
                              v
                         ┌─────────────────┐
                         │   Database      │ <-- Create user with email
                         └─────────────────┘
                              |
                              | Generate session token
                              | (payload: user_id, email, expiration)
                              |
                              | Return session token + user info
                              v
                         ┌─────────────────┐
                         │  Frontend App   │ <-- Receive token and user
                         └─────────────────┘


┌─────────────────────────────────────────────────────────────────────────────┐
│                  STEP 5: SAVE SESSION & AUTHENTICATE                        │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────┐
│  Frontend App   │
└─────────────────┘
     |
     | Save session token to local storage
     |
     | Update auth state: isAuthenticated = true, user = {id, email}
     |
     | Navigate to protected route (HomeScreen)
     |
     v
    User (Authenticated)


┌─────────────────────────────────────────────────────────────────────────────┐
│              STEP 6: SESSION PERSISTENCE (APP RESTART)                     │
└─────────────────────────────────────────────────────────────────────────────┘

    User
     |
     | Restarts app
     v
┌─────────────────┐
│  Frontend App   │
└─────────────────┘
     |
     | Load token from local storage
     |
     | Verify token with backend
     v
┌─────────────────┐
│   Backend API   │
└─────────────────┘
     |
     | Validate token signature & expiration
     |
     +---> Token Valid -----> Get user by ID
     |                        |
     |                        v
     |                   ┌─────────────────┐
     |                   │   Database      │ <-- Query user
     |                   └─────────────────┘
     |                        |
     |                        | Return user info
     |                        v
     |                   ┌─────────────────┐
     |                   │  Frontend App   │ <-- Restore auth state
     |                   └─────────────────┘
     |                        |
     |                        | Navigate to HomeScreen
     |                        |
     +---> Token Invalid/Expired -----> Clear local storage
                                        |
                                        | Navigate to LoginScreen
                                        v
                                   ┌─────────────────┐
                                   │  Frontend App   │
                                   └─────────────────┘
```

### Technical Implementation Flow

**Key Components:**

1. **Email Auth Hash Generation**
   - Generate hash using SHA-256: hash = SHA-256(email, salt)
   - Same email + same salt = same hash (deterministic)
   - Storage structure: In-memory map with email as key, hash as value, with expiration. Redis will automatically expire the key. e.g. after expiration the key will not be found.
   - Format: email => hash
   - Rate limiting: If email exists in map, don't send link again (rate limit)
   - Email auth format: `/email-auth?email={email}&secret={hash}`
   - Verification: Extract email and secret from request body, match secret with stored hash
   - Deletion: Email is only deleted by scheduled cleanup task (not manually in verification code)
   - User is NOT created until hash is verified

2. **Deep Link Format**
   - URL Scheme: `app://email-auth?email={email}&secret={hash}`
   - Purpose: Pass email and secret from email to the app (transport mechanism)
   - App extracts email and secret from URL query parameters
   - App then sends email and secret to backend in request body for verification
   - Handles both cold start (app closed) and warm start (app running)
   - Note: Deep link does NOT verify - it just passes data to app, which sends to backend

3. **Session Management**
   - Session token issued after successful email auth verification
   - Token stored in local storage for persistence
   - Token validated on app start via user info endpoint

4. **Security Measures**
   - Rate limiting (email-based): Built-in via in-memory map check
     - If email exists in map, don't send link again (rate limit)
     - One request per email until expiration (30 minutes)
     - Cleanup task removes email after expiration, allowing new requests
   - Hash matching prevents unauthorized access
   - Session token signature validation ensures token integrity


### Backend Subtasks

#### Task 1.1: Backend Project Initialization & Hello World
**Goal**: Set up basic Go project structure and simple hello world endpoint
- [ ] Create `backend/` directory
- [ ] Initialize Go module (`go mod init`)
- [ ] Create basic directory structure (`cmd/`, `internal/`, `pkg/`)
- [ ] Add `.gitignore` for Go
- [ ] Set up basic HTTP server
- [ ] Create `/hello` endpoint that returns "Hello World"
- [ ] Test server runs and `/hello` endpoint responds correctly
- [ ] Verify project structure is correct

**Acceptance**: `go mod init` runs successfully, basic directories exist, `/hello` endpoint returns "Hello World"

---

#### Task 1.2: Backend GraphQL Setup
**Goal**: Set up GraphQL server and schema
- [ ] Install GraphQL library for Go (gqlgen or graphql-go)
- [ ] Set up GraphQL server
- [ ] Create GraphQL schema file
- [ ] Configure GraphQL endpoint (e.g., `/graphql`)
- [ ] Set up GraphQL playground/explorer (optional, for development)
- [ ] Test GraphQL endpoint is accessible
- [ ] Verify GraphQL server runs correctly

**Acceptance**: GraphQL server runs, schema file created, GraphQL endpoint accessible

---

#### Task 1.3: GraphQL Type Generation
**Goal**: Generate types from GraphQL schema for both Go and TypeScript
- [ ] Set up GraphQL code generation tool for Go (gqlgen or similar)
- [ ] Set up GraphQL code generation for TypeScript (GraphQL Code Generator)
- [ ] Configure code generation config files
- [ ] Generate Go types from GraphQL schema
- [ ] Generate TypeScript types from GraphQL schema
- [ ] Verify generated types are correct
- [ ] Set up automated type generation workflow (optional)

**Acceptance**: Types generated for both Go and TypeScript, can be used in code

---

#### Task 1.4: Backend Dependencies Installation
**Goal**: Install required Go packages
- [ ] Install GraphQL library (if not already installed in Task 1.2)
- [ ] Install GORM and PostgreSQL driver
- [ ] Install JWT library
- [ ] Install environment variable loader
- [ ] Install other required dependencies
- [ ] Run `go mod tidy` to verify dependencies

**Acceptance**: All dependencies install without errors, `go.mod` file updated

---

#### Task 1.3: Backend Database Setup
**Goal**: Configure SQLite database connection
- [ ] Create database configuration file
- [ ] Set up GORM connection to SQLite
- [ ] Create database initialization function
- [ ] Test database connection
- [ ] Add database close/cleanup logic

**Acceptance**: Database connects successfully, can create tables

---

#### Task 1.9: Backend Database Schema - Users Table
**Goal**: Create users table for authentication
- [ ] Define User model struct
- [ ] Create users table migration
- [ ] Add fields: id, email, created_at, updated_at
- [ ] Test table creation
- [ ] Verify table structure

**Acceptance**: Users table exists with correct schema

---


**Goal**: Set up in-memory map for email auth hashes
- [ ] Create in-memory map structure (thread-safe)
- [ ] Define storage structure:
  - Key: email (string)
  - Value: hash (string) with expiration timestamp
- [ ] Implement store function: Store(email, hash, expires_at)
- [ ] Implement retrieve function: Get(email) -> hash
- [ ] Implement delete function: Delete(email)
- [ ] Implement exists check: Exists(email) -> bool
- [ ] Set up task scheduler/background job system for cleanup
- [ ] Test map operations (store, retrieve, delete, exists)

**Acceptance**: In-memory map configured, can store and retrieve email=>hash pairs with expiration

---

#### Task 1.7: Backend GraphQL Server Integration
**Goal**: Integrate GraphQL server with main application
- [ ] Integrate GraphQL endpoint with main server
- [ ] Set up basic middleware (logging, recovery)
- [ ] Configure CORS for React Native
- [ ] Add health check endpoint (`GET /health`)
- [ ] Test server starts and GraphQL endpoint responds
- [ ] Verify GraphQL queries work correctly

**Acceptance**: Server runs on port, GraphQL endpoint accessible, health endpoint returns 200

---

#### Task 1.8: Backend Environment Configuration
**Goal**: Set up environment variables
- [ ] Create `.env.example` file
- [ ] Add environment variables: PORT, DATABASE_URL (PostgreSQL), JWT_SECRET, GRAPHQL_ENDPOINT
- [ ] Load environment variables in server
- [ ] Add default values for development
- [ ] Test configuration loading

**Acceptance**: Environment variables load correctly, defaults work

---

#### Task 1.10: Backend Email Auth Hash Generation
**Goal**: Generate SHA-256 hash for email auths
- [ ] Create hash generation function using SHA-256
- [ ] Implement: hash = SHA-256(email, salt)
- [ ] Ensure same email + salt produces same hash (deterministic)
- [ ] Set expiration time (30 minutes)
- [ ] Store hash in in-memory map: email => hash
- [ ] Schedule background task to delete email from map after expires_at
- [ ] Test hash generation (same input = same output)
- [ ] Test hash storage and scheduled cleanup

**Acceptance**: Hash generation is deterministic, stored in map correctly, cleanup task scheduled automatically

---

#### Task 1.11: Backend Request Email Auth GraphQL Mutation
**Goal**: Create GraphQL mutation to request email auth
- [ ] Create GraphQL mutation `requestEmailAuth(email: String!): Boolean!`
- [ ] Add mutation to GraphQL schema
- [ ] Validate email format in request
- [ ] Check if email exists in in-memory map (rate limiting)
  - [ ] If exists: return error "Try after X time" (rate limit)
  - [ ] If not exists: continue
- [ ] Generate hash = SHA-256(email, salt)
- [ ] Store in in-memory map: email => hash (NO user creation yet)
- [ ] Schedule cleanup task to delete email from map after expires_at
- [ ] Generate email auth: `/email-auth?email={email}&secret={hash}`
- [ ] Return success response
- [ ] Test endpoint with valid email
- [ ] Test rate limiting (same email twice)
- [ ] Test cleanup task is scheduled correctly

**Acceptance**: Endpoint accepts email, checks rate limit, generates hash, stores in map, schedules cleanup, returns success

---

#### Task 1.12: Backend Email Service Interface
**Goal**: Create abstract email sending interface
- [ ] Define EmailSender interface
- [ ] Create mock email sender implementation
- [ ] Mock should log email content to console
- [ ] Include deep link URL in email content
- [ ] Test interface and mock implementation

**Acceptance**: Email interface exists, mock logs email with link

---

#### Task 1.13: Backend Integrate Email Service with Email Auth
**Goal**: Send email when email auth is requested
- [ ] Call email service in request-link endpoint
- [ ] Generate deep link URL: `app://email-auth?email={email}&secret={hash}`
- [ ] Format email message with link
- [ ] Send email via mock service
- [ ] Test email is sent/logged when requesting link

**Acceptance**: Requesting email auth triggers email with correct deep link format

---

#### Task 1.14: Backend Rate Limiting - Basic Setup (Optional)
**Goal**: Add basic rate limiting per IP address (optional, email-based rate limiting already implemented)
- [ ] Create rate limiter middleware
- [ ] Limit requests per IP address
- [ ] Set limit: 5 requests per 15 minutes per IP
- [ ] Return 429 status when limit exceeded
- [ ] Test rate limiting works

**Acceptance**: Rate limiting prevents excessive requests per IP, returns 429
**Note**: Email-based rate limiting is already handled by checking if email exists in map (Task 1.9)

---

#### Task 1.15: Backend Rate Limiting - Per Email (Already Implemented)
**Goal**: Rate limiting per email is already implemented via in-memory map check
- [ ] Rate limiting is handled in Task 1.9
- [ ] If email exists in map, don't send link again
- [ ] This automatically prevents multiple requests for same email
- [ ] Cleanup task removes email after expiration, allowing new requests

**Acceptance**: Email-based rate limiting works via map existence check
**Note**: This is already implemented in Task 1.9 - no separate task needed

---

#### Task 1.16: Backend Email Auth Verification GraphQL Mutation
**Goal**: Create GraphQL mutation to verify email auth
- [ ] Create GraphQL mutation `verifyEmailAuth(email: String!, secret: String!): AuthPayload!`
- [ ] Add mutation to GraphQL schema
- [ ] Define AuthPayload type (token, user)
- [ ] Extract email from request body
- [ ] Extract secret from request body
- [ ] Check if email exists in in-memory map
  - [ ] If email does NOT exist: return error "Timed out" or "Does not exist"
  - [ ] (Email either never existed or was auto-deleted by cleanup task after expiration)
- [ ] If email exists: Get hash value for email from map: Get(email) -> hash
- [ ] Compare secret from request body with hash from map
  - [ ] If secret matches hash: continue
  - [ ] If secret does NOT match hash: return "Invalid secret" error
- [ ] Do NOT delete email from map manually (deletion handled by scheduler only)
- [ ] Return error responses for invalid cases

**Acceptance**: Endpoint validates email exists in map, matches secret from request body with stored hash. Email remains in map after verification - deletion only happens via scheduled cleanup task.

---

#### Task 1.17: Backend Session Token Generation
**Goal**: Generate JWT session tokens after successful verification
- [ ] Create JWT token generation function
- [ ] Include user ID and email in token payload
- [ ] Set token expiration (7 days)
- [ ] Sign token with secret key
- [ ] Test token generation and structure

**Acceptance**: JWT tokens are generated with correct payload and signature

---

#### Task 1.18: Backend Complete Email Auth Verification Flow
**Goal**: Complete the verification endpoint with session creation
- [ ] After secret validation passes (secret matches hash)
- [ ] Email already extracted from request body (from Step 1.14)
- [ ] Email remains in in-memory map (deletion handled by scheduler)
- [ ] Check if user exists by email in database
- [ ] Create user if doesn't exist (first time login)
- [ ] Get user details (existing or newly created)
- [ ] Generate session token with user ID and email
- [ ] Return session token in response
- [ ] Return user information in response
- [ ] Test full verification flow

**Acceptance**: Verifying valid secret creates user if needed, returns session token and user info, email deleted from map

---

#### Task 1.19: Backend Auth Middleware
**Goal**: Create middleware to protect routes
- [ ] Create authentication middleware
- [ ] Extract JWT token from Authorization header
- [ ] Validate token signature and expiration
- [ ] Attach user info to request context
- [ ] Return 401 for invalid/missing tokens
- [ ] Test middleware with valid and invalid tokens

**Acceptance**: Middleware validates tokens, protects routes, adds user to context

---

#### Task 1.20: Backend Get Current User GraphQL Query
**Goal**: Create GraphQL query to get authenticated user info
- [ ] Create GraphQL query `me: User!`
- [ ] Add query to GraphQL schema
- [ ] Protect with auth middleware/resolver
- [ ] Return current user information
- [ ] Test endpoint with valid session token
- [ ] Test endpoint without token (should fail)

**Acceptance**: Endpoint returns user info when authenticated, fails when not

---

### Frontend Subtasks

#### Task 1.21: Frontend Project Initialization
**Goal**: Set up React Native project
- [ ] Create `frontend/` directory
- [ ] Initialize React Native project (Expo or bare RN)
- [ ] Verify project structure
- [ ] Test app runs on simulator/device
- [ ] Add `.gitignore` for React Native

**Acceptance**: React Native app runs successfully

---

#### Task 1.22: Frontend Dependencies Installation
**Goal**: Install required npm packages
- [ ] Install React Navigation packages
- [ ] Install AsyncStorage
- [ ] Install deep linking library
- [ ] Install Apollo Client and GraphQL dependencies
- [ ] Install GraphQL Code Generator (for type generation)
- [ ] Install TypeScript (if not already installed)
- [ ] Install state management library (if needed)
- [ ] Run `npm install` and verify

**Acceptance**: All dependencies install without errors

---

#### Task 1.23: Frontend Basic Navigation Setup
**Goal**: Set up React Navigation structure
- [ ] Create navigation container
- [ ] Set up stack navigator
- [ ] Create placeholder screens (LoginScreen, HomeScreen)
- [ ] Test navigation between screens
- [ ] Verify navigation works

**Acceptance**: Navigation works, can navigate between screens

---

#### Task 1.24: Frontend GraphQL Client Setup (Apollo Client)
**Goal**: Set up Apollo Client for GraphQL communication
- [ ] Install Apollo Client and GraphQL dependencies
- [ ] Install generated TypeScript types (from Task 1.3)
- [ ] Create Apollo Client instance
- [ ] Configure GraphQL endpoint URL from environment
- [ ] Set up Apollo Client cache
- [ ] Configure authentication headers (for future use)
- [ ] Test Apollo Client can make GraphQL queries
- [ ] Verify generated TypeScript types work correctly

**Acceptance**: Apollo Client configured, can make GraphQL queries, TypeScript types available

---

#### Task 1.25: Frontend Email Input Screen UI
**Goal**: Create email input screen
- [ ] Create EmailLoginScreen component
- [ ] Add email input field
- [ ] Add submit button
- [ ] Add basic styling
- [ ] Test UI renders correctly

**Acceptance**: Email input screen displays with input and button

---

#### Task 1.26: Frontend Email Validation
**Goal**: Add email format validation
- [ ] Add email validation function
- [ ] Validate on input change
- [ ] Show validation error message
- [ ] Disable submit button if email invalid
- [ ] Test validation works correctly

**Acceptance**: Invalid emails show error, submit disabled for invalid emails

---

#### Task 1.25: Frontend Request Email Auth API Integration
**Goal**: Connect email screen to backend API
- [ ] Create `requestMagicLink` API function
- [ ] Call backend `/api/auth/request-link` endpoint
- [ ] Handle loading state
- [ ] Handle success response
- [ ] Handle error response
- [ ] Show success message to user
- [ ] Test API call works

**Acceptance**: Submitting email calls backend, shows success message

---

#### Task 1.28: Frontend Auth Context Setup
**Goal**: Create authentication state management
- [ ] Create AuthContext
- [ ] Add auth state (user, token, isAuthenticated)
- [ ] Add auth actions (login, logout)
- [ ] Wrap app with AuthProvider
- [ ] Test context provides state

**Acceptance**: Auth context available throughout app

---

#### Task 1.29: Frontend Session Token Storage
**Goal**: Persist session token in AsyncStorage
- [ ] Create storage utility functions
- [ ] Save token to AsyncStorage on login
- [ ] Load token from AsyncStorage on app start
- [ ] Clear token on logout
- [ ] Test token persists across app restarts

**Acceptance**: Token saved and loaded correctly, persists after app restart

---

#### Task 1.30: Frontend Deep Link Configuration
**Goal**: Configure app to handle deep links
- [ ] Configure URL scheme in app config (`app://`)
- [ ] Set up deep link listener
- [ ] Parse deep link URL
- [ ] Extract email and secret from URL parameters
- [ ] URL format: `app://email-auth?email={email}&secret={hash}`
- [ ] Test deep link opens app

**Acceptance**: App handles deep links, extracts email and secret from URL

---

#### Task 1.31: Frontend Deep Link Handler - Cold Start
**Goal**: Handle deep link when app is closed
- [ ] Detect app launch from deep link
- [ ] Extract email and secret from URL on app start
- [ ] Call verify endpoint with email and secret
- [ ] Handle verification response
- [ ] Test cold start deep link flow

**Acceptance**: App opens from link, extracts email and secret, verifies, logs in user

---

#### Task 1.32: Frontend Deep Link Handler - Warm Start
**Goal**: Handle deep link when app is already running
- [ ] Listen for deep link events while app running
- [ ] Extract email and secret from deep link
- [ ] Call verify endpoint with email and secret
- [ ] Handle verification response
- [ ] Test warm start deep link flow

**Acceptance**: Clicking link while app open extracts email and secret, verifies, logs in

---

#### Task 1.31: Frontend Verify Email Auth API Integration
**Goal**: Connect deep link handler to backend
- [ ] Create `verifyMagicLink` API function
- [ ] Call backend `/api/auth/verify-link` endpoint with email and secret
- [ ] Send email and secret in request body (POST body)
- [ ] Handle success: save session token, update auth state
- [ ] Handle errors: show error message ("Timed out", "Does not exist", "Invalid secret")
- [ ] Navigate to appropriate screen after verification
- [ ] Test verification flow

**Acceptance**: Verifying email and secret saves session, updates auth state, navigates correctly

---

#### Task 1.34: Frontend Auth State Restoration
**Goal**: Restore auth state on app start
- [ ] Load token from AsyncStorage on app start
- [ ] Validate token exists
- [ ] Call GraphQL query (e.g., `me` query) to verify token still valid
- [ ] Use generated TypeScript types
- [ ] Update auth state if token valid
- [ ] Navigate to home if authenticated, login if not
- [ ] Test auth restoration

**Acceptance**: App restores session on start via GraphQL, navigates correctly based on auth state

---

#### Task 1.35: Frontend Protected Routes
**Goal**: Protect routes that require authentication
- [ ] Create protected route wrapper component
- [ ] Check auth state before rendering
- [ ] Redirect to login if not authenticated
- [ ] Show protected content if authenticated
- [ ] Test protected routes work

**Acceptance**: Protected routes redirect to login when not authenticated

---

#### Task 1.36: Phase 1 Integration Testing
**Goal**: Test complete authentication flow
- [ ] Test: Request email auth with email
- [ ] Test: Receive email with deep link
- [ ] Test: Click link (cold start)
- [ ] Test: Click link (warm start)
- [ ] Test: App restores session on restart
- [ ] Test: Protected routes work
- [ ] Test: Rate limiting works
- [ ] Test: Expired tokens fail
- [ ] Test: Used tokens fail

**Acceptance**: Complete authentication flow works end-to-end

---

## Phase 2: Onboarding & Feedback (After Authentication)

### Frontend Subtasks

#### Task 2.1: Onboarding State Management Setup
**Goal**: Create system to track onboarding completion
- [ ] Create onboarding context/state
- [ ] Add onboarding completion status
- [ ] Create storage key format: `onboarding_completed_{userId}`
- [ ] Add functions to check/set completion status
- [ ] Test onboarding state management

**Acceptance**: Onboarding state can be checked and set per user

---

#### Task 2.2: Onboarding Bottom Sheet Component - Base
**Goal**: Create reusable bottom sheet component
- [ ] Install bottom sheet library (if needed) or create mock component
- [ ] Create BaseBottomSheet component
- [ ] Add open/close functionality
- [ ] Add backdrop/overlay
- [ ] Test bottom sheet opens and closes

**Acceptance**: Bottom sheet component renders and can be opened/closed

---

#### Task 2.3: Onboarding Bottom Sheet 1 - Welcome Screen
**Goal**: Create first onboarding sheet
- [ ] Create WelcomeSheet component
- [ ] Add welcome message/content
- [ ] Add "Next" or "Continue" button
- [ ] Style according to design (or mock design)
- [ ] Test sheet displays correctly

**Acceptance**: Welcome sheet displays with content and navigation button

---

#### Task 2.4: Onboarding Bottom Sheet 2 - Feedback Form
**Goal**: Create feedback collection sheet
- [ ] Create FeedbackSheet component
- [ ] Add text input field (mock existing custom component)
- [ ] Add "Send Feedback" button
- [ ] Add placeholder text
- [ ] Style form appropriately
- [ ] Test form displays correctly

**Acceptance**: Feedback sheet displays with input field and submit button

---

#### Task 2.5: Onboarding Bottom Sheet 3 - Store Redirect
**Goal**: Create store redirect sheet
- [ ] Create StoreSheet component
- [ ] Add image/illustration placeholder
- [ ] Add "Open Store" button
- [ ] Detect platform (iOS/Android)
- [ ] Style sheet appropriately
- [ ] Test sheet displays correctly

**Acceptance**: Store sheet displays with image and platform-specific button

---

#### Task 2.6: Onboarding Sequential Flow Logic
**Goal**: Implement sequential navigation between sheets
- [ ] Create onboarding flow manager
- [ ] Track current sheet index
- [ ] Navigate from Sheet 1 to Sheet 2
- [ ] Navigate from Sheet 2 to Sheet 3
- [ ] Close flow after Sheet 3
- [ ] Test sequential navigation works

**Acceptance**: Sheets appear in sequence, navigation works correctly

---

#### Task 2.7: Onboarding Trigger After Login
**Goal**: Show onboarding immediately after successful login
- [ ] Detect successful login event
- [ ] Check if onboarding already completed for user
- [ ] Show onboarding if not completed
- [ ] Skip onboarding if already completed
- [ ] Test onboarding appears after login

**Acceptance**: Onboarding appears after login if not completed

---

#### Task 2.8: Onboarding Persistence Check
**Goal**: Ensure onboarding shows only once per user
- [ ] Check onboarding status on app start
- [ ] Check onboarding status after login
- [ ] Mark onboarding as completed after Sheet 3
- [ ] Store completion status with user ID
- [ ] Test onboarding doesn't reappear after completion

**Acceptance**: Onboarding shows once per user, persists across app restarts

---

### Backend Subtasks

#### Task 2.9: Backend Feedback Database Schema
**Goal**: Create table to store feedback submissions
- [ ] Define Feedback model struct
- [ ] Create feedback table migration
- [ ] Add fields: id, user_id, content, created_at
- [ ] Add foreign key to users table
- [ ] Test table creation

**Acceptance**: Feedback table exists with correct schema and relationships

---

#### Task 2.10: Backend Feedback Repository
**Goal**: Create data access layer for feedback
- [ ] Create feedback repository interface
- [ ] Implement Create function
- [ ] Implement GetByUserID function
- [ ] Add error handling
- [ ] Test repository functions

**Acceptance**: Feedback can be created and retrieved from database

---

#### Task 2.11: Backend Feedback Submission Endpoint
**Goal**: Create API endpoint to submit feedback
- [ ] Create `POST /api/feedback` endpoint
- [ ] Protect with auth middleware
- [ ] Extract user from context
- [ ] Validate feedback content (not empty, max length)
- [ ] Save feedback to database
- [ ] Return success response
- [ ] Test endpoint with authenticated request

**Acceptance**: Endpoint accepts feedback, saves to database, returns success

---

#### Task 2.12: Backend Slack Client Interface
**Goal**: Create abstract Slack client interface
- [ ] Define SlackClient interface
- [ ] Add PublishFeedback method signature
- [ ] Document interface contract
- [ ] Test interface definition

**Acceptance**: Slack interface defined with clear contract

---

#### Task 2.13: Backend Mock Slack Implementation
**Goal**: Create mock Slack client
- [ ] Implement SlackClient interface
- [ ] Log feedback message to console
- [ ] Format message with user info and feedback
- [ ] Simulate API call delay (optional)
- [ ] Return success/error for testing
- [ ] Test mock implementation

**Acceptance**: Mock Slack logs feedback messages correctly

---

#### Task 2.14: Backend Integrate Slack with Feedback
**Goal**: Trigger Slack publish when feedback submitted
- [ ] Inject Slack client into feedback handler
- [ ] Call Slack.PublishFeedback after saving feedback
- [ ] Handle Slack errors gracefully (don't fail feedback save)
- [ ] Log Slack publish attempts
- [ ] Test Slack is called on feedback submission

**Acceptance**: Submitting feedback triggers Slack publish, errors handled gracefully

---

### Frontend Subtasks (Continued)

#### Task 2.15: Frontend Feedback Form State Management
**Goal**: Manage feedback input state
- [ ] Add state for feedback text
- [ ] Handle input changes
- [ ] Preserve input during navigation
- [ ] Clear input after successful submission
- [ ] Test state management works

**Acceptance**: Feedback input state managed correctly, preserved during errors

---

#### Task 2.16: Frontend Feedback Submission API Integration
**Goal**: Connect feedback form to backend
- [ ] Create `submitFeedback` API function
- [ ] Call backend `/api/feedback` endpoint
- [ ] Include session token in request
- [ ] Send feedback content in request body
- [ ] Test API call works

**Acceptance**: Feedback form calls backend API correctly

---

#### Task 2.17: Frontend Feedback Submission Loading State
**Goal**: Show loading state during submission
- [ ] Add loading state to feedback form
- [ ] Disable submit button while loading
- [ ] Show loading indicator
- [ ] Prevent form interaction during submission
- [ ] Test loading state displays correctly

**Acceptance**: Submit button disabled and loading shown during submission

---

#### Task 2.18: Frontend Feedback Submission Error Handling
**Goal**: Handle submission errors gracefully
- [ ] Catch API errors
- [ ] Display error message to user
- [ ] Keep feedback text in input (don't clear)
- [ ] Re-enable submit button on error
- [ ] Allow retry without losing input
- [ ] Test error handling works

**Acceptance**: Errors shown to user, input preserved, retry works

---

#### Task 2.19: Frontend Feedback Submission Success Handling
**Goal**: Handle successful submission
- [ ] Detect successful API response
- [ ] Show success message (optional)
- [ ] Close feedback sheet only after success
- [ ] Navigate to next sheet (Store sheet)
- [ ] Test success flow works

**Acceptance**: Successful submission closes sheet and navigates correctly

---

#### Task 2.20: Frontend Prevent Double Submission
**Goal**: Prevent duplicate feedback submissions
- [ ] Add submission flag to prevent double clicks
- [ ] Disable button immediately on click
- [ ] Check flag before making API call
- [ ] Reset flag after completion/error
- [ ] Test double submission prevention

**Acceptance**: Multiple rapid clicks don't cause duplicate submissions

---

#### Task 2.21: Frontend Store Redirect Implementation
**Goal**: Open App Store or Play Store
- [ ] Detect platform (iOS/Android)
- [ ] Create App Store URL for iOS
- [ ] Create Play Store URL for Android
- [ ] Use Linking.openURL() to open store
- [ ] Handle errors if store can't open
- [ ] Test store redirect works on both platforms

**Acceptance**: Store button opens correct store for platform

---

#### Task 2.22: Frontend Complete Onboarding Flow
**Goal**: Complete onboarding marks user as completed
- [ ] Mark onboarding complete after Store sheet
- [ ] Save completion status with user ID
- [ ] Close onboarding flow
- [ ] Navigate to main app screen
- [ ] Test onboarding completion persists

**Acceptance**: Completing onboarding marks user as completed, doesn't show again

---

#### Task 2.23: Frontend Onboarding Edge Cases
**Goal**: Handle edge cases in onboarding flow
- [ ] Test: App restart during onboarding
- [ ] Test: Login with different user after completing onboarding
- [ ] Test: Onboarding doesn't show for new user if already completed for another user
- [ ] Test: Onboarding shows for new user even if completed for previous user
- [ ] Test: Navigation away from onboarding and back

**Acceptance**: All edge cases handled correctly

---

#### Task 2.24: Phase 2 Integration Testing
**Goal**: Test complete onboarding and feedback flow
- [ ] Test: Login triggers onboarding
- [ ] Test: Navigate through all 3 sheets sequentially
- [ ] Test: Submit feedback successfully
- [ ] Test: Feedback saved to backend
- [ ] Test: Slack mock receives feedback
- [ ] Test: Store redirect works
- [ ] Test: Onboarding doesn't reappear after completion
- [ ] Test: Error handling works (network errors, etc.)
- [ ] Test: Double submission prevention works

**Acceptance**: Complete Phase 2 flow works end-to-end

---

## Final Integration & Testing

#### Task 3.1: End-to-End Flow Testing
**Goal**: Test complete application flow
- [ ] Test: Full flow from email input to feedback submission
- [ ] Test: App restart at various points
- [ ] Test: Deep link handling in all scenarios
- [ ] Test: Multiple users on same device
- [ ] Test: Rate limiting prevents abuse
- [ ] Test: Error scenarios handled gracefully

**Acceptance**: Complete application works correctly in all scenarios

---

#### Task 3.2: Code Review & Cleanup
**Goal**: Clean up code and ensure quality
- [ ] Review all code for consistency
- [ ] Remove debug logs and comments
- [ ] Ensure error messages are user-friendly
- [ ] Verify no hardcoded values
- [ ] Check code follows best practices

**Acceptance**: Code is clean, consistent, and production-ready

---

## Notes

- Each task should be completed and verified before moving to the next
- Test each task independently before integration
- Mark tasks as complete using checkboxes
- If a task fails, fix it before proceeding
- Some tasks can be done in parallel (e.g., backend and frontend tasks)
- Always test on both iOS and Android for frontend tasks
