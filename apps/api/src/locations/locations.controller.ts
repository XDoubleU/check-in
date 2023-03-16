import { BadRequestException, Body, ConflictException, Controller, Delete, Get, InternalServerErrorException, NotFoundException, Param, Patch, Post, Query } from "@nestjs/common"
import { LocationsService } from "./locations.service"
import { UsersService } from "../users/users.service"
import  { type GetAllPaginatedLocationDto } from "types-custom"
import { Role , CreateLocationDto, UpdateLocationDto } from "types-custom"
import { ReqUser } from "../auth/decorators/user.decorator"
import { Roles } from "../auth/decorators/roles.decorator"
import { type LocationEntity, UserEntity } from "mikro-orm-config"

type MikroGetAllPaginatedLocationDto = Omit<GetAllPaginatedLocationDto, "locations"> & { locations: LocationEntity[] }

@Controller("locations")
export class LocationsController {
  private readonly locationsService: LocationsService
  private readonly usersService: UsersService

  public constructor(locationsService: LocationsService, usersService: UsersService) {
    this.locationsService = locationsService
    this.usersService = usersService
  }

  @Roles(Role.Admin)
  @Get()
  public async getAll(@Query("page") queryPage?: string): Promise<MikroGetAllPaginatedLocationDto> {
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
  public async getMyLocation(@ReqUser() user: UserEntity): Promise<LocationEntity> {
    if (!user.location?.id) {
      throw new BadRequestException()
    }

    const location = await this.locationsService.getById(user.location.id)
    if (!location) {
      throw new NotFoundException("Location not found")
    }

    return location
  }

  @Get(":id")
  public async get(@ReqUser() user: UserEntity, @Param("id") id: string): Promise<LocationEntity> {
    const location = await this.locationsService.getById(id)
    if (!location || (!user.roles.includes(Role.Admin) && location.user.id !== user.id)) {
      throw new NotFoundException("Location not found")
    }

    return location
  }

  @Roles(Role.Admin)
  @Post()
  public async create(@Body() createLocationDto: CreateLocationDto): Promise <LocationEntity> {
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
  public async update(@ReqUser() reqUser: UserEntity, @Param("id") id: string, @Body() updateLocationDto: UpdateLocationDto): Promise<LocationEntity> {
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
  public async delete(@Param("id") id: string): Promise<LocationEntity> {
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
