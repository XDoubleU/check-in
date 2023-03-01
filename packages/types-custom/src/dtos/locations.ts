import { IsNotEmpty, IsOptional } from "class-validator"
import { Location } from "../types"

export class CreateLocationDto {
  @IsNotEmpty()
  name: string

  @IsNotEmpty()
  capacity: number

  @IsNotEmpty()
  username: string
  
  @IsNotEmpty()
  password: string
}

export class GetAllPaginatedLocationDto {
  @IsNotEmpty()
  page: number

  @IsNotEmpty()
  totalPages: number

  @IsNotEmpty()
  locations: Location[]
}

export class UpdateLocationDto {
  @IsOptional()
  name?: string

  @IsOptional()
  capacity?: number
  
  @IsOptional()
  username?: string
  
  @IsOptional()
  password?: string
}