# Invariant: Client cannot ApproveWork on an open gig (wrong state)

# Register client
POST {{host}}/users/register
Content-Type: application/json
{
    "email": "state-client@test.com",
    "password": "pass1234",
    "role": "client",
    "name": "State Client"
}
HTTP 201

# Login client
POST {{host}}/users/login
Content-Type: application/json
{
    "email": "state-client@test.com",
    "password": "pass1234"
}
HTTP 200
[Captures]
client_token: jsonpath "$.access_token"

# Client creates gig
POST {{host}}/gigs
Authorization: Bearer {{client_token}}
Content-Type: application/json
{
    "title": "State Test Gig",
    "description": "Testing invalid state transition",
    "budget": 30000
}
HTTP 201
[Captures]
gig_id: jsonpath "$.gig.id"

# Publish gig
PUT {{host}}/gigs/{{gig_id}}/publish
Authorization: Bearer {{client_token}}
HTTP 200

# Try to approve work on an open gig → 409
POST {{host}}/gigs/{{gig_id}}/approve
Authorization: Bearer {{client_token}}
HTTP 409
