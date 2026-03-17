# Schedule: Create org → Register admin → Create workflow → Activate → Set schedule → Get schedule → Delete schedule → Get schedule (empty)

# Step 1: Create organization
POST {{host}}/organizations
Content-Type: application/json
{
    "name": "Schedule Corp",
    "plan_type": "pro",
    "credits_balance": 100
}
HTTP 201
[Captures]
org_id: jsonpath "$.organization.id"

# Step 2: Register admin
POST {{host}}/users/register
Content-Type: application/json
{
    "org_id": {{org_id}},
    "email": "sched-admin@test.com",
    "password": "pass1234",
    "role": "admin",
    "name": "Schedule Admin"
}
HTTP 201

# Step 3: Login
POST {{host}}/users/login
Content-Type: application/json
{
    "email": "sched-admin@test.com",
    "password": "pass1234"
}
HTTP 200
[Captures]
token: jsonpath "$.access_token"

# Step 4: Create workflow
POST {{host}}/workflows
Authorization: Bearer {{token}}
Content-Type: application/json
{
    "title": "Scheduled Report",
    "trigger_event": "schedule.daily"
}
HTTP 201
[Captures]
wf_id: jsonpath "$.workflow.id"

# Step 5: Activate workflow
POST {{host}}/workflows/{{wf_id}}/activate
Authorization: Bearer {{token}}
HTTP 200

# Step 6: Set schedule (cron every day at midnight)
POST {{host}}/workflows/{{wf_id}}/schedule
Authorization: Bearer {{token}}
Content-Type: application/json
{
    "cron": "0 0 * * *"
}
HTTP 200
[Asserts]
jsonpath "$.cron" == "0 0 * * *"
jsonpath "$.next_run" exists

# Step 7: Get schedule
GET {{host}}/workflows/{{wf_id}}/schedule
Authorization: Bearer {{token}}
HTTP 200
[Asserts]
jsonpath "$.cron" == "0 0 * * *"

# Step 8: Delete schedule
DELETE {{host}}/workflows/{{wf_id}}/schedule
Authorization: Bearer {{token}}
HTTP 200
[Asserts]
jsonpath "$.id" exists
