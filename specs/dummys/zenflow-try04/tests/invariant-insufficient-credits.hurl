# Invariant: Org with 0 credits cannot activate workflow

# Register admin in zero-credits org
POST {{host}}/auth/register
{
    "email": "admin@nocredits.test",
    "password": "Pass1234",
    "org_id": 3,
    "role": "admin"
}
HTTP 201

# Login
POST {{host}}/auth/login
{
    "email": "admin@nocredits.test",
    "password": "Pass1234"
}
HTTP 200
[Captures]
nocredits_token: jsonpath "$.access_token"

# Create workflow
POST {{host}}/workflows
Authorization: Bearer {{nocredits_token}}
{
    "title": "No Credits Workflow",
    "trigger_event": "test_event"
}
HTTP 201
[Captures]
nocredits_wf_id: jsonpath "$.workflow.id"

# Try to activate — should fail with 402
POST {{host}}/workflows/{{nocredits_wf_id}}/activate
Authorization: Bearer {{nocredits_token}}
HTTP 402
