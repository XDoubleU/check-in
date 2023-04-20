import { type CheckIn, type CreateCheckInDto } from "types-custom"
import { fetchHandler } from "./fetchHandler"
import Query from "./query"
import { type APIResponse } from "./types"

const CHECKIN_URL = `${process.env.NEXT_PUBLIC_API_URL ?? ""}/checkins`

export async function getDataForRangeChart(
  locationId: string,
  startDate: string,
  endDate: string
): Promise<APIResponse<unknown[]>> {
  const query = new Query({
    startDate,
    endDate
  })

  return await fetchHandler(
    `${CHECKIN_URL}/range/${locationId}${query.toString()}`
  )
}

export async function getDataForDayChart(
  locationId: string,
  date: string
): Promise<APIResponse<unknown[]>> {
  const query = new Query({
    date
  })

  return await fetchHandler(
    `${CHECKIN_URL}/day/${locationId}${query.toString()}`
  )
}

export function downloadCsvForRangeChart(
  locationId: string,
  startDate: string,
  endDate: string
): void {
  const query = new Query({
    startDate,
    endDate
  })

  window.open(`${CHECKIN_URL}/csv/range/${locationId}${query.toString()}`)
}

export function downloadCsvForDayChart(locationId: string, date: string): void {
  const query = new Query({
    date
  })

  window.open(`${CHECKIN_URL}/csv/day/${locationId}${query.toString()}`)
}

export async function createCheckIn(
  createCheckInDto: CreateCheckInDto
): Promise<APIResponse<CheckIn>> {
  return await fetchHandler(CHECKIN_URL, {
    method: "POST",
    body: JSON.stringify(createCheckInDto)
  })
}
