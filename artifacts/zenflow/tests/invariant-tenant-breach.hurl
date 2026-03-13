# Invariant: Tenant Breach — Org A user cannot access Org B workflow

# Create Org A
POST {{host}}/organizations
Content-Type: application/json
{
  "name": "OrgA",
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
  "name": "OrgB",
  "plan_type": "pro",
  "credits_balance": 100
}
HTTP 201
[Captures]
org_b_id: jsonpath "$.organization.id"

# Register user in Org A
POST {{host}}/auth/register
Content-Type: application/json
{
  "email": "usera@orga.com",
  "password": "secret123",
  "org_id": {{org_a_id}}
}
HTTP 201

# Login as Org A user
POST {{host}}/auth/login
Content-Type: application/json
{
  "email": "usera@orga.com",
  "password": "secret123"
}
HTTP 200
[Captures]
token_a: jsonpath "$.access_token"

# Register admin in Org B
POST {{host}}/auth/register
Content-Type: application/json
{
  "email": "adminb@orgb.com",
  "password": "secret123",
  "org_id": {{org_b_id}}
}
HTTP 201

# Login as Org B admin
POST {{host}}/auth/login
Content-Type: application/json
{
  "email": "adminb@orgb.com",
  "password": "secret123"
}
HTTP 200
[Captures]
token_b: jsonpath "$.access_token"

# Org B admin creates a workflow
POST {{host}}/workflows
Content-Type: application/json
Authorization: Bearer {{token_b}}
{
  "title": "Org B Workflow",
  "trigger_event": "webhook"
}
HTTP 201
[Captures]
wf_b_id: jsonpath "$.workflow.id"

# Org A user tries to access Org B workflow — should be 403
GET {{host}}/workflows/{{wf_b_id}}
Authorization: Bearer {{token_a}}
HTTP 403
