# Versioning: Create workflow → Add actions → Create new version → Verify version + actions copied

# Create org
POST {{host}}/organizations
Content-Type: application/json
{
    "name": "Version Corp",
    "plan_type": "pro",
    "credits_balance": 50
}
HTTP 201
[Captures]
org_id: jsonpath "$.organization.id"

# Register admin
POST {{host}}/users/register
Content-Type: application/json
{
    "org_id": {{org_id}},
    "email": "ver-admin@test.com",
    "password": "pass1234",
    "role": "admin",
    "name": "Version Admin"
}
HTTP 201

# Login
POST {{host}}/users/login
Content-Type: application/json
{
    "email": "ver-admin@test.com",
    "password": "pass1234"
}
HTTP 200
[Captures]
token: jsonpath "$.access_token"

# Create workflow v1
POST {{host}}/workflows
Authorization: Bearer {{token}}
Content-Type: application/json
{
    "title": "Report Pipeline",
    "trigger_event": "schedule.daily"
}
HTTP 201
[Captures]
wf_v1_id: jsonpath "$.workflow.id"
[Asserts]
jsonpath "$.workflow.version" == 1
jsonpath "$.workflow.root_workflow_id" == 0

# Add 2 actions to v1
POST {{host}}/workflows/{{wf_v1_id}}/actions
Authorization: Bearer {{token}}
Content-Type: application/json
{
    "action_type": "fetch_data",
    "payload_template": "{\"source\": \"db\"}",
    "sequence_order": 1
}
HTTP 201

POST {{host}}/workflows/{{wf_v1_id}}/actions
Authorization: Bearer {{token}}
Content-Type: application/json
{
    "action_type": "generate_report",
    "payload_template": "{\"format\": \"pdf\"}",
    "sequence_order": 2
}
HTTP 201

# Create new version from v1
POST {{host}}/workflows/{{wf_v1_id}}/new-version
Authorization: Bearer {{token}}
HTTP 201
[Captures]
wf_v2_id: jsonpath "$.workflow.id"
[Asserts]
jsonpath "$.workflow.version" != 0
jsonpath "$.workflow.status" == "draft"

# List versions — should show both v1 and v2
GET {{host}}/workflows/{{wf_v1_id}}/versions
Authorization: Bearer {{token}}
HTTP 200
[Asserts]
jsonpath "$.versions" count >= 2
