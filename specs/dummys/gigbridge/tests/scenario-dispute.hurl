# === Dispute Flow ===
# gig: under_review → disputed

# --- Client: Register + Login ---
POST {{host}}/auth/register
Content-Type: application/json
{
  "email": "client3@gigbridge-test.com",
  "password": "Pass1234!",
  "role": "client",
  "name": "Client Three"
}

HTTP 200

POST {{host}}/auth/login
Content-Type: application/json
{
  "email": "client3@gigbridge-test.com",
  "password": "Pass1234!"
}

HTTP 200
[Captures]
client_token: jsonpath "$.AccessToken"

# --- Freelancer: Register + Login ---
POST {{host}}/auth/register
Content-Type: application/json
{
  "email": "freelancer3@gigbridge-test.com",
  "password": "Pass1234!",
  "role": "freelancer",
  "name": "Freelancer Three"
}

HTTP 200

POST {{host}}/auth/login
Content-Type: application/json
{
  "email": "freelancer3@gigbridge-test.com",
  "password": "Pass1234!"
}

HTTP 200
[Captures]
freelancer_token: jsonpath "$.AccessToken"

# --- CreateGig + Publish ---
POST {{host}}/gigs
Authorization: Bearer {{client_token}}
Content-Type: application/json
{
  "title": "Build API",
  "description": "REST API development",
  "budget": 8000
}

HTTP 200
[Captures]
gig_id: jsonpath "$.gig.id"

PUT {{host}}/gigs/{{gig_id}}/publish
Authorization: Bearer {{client_token}}

HTTP 200

# --- SubmitProposal + AcceptProposal ---
POST {{host}}/gigs/{{gig_id}}/proposals
Authorization: Bearer {{freelancer_token}}
Content-Type: application/json
{
  "bid_amount": 7000
}

HTTP 200
[Captures]
proposal_id: jsonpath "$.proposal.id"

POST {{host}}/proposals/{{proposal_id}}/accept
Authorization: Bearer {{client_token}}

HTTP 200

# --- SubmitWork → under_review ---
POST {{host}}/gigs/{{gig_id}}/submit-work
Authorization: Bearer {{freelancer_token}}

HTTP 200
[Asserts]
jsonpath "$.gig.status" == "under_review"

# --- RaiseDispute (client) → gig: under_review → disputed ---
POST {{host}}/gigs/{{gig_id}}/dispute
Authorization: Bearer {{client_token}}

HTTP 200
[Asserts]
jsonpath "$.gig.status" == "disputed"
