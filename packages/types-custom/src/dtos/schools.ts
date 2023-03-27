import { type School } from "../types"
import { type Pagination } from "./shared"

export interface CreateSchoolDto {
  name: string
}

export interface GetAllPaginatedSchoolDto {
  data: School[]
  pagination: Pagination
}

export interface UpdateSchoolDto {
  name: string
}
