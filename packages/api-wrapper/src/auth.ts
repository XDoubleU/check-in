import { SignInDto } from "types"
import { fetchHandler } from "./fetchHandler"

const AUTH_URL = `${process.env.NEXT_PUBLIC_API_URL}/auth`

export async function signin(username: string, password: string): Promise<string | null> {
  const data: SignInDto = {
    username,
    password
  }

  const response = await fetchHandler(`${AUTH_URL}/signin`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(data)
  })

  if (!response) {
    return "Invalid credentials"
  }

  return null
}

export async function signOut(): Promise<void> {
  await fetchHandler(`${AUTH_URL}/signout`)
}