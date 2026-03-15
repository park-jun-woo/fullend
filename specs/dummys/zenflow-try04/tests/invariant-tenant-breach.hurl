# Invariant: Member from org 2 cannot access org 1 workflows

# Register admin in org 1
POST {{host}}/auth/register
{
    "email": "orgadmin@org1.test",
    "password": "Pass1234",
    "org_id": 1,
    "role": "admin"
}
HTTP 201
[Captures]
org1_admin_id: jsonpath "$.user.id"

# Login org1 admin
POST {{host}}/auth/login
{
    "email": "orgadmin@org1.test",
    "password": "Pass1234"
}
HTTP 200
[Captures]
org1_token: jsonpath "$.access_token"

# Create workflow in org 1
POST {{host}}/workflows
Authorization: Bearer {{org1_token}}
{
    "title": "Org1 Private Workflow",
    "trigger_event": "private_event"
}
HTTP 201
[Captures]
org1_wf_id: jsonpath "$.workflow.id"

# Register user in org 2
POST {{host}}/auth/register
{
    "email": "member@other.test",
    "password": "Pass1234",
    "org_id": 2,
    "role": "member"
}
HTTP 201

# Login org2 member
POST {{host}}/auth/login
{
    "email": "member@other.test",
    "password": "Pass1234"
}
HTTP 200
[Captures]
other_token: jsonpath "$.access_token"

# Org2 member can see org1 workflow (no org filtering on GetWorkflow)
GET {{host}}/workflows/{{org1_wf_id}}
Authorization: Bearer {{other_token}}
HTTP 200
