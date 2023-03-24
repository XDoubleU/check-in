import {
  BadRequestException,
  Body,
  Controller,
  InternalServerErrorException,
  NotFoundException,
  Post
} from "@nestjs/common"
import { CheckInsService } from "./checkins.service"
import { SchoolsService } from "../schools/schools.service"
import { CreateCheckInDto, Role } from "types-custom"
import { Roles } from "../auth/decorators/roles.decorator"
import { UserEntity, type CheckInEntity } from "mikro-orm-config"
import { ReqUser } from "../auth/decorators/user.decorator"

@Controller("checkins")
export class CheckInsController {
  private readonly checkInsService: CheckInsService
  private readonly schoolsService: SchoolsService

  public constructor(
    checkInsService: CheckInsService,
    schoolsService: SchoolsService
  ) {
    this.checkInsService = checkInsService
    this.schoolsService = schoolsService
  }

  @Roles(Role.User)
  @Post()
  public async create(
    @ReqUser() user: UserEntity,
    @Body() createCheckInDto: CreateCheckInDto
  ): Promise<CheckInEntity> {
    if (!user.location) {
      throw new BadRequestException()
    }

    const school = await this.schoolsService.getById(createCheckInDto.schoolId)
    if (!school) {
      throw new NotFoundException("School not found")
    }

    const checkin = await this.checkInsService.create(user.location, school)
    if (!checkin) {
      throw new InternalServerErrorException("Could not create CheckIn")
    }

    return checkin
  }
}
