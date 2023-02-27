import { BadRequestException, Body, ConflictException, Controller, Delete, Get, InternalServerErrorException, NotFoundException, Param, Patch, Post, Query } from "@nestjs/common"
import { LocationsService } from "./locations.service"
import { UsersService } from "../users/users.service"
import { CreateLocationDto, GetAllPaginatedLocationDto, Location, Role, UpdateLocationDto, User } from "types"
import { ReqUser } from "../auth/decorators/user.decorator"
import { Roles } from "../auth/decorators/roles.decorator"

@Controller("locations")
export class LocationsController {
  constructor(
    private readonly locationsService: LocationsService,
    private readonly usersService: UsersService
  ) {}

  @Roles(Role.Admin)
  @Get()
  async getAll(@Query("page") queryPage?: string): Promise<GetAllPaginatedLocationDto> {
    const pageSize = 3
    const page = queryPage ? parseInt(queryPage) : 1
    const count = await this.locationsService.getTotalCount()
    const locations = await this.locationsService.getAllPaged(page, pageSize)

    return {
      page: page,
      totalPages: Math.ceil(count/pageSize),
      locations: locations
    }
  }

  @Roles(Role.User)
  @Get("me")
  async getMyLocation(@ReqUser() user: User): Promise<Location> {
    if (!user.locationId) {
      throw new BadRequestException()
    }

    const location = await this.locationsService.getById(user.locationId)
    if (!location) {
      throw new NotFoundException("Location not found")
    }

    return location
  }

  @Get(":id")
  async get(@ReqUser() user: User, @Param("id") id: string): Promise<Location> {
    const location = await this.locationsService.getById(id)
    if (!location || (!user.roles.includes(Role.Admin) && location.user.id !== user.id)) {
      throw new NotFoundException("Location not found")
    }

    return location
  }

  @Roles(Role.Admin)
  @Post()
  async create(@Body() createLocationDto: CreateLocationDto): Promise <Location> {
    const existingLocation = await this.locationsService.getByName(createLocationDto.name)
    if (existingLocation) {
      throw new ConflictException("Location with this name already exists")
    }

    const existingUser = await this.usersService.getByUserName(createLocationDto.username)
    if (existingUser) {
      throw new ConflictException("User with this username already exists")
    }

    const user = await this.usersService.create(createLocationDto.username, createLocationDto.password)
    return await this.locationsService.create(createLocationDto.name, createLocationDto.capacity, user)
  }

  @Patch(":id")
  async update(@ReqUser() reqUser: User, @Param("id") id: string, @Body() updateLocationDto: UpdateLocationDto): Promise<Location> {
    const location = await this.locationsService.getById(id)
    if (!location || (!reqUser.roles.includes(Role.Admin) && location.user.id !== reqUser.id)) {
      throw new NotFoundException("Location not found")
    }
    
    if (updateLocationDto.name) {
      const existingLocation = await this.locationsService.getByName(updateLocationDto.name)
      if (existingLocation) {
        throw new ConflictException("Location with this name already exists")
      }
    }

    if (updateLocationDto.username) {
      const existingUser = await this.usersService.getByUserName(updateLocationDto.username)
      if (existingUser) {
        throw new ConflictException("User with this username already exists")
      }
    }

    const user = await this.usersService.getById(location.user.id)
    if (!user) {
      throw new InternalServerErrorException("User from location couldn't be fetched")
    }

    await this.usersService.update(user, updateLocationDto.username, updateLocationDto.password)

    return await this.locationsService.update(location, updateLocationDto.name, updateLocationDto.capacity)
  }

  @Roles(Role.Admin)
  @Delete(":id")
  async delete(@Param("id") id: string): Promise<Location> {
    let location = await this.locationsService.getById(id)
    if (!location) {
      throw new NotFoundException("Location not found")
    }

    const user = await this.usersService.getById(location.user.id)
    if (!user) {
      throw new InternalServerErrorException("User from location couldn't be fetched")
    }

    location = await this.locationsService.delete(location)
    await this.usersService.delete(user)

    return location
  }
}
