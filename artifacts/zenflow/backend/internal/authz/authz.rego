package authz

# @ownership workflow: workflows.org_id

default allow = false

allow if {
    input.action == "CreateWorkflow"
    input.resource == "workflow"
    input.claims.role == "admin"
}

allow if {
    input.action == "ListWorkflows"
    input.resource == "workflow"
}

allow if {
    input.action == "ActivateWorkflow"
    input.resource == "workflow"
    input.claims.role == "admin"
}
