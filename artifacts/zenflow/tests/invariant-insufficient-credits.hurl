# Invariant: Insufficient Credits — 0 credits org cannot activate workflow

# Create org with 0 credits
POST {{host}}/organizations
Content-Type: application/json
{
  "name": "BrokeOrg",
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
  "email": "admin@brokeorg.com",
  "password": "secret123",
  "org_id": {{org_id}}
}
HTTP 201

# Login
POST {{host}}/auth/login
Content-Type: application/json
{
  "email": "admin@brokeorg.com",
  "password": "secret123"
}
HTTP 200
[Captures]
token: jsonpath "$.access_token"

# Create workflow
POST {{host}}/workflows
Content-Type: application/json
Authorization: Bearer {{token}}
{
  "title": "Broke Workflow",
  "trigger_event": "webhook"
}
HTTP 201
[Captures]
wf_id: jsonpath "$.workflow.id"

# Try to activate — should fail with 402
POST {{host}}/workflows/{{wf_id}}/activate
Content-Type: application/json
Authorization: Bearer {{token}}
HTTP 402
