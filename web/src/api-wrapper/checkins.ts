import { fetchHandler } from "./fetchHandler"
import { type APIResponse } from "./types"
import { type School, type CheckIn, type CheckInDto } from "./types/apiTypes"

const CHECKINS_ENDPOINT = "checkins"

export async function getAllSchoolsSortedForLocation(): Promise<
  APIResponse<School[]>
> {
  return await fetchHandler(`${CHECKINS_ENDPOINT}/schools`)
}

export async function createCheckIn(
  createCheckInDto: CheckInDto
): Promise<APIResponse<CheckIn>> {
  return await fetchHandler(CHECKINS_ENDPOINT, {
    method: "POST",
    body: JSON.stringify(createCheckInDto)
  })
}
