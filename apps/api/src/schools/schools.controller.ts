import { Body, ConflictException, Controller, Delete, Get, NotFoundException, Param, ParseIntPipe, Patch, Post, Query } from "@nestjs/common"
import { SchoolsService } from "./schools.service"
import { CreateSchoolDto, GetAllPaginatedSchoolDto, Role, School, UpdateSchoolDto, User } from "types"
import { Roles } from "../auth/decorators/roles.decorator"
import { ReqUser } from "../auth/decorators/user.decorator"

@Controller("schools")
export class SchoolsController {
  constructor(private readonly schoolsService: SchoolsService) {}

  @Roles(Role.User)
  @Get("all")
  async getAll(@ReqUser() user: User): Promise<School[]> {
    // TODO: test role access
    return await this.schoolsService.getAll(user.locationId as string)
  }

  @Roles(Role.Admin)
  @Get()
  async getAllPaged(@Query("page") queryPage?: string): Promise<GetAllPaginatedSchoolDto> {
    // TODO: test role access
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
  async create(@Body() createSchoolDto: CreateSchoolDto): Promise <School> {
    // TODO: test role access
    const existingSchool = await this.schoolsService.getByName(createSchoolDto.name)
    if (existingSchool) {
      throw new ConflictException("School with this name already exists")
    }

    return await this.schoolsService.create(createSchoolDto.name)
  }

  @Roles(Role.Admin)
  @Patch(":id")
  async update(@Param("id", ParseIntPipe) id: number, @Body() updateSchoolDto: UpdateSchoolDto): Promise<School> {
    // TODO: test role access
    // TODO: except id === 1
    const school = await this.schoolsService.getById(id)
    if (!school) {
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
  async delete(@Param("id", ParseIntPipe) id: number): Promise<School> {
    // TODO: test role access
    // TODO: except id === 1
    const school = await this.schoolsService.getById(id)
    if (!school) {
      throw new NotFoundException("School not found")
    }

    return await this.schoolsService.delete(school)
  }
}
