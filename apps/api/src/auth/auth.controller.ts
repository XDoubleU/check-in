import { Body, Controller, Get, Post, Req, UnauthorizedException, UseGuards } from "@nestjs/common"
import { AuthService, Tokens } from "./auth.service"
import { LoginDto } from "types"
import { RefreshTokenGuard } from "./guards/refreshToken.guard"
import { UsersService } from "../users/users.service"
import { Request } from "express"
import { JwtPayload } from "./strategies/accessToken.strategy"
import { JwtPayloadWithRefresh } from "./strategies/refreshToken.strategy"
import { Public } from "./decorators/public.decorator"

export type AccessTokenRequest = Request & { user: JwtPayload }
export type RefreshTokenRequest = Request & { user: JwtPayloadWithRefresh }

@Controller("auth")
export class AuthController {
  constructor(
    private usersService: UsersService,
    private authService: AuthService
  ) {}
  
  @Public()
  @Post("login")
  async login(@Body() loginDto: LoginDto): Promise<Tokens> {
    const result = await this.authService.login(loginDto.username, loginDto.password)
    if (result === null) {
      throw new UnauthorizedException("Invalid credentials")
    }
    return result
  }

  @Public()
  @UseGuards(RefreshTokenGuard)
  @Get("refresh")
  async refreshTokens(@Req() req: RefreshTokenRequest): Promise<Tokens> {
    const user = await this.usersService.getById(req.user.sub)
    if (user === null) {
      throw new UnauthorizedException("Invalid refreshtoken")
    }

    return await this.authService.refreshTokens(user)
  }
}
