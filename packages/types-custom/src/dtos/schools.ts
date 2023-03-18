import { type School } from "../types"

export interface CreateSchoolDto {
  name: string
}

export interface GetAllPaginatedSchoolDto {
  page: number
  totalPages: number
  schools: School[]
}

export interface UpdateSchoolDto {
  name: string
}
