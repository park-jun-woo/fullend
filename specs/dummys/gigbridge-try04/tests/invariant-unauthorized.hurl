# Invariant: Freelancer B cannot SubmitWork on Freelancer A's gig

# Register client
POST {{host}}/users/register
Content-Type: application/json
{
    "email": "inv-client@test.com",
    "password": "pass1234",
    "role": "client",
    "name": "Inv Client"
}
HTTP 201

# Register freelancer A
POST {{host}}/users/register
Content-Type: application/json
{
    "email": "inv-freelancer-a@test.com",
    "password": "pass1234",
    "role": "freelancer",
    "name": "Freelancer A"
}
HTTP 201

# Register freelancer B
POST {{host}}/users/register
Content-Type: application/json
{
    "email": "inv-freelancer-b@test.com",
    "password": "pass1234",
    "role": "freelancer",
    "name": "Freelancer B"
}
HTTP 201

# Login client
POST {{host}}/users/login
Content-Type: application/json
{
    "email": "inv-client@test.com",
    "password": "pass1234"
}
HTTP 200
[Captures]
client_token: jsonpath "$.access_token"

# Login freelancer A
POST {{host}}/users/login
Content-Type: application/json
{
    "email": "inv-freelancer-a@test.com",
    "password": "pass1234"
}
HTTP 200
[Captures]
freelancer_a_token: jsonpath "$.access_token"

# Login freelancer B
POST {{host}}/users/login
Content-Type: application/json
{
    "email": "inv-freelancer-b@test.com",
    "password": "pass1234"
}
HTTP 200
[Captures]
freelancer_b_token: jsonpath "$.access_token"

# Client creates and publishes gig
POST {{host}}/gigs
Authorization: Bearer {{client_token}}
Content-Type: application/json
{
    "title": "Invariant Test Gig",
    "description": "Testing unauthorized access",
    "budget": 50000
}
HTTP 201
[Captures]
gig_id: jsonpath "$.gig.id"

PUT {{host}}/gigs/{{gig_id}}/publish
Authorization: Bearer {{client_token}}
HTTP 200

# Freelancer A submits proposal
POST {{host}}/gigs/{{gig_id}}/proposals
Authorization: Bearer {{freelancer_a_token}}
Content-Type: application/json
{
    "bid_amount": 40000
}
HTTP 201
[Captures]
proposal_id: jsonpath "$.proposal.id"

# Client accepts proposal (assigns freelancer A)
POST {{host}}/proposals/{{proposal_id}}/accept
Authorization: Bearer {{client_token}}
HTTP 200

# Freelancer B tries to submit work on Freelancer A's gig → 403
POST {{host}}/gigs/{{gig_id}}/submit-work
Authorization: Bearer {{freelancer_b_token}}
HTTP 403
