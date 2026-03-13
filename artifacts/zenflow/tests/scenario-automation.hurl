# Happy Path: Admin creates workflow, adds actions, activates, executes

# Step 1: Create organization
POST {{host}}/organizations
Content-Type: application/json
{
  "name": "TestOrg",
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
  "email": "admin@testorg.com",
  "password": "secret123",
  "org_id": {{org_id}}
}
HTTP 201

# Step 3: Login
POST {{host}}/auth/login
Content-Type: application/json
{
  "email": "admin@testorg.com",
  "password": "secret123"
}
HTTP 200
[Captures]
token: jsonpath "$.access_token"

# Step 4: Create workflow
POST {{host}}/workflows
Content-Type: application/json
Authorization: Bearer {{token}}
{
  "title": "Email Campaign",
  "trigger_event": "webhook_received"
}
HTTP 201
[Captures]
wf_id: jsonpath "$.workflow.id"

# Step 5: Add action 1
POST {{host}}/workflows/{{wf_id}}/actions
Content-Type: application/json
Authorization: Bearer {{token}}
{
  "action_type": "send_email",
  "sequence_order": 1
}
HTTP 201

# Step 6: Add action 2
POST {{host}}/workflows/{{wf_id}}/actions
Content-Type: application/json
Authorization: Bearer {{token}}
{
  "action_type": "http_request",
  "sequence_order": 2
}
HTTP 201

# Step 7: Activate workflow
POST {{host}}/workflows/{{wf_id}}/activate
Content-Type: application/json
Authorization: Bearer {{token}}
HTTP 200

# Step 8: Execute workflow
POST {{host}}/workflows/{{wf_id}}/execute
Content-Type: application/json
Authorization: Bearer {{token}}
HTTP 200
[Asserts]
jsonpath "$.execution.status" == "completed"
jsonpath "$.execution.credits_spent" == 1
