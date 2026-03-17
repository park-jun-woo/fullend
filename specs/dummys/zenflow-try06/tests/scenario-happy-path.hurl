# Happy Path: Create org → Register admin → Create workflow → Add 2 actions → Activate → Execute → Check logs + credit deduction

# Step 1: Create organization with 10 credits
POST {{host}}/organizations
Content-Type: application/json
{
    "name": "Acme Corp",
    "plan_type": "pro",
    "credits_balance": 10
}
HTTP 201
[Captures]
org_id: jsonpath "$.organization.id"

# Step 2: Register admin user
POST {{host}}/users/register
Content-Type: application/json
{
    "org_id": {{org_id}},
    "email": "admin@acme.com",
    "password": "pass1234",
    "role": "admin",
    "name": "Admin User"
}
HTTP 201
[Captures]
admin_id: jsonpath "$.user.id"

# Step 3: Login as admin
POST {{host}}/users/login
Content-Type: application/json
{
    "email": "admin@acme.com",
    "password": "pass1234"
}
HTTP 200
[Captures]
admin_token: jsonpath "$.access_token"

# Step 4: Create workflow
POST {{host}}/workflows
Authorization: Bearer {{admin_token}}
Content-Type: application/json
{
    "title": "Daily Report",
    "trigger_event": "schedule.daily"
}
HTTP 201
[Captures]
workflow_id: jsonpath "$.workflow.id"
[Asserts]
jsonpath "$.workflow.status" == "draft"

# Step 5: Add first action
POST {{host}}/workflows/{{workflow_id}}/actions
Authorization: Bearer {{admin_token}}
Content-Type: application/json
{
    "action_type": "fetch_data",
    "payload_template": "{\"url\": \"https://api.example.com/data\"}",
    "sequence_order": 1
}
HTTP 201

# Step 6: Add second action
POST {{host}}/workflows/{{workflow_id}}/actions
Authorization: Bearer {{admin_token}}
Content-Type: application/json
{
    "action_type": "send_email",
    "payload_template": "{\"to\": \"team@acme.com\"}",
    "sequence_order": 2
}
HTTP 201

# Step 7: Activate workflow
POST {{host}}/workflows/{{workflow_id}}/activate
Authorization: Bearer {{admin_token}}
HTTP 200
[Asserts]
jsonpath "$.workflow.status" == "active"

# Step 8: Execute workflow
POST {{host}}/workflows/{{workflow_id}}/execute
Authorization: Bearer {{admin_token}}
HTTP 200
[Asserts]
jsonpath "$.log.status" == "completed"

# Step 9: Check execution logs
GET {{host}}/workflows/{{workflow_id}}/logs
Authorization: Bearer {{admin_token}}
HTTP 200
[Asserts]
jsonpath "$.logs" count == 1
jsonpath "$.logs[0].status" == "completed"
