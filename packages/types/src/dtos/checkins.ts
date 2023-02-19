import { IsNotEmpty } from "class-validator"

export class CreateCheckInDto {
  @IsNotEmpty()
  locationId: string
  
  @IsNotEmpty()
  schoolId: number
}