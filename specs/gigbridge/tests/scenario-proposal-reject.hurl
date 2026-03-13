# === Proposal Rejection Flow ===
# proposal: pending → rejected

# --- Client: Register + Login ---
POST {{host}}/auth/register
Content-Type: application/json
{
  "email": "client2@gigbridge-test.com",
  "password": "Pass1234!",
  "role": "client",
  "name": "Client Two"
}

HTTP 200

POST {{host}}/auth/login
Content-Type: application/json
{
  "email": "client2@gigbridge-test.com",
  "password": "Pass1234!"
}

HTTP 200
[Captures]
client_token: jsonpath "$.AccessToken"

# --- Freelancer: Register + Login ---
POST {{host}}/auth/register
Content-Type: application/json
{
  "email": "freelancer2@gigbridge-test.com",
  "password": "Pass1234!",
  "role": "freelancer",
  "name": "Freelancer Two"
}

HTTP 200

POST {{host}}/auth/login
Content-Type: application/json
{
  "email": "freelancer2@gigbridge-test.com",
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
  "title": "Design Logo",
  "description": "Need a logo design",
  "budget": 3000
}

HTTP 200
[Captures]
gig_id: jsonpath "$.gig.id"

PUT {{host}}/gigs/{{gig_id}}/publish
Authorization: Bearer {{client_token}}

HTTP 200

# --- SubmitProposal (freelancer) ---
POST {{host}}/gigs/{{gig_id}}/proposals
Authorization: Bearer {{freelancer_token}}
Content-Type: application/json
{
  "bid_amount": 2500
}

HTTP 200
[Captures]
proposal_id: jsonpath "$.proposal.id"
[Asserts]
jsonpath "$.proposal.status" == "pending"

# --- RejectProposal (client) → proposal: pending → rejected ---
POST {{host}}/proposals/{{proposal_id}}/reject
Authorization: Bearer {{client_token}}

HTTP 200
[Asserts]
jsonpath "$.proposal.status" == "rejected"
