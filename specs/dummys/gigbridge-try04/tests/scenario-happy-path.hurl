# Happy Path: Client creates gig → Freelancer proposes → Client accepts → Escrow → Freelancer submits → Client approves → Funds released

# Step 1: Register client
POST {{host}}/users/register
Content-Type: application/json
{
    "email": "client@test.com",
    "password": "pass1234",
    "role": "client",
    "name": "Test Client"
}
HTTP 201
[Captures]
client_id: jsonpath "$.user.id"

# Step 2: Register freelancer
POST {{host}}/users/register
Content-Type: application/json
{
    "email": "freelancer@test.com",
    "password": "pass1234",
    "role": "freelancer",
    "name": "Test Freelancer"
}
HTTP 201
[Captures]
freelancer_id: jsonpath "$.user.id"

# Step 3: Login as client
POST {{host}}/users/login
Content-Type: application/json
{
    "email": "client@test.com",
    "password": "pass1234"
}
HTTP 200
[Captures]
client_token: jsonpath "$.access_token"

# Step 4: Login as freelancer
POST {{host}}/users/login
Content-Type: application/json
{
    "email": "freelancer@test.com",
    "password": "pass1234"
}
HTTP 200
[Captures]
freelancer_token: jsonpath "$.access_token"

# Step 5: Client creates gig
POST {{host}}/gigs
Authorization: Bearer {{client_token}}
Content-Type: application/json
{
    "title": "Build a website",
    "description": "Need a modern website",
    "budget": 100000
}
HTTP 201
[Captures]
gig_id: jsonpath "$.gig.id"
[Asserts]
jsonpath "$.gig.status" == "draft"

# Step 6: Client publishes gig
PUT {{host}}/gigs/{{gig_id}}/publish
Authorization: Bearer {{client_token}}
HTTP 200
[Asserts]
jsonpath "$.gig.status" == "open"

# Step 7: Freelancer submits proposal
POST {{host}}/gigs/{{gig_id}}/proposals
Authorization: Bearer {{freelancer_token}}
Content-Type: application/json
{
    "bid_amount": 80000
}
HTTP 201
[Captures]
proposal_id: jsonpath "$.proposal.id"
[Asserts]
jsonpath "$.proposal.status" == "pending"

# Step 8: Client accepts proposal (escrow hold)
POST {{host}}/proposals/{{proposal_id}}/accept
Authorization: Bearer {{client_token}}
HTTP 200
[Asserts]
jsonpath "$.gig.status" == "in_progress"

# Step 9: Freelancer submits work
POST {{host}}/gigs/{{gig_id}}/submit-work
Authorization: Bearer {{freelancer_token}}
HTTP 200
[Asserts]
jsonpath "$.gig.status" == "under_review"

# Step 10: Client approves work (funds released)
POST {{host}}/gigs/{{gig_id}}/approve
Authorization: Bearer {{client_token}}
HTTP 200
[Asserts]
jsonpath "$.gig.status" == "completed"
