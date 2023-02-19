import { CreateLocationDto, GetAllPaginatedLocationDto, Location, UpdateLocationDto } from "types"
import Query from "./query"
import { fetchHandler } from "./fetchHandler"

const LOCATIONS_URL = `${process.env.NEXT_PUBLIC_API_URL}/locations`

export async function getAllLocations(page?: number): Promise<GetAllPaginatedLocationDto | null> {
  const query = new Query({
    page
  })

  const response = await fetchHandler(LOCATIONS_URL + query)
  if (!response) {
    return null
  }

  return (await response.json()) as GetAllPaginatedLocationDto
}

export async function getMyLocation(): Promise<Location | null> {
  const response = await fetchHandler(`${LOCATIONS_URL}/me`)
  if (!response) {
    return null
  }

  return (await response.json()) as Location
}

export async function getLocation(id: string): Promise<Location | null> {
  const response = await fetchHandler(`${LOCATIONS_URL}/${id}`)
  if (!response) {
    return null
  }

  return (await response.json()) as Location
}

export async function createLocation(name: string, capacity: number, username: string, password: string): Promise<Location | null> {
  const data: CreateLocationDto = {
    name,
    capacity,
    username,
    password
  }

  const response = await fetchHandler(LOCATIONS_URL, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(data)
  })

  if (!response) {
    return null
  }

  return (await response.json()) as Location
}

export async function updateLocation(id: string, name?: string, capacity?: number, username?: string, password?: string): Promise<Location | null> {
  const data: UpdateLocationDto = {
    name,
    capacity,
    username,
    password
  }

  const response = await fetchHandler(`${LOCATIONS_URL}/${id}`, {
    method: "PATCH",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(data)
  })

  if (!response) {
    return null
  }

  return (await response.json()) as Location
}

export async function deleteLocation(id: string): Promise<Location | null> {
  const response = await fetchHandler(`${LOCATIONS_URL}/${id}`, {
    method: "DELETE",
    headers: {
      "Content-Type": "application/json",
    }
  })

  if (!response) {
    return null
  }

  return (await response.json()) as Location
}