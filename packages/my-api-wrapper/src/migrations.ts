import { fetchHandler } from "./fetchHandler"
import type APIResponse from "./types/apiResponse"

const MIGRATIONS_URL = `${process.env.NEXT_PUBLIC_API_URL ?? ""}/migrations`

export async function applyMigrationsUp(): Promise<APIResponse<string>> {
  return await fetchHandler(`${MIGRATIONS_URL}/up`)
}

export async function applyMigrationsDown(): Promise<APIResponse<string>> {
  return await fetchHandler(`${MIGRATIONS_URL}/down`)
}

export async function applySeeder(): Promise<APIResponse<void>> {
  return await fetchHandler(`${MIGRATIONS_URL}/seed`)
}
