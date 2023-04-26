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
import {
  type CreateUserDto,
  type GetAllPaginatedUserDto,
  Role,
  type UpdateUserDto
} from "types-custom"
import { Roles } from "../auth/decorators/roles.decorator"
import { ReqUser } from "../auth/decorators/user.decorator"
import { UsersService } from "./users.service"
import { UserEntity } from "../entities"

type MikroGetAllPaginatedUserDto = Omit<GetAllPaginatedUserDto, "data"> & {
  data: UserEntity[]
}

const NOT_FOUND = "User not found"

@Controller("users")
export class UsersController {
  private readonly usersService: UsersService

  public constructor(usersService: UsersService) {
    this.usersService = usersService
  }

  @Get("me")
  public getUserInfo(@ReqUser() user: UserEntity): UserEntity {
    return user
  }

  @Roles(Role.Manager)
  @Get(":id")
  public async get(@Param("id") id: string): Promise<UserEntity> {
    const user = await this.usersService.getById(id)
    if (!user) {
      throw new NotFoundException(NOT_FOUND)
    }

    return user
  }

  @Roles(Role.Admin)
  @Get()
  public async getAllManagersPaged(
    @Query("page") queryPage?: string
  ): Promise<MikroGetAllPaginatedUserDto> {
    const pageSize = 4
    const current = queryPage ? parseInt(queryPage) : 1
    const amountOfUsers = await this.usersService.getManagerCount()
    const users = await this.usersService.getAllManagersPaged(current, pageSize)

    return {
      data: users,
      pagination: {
        current,
        total: Math.ceil(amountOfUsers / pageSize)
      }
    }
  }

  @Roles(Role.Admin)
  @Post()
  public async create(
    @Body() createUserDto: CreateUserDto
  ): Promise<UserEntity> {
    const existingUser = await this.usersService.getByUserName(
      createUserDto.username
    )
    if (existingUser) {
      throw new ConflictException("User with this username already exists")
    }

    return await this.usersService.create(
      createUserDto.username,
      createUserDto.password,
      Role.Manager
    )
  }

  @Roles(Role.Admin)
  @Patch(":id")
  public async update(
    @Param("id") id: string,
    @Body() updateUserDto: UpdateUserDto
  ): Promise<UserEntity> {
    const user = await this.usersService.getById(id)
    if (!user) {
      throw new NotFoundException(NOT_FOUND)
    }

    const existingUser = await this.usersService.getByUserName(
      updateUserDto.username
    )
    if (existingUser) {
      throw new ConflictException("User with this username already exists")
    }

    return await this.usersService.update(
      user,
      updateUserDto.username,
      updateUserDto.password
    )
  }

  @Roles(Role.Admin)
  @Delete(":id")
  public async delete(@Param("id") id: string): Promise<UserEntity> {
    const user = await this.usersService.getById(id)
    if (!user) {
      throw new NotFoundException(NOT_FOUND)
    }

    return await this.usersService.delete(user)
  }
}
