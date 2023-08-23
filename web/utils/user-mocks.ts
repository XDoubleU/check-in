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
      id: "userId",
      username: "default",
      role: "default",
      location: {
        id: "locationId",
        name: "location",
        normalizedName: "location",
        available: 2,
        capacity: 10,
        timeZone: "Europe/Brussels",
        userId: "userId",
        yesterdayFullAt: ""
      }
    }
  })
}

export async function managerUserMock(): Promise<APIResponse<User>> {
  return Promise.resolve({
    ok: true,
    data: {
      id: "userId",
      username: "manager",
      role: "manager"
    }
  })
}

export async function adminUserMock(): Promise<APIResponse<User>> {
  return Promise.resolve({
    ok: true,
    data: {
      id: "userId",
      username: "admin",
      role: "admin"
    }
  })
}
