/* eslint-disable sonarjs/no-duplicate-string */
import { fetchHandler } from "./fetchHandler"
import { type APIResponse } from "./types"
import { validate as isValidUUID } from "uuid"
import {
  type CheckInsLocationEntryRawMap,
  type CreateLocationDto,
  type Location,
  type PaginatedLocationsDto,
  type UpdateLocationDto
} from "./types/apiTypes"
import queryString from "query-string"
import { API_DATE_FORMAT } from "api-wrapper/types/apiTypes"
import { format } from "date-fns"

const LOCATIONS_ENDPOINT = "locations"

export async function getDataForRangeChart(
  locationId: string,
  startDate: Date,
  endDate: Date
): Promise<APIResponse<CheckInsLocationEntryRawMap>> {
  if (!isValidUUID(locationId)) {
    return {
      ok: false,
      message: "Invalid UUID"
    }
  }

  return await fetchHandler(
    `${LOCATIONS_ENDPOINT}/${locationId}/checkins/range`,
    undefined,
    undefined,
    {
      startDate: format(startDate, API_DATE_FORMAT),
      endDate: format(endDate, API_DATE_FORMAT),
      returnType: "raw"
    }
  )
}

export async function getDataForDayChart(
  locationId: string,
  date: Date
): Promise<APIResponse<CheckInsLocationEntryRawMap>> {
  if (!isValidUUID(locationId)) {
    return {
      ok: false,
      message: "Invalid UUID"
    }
  }

  return await fetchHandler(
    `${LOCATIONS_ENDPOINT}/${locationId}/checkins/day`,
    undefined,
    undefined,
    {
      date: format(date, API_DATE_FORMAT),
      returnType: "raw"
    }
  )
}

export function downloadCSVForRangeChart(
  locationId: string,
  startDate: Date,
  endDate: Date
): void {
  if (!isValidUUID(locationId)) {
    return
  }

  const query = queryString.stringify({
    startDate: format(startDate, API_DATE_FORMAT),
    endDate: format(endDate, API_DATE_FORMAT),
    returnType: "csv"
  })

  open(
    `${
      process.env.NEXT_PUBLIC_API_URL ?? ""
    }/${LOCATIONS_ENDPOINT}/${locationId}/checkins/range?${query}`
  )
}

export function downloadCSVForDayChart(locationId: string, date: Date): void {
  if (!isValidUUID(locationId)) {
    return
  }

  const query = queryString.stringify({
    date: format(date, API_DATE_FORMAT),
    returnType: "csv"
  })

  open(
    `${
      process.env.NEXT_PUBLIC_API_URL ?? ""
    }/${LOCATIONS_ENDPOINT}/${locationId}/checkins/day?${query}`
  )
}

export async function getAllLocations(
  page?: number
): Promise<APIResponse<PaginatedLocationsDto>> {
  return await fetchHandler(`${LOCATIONS_ENDPOINT}`, undefined, undefined, {
    page
  })
}

export async function getLocation(id: string): Promise<APIResponse<Location>> {
  if (!isValidUUID(id)) {
    return {
      ok: false,
      message: "Invalid UUID"
    }
  }

  return await fetchHandler(`${LOCATIONS_ENDPOINT}/${id}`)
}

export async function createLocation(
  createLocationDto: CreateLocationDto
): Promise<APIResponse<Location>> {
  return await fetchHandler(`${LOCATIONS_ENDPOINT}`, "POST", createLocationDto)
}

export async function updateLocation(
  id: string,
  updateLocationDto: UpdateLocationDto
): Promise<APIResponse<Location>> {
  if (!isValidUUID(id)) {
    return {
      ok: false,
      message: "Invalid UUID"
    }
  }

  return await fetchHandler(
    `${LOCATIONS_ENDPOINT}/${id}`,
    "PATCH",
    updateLocationDto
  )
}

export async function deleteLocation(
  id: string
): Promise<APIResponse<Location>> {
  if (!isValidUUID(id)) {
    return {
      ok: false,
      message: "Invalid UUID"
    }
  }

  return await fetchHandler(`${LOCATIONS_ENDPOINT}/${id}`, "DELETE")
}
