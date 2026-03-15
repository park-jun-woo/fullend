# Happy Path: Admin creates workflow, adds actions, activates, executes

# Register admin
POST {{host}}/auth/register
{
    "email": "admin@zenflow.test",
    "password": "Pass1234",
    "org_id": 1,
    "role": "admin"
}
HTTP 201
[Captures]
admin_id: jsonpath "$.user.id"

# Login admin
POST {{host}}/auth/login
{
    "email": "admin@zenflow.test",
    "password": "Pass1234"
}
HTTP 200
[Captures]
admin_token: jsonpath "$.access_token"

# Create workflow
POST {{host}}/workflows
Authorization: Bearer {{admin_token}}
{
    "title": "Daily Report",
    "trigger_event": "daily_cron"
}
HTTP 201
[Captures]
workflow_id: jsonpath "$.workflow.id"
[Asserts]
jsonpath "$.workflow.status" == "draft"

# Add action 1
POST {{host}}/workflows/{{workflow_id}}/actions
Authorization: Bearer {{admin_token}}
{
    "action_type": "fetch_data",
    "payload_template": "SELECT * FROM reports",
    "sequence_order": 1
}
HTTP 201

# Add action 2
POST {{host}}/workflows/{{workflow_id}}/actions
Authorization: Bearer {{admin_token}}
{
    "action_type": "send_email",
    "payload_template": "Send report to team",
    "sequence_order": 2
}
HTTP 201

# Activate workflow
POST {{host}}/workflows/{{workflow_id}}/activate
Authorization: Bearer {{admin_token}}
HTTP 200
[Asserts]
jsonpath "$.workflow.status" == "active"

# Execute workflow
POST {{host}}/workflows/{{workflow_id}}/execute
Authorization: Bearer {{admin_token}}
HTTP 200
[Asserts]
jsonpath "$.execution_log.status" == "completed"
jsonpath "$.execution_log.credits_spent" == 1
