# Invariant: Cannot approve work on a gig that is still 'open'

# Register client
POST {{host}}/auth/register
Content-Type: application/json
{
    "email": "state-client@test.com",
    "password": "pass1234",
    "role": "client",
    "name": "State Client"
}
HTTP 201

# Login client
POST {{host}}/auth/login
Content-Type: application/json
{
    "email": "state-client@test.com",
    "password": "pass1234"
}
HTTP 200
[Captures]
state_client_token: jsonpath "$.access_token"

# Create gig
POST {{host}}/gigs
Content-Type: application/json
Authorization: Bearer {{state_client_token}}
{
    "title": "State Test Gig",
    "description": "Test",
    "budget": 300
}
HTTP 201
[Captures]
state_gig_id: jsonpath "$.gig.id"

# Publish gig
PUT {{host}}/gigs/{{state_gig_id}}/publish
Content-Type: application/json
Authorization: Bearer {{state_client_token}}
HTTP 200

# Try to approve work on 'open' gig → 409
POST {{host}}/gigs/{{state_gig_id}}/approve
Content-Type: application/json
Authorization: Bearer {{state_client_token}}
HTTP 409
