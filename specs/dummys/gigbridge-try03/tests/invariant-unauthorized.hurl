# Invariant: Freelancer B cannot submit work on Freelancer A's gig

# Register client
POST {{host}}/auth/register
Content-Type: application/json
{
    "email": "inv-client@test.com",
    "password": "pass1234",
    "role": "client",
    "name": "Inv Client"
}
HTTP 201
[Captures]
inv_client_id: jsonpath "$.user.id"

# Login client
POST {{host}}/auth/login
Content-Type: application/json
{
    "email": "inv-client@test.com",
    "password": "pass1234"
}
HTTP 200
[Captures]
inv_client_token: jsonpath "$.access_token"

# Register freelancer A
POST {{host}}/auth/register
Content-Type: application/json
{
    "email": "inv-fl-a@test.com",
    "password": "pass1234",
    "role": "freelancer",
    "name": "Freelancer A"
}
HTTP 201

# Login freelancer A
POST {{host}}/auth/login
Content-Type: application/json
{
    "email": "inv-fl-a@test.com",
    "password": "pass1234"
}
HTTP 200
[Captures]
inv_fl_a_token: jsonpath "$.access_token"

# Register freelancer B
POST {{host}}/auth/register
Content-Type: application/json
{
    "email": "inv-fl-b@test.com",
    "password": "pass1234",
    "role": "freelancer",
    "name": "Freelancer B"
}
HTTP 201

# Login freelancer B
POST {{host}}/auth/login
Content-Type: application/json
{
    "email": "inv-fl-b@test.com",
    "password": "pass1234"
}
HTTP 200
[Captures]
inv_fl_b_token: jsonpath "$.access_token"

# Client creates + publishes gig
POST {{host}}/gigs
Content-Type: application/json
Authorization: Bearer {{inv_client_token}}
{
    "title": "Inv Gig",
    "description": "Test",
    "budget": 500
}
HTTP 201
[Captures]
inv_gig_id: jsonpath "$.gig.id"

PUT {{host}}/gigs/{{inv_gig_id}}/publish
Content-Type: application/json
Authorization: Bearer {{inv_client_token}}
HTTP 200

# Freelancer A submits proposal
POST {{host}}/gigs/{{inv_gig_id}}/proposals
Content-Type: application/json
Authorization: Bearer {{inv_fl_a_token}}
{
    "bid_amount": 450
}
HTTP 201
[Captures]
inv_proposal_id: jsonpath "$.proposal.id"

# Client accepts proposal
POST {{host}}/proposals/{{inv_proposal_id}}/accept
Content-Type: application/json
Authorization: Bearer {{inv_client_token}}
HTTP 200

# Freelancer B tries to submit work on Freelancer A's gig → 403
POST {{host}}/gigs/{{inv_gig_id}}/submit-work
Content-Type: application/json
Authorization: Bearer {{inv_fl_b_token}}
HTTP 403
