import {
  type SchoolDto,
  type PaginatedSchoolsDto,
  type School
} from "./types/apiTypes"
import { fetchHandler } from "./fetchHandler"
import { type APIResponse } from "./types"

const SCHOOLS_ENDPOINT = "schools"

export async function getAllSchoolsPaged(
  page?: number
): Promise<APIResponse<PaginatedSchoolsDto>> {
  return await fetchHandler(`${SCHOOLS_ENDPOINT}`, undefined, { page })
}

export async function createSchool(
  createSchoolDto: SchoolDto
): Promise<APIResponse<School>> {
  return await fetchHandler(SCHOOLS_ENDPOINT, {
    method: "POST",
    body: JSON.stringify(createSchoolDto)
  })
}

export async function updateSchool(
  id: number,
  updateSchoolDto: SchoolDto
): Promise<APIResponse<School>> {
  return await fetchHandler(`${SCHOOLS_ENDPOINT}/${id}`, {
    method: "PATCH",
    body: JSON.stringify(updateSchoolDto)
  })
}

export async function deleteSchool(id: number): Promise<APIResponse<School>> {
  return await fetchHandler(`${SCHOOLS_ENDPOINT}/${id}`, {
    method: "DELETE"
  })
}
