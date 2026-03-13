package authz

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
    input.action == "GetWorkflow"
    input.resource == "workflow"
}

allow if {
    input.action == "ActivateWorkflow"
    input.resource == "workflow"
    input.claims.role == "admin"
}

allow if {
    input.action == "PauseWorkflow"
    input.resource == "workflow"
    input.claims.role == "admin"
}

allow if {
    input.action == "ArchiveWorkflow"
    input.resource == "workflow"
    input.claims.role == "admin"
}

allow if {
    input.action == "ExecuteWorkflow"
    input.resource == "workflow"
}

allow if {
    input.action == "CreateAction"
    input.resource == "action"
    input.claims.role == "admin"
}

allow if {
    input.action == "ListActions"
    input.resource == "action"
}
