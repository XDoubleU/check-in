import { type User, type SignInDto } from "types-custom"
import { fetchHandlerNoRefresh } from "./fetchHandler"
import type APIResponse from "./types/apiResponse"

const AUTH_URL = `${process.env.NEXT_PUBLIC_API_URL ?? ""}/auth`

export async function signin(signInDto: SignInDto): Promise<APIResponse<User>> {
  return await fetchHandlerNoRefresh(`${AUTH_URL}/signin`, {
    method: "POST",
    body: JSON.stringify(signInDto)
  })
}

export async function signOut(): Promise<void> {
  await fetchHandlerNoRefresh(`${AUTH_URL}/signout`)
}
