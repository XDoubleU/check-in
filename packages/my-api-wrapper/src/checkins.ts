import { type CheckIn, type CreateCheckInDto } from "types-custom"
import { fetchHandler } from "./fetchHandler"
import type APIResponse from "./types/apiResponse"

const CHECKIN_URL = `${process.env.NEXT_PUBLIC_API_URL ?? ""}/checkins`

export async function createCheckIn(
  createCheckInDto: CreateCheckInDto
): Promise<APIResponse<CheckIn>> {
  return await fetchHandler(CHECKIN_URL, {
    method: "POST",
    body: JSON.stringify(createCheckInDto)
  })
}
