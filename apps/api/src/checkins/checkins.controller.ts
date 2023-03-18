import {
  Body,
  Controller,
  InternalServerErrorException,
  NotFoundException,
  Post
} from "@nestjs/common"
import { CheckInsService } from "./checkins.service"
import { LocationsService } from "../locations/locations.service"
import { SchoolsService } from "../schools/schools.service"
import { CreateCheckInDto, Role } from "types-custom"
import { Roles } from "../auth/decorators/roles.decorator"
import { type CheckInEntity } from "mikro-orm-config"

@Controller("checkins")
export class CheckInsController {
  private readonly checkInsService: CheckInsService
  private readonly locationsService: LocationsService
  private readonly schoolsService: SchoolsService

  public constructor(
    checkInsService: CheckInsService,
    locationsService: LocationsService,
    schoolsService: SchoolsService
  ) {
    this.checkInsService = checkInsService
    this.locationsService = locationsService
    this.schoolsService = schoolsService
  }

  @Roles(Role.User)
  @Post()
  public async create(
    @Body() createCheckInDto: CreateCheckInDto
  ): Promise<CheckInEntity> {
    const location = await this.locationsService.getById(
      createCheckInDto.locationId
    )
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
