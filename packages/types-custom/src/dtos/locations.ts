import { type Location } from "../types"

export interface CreateLocationDto {
  name: string
  capacity: number
  username: string
  password: string
}

export interface GetAllPaginatedLocationDto {
  page: number
  totalPages: number
  locations: Location[]
}

export interface UpdateLocationDto {
  name?: string
  capacity?: number
  username?: string
  password?: string
}