// fullend:gen ssot=frontend/workflow-detail.html contract=5182303
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { useParams } from 'react-router-dom'
import { useForm } from 'react-hook-form'
import { api } from '../api'

export default function WorkflowDetail() {
  const { id } = useParams()
  const queryClient = useQueryClient()

  const { data: getWorkflowData, isLoading: getWorkflowDataLoading, error: getWorkflowDataError } = useQuery({
    queryKey: ['GetWorkflow', id],
    queryFn: () => api.GetWorkflow({ id: id }),
  })

  const addActionForm = useForm()
  const addActionMutation = useMutation({
    mutationFn: (data: any) => api.AddAction({ ...data, id: id }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['GetWorkflow'] })
    },
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

  return (
    <div>
      <title>ZenFlow Workflow Detail</title>
      {getWorkflowDataLoading && <div>로딩 중...</div>}
      {getWorkflowDataError && <div>오류가 발생했습니다</div>}
      {getWorkflowData && (
        <div>
          <h2>{getWorkflowData.workflow.title}</h2>
          <p>{getWorkflowData.workflow.status}</p>
          <p>{getWorkflowData.workflow.trigger_event}</p>
          <button onClick={() => activateWorkflowMutation.mutate({})}>Activate</button>
          <button onClick={() => pauseWorkflowMutation.mutate({})}>Pause</button>
          <button onClick={() => archiveWorkflowMutation.mutate({})}>Archive</button>
          <button onClick={() => executeWorkflowMutation.mutate({})}>Execute</button>
        </div>
      )}
      <form onSubmit={addActionForm.handleSubmit((data) => addActionMutation.mutate(data))}>
        <input placeholder="Action Type" {...addActionForm.register('action_type')} />
        <input placeholder="Payload Template" {...addActionForm.register('payload_template')} />
        <input type="number" placeholder="Order" {...addActionForm.register('sequence_order', { valueAsNumber: true })} />
        <button type="submit">Add Action</button>
      </form>
    </div>
  )
}
