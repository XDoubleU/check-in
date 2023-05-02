import {
  type CreateLocationDto,
  type GetAllPaginatedLocationDto,
  type Location,
  type UpdateLocationDto
} from "types-custom"
import Query from "./query"
import { fetchHandler } from "./fetchHandler"
import { type APIResponse } from "./types"

const LOCATIONS_ENDPOINT = "locations"

export async function getMyLocation(): Promise<APIResponse<Location>> {
  return await fetchHandler(`${LOCATIONS_ENDPOINT}/me`)
}

export async function getAllLocations(
  page?: number
): Promise<APIResponse<GetAllPaginatedLocationDto>> {
  const query = new Query({
    page
  })

  return await fetchHandler(`${LOCATIONS_ENDPOINT}${query.toString()}`)
}

export async function getLocation(id: string): Promise<APIResponse<Location>> {
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
  return await fetchHandler(`${LOCATIONS_ENDPOINT}/${id}`, {
    method: "PATCH",
    body: JSON.stringify(updateLocationDto)
  })
}

export async function deleteLocation(
  id: string
): Promise<APIResponse<Location>> {
  return await fetchHandler(`${LOCATIONS_ENDPOINT}/${id}`, {
    method: "DELETE"
  })
}
