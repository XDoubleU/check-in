import { MouseEventHandler } from "react"
import { Pagination } from "react-bootstrap"

type CustomPaginationProps = {
  current: number,
  total: number,
  pageSize: number,
  onClick: MouseEventHandler<HTMLElement>
}

export default function CustomPagination({current, total, pageSize, onClick}: CustomPaginationProps) {
  const pages = []
  for (let number = 1; number <= Math.ceil(total/pageSize); number++) {
    pages.push(
      <Pagination.Item key={number} active={number === current} onClick={onClick}>
        {number}
      </Pagination.Item>,
    )
  }
  
  return <Pagination className="justify-content-center">{pages}</Pagination>
}