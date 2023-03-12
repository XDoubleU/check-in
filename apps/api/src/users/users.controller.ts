import { Controller, Get } from "@nestjs/common"
import { UserEntity } from "mikro-orm-config"
import { ReqUser } from "../auth/decorators/user.decorator"

@Controller("users")
export class UsersController {
  @Get("me")
  getUserInfo(@ReqUser() user: UserEntity): UserEntity {
    return user
  }
}
