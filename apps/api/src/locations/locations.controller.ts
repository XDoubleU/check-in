import { BadRequestException, Body, ConflictException, Controller, Delete, Get, NotFoundException, Param, Patch, Post, Query } from "@nestjs/common"
import { LocationsService } from "./locations.service"
import { UsersService } from "../users/users.service"
import { CreateLocationDto, GetAllPaginatedLocationDto, Location, UpdateLocationDto } from "types"

@Controller("locations")
export class LocationsController {
  constructor(
    private readonly locationsService: LocationsService,
    private readonly usersService: UsersService
  ) {}

  @Get()
  async getAll(@Query("page") pageQ?: number): Promise<GetAllPaginatedLocationDto> {
    const pageSize = 3
    const page = pageQ ?? 1
    const count = await this.locationsService.getTotalCount()
    const locations = await this.locationsService.getAll(page, pageSize)

    return {
      page: page,
      totalPages: count/pageSize,
      locations: locations
    }
  }

  @Get(":id")
  async get(@Param("id") id: string): Promise<Location> {
    const location = await this.locationsService.getById(id)
    if (location === null) {
      throw new NotFoundException("Location not found")
    }

    return location
  }

  @Post()
  async create(@Body() createLocationDto: CreateLocationDto): Promise <Location> {
    const existingLocation = await this.locationsService.getByName(createLocationDto.name)
    if (existingLocation !== null) {
      throw new ConflictException("Location with this name already exists")
    }

    const existingUser = await this.usersService.getByUserName(createLocationDto.username)
    if (existingUser !== null) {
      throw new ConflictException("User with this username already exists")
    }

    const user = await this.usersService.create(createLocationDto.name, createLocationDto.password)
    if (user === null) {
      throw new BadRequestException()
    }

    const location = await this.locationsService.create(createLocationDto.name, createLocationDto.capacity, user)
    if (location === null) {
      throw new BadRequestException()
    }

    return location
  }

  @Patch(":id")
  async update(@Param("id") id: string, @Body() updateLocationDto: UpdateLocationDto): Promise<Location> {
    // TODO: check if user.username and location.name are already used
    
    let location = await this.locationsService.getById(id)
    if (location === null) {
      throw new NotFoundException("Location not found")
    }

    let user = await this.usersService.getById(location.userId)
    if (user === null) {
      throw new NotFoundException("User not found")
    }

    user = await this.usersService.update(user, updateLocationDto.username, updateLocationDto.password)
    if (user === null) {
      throw new BadRequestException()
    }

    location = await this.locationsService.update(location, updateLocationDto.name, updateLocationDto.capacity)
    if (location === null) {
      throw new BadRequestException()
    }

    return location
  }

  @Delete(":id")
  async delete(@Param("id") id: string): Promise<Location> {
    let location = await this.locationsService.getById(id)
    if (location === null) {
      throw new NotFoundException("Location not found")
    }

    let user = await this.usersService.getById(location.userId)
    if (user === null) {
      throw new NotFoundException("User not found")
    }

    location = await this.locationsService.delete(location)
    if (location === null) {
      throw new BadRequestException()
    }
    
    user = await this.usersService.delete(user)
    if (user === null) {
      throw new BadRequestException()
    }

    return location
  }
}
