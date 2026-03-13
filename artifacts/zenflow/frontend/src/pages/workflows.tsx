// fullend:gen ssot=frontend/workflows.html contract=d343ec5
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { useForm } from 'react-hook-form'
import { api } from '../api'

export default function Workflows() {
  const queryClient = useQueryClient()

  const { data: listWorkflowsData, isLoading: listWorkflowsDataLoading, error: listWorkflowsDataError } = useQuery({
    queryKey: ['ListWorkflows'],
    queryFn: () => api.ListWorkflows(),
  })

  const createWorkflowForm = useForm()
  const createWorkflowMutation = useMutation({
    mutationFn: (data: any) => api.CreateWorkflow(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['ListWorkflows'] })
    },
  })

  return (
    <div>
      <title>ZenFlow - Workflows</title>
      <h1>Workflows</h1>
      {listWorkflowsDataLoading && <div>로딩 중...</div>}
      {listWorkflowsDataError && <div>오류가 발생했습니다</div>}
      {listWorkflowsData && (
        <section>
          <ul>
            <li>
              {listWorkflowsData.workflows?.map((item: any, index: number) => (
                <span key={index}>
                </span>
              ))}
            </li>
          </ul>
          {listWorkflowsData.workflows?.length === 0 && <p>No workflows yet.</p>}
        </section>
      )}
      <h2>Create Workflow</h2>
      <form onSubmit={createWorkflowForm.handleSubmit((data) => createWorkflowMutation.mutate(data))}>
        <input type="text" placeholder="Title" {...createWorkflowForm.register('title')} />
        <input type="text" placeholder="Trigger Event" {...createWorkflowForm.register('trigger_event')} />
        <button type="submit">Create</button>
      </form>
    </div>
  )
}
