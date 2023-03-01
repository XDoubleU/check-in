import { Body, Controller, NotFoundException, Post } from "@nestjs/common"
import { CheckInsService } from "./checkins.service"
import { LocationsService } from "../locations/locations.service"
import { SchoolsService } from "../schools/schools.service"
import { CheckIn, CreateCheckInDto, Role } from "types"
import { Roles } from "../auth/decorators/roles.decorator"

@Controller("checkins")
export class CheckInsController {
  constructor(
    private readonly checkInsService: CheckInsService,
    private readonly locationsService: LocationsService,
    private readonly schoolsService: SchoolsService
  ) {}

  @Roles(Role.User)
  @Post()
  async create(@Body() createCheckInDto: CreateCheckInDto): Promise<CheckIn> {
    const location = await this.locationsService.getById(createCheckInDto.locationId)
    if (!location) {
      throw new NotFoundException("Location not found")
    }

    const school = await this.schoolsService.getById(createCheckInDto.schoolId)
    if (!school) {
      throw new NotFoundException("School not found")
    }
    
    return await this.checkInsService.create(location, school)
  }
}
