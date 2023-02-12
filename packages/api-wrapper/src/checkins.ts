import { CheckIn } from "database"
import { CreateCheckInDto } from "dtos"

const CHECKIN_URL = `${process.env.API_URL}/checkins`

export async function createCheckIn(locationId: string, schoolId: number): Promise<CheckIn> {
  const data: CreateCheckInDto = {
    locationId,
    schoolId
  }

  const response = await fetch(CHECKIN_URL, {
    method: "POST",
    body: JSON.stringify(data)
  })

  return (await response.json()) as CheckIn
}