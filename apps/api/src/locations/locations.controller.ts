/* eslint-disable sonarjs/no-duplicate-string */
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
  Query,
  Res
} from "@nestjs/common"
import { LocationsService } from "./locations.service"
import { UsersService } from "../users/users.service"
import {
  type GetAllPaginatedLocationDto,
  Role,
  type CreateLocationDto,
  type UpdateLocationDto
} from "types-custom"
import { ReqUser } from "../auth/decorators/user.decorator"
import { Roles } from "../auth/decorators/roles.decorator"
import { type LocationEntity, UserEntity } from "mikro-orm-config"
import { Public } from "../auth/decorators/public.decorator"
import { type Response } from "express"
import { convertToLocationUpdateEventDto } from "../helpers/conversion"

type MikroGetAllPaginatedLocationDto = Omit<
  GetAllPaginatedLocationDto,
  "data"
> & { data: LocationEntity[] }

@Controller("locations")
export class LocationsController {
  private readonly locationsService: LocationsService
  private readonly usersService: UsersService

  public constructor(
    locationsService: LocationsService,
    usersService: UsersService
  ) {
    this.locationsService = locationsService
    this.usersService = usersService
  }

  @Get("ws")
  @Public()
  public async getInitStateSse(@Res() res: Response): Promise<void> {
    res.set("Access-Control-Allow-Origin", "*")

    const locations = await this.locationsService.getAll()
    const data = locations.map((location) => {
      return convertToLocationUpdateEventDto(location)
    })

    res.json(data)
  }

  @Roles(Role.User)
  @Get("me")
  public getLocationOfUser(@ReqUser() user: UserEntity): LocationEntity | undefined {    
    return user.location
  }

  @Roles(Role.Manager)
  @Get()
  public async getAll(
    @Query("page") queryPage?: string
  ): Promise<MikroGetAllPaginatedLocationDto> {
    const pageSize = 3
    const current = queryPage ? parseInt(queryPage) : 1
    const amountOfLocations = await this.locationsService.getTotalCount()
    const locations = await this.locationsService.getAllPaged(current, pageSize)

    return {
      data: locations,
      pagination: {
        current,
        total: Math.ceil(amountOfLocations / pageSize)
      }
    }
  }

  @Get(":id")
  public async get(
    @ReqUser() user: UserEntity,
    @Param("id") id: string
  ): Promise<LocationEntity> {
    const location = await this.locationsService.getLocation(id, user)
    if (!location) {
      throw new NotFoundException("Location not found")
    }

    return location
  }

  @Roles(Role.Manager)
  @Post()
  public async create(
    @Body() createLocationDto: CreateLocationDto
  ): Promise<LocationEntity> {
    const existingLocation = await this.locationsService.getByName(
      createLocationDto.name
    )
    if (existingLocation) {
      throw new ConflictException("Location with this name already exists")
    }

    const existingUser = await this.usersService.getByUserName(
      createLocationDto.username
    )
    if (existingUser) {
      throw new ConflictException("User with this username already exists")
    }

    const user = await this.usersService.create(
      createLocationDto.username,
      createLocationDto.password,
      Role.User
    )
    return await this.locationsService.create(
      createLocationDto.name,
      createLocationDto.capacity,
      user
    )
  }

  @Patch(":id")
  public async update(
    @ReqUser() reqUser: UserEntity,
    @Param("id") id: string,
    @Body() updateLocationDto: UpdateLocationDto
  ): Promise<LocationEntity> {
    const location = await this.locationsService.getLocation(id, reqUser)
    if (!location) {
      throw new NotFoundException("Location not found")
    }

    if (updateLocationDto.name) {
      const existingLocation = await this.locationsService.getByName(
        updateLocationDto.name
      )
      if (existingLocation) {
        throw new ConflictException("Location with this name already exists")
      }
    }

    if (updateLocationDto.username) {
      const existingUser = await this.usersService.getByUserName(
        updateLocationDto.username
      )
      if (existingUser) {
        throw new ConflictException("User with this username already exists")
      }
    }

    const user = await this.usersService.getById(location.user.id)

    await this.usersService.update(
      // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
      user!,
      updateLocationDto.username,
      updateLocationDto.password
    )

    return await this.locationsService.update(
      location,
      updateLocationDto.name,
      updateLocationDto.capacity
    )
  }

  @Roles(Role.Manager)
  @Delete(":id")
  public async delete(@Param("id") id: string): Promise<LocationEntity> {
    let location = await this.locationsService.getById(id)
    if (!location) {
      throw new NotFoundException("Location not found")
    }

    const user = await this.usersService.getById(location.user.id)

    location = await this.locationsService.delete(location)
    // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
    await this.usersService.delete(user!)

    return location
  }
}
