import {
  Body,
  ConflictException,
  Controller,
  Delete,
  Get,
  NotFoundException,
  Param,
  Patch,
  Post,
  Query
} from "@nestjs/common"
import { SchoolsService } from "./schools.service"
import {
  type GetAllPaginatedSchoolDto,
  Role,
  type CreateSchoolDto,
  type UpdateSchoolDto
} from "types-custom"
import { Roles } from "../auth/decorators/roles.decorator"
import { ReqUser } from "../auth/decorators/user.decorator"
import { type SchoolEntity, UserEntity } from "../entities"

type MikroGetAllPaginatedSchoolDto = Omit<GetAllPaginatedSchoolDto, "data"> & {
  data: SchoolEntity[]
}

@Controller("schools")
export class SchoolsController {
  private readonly schoolsService: SchoolsService

  public constructor(schoolsService: SchoolsService) {
    this.schoolsService = schoolsService
  }

  @Roles(Role.User)
  @Get("location")
  public async getAllForLocation(
    @ReqUser() user: UserEntity
  ): Promise<SchoolEntity[]> {
    // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
    return await this.schoolsService.getAllForLocation(user.location!.id)
  }

  @Roles(Role.Manager)
  @Get()
  public async getAllPaged(
    @Query("page") queryPage?: string
  ): Promise<MikroGetAllPaginatedSchoolDto> {
    const pageSize = 4
    const current = queryPage ? parseInt(queryPage) : 1
    const amountOfSchools = await this.schoolsService.getTotalCount()
    const schools = await this.schoolsService.getAllPaged(current, pageSize)

    return {
      data: schools,
      pagination: {
        current,
        total: Math.ceil(amountOfSchools / pageSize)
      }
    }
  }

  @Roles(Role.Manager)
  @Post()
  public async create(
    @Body() createSchoolDto: CreateSchoolDto
  ): Promise<SchoolEntity> {
    const existingSchool = await this.schoolsService.getByName(
      createSchoolDto.name
    )
    if (existingSchool) {
      throw new ConflictException("School with this name already exists")
    }

    return await this.schoolsService.create(createSchoolDto.name)
  }

  @Roles(Role.Manager)
  @Patch(":id")
  public async update(
    @Param("id") id: string,
    @Body() updateSchoolDto: UpdateSchoolDto
  ): Promise<SchoolEntity> {
    const school = await this.schoolsService.getById(parseInt(id))
    if (!school || parseInt(id) === 1) {
      throw new NotFoundException("School not found")
    }

    const existingSchool = await this.schoolsService.getByName(
      updateSchoolDto.name
    )
    if (existingSchool) {
      throw new ConflictException("School with this name already exists")
    }

    return await this.schoolsService.update(school, updateSchoolDto.name)
  }

  @Roles(Role.Manager)
  @Delete(":id")
  public async delete(@Param("id") id: string): Promise<SchoolEntity> {
    const school = await this.schoolsService.getById(parseInt(id))
    if (!school || parseInt(id) === 1) {
      throw new NotFoundException("School not found")
    }

    return await this.schoolsService.delete(school)
  }
}
