# Template Marketplace: Org A publishes → duplicate blocked → Org B browses → Org B clones

# Create Org A
POST {{host}}/organizations
Content-Type: application/json
{
    "name": "Publisher Corp",
    "plan_type": "pro",
    "credits_balance": 50
}
HTTP 201
[Captures]
org_a_id: jsonpath "$.organization.id"

# Create Org B
POST {{host}}/organizations
Content-Type: application/json
{
    "name": "Consumer Corp",
    "plan_type": "pro",
    "credits_balance": 50
}
HTTP 201
[Captures]
org_b_id: jsonpath "$.organization.id"

# Register Org A admin
POST {{host}}/users/register
Content-Type: application/json
{
    "org_id": {{org_a_id}},
    "email": "pub-admin@test.com",
    "password": "pass1234",
    "role": "admin",
    "name": "Publisher Admin"
}
HTTP 201

# Register Org B admin
POST {{host}}/users/register
Content-Type: application/json
{
    "org_id": {{org_b_id}},
    "email": "con-admin@test.com",
    "password": "pass1234",
    "role": "admin",
    "name": "Consumer Admin"
}
HTTP 201

# Login Org A
POST {{host}}/users/login
Content-Type: application/json
{
    "email": "pub-admin@test.com",
    "password": "pass1234"
}
HTTP 200
[Captures]
token_a: jsonpath "$.access_token"

# Login Org B
POST {{host}}/users/login
Content-Type: application/json
{
    "email": "con-admin@test.com",
    "password": "pass1234"
}
HTTP 200
[Captures]
token_b: jsonpath "$.access_token"

# Org A creates workflow with actions
POST {{host}}/workflows
Authorization: Bearer {{token_a}}
Content-Type: application/json
{
    "title": "ETL Pipeline",
    "trigger_event": "schedule.hourly"
}
HTTP 201
[Captures]
wf_id: jsonpath "$.workflow.id"

POST {{host}}/workflows/{{wf_id}}/actions
Authorization: Bearer {{token_a}}
Content-Type: application/json
{
    "action_type": "extract",
    "payload_template": "{\"source\": \"s3\"}",
    "sequence_order": 1
}
HTTP 201

POST {{host}}/workflows/{{wf_id}}/actions
Authorization: Bearer {{token_a}}
Content-Type: application/json
{
    "action_type": "transform",
    "payload_template": "{\"format\": \"parquet\"}",
    "sequence_order": 2
}
HTTP 201

# Org A publishes workflow as template
POST {{host}}/templates
Authorization: Bearer {{token_a}}
Content-Type: application/json
{
    "source_workflow_id": {{wf_id}},
    "title": "ETL Pipeline Template",
    "description": "Hourly ETL from S3",
    "category": "data"
}
HTTP 201
[Captures]
tmpl_id: jsonpath "$.template.id"
[Asserts]
jsonpath "$.template.category" == "data"
jsonpath "$.template.clone_count" == 0

# Org A tries duplicate publish → 409 (@exists guard)
POST {{host}}/templates
Authorization: Bearer {{token_a}}
Content-Type: application/json
{
    "source_workflow_id": {{wf_id}},
    "title": "Duplicate",
    "description": "Should fail",
    "category": "data"
}
HTTP 409

# Org B browses templates (cursor pagination, public)
GET {{host}}/templates?limit=10
HTTP 200
[Asserts]
jsonpath "$.items" count >= 1
jsonpath "$.has_next" exists

# Org B gets template detail
GET {{host}}/templates/{{tmpl_id}}
HTTP 200
[Asserts]
jsonpath "$.template.title" == "ETL Pipeline Template"
jsonpath "$.org_name" == "Publisher Corp"

# Org B clones template
POST {{host}}/templates/{{tmpl_id}}/clone
Authorization: Bearer {{token_b}}
HTTP 201
[Captures]
cloned_wf_id: jsonpath "$.workflow.id"
[Asserts]
jsonpath "$.workflow.status" == "draft"
jsonpath "$.template_id" == {{tmpl_id}}

# Verify Org B now has the workflow
GET {{host}}/workflows/{{cloned_wf_id}}
Authorization: Bearer {{token_b}}
HTTP 200
[Asserts]
jsonpath "$.workflow.title" == "ETL Pipeline"
