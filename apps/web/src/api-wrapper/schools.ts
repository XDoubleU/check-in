import {
  type CreateSchoolDto,
  type GetAllPaginatedSchoolDto,
  type School,
  type UpdateSchoolDto
} from "types-custom"
import Query from "./query"
import { fetchHandler } from "./fetchHandler"
import { type APIResponse } from "./types"

const SCHOOLS_URL = `${process.env.NEXT_PUBLIC_API_URL ?? ""}/schools`

export async function getAllSchoolsSortedForLocation(): Promise<
  APIResponse<School[]>
> {
  return await fetchHandler(`${SCHOOLS_URL}/location`)
}

export async function getAllSchoolsPaged(
  page?: number
): Promise<APIResponse<GetAllPaginatedSchoolDto>> {
  const query = new Query({
    page
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
