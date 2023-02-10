import { Pagination } from "react-bootstrap"

export type CustomPaginationProps = {
  current: number,
  total: number,
  pageSize: number
}

export default function CustomPagination({current, total, pageSize}: CustomPaginationProps) {
  const pages = []
  for (let number = 1; number <= Math.ceil(total/pageSize); number++) {
    pages.push(
      <Pagination.Item key={number} active={number === current} href={`?page=${number}`}>
        {number}
      </Pagination.Item>,
    )
  }
  
  return <Pagination className="justify-content-center">{pages}</Pagination>
}