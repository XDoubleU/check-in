import { fetchHandler } from "./fetchHandler"
import { type APIResponse } from "./types"
import { validate as isValidUUID } from "uuid"
import {
  type CheckIn,
  type CheckInsGraphDto,
  type CreateLocationDto,
  type Location,
  type PaginatedLocationsDto,
  type UpdateLocationDto
} from "./types/apiTypes"
import queryString from "query-string"
import { DATE_FORMAT } from "api-wrapper/types/apiTypes"
import moment, { type Moment } from "moment"

const LOCATIONS_ENDPOINT = "locations"

const INVALID_UUID = "Invalid UUID"

function areValidUUIDs(strings: string[]): boolean {
  for (const str of strings) {
    if (!isValidUUID(str)) {
      return false
    }
  }

  return true
}

export async function getDataForRangeChart(
  locationIds: string[],
  startDate: Moment,
  endDate: Moment
): Promise<APIResponse<CheckInsGraphDto>> {
  if (!areValidUUIDs(locationIds)) {
    return {
      ok: false,
      message: INVALID_UUID
    }
  }

  return await fetchHandler(
    `all-locations/checkins/range`,
    undefined,
    undefined,
    {
      ids: locationIds,
      startDate: moment(startDate).format(DATE_FORMAT),
      endDate: moment(endDate).format(DATE_FORMAT),
      returnType: "raw"
    }
  )
}

export async function getDataForDayChart(
  locationIds: string[],
  date: Moment
): Promise<APIResponse<CheckInsGraphDto>> {
  if (!areValidUUIDs(locationIds)) {
    return {
      ok: false,
      message: INVALID_UUID
    }
  }

  return await fetchHandler(
    `all-locations/checkins/day`,
    undefined,
    undefined,
    {
      ids: locationIds,
      date: moment(date).format(DATE_FORMAT),
      returnType: "raw"
    }
  )
}

export function downloadCSVForRangeChart(
  locationIds: string[],
  startDate: Moment,
  endDate: Moment
): void {
  if (!areValidUUIDs(locationIds)) {
    return
  }

  const query = queryString.stringify(
    {
      ids: locationIds,
      startDate: moment(startDate).format(DATE_FORMAT),
      endDate: moment(endDate).format(DATE_FORMAT),
      returnType: "csv"
    },
    { arrayFormat: "comma" }
  )

  open(
    `${
      process.env.NEXT_PUBLIC_API_URL ?? ""
    }/all-locations/checkins/range?${query}`
  )
}

export function downloadCSVForDayChart(
  locationIds: string[],
  date: Moment
): void {
  if (!areValidUUIDs(locationIds)) {
    return
  }

  const query = queryString.stringify(
    {
      ids: locationIds,
      date: moment(date).format(DATE_FORMAT),
      returnType: "csv"
    },
    { arrayFormat: "comma" }
  )

  open(
    `${
      process.env.NEXT_PUBLIC_API_URL ?? ""
    }/all-locations/checkins/day?${query}`
  )
}

export async function getCheckInsToday(
  locationId: string
): Promise<APIResponse<CheckIn[]>> {
  return await fetchHandler(`${LOCATIONS_ENDPOINT}/${locationId}/checkins`)
}

export async function deleteCheckIn(
  locationId: string,
  checkInId: number
): Promise<APIResponse<CheckIn>> {
  return await fetchHandler(
    `${LOCATIONS_ENDPOINT}/${locationId}/checkins/${checkInId.toString()}`,
    "DELETE"
  )
}

export async function getAllLocations(): Promise<APIResponse<Location[]>> {
  return await fetchHandler("all-locations")
}

export async function getAllLocationsPaged(
  page?: number
): Promise<APIResponse<PaginatedLocationsDto>> {
  return await fetchHandler(LOCATIONS_ENDPOINT, undefined, undefined, {
    page
  })
}

export async function getLocation(id: string): Promise<APIResponse<Location>> {
  if (!isValidUUID(id)) {
    return {
      ok: false,
      message: INVALID_UUID
    }
  }

  return await fetchHandler(`${LOCATIONS_ENDPOINT}/${id}`)
}

export async function createLocation(
  createLocationDto: CreateLocationDto
): Promise<APIResponse<Location>> {
  return await fetchHandler(LOCATIONS_ENDPOINT, "POST", createLocationDto)
}

export async function updateLocation(
  id: string,
  updateLocationDto: UpdateLocationDto
): Promise<APIResponse<Location>> {
  if (!isValidUUID(id)) {
    return {
      ok: false,
      message: INVALID_UUID
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
      message: INVALID_UUID
    }
  }

  return await fetchHandler(`${LOCATIONS_ENDPOINT}/${id}`, "DELETE")
}
