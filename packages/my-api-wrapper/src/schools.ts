import { type CreateSchoolDto, type GetAllPaginatedSchoolDto, type School, type UpdateSchoolDto } from "types-custom"
import Query from "./query"
import { fetchHandler } from "./fetchHandler"

const SCHOOLS_URL = `${process.env.NEXT_PUBLIC_API_URL ?? ""}/schools`

export async function getAllSchools(page?: number, pageSize?: number): Promise<GetAllPaginatedSchoolDto | null> {
  const query = new Query({
    page,
    pageSize
  })

  const response = await fetchHandler(`${SCHOOLS_URL}${query.toString()}`)
  if (!response) {
    return null
  }

  return (await response.json()) as GetAllPaginatedSchoolDto
}

export async function createSchool(name: string): Promise<School | null> {
  const data: CreateSchoolDto = {
    name
  }

  const response = await fetchHandler(SCHOOLS_URL, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(data)
  })

  if (!response) {
    return null
  }

  return (await response.json()) as School
}

export async function updateSchool(id: number, name: string): Promise<School | null> {
  const data: UpdateSchoolDto = {
    name
  }

  const response = await fetchHandler(`${SCHOOLS_URL}/${id}`, {
    method: "PATCH",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(data)
  })

  if (!response) {
    return null
  }

  return (await response.json()) as School
}

export async function deleteSchool(id: number): Promise<School | null> {
  const response = await fetchHandler(`${SCHOOLS_URL}/${id}`, {
    method: "DELETE",
    headers: {
      "Content-Type": "application/json",
    }
  })

  if (!response) {
    return null
  }

  return (await response.json()) as School
}