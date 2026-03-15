# === Happy Path: Full Gig Lifecycle ===
# draft → open → in_progress → under_review → completed

# --- Client: Register ---
POST {{host}}/auth/register
Content-Type: application/json
{
  "email": "client@gigbridge-test.com",
  "password": "Pass1234!",
  "role": "client",
  "name": "Test Client"
}

HTTP 200
[Asserts]
jsonpath "$.user.id" exists
jsonpath "$.user.role" == "client"

# --- Client: Login ---
POST {{host}}/auth/login
Content-Type: application/json
{
  "email": "client@gigbridge-test.com",
  "password": "Pass1234!"
}

HTTP 200
[Captures]
client_token: jsonpath "$.AccessToken"

# --- Freelancer: Register ---
POST {{host}}/auth/register
Content-Type: application/json
{
  "email": "freelancer@gigbridge-test.com",
  "password": "Pass1234!",
  "role": "freelancer",
  "name": "Test Freelancer"
}

HTTP 200
[Asserts]
jsonpath "$.user.role" == "freelancer"

# --- Freelancer: Login ---
POST {{host}}/auth/login
Content-Type: application/json
{
  "email": "freelancer@gigbridge-test.com",
  "password": "Pass1234!"
}

HTTP 200
[Captures]
freelancer_token: jsonpath "$.AccessToken"

# --- CreateGig (client) ---
POST {{host}}/gigs
Authorization: Bearer {{client_token}}
Content-Type: application/json
{
  "title": "Build Landing Page",
  "description": "Need a responsive landing page",
  "budget": 5000
}

HTTP 200
[Captures]
gig_id: jsonpath "$.gig.id"
[Asserts]
jsonpath "$.gig.status" == "draft"
jsonpath "$.gig.title" == "Build Landing Page"
jsonpath "$.gig.budget" == 5000

# --- PublishGig (client) ---
PUT {{host}}/gigs/{{gig_id}}/publish
Authorization: Bearer {{client_token}}

HTTP 200
[Asserts]
jsonpath "$.gig.status" == "open"

# --- GetGig (public) ---
GET {{host}}/gigs/{{gig_id}}

HTTP 200
[Asserts]
jsonpath "$.gig.status" == "open"

# --- ListGigs (public) ---
GET {{host}}/gigs

HTTP 200
[Asserts]
jsonpath "$.items" isCollection
jsonpath "$.total" >= 1

# --- SubmitProposal (freelancer) ---
POST {{host}}/gigs/{{gig_id}}/proposals
Authorization: Bearer {{freelancer_token}}
Content-Type: application/json
{
  "bid_amount": 4500
}

HTTP 200
[Captures]
proposal_id: jsonpath "$.proposal.id"
[Asserts]
jsonpath "$.proposal.status" == "pending"
jsonpath "$.proposal.bid_amount" == 4500

# --- AcceptProposal (client) → gig: open → in_progress, proposal: pending → accepted ---
POST {{host}}/proposals/{{proposal_id}}/accept
Authorization: Bearer {{client_token}}

HTTP 200
[Asserts]
jsonpath "$.proposal.status" == "accepted"

# --- Verify gig is now in_progress ---
GET {{host}}/gigs/{{gig_id}}

HTTP 200
[Asserts]
jsonpath "$.gig.status" == "in_progress"

# --- SubmitWork (freelancer) → gig: in_progress → under_review ---
POST {{host}}/gigs/{{gig_id}}/submit-work
Authorization: Bearer {{freelancer_token}}

HTTP 200
[Asserts]
jsonpath "$.gig.status" == "under_review"

# --- ApproveWork (client) → gig: under_review → completed ---
POST {{host}}/gigs/{{gig_id}}/approve
Authorization: Bearer {{client_token}}

HTTP 200
[Asserts]
jsonpath "$.gig.status" == "completed"
