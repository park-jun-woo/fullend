// fullend:gen ssot=frontend/gig-detail.html contract=b640e8c
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { useParams } from 'react-router-dom'
import { api } from '../api'

export default function GigDetail() {
  const { id } = useParams()
  const queryClient = useQueryClient()

  const { data: getGigData, isLoading: getGigDataLoading, error: getGigDataError } = useQuery({
    queryKey: ['GetGig', id],
    queryFn: () => api.GetGig({ id: id }),
  })

  const publishGigMutation = useMutation({
    mutationFn: (data: any) => api.PublishGig({ ...data, id: id }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['GetGig'] })
    },
  })

  const submitWorkMutation = useMutation({
    mutationFn: (data: any) => api.SubmitWork({ ...data, id: id }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['GetGig'] })
    },
  })

  const approveWorkMutation = useMutation({
    mutationFn: (data: any) => api.ApproveWork({ ...data, id: id }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['GetGig'] })
    },
  })

  return (
    <div>
      <title>Gig Detail</title>
      {getGigDataLoading && <div>로딩 중...</div>}
      {getGigDataError && <div>오류가 발생했습니다</div>}
      {getGigData && (
        <div>
          <h1>{getGigData.gig.title}</h1>
          <p>{getGigData.gig.description}</p>
          <span>{getGigData.gig.budget}</span>
          <span>{getGigData.gig.status}</span>
          <button onClick={() => publishGigMutation.mutate({})}>Publish</button>
          <button onClick={() => submitWorkMutation.mutate({})}>Submit Work</button>
          <button onClick={() => approveWorkMutation.mutate({})}>Approve</button>
        </div>
      )}
    </div>
  )
}
