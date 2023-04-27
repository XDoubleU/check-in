import {
  Body,
  Controller,
  Get,
  NotFoundException,
  Param,
  ParseUUIDPipe,
  Post,
  Query,
  Res
} from "@nestjs/common"
import { CheckInsService } from "./checkins.service"
import { SchoolsService } from "../schools/schools.service"
import { type CreateCheckInDto, DATE_FORMAT, Role } from "types-custom"
import { Roles } from "../auth/decorators/roles.decorator"
import { ReqUser } from "../auth/decorators/user.decorator"
import { LocationsService } from "../locations/locations.service"
import { endOfDay, format, startOfDay } from "date-fns"
import {
  convertDatetime,
  convertDayData,
  convertRangeData
} from "../helpers/dataConverters"
import { Parser } from "json2csv"
import { type Response } from "express"
import { type CheckInEntity, UserEntity } from "../entities"
import { ParseDatePipe } from "../pipes/parse-date.pipe"

@Controller("checkins")
export class CheckInsController {
  private readonly checkInsService: CheckInsService
  private readonly schoolsService: SchoolsService
  private readonly locationsService: LocationsService

  public constructor(
    checkInsService: CheckInsService,
    schoolsService: SchoolsService,
    locationsService: LocationsService
  ) {
    this.checkInsService = checkInsService
    this.schoolsService = schoolsService
    this.locationsService = locationsService
  }

  @Get("range/:locationId")
  public async getDataForRangeChart(
    @ReqUser() user: UserEntity,
    @Param("locationId", ParseUUIDPipe) locationId: string,
    @Query("startDate", ParseDatePipe) queryStartDate: Date,
    @Query("endDate", ParseDatePipe) queryEndDate: Date
  ): Promise<unknown[]> {
    const location = await this.locationsService.getLocation(locationId, user)
    if (!location) {
      throw new NotFoundException("Location not found")
    }

    const startDate = startOfDay(new Date(queryStartDate))
    const endDate = endOfDay(new Date(queryEndDate))

    return convertRangeData(
      await this.checkInsService.getAll(location, startDate, endDate),
      await this.schoolsService.getAll()
    )
  }

  @Get("csv/range/:locationId")
  public async getCsvForRangeChart(
    @Res() res: Response,
    @ReqUser() user: UserEntity,
    @Param("locationId", ParseUUIDPipe) locationId: string,
    @Query("startDate", ParseDatePipe) queryStartDate: Date,
    @Query("endDate", ParseDatePipe) queryEndDate: Date
  ): Promise<void> {
    const data = await this.getDataForRangeChart(
      user,
      locationId,
      queryStartDate,
      queryEndDate
    )

    const json = convertDatetime(data, DATE_FORMAT)
    const parser = new Parser({
      fields: Object.getOwnPropertyNames(data[0])
    })
    const csv = parser.parse(json)

    res.header("Content-Type", "text/csv")
    res.attachment(`Range-${format(new Date(), "yyyyMMddHHmmss")}.csv`)
    res.send(csv)
  }

  @Get("day/:locationId")
  public async getDataForDayChart(
    @ReqUser() user: UserEntity,
    @Param("locationId", ParseUUIDPipe) locationId: string,
    @Query("date", ParseDatePipe) queryDate: Date
  ): Promise<unknown[]> {
    const location = await this.locationsService.getLocation(locationId, user)
    if (!location) {
      throw new NotFoundException("Location not found")
    }

    const startDate = startOfDay(new Date(queryDate))
    const endDate = endOfDay(new Date(queryDate))

    return convertDayData(
      await this.checkInsService.getAll(location, startDate, endDate),
      await this.schoolsService.getAll()
    )
  }

  @Get("csv/day/:locationId")
  public async getCsvForDayChart(
    @Res() res: Response,
    @ReqUser() user: UserEntity,
    @Param("locationId", ParseUUIDPipe) locationId: string,
    @Query("date", ParseDatePipe) queryDate: Date
  ): Promise<void> {
    const data = await this.getDataForDayChart(user, locationId, queryDate)

    const json = convertDatetime(data, "yyyy-MM-dd-HH-mm")
    const parser = new Parser({
      fields: Object.getOwnPropertyNames(data[0])
    })
    const csv = parser.parse(json)

    res.header("Content-Type", "text/csv")
    res.attachment(`Day-${format(new Date(), "yyyyMMddHHmmss")}.csv`)
    res.send(csv)
  }

  @Roles(Role.User)
  @Post()
  public async create(
    @ReqUser() user: UserEntity,
    @Body() createCheckInDto: CreateCheckInDto
  ): Promise<CheckInEntity> {
    const school = await this.schoolsService.getById(createCheckInDto.schoolId)
    if (!school) {
      throw new NotFoundException("School not found")
    }

    // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
    return await this.checkInsService.create(user.location!, school)
  }
}
