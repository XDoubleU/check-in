import { Body, ConflictException, Controller, Delete, Get, NotFoundException, Param, Patch, Post, Query } from "@nestjs/common"
import { SchoolsService } from "./schools.service"
import type { CreateSchoolDto, GetAllPaginatedSchoolDto, UpdateSchoolDto } from "types-custom"
import { Role } from "types-custom"
import { Roles } from "../auth/decorators/roles.decorator"
import { ReqUser } from "../auth/decorators/user.decorator"
import { SchoolEntity, UserEntity } from "mikro-orm-config"

type MikroGetAllPaginatedSchoolDto = Omit<GetAllPaginatedSchoolDto, "schools"> & { schools: SchoolEntity[] }

@Controller("schools")
export class SchoolsController {
  constructor(private readonly schoolsService: SchoolsService) {}

  @Roles(Role.User)
  @Get("all")
  async getAll(@ReqUser() user: UserEntity): Promise<SchoolEntity[]> {
    if (!user.location?.id) {
      return []
    }

    return await this.schoolsService.getAll(user.location.id)
  }

  @Roles(Role.Admin)
  @Get()
  async getAllPaged(@Query("page") queryPage?: string): Promise<MikroGetAllPaginatedSchoolDto> {
    const pageSize = 4
    const page = queryPage ? parseInt(queryPage) : 1
    const count = await this.schoolsService.getTotalCount()
    const schools = await this.schoolsService.getAllPaged(page, pageSize)

    return {
      page: page,
      totalPages: Math.ceil(count/pageSize),
      schools: schools
    }
  }

  @Roles(Role.Admin)
  @Post()
  async create(@Body() createSchoolDto: CreateSchoolDto): Promise<SchoolEntity> {
    const existingSchool = await this.schoolsService.getByName(createSchoolDto.name)
    if (existingSchool) {
      throw new ConflictException("School with this name already exists")
    }

    return await this.schoolsService.create(createSchoolDto.name)
  }

  @Roles(Role.Admin)
  @Patch(":id")
  async update(@Param("id") id: string, @Body() updateSchoolDto: UpdateSchoolDto): Promise<SchoolEntity> {
    const school = await this.schoolsService.getById(parseInt(id))
    if (!school || parseInt(id) === 1) {
      throw new NotFoundException("School not found")
    }

    const existingSchool = await this.schoolsService.getByName(updateSchoolDto.name)
    if (existingSchool) {
      throw new ConflictException("School with this name already exists")
    }

    return await this.schoolsService.update(school, updateSchoolDto.name)
  }

  @Roles(Role.Admin)
  @Delete(":id")
  async delete(@Param("id") id: string): Promise<SchoolEntity> {
    const school = await this.schoolsService.getById(parseInt(id))
    if (!school || parseInt(id) === 1) {
      throw new NotFoundException("School not found")
    }

    return await this.schoolsService.delete(school)
  }
}
