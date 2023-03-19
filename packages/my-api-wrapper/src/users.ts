import { type User } from "types-custom"
import { fetchHandler } from "./fetchHandler"
import type APIResponse from "./types/apiResponse"

const USERS_URL = `${process.env.NEXT_PUBLIC_API_URL ?? ""}/users`

export async function getMyUser(): Promise<APIResponse<User>> {
  return await fetchHandler(`${USERS_URL}/me`)
}
