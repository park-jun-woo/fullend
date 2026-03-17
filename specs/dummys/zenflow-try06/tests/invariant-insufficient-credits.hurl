# Invariant: Organization with depleted credits cannot activate a new workflow → 402

# Create org with 1 credit
POST {{host}}/organizations
Content-Type: application/json
{
    "name": "Low Credit Corp",
    "plan_type": "free",
    "credits_balance": 1
}
HTTP 201
[Captures]
org_id: jsonpath "$.organization.id"

# Register admin
POST {{host}}/users/register
Content-Type: application/json
{
    "org_id": {{org_id}},
    "email": "low-admin@test.com",
    "password": "pass1234",
    "role": "admin",
    "name": "Low Admin"
}
HTTP 201

# Login
POST {{host}}/users/login
Content-Type: application/json
{
    "email": "low-admin@test.com",
    "password": "pass1234"
}
HTTP 200
[Captures]
token: jsonpath "$.access_token"

# Create workflow A
POST {{host}}/workflows
Authorization: Bearer {{token}}
Content-Type: application/json
{
    "title": "Workflow A",
    "trigger_event": "manual"
}
HTTP 201
[Captures]
wf_a_id: jsonpath "$.workflow.id"

# Activate workflow A (1 credit > 0, succeeds)
POST {{host}}/workflows/{{wf_a_id}}/activate
Authorization: Bearer {{token}}
HTTP 200

# Execute workflow A (deducts 1 credit → 0 remaining)
POST {{host}}/workflows/{{wf_a_id}}/execute
Authorization: Bearer {{token}}
HTTP 200
[Asserts]
jsonpath "$.credits_remaining" == 0

# Create workflow B
POST {{host}}/workflows
Authorization: Bearer {{token}}
Content-Type: application/json
{
    "title": "Workflow B",
    "trigger_event": "manual"
}
HTTP 201
[Captures]
wf_b_id: jsonpath "$.workflow.id"

# Try to activate workflow B with 0 credits → 402
POST {{host}}/workflows/{{wf_b_id}}/activate
Authorization: Bearer {{token}}
HTTP 402
