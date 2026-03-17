// fullend:gen ssot=frontend/workflows.html contract=fad398f
import { useQuery } from '@tanstack/react-query'
import { api } from '../api'

export default function Workflows() {

  const { data: listWorkflowsData, isLoading: listWorkflowsDataLoading, error: listWorkflowsDataError } = useQuery({
    queryKey: ['ListWorkflows'],
    queryFn: () => api.ListWorkflows(),
  })

  return (
    <div>
      <title>Workflows</title>
      {listWorkflowsDataLoading && <div>로딩 중...</div>}
      {listWorkflowsDataError && <div>오류가 발생했습니다</div>}
      {listWorkflowsData && (
        <section>
          <ul>
            {listWorkflowsData.workflows?.map((item: any, index: number) => (
              <li key={index}>
                <h3>{item.title}</h3>
                <span>{item.status}</span>
                <span>{item.trigger_event}</span>
              </li>
            ))}
          </ul>
        </section>
      )}
    </div>
  )
}
