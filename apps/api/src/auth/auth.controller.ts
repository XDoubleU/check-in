import { Body, Controller, Get, Post, Res, UnauthorizedException, UseGuards } from "@nestjs/common"
import { AuthService } from "./auth.service"
import { RefreshTokenGuard } from "./guards/refreshToken.guard"
import type { Response } from "express"
import { Public } from "./decorators/public.decorator"
import { ReqUser } from "./decorators/user.decorator"
import type { SignInDto } from "types-custom"
import { UserEntity } from "mikro-orm-config"

@Controller("auth")
export class AuthController {
  constructor(
    private authService: AuthService
  ) {}
  
  @Public()
  @Post("signin")
  async signin(@Body() signinDto: SignInDto, @Res({ passthrough: true }) res: Response): Promise<void> {
    const tokens = await this.authService.signin(signinDto.username, signinDto.password)
    if (!tokens) {
      throw new UnauthorizedException("Invalid credentials")
    }

    this.authService.setTokensAsCookies(tokens, res)
    res.status(200)
  }

  @Get("signout")
  signout(@Res({ passthrough: true }) res: Response): void {
    res.clearCookie("accessToken")
    res.clearCookie("refreshToken", {
      path: "/auth/refresh"
    })
  }

  @Public()
  @UseGuards(RefreshTokenGuard)
  @Get("refresh")
  async refreshTokens(@ReqUser() user: UserEntity, @Res({ passthrough: true }) res: Response): Promise<void> {
    const tokens = await this.authService.refreshTokens(user)
    this.authService.setTokensAsCookies(tokens, res)
  }
}
