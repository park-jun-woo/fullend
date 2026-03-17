// fullend:gen ssot=frontend/workflow-detail.html contract=8bfe0d9
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { useParams } from 'react-router-dom'
import { api } from '../api'

export default function WorkflowDetail() {
  const { id } = useParams()
  const queryClient = useQueryClient()

  const { data: getWorkflowData, isLoading: getWorkflowDataLoading, error: getWorkflowDataError } = useQuery({
    queryKey: ['GetWorkflow', id],
    queryFn: () => api.GetWorkflow({ id: id }),
  })

  const activateWorkflowMutation = useMutation({
    mutationFn: (data: any) => api.ActivateWorkflow({ ...data, id: id }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['GetWorkflow'] })
    },
  })

  const pauseWorkflowMutation = useMutation({
    mutationFn: (data: any) => api.PauseWorkflow({ ...data, id: id }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['GetWorkflow'] })
    },
  })

  const archiveWorkflowMutation = useMutation({
    mutationFn: (data: any) => api.ArchiveWorkflow({ ...data, id: id }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['GetWorkflow'] })
    },
  })

  const executeWorkflowMutation = useMutation({
    mutationFn: (data: any) => api.ExecuteWorkflow({ ...data, id: id }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['GetWorkflow'] })
    },
  })

  const createWorkflowVersionMutation = useMutation({
    mutationFn: (data: any) => api.CreateWorkflowVersion({ ...data, id: id }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['GetWorkflow'] })
    },
  })

  return (
    <div>
      <title>Workflow Detail</title>
      {getWorkflowDataLoading && <div>로딩 중...</div>}
      {getWorkflowDataError && <div>오류가 발생했습니다</div>}
      {getWorkflowData && (
        <section>
          <h2>{getWorkflowData.workflow.title}</h2>
          <span>{getWorkflowData.workflow.status}</span>
          <span>{getWorkflowData.workflow.trigger_event}</span>
          <span>{getWorkflowData.workflow.version}</span>
          <button onClick={() => activateWorkflowMutation.mutate({})}>Activate</button>
          <button onClick={() => activateWorkflowMutation.mutate({})}>Resume</button>
          <button onClick={() => pauseWorkflowMutation.mutate({})}>Pause</button>
          <button onClick={() => archiveWorkflowMutation.mutate({})}>Archive</button>
          <button onClick={() => executeWorkflowMutation.mutate({})}>Execute</button>
          <button onClick={() => createWorkflowVersionMutation.mutate({})}>New Version</button>
        </section>
      )}
    </div>
  )
}
