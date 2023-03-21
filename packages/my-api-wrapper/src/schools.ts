import {
  type CreateSchoolDto,
  type GetAllPaginatedSchoolDto,
  type School,
  type UpdateSchoolDto
} from "types-custom"
import Query from "./query"
import { fetchHandler } from "./fetchHandler"
import type APIResponse from "./types/apiResponse"

const SCHOOLS_URL = `${process.env.NEXT_PUBLIC_API_URL ?? ""}/schools`

export async function getAllSchoolsCheckIn(): Promise<APIResponse<School[]>> {
  return await fetchHandler(`${SCHOOLS_URL}/all`)
}

export async function getAllSchools(
  page?: number,
  pageSize?: number
): Promise<APIResponse<GetAllPaginatedSchoolDto>> {
  const query = new Query({
    page,
    pageSize
  })

  return await fetchHandler(`${SCHOOLS_URL}${query.toString()}`)
}

export async function createSchool(
  createSchoolDto: CreateSchoolDto
): Promise<APIResponse<School>> {
  return await fetchHandler(SCHOOLS_URL, {
    method: "POST",
    body: JSON.stringify(createSchoolDto)
  })
}

export async function updateSchool(
  id: number,
  updateSchoolDto: UpdateSchoolDto
): Promise<APIResponse<School>> {
  return await fetchHandler(`${SCHOOLS_URL}/${id}`, {
    method: "PATCH",
    body: JSON.stringify(updateSchoolDto)
  })
}

export async function deleteSchool(id: number): Promise<APIResponse<School>> {
  return await fetchHandler(`${SCHOOLS_URL}/${id}`, {
    method: "DELETE"
  })
}
