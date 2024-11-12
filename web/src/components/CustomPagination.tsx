import { Pagination } from "react-bootstrap"

export interface CustomPaginationProps {
  current: number
  total: number
}

export default function CustomPagination({
  current,
  total
}: Readonly<CustomPaginationProps>) {
  const pages = []
  for (let number = 1; number <= total; number++) {
    pages.push(
      <Pagination.Item
        key={number}
        active={number === current}
        href={`?page=${number.toString()}`}
      >
        {number}
      </Pagination.Item>
    )
  }

  return <Pagination className="justify-content-center">{pages}</Pagination>
}
