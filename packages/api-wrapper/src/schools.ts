import { CreateSchoolDto, GetAllPaginatedSchoolDto, School, UpdateSchoolDto } from "types"
import Query from "./query"

const SCHOOLS_URL = `${process.env.API_URL}/schools`

export async function getAllSchools(page?: number, pageSize?: number): Promise<GetAllPaginatedSchoolDto> {
  const query = new Query({
    page,
    pageSize
  })

  const response = await fetch(SCHOOLS_URL + query)

  return (await response.json()) as GetAllPaginatedSchoolDto
}

export async function createSchool(name: string): Promise<School> {
  const data: CreateSchoolDto = {
    name
  }

  const response = await fetch(SCHOOLS_URL, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(data)
  })

  return (await response.json()) as School
}

export async function updateSchool(id: number, name: string): Promise<School> {
  const data: UpdateSchoolDto = {
    name
  }

  const response = await fetch(`${SCHOOLS_URL}/${id}`, {
    method: "PATCH",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(data)
  })

  return (await response.json()) as School
}

export async function deleteSchool(id: number): Promise<School> {
  const response = await fetch(`${SCHOOLS_URL}/${id}`, {
    method: "DELETE",
    headers: {
      "Content-Type": "application/json",
    }
  })

  return (await response.json()) as School
}