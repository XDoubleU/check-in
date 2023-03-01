import { Controller, Get } from "@nestjs/common"
import { User } from "types"
import { ReqUser } from "../auth/decorators/user.decorator"

@Controller("users")
export class UsersController {
  @Get("me")
  async getUserInfo(@ReqUser() user: User): Promise<User> {
    return user
  }
}
