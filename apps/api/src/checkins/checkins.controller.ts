import { BadRequestException, Body, Controller, NotFoundException, Post } from "@nestjs/common"
import { CheckIn } from "database"
import { CheckInsService } from "./checkins.service"
import { LocationsService } from "../locations/locations.service"
import { SchoolsService } from "../schools/schools.service"
import { CreateCheckInDto } from "dtos"

@Controller("checkins")
export class CheckInsController {
  constructor(
    private readonly checkInsService: CheckInsService,
    private readonly locationsService: LocationsService,
    private readonly schoolsService: SchoolsService
  ) {}
  // TODO: get available places using websockets

  @Post()
  async create(@Body() createCheckInDto: CreateCheckInDto): Promise<CheckIn> {
    const location = await this.locationsService.getById(createCheckInDto.locationId)
    if (location === null) {
      throw new NotFoundException("Location not found")
    }

    const school = await this.schoolsService.getById(createCheckInDto.schoolId)
    if (school === null) {
      throw new NotFoundException("School not found")
    }

    const checkIn = await this.checkInsService.create(location, school)
    if (checkIn === null) {
      throw new BadRequestException()
    }
    
    return checkIn
  }
}
