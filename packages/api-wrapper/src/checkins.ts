import { CheckIn, CreateCheckInDto } from "types"

const CHECKIN_URL = `${process.env.API_URL}/checkins`

export async function createCheckIn(locationId: string, schoolId: number): Promise<CheckIn> {
  const data: CreateCheckInDto = {
    locationId,
    schoolId
  }

  const response = await fetch(CHECKIN_URL, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(data)
  })

  return (await response.json()) as CheckIn
}