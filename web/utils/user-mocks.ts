import { type APIResponse } from "api-wrapper"
import { type User } from "api-wrapper/types/apiTypes"

export async function noUserMock(): Promise<APIResponse<User>> {
  return Promise.resolve({
    ok: false
  })
}

export async function defaultUserMock(): Promise<APIResponse<User>> {
  return Promise.resolve({
    ok: true,
    data: {
      id: "id",
      username: "default",
      role: "default"
    }
  })
}

export async function managerUserMock(): Promise<APIResponse<User>> {
  return Promise.resolve({
    ok: true,
    data: {
      id: "id",
      username: "manager",
      role: "manager"
    }
  })
}

export async function adminUserMock(): Promise<APIResponse<User>> {
  return Promise.resolve({
    ok: true,
    data: {
      id: "id",
      username: "admin",
      role: "admin"
    }
  })
}