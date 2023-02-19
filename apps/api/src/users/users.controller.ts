import { Controller, Get, NotFoundException, Req } from "@nestjs/common"
import { UsersService } from "./users.service"
import { TokenRequest } from "../auth/auth.controller"
import { User } from "types"

@Controller("users")
export class UsersController {
  constructor(
    private readonly usersService: UsersService
  ) {}

  @Get("me")
  async getUserInfo(@Req() req: TokenRequest): Promise<User> {
    const user = await this.usersService.getById(req.user.sub)
    if (!user) {
      throw new NotFoundException()
    }
    
    return user
  }
}
