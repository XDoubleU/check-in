import { Body, Controller, InternalServerErrorException, NotFoundException, Post } from "@nestjs/common"
import { CheckInsService } from "./checkins.service"
import { LocationsService } from "../locations/locations.service"
import { SchoolsService } from "../schools/schools.service"
import type { CreateCheckInDto } from "types-custom"
import { Role } from "types-custom"
import { Roles } from "../auth/decorators/roles.decorator"
import { CheckInEntity } from "mikro-orm-config"

@Controller("checkins")
export class CheckInsController {
  constructor(
    private readonly checkInsService: CheckInsService,
    private readonly locationsService: LocationsService,
    private readonly schoolsService: SchoolsService
  ) {}

  @Roles(Role.User)
  @Post()
  async create(@Body() createCheckInDto: CreateCheckInDto): Promise<CheckInEntity> {
    const location = await this.locationsService.getById(createCheckInDto.locationId)
    if (!location) {
      throw new NotFoundException("Location not found")
    }

    const school = await this.schoolsService.getById(createCheckInDto.schoolId)
    if (!school) {
      throw new NotFoundException("School not found")
    }
    
    const checkin = await this.checkInsService.create(location, school)
    if (!checkin) {
      throw new InternalServerErrorException("Could not create CheckIn")
    }

    return checkin
  }
}
