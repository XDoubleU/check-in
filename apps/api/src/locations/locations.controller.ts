import { BadRequestException, Body, ConflictException, Controller, Delete, Get, NotFoundException, Param, Patch, Post, Query } from "@nestjs/common"
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
    // TODO: test role access
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
    // TODO: test role access
    if (!user.locationId) {
      throw new NotFoundException()
    }
    
    const location = await this.locationsService.getById(user.locationId)
    if (!location) {
      throw new NotFoundException()
    }

    return location
  }

  @Get(":id")
  async get(@ReqUser() user: User, @Param("id") id: string): Promise<Location> {
    const location = await this.locationsService.getById(id)
    if (!location || (user.role !== Role.Admin && location.user.id !== user.id)) {
      throw new NotFoundException("Location not found")
    }

    return location
  }

  @Roles(Role.Admin)
  @Post()
  async create(@Body() createLocationDto: CreateLocationDto): Promise <Location> {
    // TODO: test role access
    const existingLocation = await this.locationsService.getByName(createLocationDto.name)
    if (existingLocation) {
      throw new ConflictException("Location with this name already exists")
    }

    const existingUser = await this.usersService.getByUserName(createLocationDto.username)
    if (existingUser) {
      throw new ConflictException("User with this username already exists")
    }

    const user = await this.usersService.create(createLocationDto.username, createLocationDto.password)
    if (!user) {
      throw new BadRequestException()
    }

    const location = await this.locationsService.create(createLocationDto.name, createLocationDto.capacity, user)
    if (!location) {
      throw new BadRequestException()
    }

    return location
  }

  @Patch(":id")
  async update(@ReqUser() reqUser: User, @Param("id") id: string, @Body() updateLocationDto: UpdateLocationDto): Promise<Location> {
    if (updateLocationDto.name) {
      const existingLocation = await this.locationsService.getByName(updateLocationDto.name)
      if (existingLocation) {
        throw new ConflictException()
      }
    }

    if (updateLocationDto.username) {
      const existingUser = await this.usersService.getByUserName(updateLocationDto.username)
      if (existingUser) {
        throw new ConflictException()
      }
    }
    
    let location = await this.locationsService.getById(id)
    if (!location || (reqUser.role !== Role.Admin && location.user.id !== reqUser.id)) {
      throw new NotFoundException("Location not found")
    }

    let user = await this.usersService.getById(location.user.id)
    if (!user) {
      throw new NotFoundException("User not found")
    }

    user = await this.usersService.update(user, updateLocationDto.username, updateLocationDto.password)
    if (!user) {
      throw new BadRequestException()
    }

    location = await this.locationsService.update(location, updateLocationDto.name, updateLocationDto.capacity)
    if (!location) {
      throw new BadRequestException()
    }

    return location
  }

  @Roles(Role.Admin)
  @Delete(":id")
  async delete(@Param("id") id: string): Promise<Location> {
    // TODO: test role access
    let location = await this.locationsService.getById(id)
    if (!location) {
      throw new NotFoundException("Location not found")
    }

    let user = await this.usersService.getById(location.user.id)
    if (!user) {
      throw new NotFoundException("User not found")
    }

    location = await this.locationsService.delete(location)
    if (!location) {
      throw new BadRequestException()
    }
    
    user = await this.usersService.delete(user)
    if (!user) {
      throw new BadRequestException()
    }

    return location
  }
}
