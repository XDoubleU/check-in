import { type Location } from "../types"
import { type Pagination } from "./shared"

export interface CreateLocationDto {
  name: string
  capacity: number
  username: string
  password: string
}

export interface GetAllPaginatedLocationDto {
  data: Location[]
  pagination: Pagination
}

export interface UpdateLocationDto {
  name?: string
  capacity?: number
  username?: string
  password?: string
}
