import { type CheckIn, type CreateCheckInDto } from "types-custom"
import { fetchHandler } from "./fetchHandler"

const CHECKIN_URL = `${process.env.NEXT_PUBLIC_API_URL ?? ""}/checkins`

export async function createCheckIn(locationId: string, schoolId: number): Promise<CheckIn | null> {
  const data: CreateCheckInDto = {
    locationId,
    schoolId
  }

  const response = await fetchHandler(CHECKIN_URL, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(data)
  })
  if (!response){
    return null
  }

  return (await response.json()) as CheckIn
}