import { fetchHandlerNoRefresh } from "./fetchHandler"
import { type APIResponse } from "./types"
import { type SignInDto, type User } from "./types/apiTypes"

const AUTH_ENDPOINT = "auth"

export async function signIn(signInDto: SignInDto): Promise<APIResponse<User>> {
  return await fetchHandlerNoRefresh(`${AUTH_ENDPOINT}/signin`, {
    method: "POST",
    body: JSON.stringify(signInDto)
  })
}

export async function signOut(): Promise<void> {
  await fetchHandlerNoRefresh(`${AUTH_ENDPOINT}/signout`)
}
