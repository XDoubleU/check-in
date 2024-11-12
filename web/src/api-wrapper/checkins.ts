import { fetchHandler } from "./fetchHandler"
import { type APIResponse } from "./types"
import {
  type School,
  type CheckIn,
  type CreateCheckInDto
} from "./types/apiTypes"

const CHECKINS_ENDPOINT = "checkins"

export async function getAllSchoolsSortedForLocation(): Promise<
  APIResponse<School[]>
> {
  return await fetchHandler(`${CHECKINS_ENDPOINT}/schools`)
}

export async function createCheckIn(
  createCheckInDto: CreateCheckInDto
): Promise<APIResponse<CheckIn>> {
  return await fetchHandler(CHECKINS_ENDPOINT, "POST", createCheckInDto)
}
