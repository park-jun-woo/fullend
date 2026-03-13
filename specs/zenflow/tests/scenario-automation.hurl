# Happy Path: Admin creates workflow, adds actions, activates, executes

# Step 1: Create Organization
POST {{host}}/organizations
Content-Type: application/json
{
  "name": "Acme Corp",
  "plan_type": "pro",
  "credits_balance": 100
}
HTTP 201
[Captures]
org_id: jsonpath "$.organization.id"

# Step 2: Register admin user
POST {{host}}/auth/register
Content-Type: application/json
{
  "org_id": {{org_id}},
  "email": "admin@acme.com",
  "password": "secret123",
  "role": "admin"
}
HTTP 201
[Captures]
user_id: jsonpath "$.user.id"

# Step 3: Login
POST {{host}}/auth/login
Content-Type: application/json
{
  "email": "admin@acme.com",
  "password": "secret123"
}
HTTP 200
[Captures]
token: jsonpath "$.access_token"

# Step 4: Create Workflow
POST {{host}}/workflows
Content-Type: application/json
Authorization: Bearer {{token}}
{
  "title": "Email Automation",
  "trigger_event": "webhook_received"
}
HTTP 201
[Captures]
workflow_id: jsonpath "$.workflow.id"

# Step 5: Add Action 1
POST {{host}}/workflows/{{workflow_id}}/actions
Content-Type: application/json
Authorization: Bearer {{token}}
{
  "action_type": "send_email",
  "payload_template": {"to": "user@example.com"},
  "sequence_order": 1
}
HTTP 201

# Step 6: Add Action 2
POST {{host}}/workflows/{{workflow_id}}/actions
Content-Type: application/json
Authorization: Bearer {{token}}
{
  "action_type": "http_request",
  "payload_template": {"url": "https://api.example.com/notify"},
  "sequence_order": 2
}
HTTP 201

# Step 7: Activate Workflow
POST {{host}}/workflows/{{workflow_id}}/activate
Content-Type: application/json
Authorization: Bearer {{token}}
HTTP 200
[Asserts]
jsonpath "$.workflow.status" == "active"

# Step 8: Execute Workflow
POST {{host}}/workflows/{{workflow_id}}/execute
Content-Type: application/json
Authorization: Bearer {{token}}
HTTP 200
[Asserts]
jsonpath "$.execution.status" == "completed"
jsonpath "$.execution.credits_spent" == 1

# Step 9: Verify actions listed
GET {{host}}/workflows/{{workflow_id}}/actions
Authorization: Bearer {{token}}
HTTP 200
[Asserts]
jsonpath "$.actions" count == 2
