// fullend:gen ssot=frontend/workflows.html contract=3ac0d73
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
      <title>ZenFlow Workflows</title>
      {listWorkflowsDataLoading && <div>로딩 중...</div>}
      {listWorkflowsDataError && <div>오류가 발생했습니다</div>}
      {listWorkflowsData && (
        <div>
          <ul>
            {listWorkflowsData.workflows?.map((item: any, index: number) => (
              <li key={index}>
                <span>{item.title}</span>
                <span>{item.status}</span>
              </li>
            ))}
          </ul>
          {listWorkflowsData.workflows?.length === 0 && <div>No workflows yet.</div>}
        </div>
      )}
      <form onSubmit={createWorkflowForm.handleSubmit((data) => createWorkflowMutation.mutate(data))}>
        <input placeholder="Title" {...createWorkflowForm.register('title')} />
        <input placeholder="Trigger Event" {...createWorkflowForm.register('trigger_event')} />
        <button type="submit">Create Workflow</button>
      </form>
    </div>
  )
}
