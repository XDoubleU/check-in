import { Controller, Get, NotFoundException, Param } from "@nestjs/common"
import { UserEntity } from "mikro-orm-config"
import { Role } from "types-custom"
import { Roles } from "../auth/decorators/roles.decorator"
import { ReqUser } from "../auth/decorators/user.decorator"
import { UsersService } from "./users.service"

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

  @Roles(Role.Admin)
  @Get(":id")
  public async get(@Param("id") id: string): Promise<UserEntity> {
    const user = await this.usersService.getById(id)
    if (!user) {
      throw new NotFoundException("User not found")
    }

    return user
  }
}
