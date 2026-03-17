package authz

import rego.v1

# @ownership workflow: workflows.org_id

default allow := false

# CreateOrganization - public (no auth needed, but endpoint is unprotected)

# CreateWorkflow - admin only
allow if {
    input.action == "CreateWorkflow"
    input.resource == "workflow"
    input.claims.role == "admin"
}

# ListWorkflows - any authenticated user (org isolation via query)
allow if {
    input.action == "ListWorkflows"
    input.resource == "workflow"
}

# GetWorkflow - any authenticated user (org isolation via query)
allow if {
    input.action == "GetWorkflow"
    input.resource == "workflow"
}

# ActivateWorkflow - admin + org match
allow if {
    input.action == "ActivateWorkflow"
    input.resource == "workflow"
    input.claims.role == "admin"
}

# PauseWorkflow - admin + org match
allow if {
    input.action == "PauseWorkflow"
    input.resource == "workflow"
    input.claims.role == "admin"
}

# ArchiveWorkflow - admin + org match
allow if {
    input.action == "ArchiveWorkflow"
    input.resource == "workflow"
    input.claims.role == "admin"
}

# ExecuteWorkflow - any authenticated member (org isolation via query)
allow if {
    input.action == "ExecuteWorkflow"
    input.resource == "workflow"
}

# AddAction - admin only
allow if {
    input.action == "AddAction"
    input.resource == "workflow"
    input.claims.role == "admin"
}

# CreateWorkflowVersion - admin only
allow if {
    input.action == "CreateWorkflowVersion"
    input.resource == "workflow"
    input.claims.role == "admin"
}

# ListWorkflowVersions - any authenticated user
allow if {
    input.action == "ListWorkflowVersions"
    input.resource == "workflow"
}

# ExecuteWithReport - any authenticated member
allow if {
    input.action == "ExecuteWithReport"
    input.resource == "workflow"
}

# GetExecutionReport - any authenticated user
allow if {
    input.action == "GetExecutionReport"
    input.resource == "workflow"
}

# PublishTemplate - admin only
allow if {
    input.action == "PublishTemplate"
    input.resource == "template"
    input.claims.role == "admin"
}

# CloneTemplate - any authenticated user
allow if {
    input.action == "CloneTemplate"
    input.resource == "template"
}

# CreateWebhook - admin only
allow if {
    input.action == "CreateWebhook"
    input.resource == "webhook"
    input.claims.role == "admin"
}

# ListWebhooks - any authenticated user
allow if {
    input.action == "ListWebhooks"
    input.resource == "webhook"
}

# DeleteWebhook - admin only
allow if {
    input.action == "DeleteWebhook"
    input.resource == "webhook"
    input.claims.role == "admin"
}

# SetSchedule - admin only
allow if {
    input.action == "SetSchedule"
    input.resource == "workflow"
    input.claims.role == "admin"
}

# GetSchedule - any authenticated user
allow if {
    input.action == "GetSchedule"
    input.resource == "workflow"
}

# DeleteSchedule - admin only
allow if {
    input.action == "DeleteSchedule"
    input.resource == "workflow"
    input.claims.role == "admin"
}

# ListExecutionLogs - any authenticated user
allow if {
    input.action == "ListExecutionLogs"
    input.resource == "workflow"
}
