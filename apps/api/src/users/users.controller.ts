import { Controller, Get } from "@nestjs/common"
import { User } from "types-custom"
import { ReqUser } from "../auth/decorators/user.decorator"

@Controller("users")
export class UsersController {
  @Get("me")
  getUserInfo(@ReqUser() user: User): User {
    return user
  }
}
