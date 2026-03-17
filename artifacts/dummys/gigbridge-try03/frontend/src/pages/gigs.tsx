// fullend:gen ssot=frontend/gigs.html contract=30f9b7a
import { useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import { api } from '../api'

export default function Gigs() {

  const [page, setPage] = useState(1)
  const [limit] = useState(20)
  const [sortBy, setSortBy] = useState('created_at')
  const [sortDir, setSortDir] = useState<'asc' | 'desc'>('desc')
  const [filters, setFilters] = useState<Record<string, string>>({})

  const { data: listGigsData, isLoading: listGigsDataLoading, error: listGigsDataError } = useQuery({
    queryKey: ['ListGigs', page, limit, sortBy, sortDir, filters],
    queryFn: () => api.ListGigs({ page, limit, sortBy, sortDir, ...filters }),
  })

  return (
    <div>
      <title>Gigs</title>
      {listGigsDataLoading && <div>로딩 중...</div>}
      {listGigsDataError && <div>오류가 발생했습니다</div>}
      {listGigsData && (
        <div>
          <div className="flex gap-2 mb-4">
            <input placeholder="status" value={filters.status ?? ''} className="px-3 py-2 border rounded" onChange={(e) => setFilters(f => ({ ...f, status: e.target.value }))} />
          </div>
          <div className="flex gap-2 mb-4">
            <button onClick={() => { setSortBy('created_at'); setSortDir(d => d === 'asc' ? 'desc' : 'asc') }}>
              created_at {sortBy === 'created_at' ? (sortDir === 'asc' ? '↑' : '↓') : ''}
            </button>
          </div>
          <div>
            {listGigsData.items?.map((item: any, index: number) => (
              <h2 key={index}>
              </h2>
            ))}
          </div>
          <div className="flex justify-between items-center mt-4">
            <button disabled={page <= 1} onClick={() => setPage(p => p - 1)}>이전</button>
            <span>{page} / {Math.ceil((listGigsData?.total ?? 0) / limit)}</span>
            <button disabled={!listGigsData?.total || page * limit >= listGigsData.total} onClick={() => setPage(p => p + 1)}>다음</button>
          </div>
        </div>
      )}
    </div>
  )
}
