# Invariant: Tenant Breach — Org A user cannot access Org B workflow

# Setup Org A
POST {{host}}/organizations
Content-Type: application/json
{
  "name": "Org A",
  "plan_type": "pro",
  "credits_balance": 50
}
HTTP 201
[Captures]
org_a_id: jsonpath "$.organization.id"

# Setup Org B
POST {{host}}/organizations
Content-Type: application/json
{
  "name": "Org B",
  "plan_type": "pro",
  "credits_balance": 50
}
HTTP 201
[Captures]
org_b_id: jsonpath "$.organization.id"

# Register Org A admin
POST {{host}}/auth/register
Content-Type: application/json
{
  "org_id": {{org_a_id}},
  "email": "admin-a@orga.com",
  "password": "secret123",
  "role": "admin"
}
HTTP 201

# Register Org B admin
POST {{host}}/auth/register
Content-Type: application/json
{
  "org_id": {{org_b_id}},
  "email": "admin-b@orgb.com",
  "password": "secret123",
  "role": "admin"
}
HTTP 201

# Login as Org A admin
POST {{host}}/auth/login
Content-Type: application/json
{
  "email": "admin-a@orga.com",
  "password": "secret123"
}
HTTP 200
[Captures]
token_a: jsonpath "$.access_token"

# Login as Org B admin
POST {{host}}/auth/login
Content-Type: application/json
{
  "email": "admin-b@orgb.com",
  "password": "secret123"
}
HTTP 200
[Captures]
token_b: jsonpath "$.access_token"

# Org A creates a workflow
POST {{host}}/workflows
Content-Type: application/json
Authorization: Bearer {{token_a}}
{
  "title": "Org A Workflow",
  "trigger_event": "webhook_received"
}
HTTP 201
[Captures]
workflow_a_id: jsonpath "$.workflow.id"

# Org B tries to access Org A's workflow — should get 404 (tenant isolation)
GET {{host}}/workflows/{{workflow_a_id}}
Authorization: Bearer {{token_b}}
HTTP 404
