import { IsNotEmpty } from "class-validator"
import { School } from "../types"

export class CreateSchoolDto {
  @IsNotEmpty()
  name: string
}

export class GetAllPaginatedSchoolDto {
  @IsNotEmpty()
  page: number

  @IsNotEmpty()
  totalPages: number

  @IsNotEmpty()
  schools: School[]
}

export class UpdateSchoolDto {
  @IsNotEmpty()
  name: string
}