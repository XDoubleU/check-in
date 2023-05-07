import {
  type CreateSchoolDto,
  type GetAllPaginatedSchoolDto,
  type School,
  type UpdateSchoolDto
} from "types-custom"
import { fetchHandler } from "./fetchHandler"
import { type APIResponse } from "./types"

const SCHOOLS_ENDPOINT = "schools"

export async function getAllSchoolsSortedForLocation(): Promise<
  APIResponse<School[]>
> {
  return await fetchHandler(`${SCHOOLS_ENDPOINT}/location`)
}

export async function getAllSchoolsPaged(
  page?: number
): Promise<APIResponse<GetAllPaginatedSchoolDto>> {
  return await fetchHandler(`${SCHOOLS_ENDPOINT}`, undefined, { page })
}

export async function createSchool(
  createSchoolDto: CreateSchoolDto
): Promise<APIResponse<School>> {
  return await fetchHandler(SCHOOLS_ENDPOINT, {
    method: "POST",
    body: JSON.stringify(createSchoolDto)
  })
}

export async function updateSchool(
  id: number,
  updateSchoolDto: UpdateSchoolDto
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
