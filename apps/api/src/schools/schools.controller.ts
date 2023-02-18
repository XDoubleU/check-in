import { BadRequestException, Body, ConflictException, Controller, Delete, Get, NotFoundException, Param, Patch, Post, Query } from "@nestjs/common"
import { SchoolsService } from "./schools.service"
import { CreateSchoolDto, GetAllPaginatedSchoolDto, School, UpdateSchoolDto } from "types"

@Controller("schools")
export class SchoolsController {
  constructor(private readonly schoolsService: SchoolsService) {}

  @Get()
  async getAll(@Query("page") queryPage?: string, @Query("pageSize") queryPageSize?: number): Promise<GetAllPaginatedSchoolDto> {
    const pageSize = queryPageSize ?? 4
    const page = queryPage ? parseInt(queryPage) : 1
    const count = await this.schoolsService.getTotalCount()
    const schools = await this.schoolsService.getAll(page, pageSize)

    return {
      page: page,
      totalPages: Math.ceil(count/pageSize),
      schools: schools
    }
  }

  @Post()
  async create(@Body() createSchoolDto: CreateSchoolDto): Promise <School> {
    const existingSchool = await this.schoolsService.getByName(createSchoolDto.name)
    if (existingSchool !== null) {
      throw new ConflictException("School with this name already exists")
    }

    const school = await this.schoolsService.create(createSchoolDto.name)
    if (school === null) {
      throw new BadRequestException()
    }

    return school
  }

  @Patch(":id")
  async update(@Param("id") id: number, @Body() updateSchoolDto: UpdateSchoolDto): Promise<School> {
    let school = await this.schoolsService.getById(id)
    if (school === null) {
      throw new NotFoundException("School not found")
    }

    school = await this.schoolsService.update(school, updateSchoolDto.name)
    if (school === null) {
      throw new BadRequestException()
    }

    return school
  }

  @Delete(":id")
  async delete(@Param("id") id: number): Promise<School> {
    let school = await this.schoolsService.getById(id)
    if (school === null) {
      throw new NotFoundException("School not found")
    }

    school = await this.schoolsService.delete(school)
    if (school === null) {
      throw new BadRequestException()
    }

    return school
  }
}
