import { fetchHandler } from "./fetchHandler"
import { type APIResponse } from "./types"

const MIGRATIONS_ENDPOINT = "migrations"

export async function applyMigrationsUp(): Promise<APIResponse<string>> {
  return await fetchHandler(`${MIGRATIONS_ENDPOINT}/up`)
}

export async function applyMigrationsDown(): Promise<APIResponse<string>> {
  return await fetchHandler(`${MIGRATIONS_ENDPOINT}/down`)
}

export async function applySeeder(): Promise<APIResponse<void>> {
  return await fetchHandler(`${MIGRATIONS_ENDPOINT}/seed`)
}
