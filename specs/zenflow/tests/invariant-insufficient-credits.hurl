# Invariant: Insufficient Credits — 0 credits org cannot activate workflow

# Create Organization with 0 credits
POST {{host}}/organizations
Content-Type: application/json
{
  "name": "Broke Corp",
  "plan_type": "free",
  "credits_balance": 0
}
HTTP 201
[Captures]
org_id: jsonpath "$.organization.id"

# Register admin
POST {{host}}/auth/register
Content-Type: application/json
{
  "org_id": {{org_id}},
  "email": "admin@broke.com",
  "password": "secret123",
  "role": "admin"
}
HTTP 201

# Login
POST {{host}}/auth/login
Content-Type: application/json
{
  "email": "admin@broke.com",
  "password": "secret123"
}
HTTP 200
[Captures]
token: jsonpath "$.access_token"

# Create Workflow
POST {{host}}/workflows
Content-Type: application/json
Authorization: Bearer {{token}}
{
  "title": "Unpaid Workflow",
  "trigger_event": "timer"
}
HTTP 201
[Captures]
workflow_id: jsonpath "$.workflow.id"

# Try to activate — should fail with 402
POST {{host}}/workflows/{{workflow_id}}/activate
Content-Type: application/json
Authorization: Bearer {{token}}
HTTP 402
