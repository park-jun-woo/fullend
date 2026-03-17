# File Report: Execute workflow with report → Check report key → Get report metadata

# Create org
POST {{host}}/organizations
Content-Type: application/json
{
    "name": "Report Corp",
    "plan_type": "pro",
    "credits_balance": 10
}
HTTP 201
[Captures]
org_id: jsonpath "$.organization.id"

# Register admin
POST {{host}}/users/register
Content-Type: application/json
{
    "org_id": {{org_id}},
    "email": "rpt-admin@test.com",
    "password": "pass1234",
    "role": "admin",
    "name": "Report Admin"
}
HTTP 201

# Login
POST {{host}}/users/login
Content-Type: application/json
{
    "email": "rpt-admin@test.com",
    "password": "pass1234"
}
HTTP 200
[Captures]
token: jsonpath "$.access_token"

# Create workflow
POST {{host}}/workflows
Authorization: Bearer {{token}}
Content-Type: application/json
{
    "title": "Report Test WF",
    "trigger_event": "manual"
}
HTTP 201
[Captures]
wf_id: jsonpath "$.workflow.id"

# Activate
POST {{host}}/workflows/{{wf_id}}/activate
Authorization: Bearer {{token}}
HTTP 200

# Execute with report
POST {{host}}/workflows/{{wf_id}}/execute-with-report
Authorization: Bearer {{token}}
HTTP 200
[Captures]
log_id: jsonpath "$.log.id"
report_key: jsonpath "$.report_key"
[Asserts]
jsonpath "$.log.status" == "completed"
jsonpath "$.report_key" startsWith "reports/"

# Get execution report metadata
GET {{host}}/execution-logs/{{log_id}}/report
Authorization: Bearer {{token}}
HTTP 200
[Asserts]
jsonpath "$.report_key" == {{report_key}}
jsonpath "$.status" == "completed"
