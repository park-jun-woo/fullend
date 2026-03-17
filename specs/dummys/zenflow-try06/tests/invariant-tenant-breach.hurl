# Invariant: Org A user cannot access Org B's workflow

# Create Org A
POST {{host}}/organizations
Content-Type: application/json
{
    "name": "Org A",
    "plan_type": "pro",
    "credits_balance": 100
}
HTTP 201
[Captures]
org_a_id: jsonpath "$.organization.id"

# Create Org B
POST {{host}}/organizations
Content-Type: application/json
{
    "name": "Org B",
    "plan_type": "pro",
    "credits_balance": 100
}
HTTP 201
[Captures]
org_b_id: jsonpath "$.organization.id"

# Register admin for Org A
POST {{host}}/users/register
Content-Type: application/json
{
    "org_id": {{org_a_id}},
    "email": "admin-a@test.com",
    "password": "pass1234",
    "role": "admin",
    "name": "Admin A"
}
HTTP 201

# Register admin for Org B
POST {{host}}/users/register
Content-Type: application/json
{
    "org_id": {{org_b_id}},
    "email": "admin-b@test.com",
    "password": "pass1234",
    "role": "admin",
    "name": "Admin B"
}
HTTP 201

# Login Org A admin
POST {{host}}/users/login
Content-Type: application/json
{
    "email": "admin-a@test.com",
    "password": "pass1234"
}
HTTP 200
[Captures]
token_a: jsonpath "$.access_token"

# Login Org B admin
POST {{host}}/users/login
Content-Type: application/json
{
    "email": "admin-b@test.com",
    "password": "pass1234"
}
HTTP 200
[Captures]
token_b: jsonpath "$.access_token"

# Org B creates a workflow
POST {{host}}/workflows
Authorization: Bearer {{token_b}}
Content-Type: application/json
{
    "title": "Org B Workflow",
    "trigger_event": "webhook"
}
HTTP 201
[Captures]
wf_b_id: jsonpath "$.workflow.id"

# Org A tries to access Org B's workflow → 404 (org isolation via query)
GET {{host}}/workflows/{{wf_b_id}}
Authorization: Bearer {{token_a}}
HTTP 404
