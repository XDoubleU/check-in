import {
  type CreateUserDto,
  type GetAllPaginatedUserDto,
  type UpdateUserDto,
  type User
} from "types-custom"
import { fetchHandler } from "./fetchHandler"
import Query from "./query"
import type APIResponse from "./types/apiResponse"

const USERS_URL = `${process.env.NEXT_PUBLIC_API_URL ?? ""}/users`

export async function getMyUser(): Promise<APIResponse<User>> {
  return await fetchHandler(`${USERS_URL}/me`)
}

export async function getUser(id: string): Promise<APIResponse<User>> {
  return await fetchHandler(`${USERS_URL}/${id}`)
}

export async function getAllUsersPaged(
  page?: number
): Promise<APIResponse<GetAllPaginatedUserDto>> {
  const query = new Query({
    page
  })

  return await fetchHandler(`${USERS_URL}${query.toString()}`)
}

export async function createUser(
  createUserDto: CreateUserDto
): Promise<APIResponse<User>> {
  return await fetchHandler(USERS_URL, {
    method: "POST",
    body: JSON.stringify(createUserDto)
  })
}

export async function updateUser(
  id: string,
  updateUserDto: UpdateUserDto
): Promise<APIResponse<User>> {
  return await fetchHandler(`${USERS_URL}/${id}`, {
    method: "PATCH",
    body: JSON.stringify(updateUserDto)
  })
}

export async function deleteUser(id: string): Promise<APIResponse<User>> {
  return await fetchHandler(`${USERS_URL}/${id}`, {
    method: "DELETE"
  })
}
