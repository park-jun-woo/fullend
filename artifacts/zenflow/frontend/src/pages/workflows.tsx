// fullend:gen ssot=frontend/workflows.html contract=8080435
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
      <h1>Workflows</h1>
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
        </div>
      )}
    </div>
  )
}
