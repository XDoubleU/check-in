import { type User } from "../types"
import { type Pagination } from "./shared"

export interface GetAllPaginatedUserDto {
  data: User[]
  pagination: Pagination
}

export interface CreateUserDto {
  username: string
  password: string
}

export interface UpdateUserDto {
  username?: string
  password?: string
}
