import { fetchHandler } from "./fetchHandler"
import { type APIResponse } from "./types"
import { validate as isValidUUID } from "uuid"
import {
  type CreateUserDto,
  type PaginatedUsersDto,
  type UpdateUserDto,
  type User
} from "./types/apiTypes"

const USERS_ENDPOINT = "users"

export async function getMyUser(): Promise<APIResponse<User>> {
  return await fetchHandler(`current-user`)
}

export async function getUser(id: string): Promise<APIResponse<User>> {
  return await fetchHandler(`${USERS_ENDPOINT}/${id}`)
}

export async function getAllUsersPaged(
  page?: number
): Promise<APIResponse<PaginatedUsersDto>> {
  return await fetchHandler(USERS_ENDPOINT, undefined, undefined, { page })
}

export async function createUser(
  createUserDto: CreateUserDto
): Promise<APIResponse<User>> {
  return await fetchHandler(USERS_ENDPOINT, "POST", createUserDto)
}

export async function updateUser(
  id: string,
  updateUserDto: UpdateUserDto
): Promise<APIResponse<User>> {
  if (!isValidUUID(id)) {
    return {
      ok: false,
      message: "Invalid UUID"
    }
  }

  return await fetchHandler(`${USERS_ENDPOINT}/${id}`, "PATCH", updateUserDto)
}

export async function deleteUser(id: string): Promise<APIResponse<User>> {
  if (!isValidUUID(id)) {
    return {
      ok: false,
      message: "Invalid UUID"
    }
  }

  return await fetchHandler(`${USERS_ENDPOINT}/${id}`, "DELETE")
}
