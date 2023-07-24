/* eslint-disable sonarjs/no-duplicate-string */
import { fetchHandler } from "./fetchHandler"
import { type APIResponse } from "./types"
import { validate as isValidUUID } from "uuid"
import {
  type CreateLocationDto,
  type Location,
  type PaginatedLocationsDto,
  type UpdateLocationDto
} from "./types/apiTypes"

const LOCATIONS_ENDPOINT = "locations"

export async function getMyLocation(): Promise<APIResponse<Location>> {
  return await fetchHandler(`${LOCATIONS_ENDPOINT}/me`)
}

export async function getAllLocations(
  page?: number
): Promise<APIResponse<PaginatedLocationsDto>> {
  return await fetchHandler(`${LOCATIONS_ENDPOINT}`, undefined, { page })
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
  return await fetchHandler(LOCATIONS_ENDPOINT, {
    method: "POST",
    body: JSON.stringify(createLocationDto)
  })
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

  return await fetchHandler(`${LOCATIONS_ENDPOINT}/${id}`, {
    method: "PATCH",
    body: JSON.stringify(updateLocationDto)
  })
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

  return await fetchHandler(`${LOCATIONS_ENDPOINT}/${id}`, {
    method: "DELETE"
  })
}
