# Webhook: Create org → Register webhook → Execute workflow → Verify webhook list

# Create org with credits
POST {{host}}/organizations
Content-Type: application/json
{
    "name": "Webhook Corp",
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
    "email": "wh-admin@test.com",
    "password": "pass1234",
    "role": "admin",
    "name": "Webhook Admin"
}
HTTP 201

# Login
POST {{host}}/users/login
Content-Type: application/json
{
    "email": "wh-admin@test.com",
    "password": "pass1234"
}
HTTP 200
[Captures]
token: jsonpath "$.access_token"

# Register webhook
POST {{host}}/webhooks
Authorization: Bearer {{token}}
Content-Type: application/json
{
    "url": "https://hooks.example.com/zenflow",
    "event_type": "workflow.executed"
}
HTTP 201
[Captures]
webhook_id: jsonpath "$.webhook.id"
[Asserts]
jsonpath "$.webhook.event_type" == "workflow.executed"

# List webhooks
GET {{host}}/webhooks
Authorization: Bearer {{token}}
HTTP 200
[Asserts]
jsonpath "$.webhooks" count == 1

# Create and activate workflow
POST {{host}}/workflows
Authorization: Bearer {{token}}
Content-Type: application/json
{
    "title": "Webhook Test WF",
    "trigger_event": "manual"
}
HTTP 201
[Captures]
wf_id: jsonpath "$.workflow.id"

POST {{host}}/workflows/{{wf_id}}/activate
Authorization: Bearer {{token}}
HTTP 200

# Execute workflow (triggers @publish "workflow.executed")
POST {{host}}/workflows/{{wf_id}}/execute
Authorization: Bearer {{token}}
HTTP 200
[Asserts]
jsonpath "$.log.status" == "completed"

# Delete webhook
DELETE {{host}}/webhooks/{{webhook_id}}
Authorization: Bearer {{token}}
HTTP 200
[Asserts]
jsonpath "$.id" exists

# Verify webhook deleted
GET {{host}}/webhooks
Authorization: Bearer {{token}}
HTTP 200
[Asserts]
jsonpath "$.webhooks" count == 0
