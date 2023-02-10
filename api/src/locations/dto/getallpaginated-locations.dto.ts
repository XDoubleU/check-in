import { Location } from "@prisma/client"

export class GetAllPaginatedLocationDto {
  page: number
  totalPages: number
  locations: Location[]
}