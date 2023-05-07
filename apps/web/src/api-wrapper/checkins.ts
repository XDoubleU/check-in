import { type CheckIn, type CreateCheckInDto } from "types-custom"
import { fetchHandler } from "./fetchHandler"
import { type APIResponse } from "./types"
import queryString from "query-string"
import { validate as isValidUUID } from "uuid"

const CHECKINS_ENDPOINT = "checkins"

export async function getDataForRangeChart(
  locationId: string,
  startDate: string,
  endDate: string
): Promise<APIResponse<unknown[]>> {
  if (!isValidUUID(locationId)) {
    return {
      ok: false,
      message: "Invalid UUID"
    }
  }

  return await fetchHandler(
    `${CHECKINS_ENDPOINT}/range/${locationId}`,
    undefined,
    {
      startDate,
      endDate
    }
  )
}

export async function getDataForDayChart(
  locationId: string,
  date: string
): Promise<APIResponse<unknown[]>> {
  if (!isValidUUID(locationId)) {
    return {
      ok: false,
      message: "Invalid UUID"
    }
  }

  return await fetchHandler(
    `${CHECKINS_ENDPOINT}/day/${locationId}`,
    undefined,
    {
      date
    }
  )
}

export function downloadCsvForRangeChart(
  locationId: string,
  startDate: string,
  endDate: string
): void {
  if (!isValidUUID(locationId)) {
    return
  }

  const query = queryString.stringify({
    startDate,
    endDate
  })

  open(
    `${
      process.env.NEXT_PUBLIC_API_URL ?? ""
    }/${CHECKINS_ENDPOINT}/csv/range/${locationId}?${query}`
  )
}

export function downloadCsvForDayChart(locationId: string, date: string): void {
  if (!isValidUUID(locationId)) {
    return
  }

  const query = queryString.stringify({
    date
  })

  open(
    `${
      process.env.NEXT_PUBLIC_API_URL ?? ""
    }/${CHECKINS_ENDPOINT}/csv/day/${locationId}?${query}`
  )
}

export async function createCheckIn(
  createCheckInDto: CreateCheckInDto
): Promise<APIResponse<CheckIn>> {
  return await fetchHandler(CHECKINS_ENDPOINT, {
    method: "POST",
    body: JSON.stringify(createCheckInDto)
  })
}
