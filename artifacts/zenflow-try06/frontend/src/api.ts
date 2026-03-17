const BASE = '/api'

async function activateWorkflow(params?: Record<string, any>) {
  const id = params?.id
  const exclude = new Set(['id'])
  const body: Record<string, any> = {}
  if (params) {
    for (const [k, v] of Object.entries(params)) {
      if (!exclude.has(k)) body[k] = v
    }
  }
  const res = await fetch(`${BASE}/workflows/${id}/activate`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(body),
  })
  return res.json()
}

async function addAction(params?: Record<string, any>) {
  const id = params?.id
  const exclude = new Set(['id'])
  const body: Record<string, any> = {}
  if (params) {
    for (const [k, v] of Object.entries(params)) {
      if (!exclude.has(k)) body[k] = v
    }
  }
  const res = await fetch(`${BASE}/workflows/${id}/actions`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(body),
  })
  return res.json()
}

async function archiveWorkflow(params?: Record<string, any>) {
  const id = params?.id
  const exclude = new Set(['id'])
  const body: Record<string, any> = {}
  if (params) {
    for (const [k, v] of Object.entries(params)) {
      if (!exclude.has(k)) body[k] = v
    }
  }
  const res = await fetch(`${BASE}/workflows/${id}/archive`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(body),
  })
  return res.json()
}

async function cloneTemplate(params?: Record<string, any>) {
  const id = params?.id
  const exclude = new Set(['id'])
  const body: Record<string, any> = {}
  if (params) {
    for (const [k, v] of Object.entries(params)) {
      if (!exclude.has(k)) body[k] = v
    }
  }
  const res = await fetch(`${BASE}/templates/${id}/clone`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(body),
  })
  return res.json()
}

async function createOrganization(params?: Record<string, any>) {
  const body = params ?? {}
  const res = await fetch(`${BASE}/organizations`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(body),
  })
  return res.json()
}

async function createWebhook(params?: Record<string, any>) {
  const body = params ?? {}
  const res = await fetch(`${BASE}/webhooks`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(body),
  })
  return res.json()
}

async function createWorkflow(params?: Record<string, any>) {
  const body = params ?? {}
  const res = await fetch(`${BASE}/workflows`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(body),
  })
  return res.json()
}

async function createWorkflowVersion(params?: Record<string, any>) {
  const id = params?.id
  const exclude = new Set(['id'])
  const body: Record<string, any> = {}
  if (params) {
    for (const [k, v] of Object.entries(params)) {
      if (!exclude.has(k)) body[k] = v
    }
  }
  const res = await fetch(`${BASE}/workflows/${id}/new-version`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(body),
  })
  return res.json()
}

async function deleteWebhook(params?: Record<string, any>) {
  const id = params?.id
  const exclude = new Set(['id'])
  const body: Record<string, any> = {}
  if (params) {
    for (const [k, v] of Object.entries(params)) {
      if (!exclude.has(k)) body[k] = v
    }
  }
  const res = await fetch(`${BASE}/webhooks/${id}`, {
    method: 'DELETE',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(body),
  })
  return res.json()
}

async function executeWithReport(params?: Record<string, any>) {
  const id = params?.id
  const exclude = new Set(['id'])
  const body: Record<string, any> = {}
  if (params) {
    for (const [k, v] of Object.entries(params)) {
      if (!exclude.has(k)) body[k] = v
    }
  }
  const res = await fetch(`${BASE}/workflows/${id}/execute-with-report`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(body),
  })
  return res.json()
}

async function executeWorkflow(params?: Record<string, any>) {
  const id = params?.id
  const exclude = new Set(['id'])
  const body: Record<string, any> = {}
  if (params) {
    for (const [k, v] of Object.entries(params)) {
      if (!exclude.has(k)) body[k] = v
    }
  }
  const res = await fetch(`${BASE}/workflows/${id}/execute`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(body),
  })
  return res.json()
}

async function getExecutionReport(params?: Record<string, any>) {
  const id = params?.id
  const query = new URLSearchParams()
  if (params) {
    const exclude = new Set(['id'])
    for (const [k, v] of Object.entries(params)) {
      if (v != null && !exclude.has(k)) query.set(k, String(v))
    }
  }
  const qs = query.toString()
  const res = await fetch(`${BASE}/execution-logs/${id}/report${qs ? '?' + qs : ''}`)
  return res.json()
}

async function getTemplate(params?: Record<string, any>) {
  const id = params?.id
  const query = new URLSearchParams()
  if (params) {
    const exclude = new Set(['id'])
    for (const [k, v] of Object.entries(params)) {
      if (v != null && !exclude.has(k)) query.set(k, String(v))
    }
  }
  const qs = query.toString()
  const res = await fetch(`${BASE}/templates/${id}${qs ? '?' + qs : ''}`)
  return res.json()
}

async function getWorkflow(params?: Record<string, any>) {
  const id = params?.id
  const query = new URLSearchParams()
  if (params) {
    const exclude = new Set(['id'])
    for (const [k, v] of Object.entries(params)) {
      if (v != null && !exclude.has(k)) query.set(k, String(v))
    }
  }
  const qs = query.toString()
  const res = await fetch(`${BASE}/workflows/${id}${qs ? '?' + qs : ''}`)
  return res.json()
}

async function listExecutionLogs(params?: Record<string, any>) {
  const id = params?.id
  const query = new URLSearchParams()
  if (params) {
    const exclude = new Set(['id'])
    for (const [k, v] of Object.entries(params)) {
      if (v != null && !exclude.has(k)) query.set(k, String(v))
    }
  }
  const qs = query.toString()
  const res = await fetch(`${BASE}/workflows/${id}/logs${qs ? '?' + qs : ''}`)
  return res.json()
}

async function listTemplates(params?: Record<string, any>) {
  const query = new URLSearchParams()
  if (params) {
    for (const [k, v] of Object.entries(params)) {
      if (v != null) query.set(k, String(v))
    }
  }
  const qs = query.toString()
  const res = await fetch(`${BASE}/templates${qs ? '?' + qs : ''}`)
  return res.json()
}

async function listWebhooks(params?: Record<string, any>) {
  const query = new URLSearchParams()
  if (params) {
    for (const [k, v] of Object.entries(params)) {
      if (v != null) query.set(k, String(v))
    }
  }
  const qs = query.toString()
  const res = await fetch(`${BASE}/webhooks${qs ? '?' + qs : ''}`)
  return res.json()
}

async function listWorkflowVersions(params?: Record<string, any>) {
  const id = params?.id
  const query = new URLSearchParams()
  if (params) {
    const exclude = new Set(['id'])
    for (const [k, v] of Object.entries(params)) {
      if (v != null && !exclude.has(k)) query.set(k, String(v))
    }
  }
  const qs = query.toString()
  const res = await fetch(`${BASE}/workflows/${id}/versions${qs ? '?' + qs : ''}`)
  return res.json()
}

async function listWorkflows(params?: Record<string, any>) {
  const query = new URLSearchParams()
  if (params) {
    for (const [k, v] of Object.entries(params)) {
      if (v != null) query.set(k, String(v))
    }
  }
  const qs = query.toString()
  const res = await fetch(`${BASE}/workflows${qs ? '?' + qs : ''}`)
  return res.json()
}

async function login(params?: Record<string, any>) {
  const body = params ?? {}
  const res = await fetch(`${BASE}/users/login`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(body),
  })
  return res.json()
}

async function pauseWorkflow(params?: Record<string, any>) {
  const id = params?.id
  const exclude = new Set(['id'])
  const body: Record<string, any> = {}
  if (params) {
    for (const [k, v] of Object.entries(params)) {
      if (!exclude.has(k)) body[k] = v
    }
  }
  const res = await fetch(`${BASE}/workflows/${id}/pause`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(body),
  })
  return res.json()
}

async function publishTemplate(params?: Record<string, any>) {
  const body = params ?? {}
  const res = await fetch(`${BASE}/templates`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(body),
  })
  return res.json()
}

async function register(params?: Record<string, any>) {
  const body = params ?? {}
  const res = await fetch(`${BASE}/users/register`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(body),
  })
  return res.json()
}

export const api = {
  ActivateWorkflow: activateWorkflow,
  AddAction: addAction,
  ArchiveWorkflow: archiveWorkflow,
  CloneTemplate: cloneTemplate,
  CreateOrganization: createOrganization,
  CreateWebhook: createWebhook,
  CreateWorkflow: createWorkflow,
  CreateWorkflowVersion: createWorkflowVersion,
  DeleteWebhook: deleteWebhook,
  ExecuteWithReport: executeWithReport,
  ExecuteWorkflow: executeWorkflow,
  GetExecutionReport: getExecutionReport,
  GetTemplate: getTemplate,
  GetWorkflow: getWorkflow,
  ListExecutionLogs: listExecutionLogs,
  ListTemplates: listTemplates,
  ListWebhooks: listWebhooks,
  ListWorkflowVersions: listWorkflowVersions,
  ListWorkflows: listWorkflows,
  Login: login,
  PauseWorkflow: pauseWorkflow,
  PublishTemplate: publishTemplate,
  Register: register
}
