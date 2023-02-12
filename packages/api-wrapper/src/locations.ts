import { CreateLocationDto, GetAllPaginatedLocationDto, UpdateLocationDto } from "types"
import Query from "./query"

const LOCATIONS_URL = `${process.env.API_URL}/locations`

export async function getAllLocations(page?: number): Promise<GetAllPaginatedLocationDto> {
  const query = new Query({
    page
  })

  const response = await fetch(LOCATIONS_URL + query)

  return (await response.json()) as GetAllPaginatedLocationDto
}

export async function getLocation(id: string): Promise<Location> {
  const response = await fetch(`${LOCATIONS_URL}/${id}`)
  return (await response.json()) as Location
}

export async function createLocation(name: string, capacity: number, username: string, password: string): Promise<Location> {
  const data: CreateLocationDto = {
    name,
    capacity,
    username,
    password
  }

  const response = await fetch(LOCATIONS_URL, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(data)
  })

  return (await response.json()) as Location
}

export async function updateLocation(id: string, name?: string, capacity?: number, username?: string, password?: string): Promise<Location> {
  const data: UpdateLocationDto = {
    name,
    capacity,
    username,
    password
  }

  const response = await fetch(`${LOCATIONS_URL}/${id}`, {
    method: "PATCH",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(data)
  })

  return (await response.json()) as Location
}

export async function deleteLocation(id: string): Promise<Location> {
  const response = await fetch(`${LOCATIONS_URL}/${id}`, {
    method: "DELETE",
    headers: {
      "Content-Type": "application/json",
    }
  })

  return (await response.json()) as Location
}