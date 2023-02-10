import { School } from "@prisma/client"

export class GetAllPaginatedSchoolDto {
  page: number
  totalPages: number
  schools: School[]
}