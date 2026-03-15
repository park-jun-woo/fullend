// fullend:gen ssot=frontend/workflow-detail.html contract=f328c60
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

  const executeWorkflowMutation = useMutation({
    mutationFn: (data: any) => api.ExecuteWorkflow({ ...data, id: id }),
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
          {getWorkflowDataLoading && <div>Loading...</div>}
          <h1>{getWorkflowData.workflow.title}</h1>
          <p>{getWorkflowData.workflow.status}</p>
          <p>{getWorkflowData.workflow.trigger_event}</p>
        </section>
      )}
      <form><button onClick={() => activateWorkflowMutation.mutate({})}>Activate</button></form>
      <form><button onClick={() => pauseWorkflowMutation.mutate({})}>Pause</button></form>
      <form><button onClick={() => executeWorkflowMutation.mutate({})}>Execute</button></form>
    </div>
  )
}
