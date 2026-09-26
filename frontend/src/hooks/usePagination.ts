// 分页 Hook：page/pageSize/total 与翻页动作
import { useCallback, useState } from 'react'

export function usePagination(initialSize = 10) {
  const [page, setPage] = useState(1)
  const [pageSize, setPageSize] = useState(initialSize)
  const [total, setTotal] = useState(0)

  const onPageChange = useCallback((p: number) => {
    setPage(p)
    window.scrollTo({ top: 0, behavior: 'smooth' })
  }, [])

  return { page, pageSize, total, setTotal, onPageChange, setPage, setPageSize }
}
