# Happy Path: Client creates gig → Freelancer submits proposal → Client accepts → Escrow → Freelancer submits work → Client approves → Funds released

# 1. Register client
POST {{host}}/auth/register
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

# 2. Login client
POST {{host}}/auth/login
Content-Type: application/json
{
    "email": "client@test.com",
    "password": "pass1234"
}
HTTP 200
[Captures]
client_token: jsonpath "$.access_token"

# 3. Register freelancer
POST {{host}}/auth/register
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

# 4. Login freelancer
POST {{host}}/auth/login
Content-Type: application/json
{
    "email": "freelancer@test.com",
    "password": "pass1234"
}
HTTP 200
[Captures]
freelancer_token: jsonpath "$.access_token"

# 5. Client creates gig
POST {{host}}/gigs
Content-Type: application/json
Authorization: Bearer {{client_token}}
{
    "title": "Build a website",
    "description": "Need a landing page",
    "budget": 1000
}
HTTP 201
[Captures]
gig_id: jsonpath "$.gig.id"

# 6. Client publishes gig
PUT {{host}}/gigs/{{gig_id}}/publish
Content-Type: application/json
Authorization: Bearer {{client_token}}
HTTP 200

# 7. Freelancer submits proposal
POST {{host}}/gigs/{{gig_id}}/proposals
Content-Type: application/json
Authorization: Bearer {{freelancer_token}}
{
    "bid_amount": 900
}
HTTP 201
[Captures]
proposal_id: jsonpath "$.proposal.id"

# 8. Client accepts proposal (escrow hold)
POST {{host}}/proposals/{{proposal_id}}/accept
Content-Type: application/json
Authorization: Bearer {{client_token}}
HTTP 200

# 9. Freelancer submits work
POST {{host}}/gigs/{{gig_id}}/submit-work
Content-Type: application/json
Authorization: Bearer {{freelancer_token}}
HTTP 200

# 10. Client approves work (funds released)
POST {{host}}/gigs/{{gig_id}}/approve
Content-Type: application/json
Authorization: Bearer {{client_token}}
HTTP 200
