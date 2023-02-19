import { Body, Controller, Get, Post, Req, Res, UnauthorizedException, UseGuards } from "@nestjs/common"
import { AuthService } from "./auth.service"
import { SignInDto } from "types"
import { RefreshTokenGuard } from "./guards/refreshToken.guard"
import { UsersService } from "../users/users.service"
import { Request, Response } from "express"
import { JwtPayload } from "./strategies/accessToken.strategy"
import { Public } from "./decorators/public.decorator"

export type TokenRequest = Request & { user: JwtPayload }

@Controller("auth")
export class AuthController {
  constructor(
    private usersService: UsersService,
    private authService: AuthService
  ) {}
  
  @Public()
  @Post("signin")
  async signin(@Body() signinDto: SignInDto, @Res({ passthrough: true }) res: Response): Promise<void> {
    const tokens = await this.authService.signin(signinDto.username, signinDto.password)
    if (tokens === null) {
      throw new UnauthorizedException("Invalid credentials")
    }

    this.authService.setTokensAsCookies(tokens, res)
    res.status(200)
  }

  @Get("signout")
  async signout(@Res({ passthrough: true }) res: Response): Promise<void> {
    res.clearCookie("accessToken")
    res.clearCookie("refreshToken", {
      path: "/auth/refresh"
    })
  }

  @Public()
  @UseGuards(RefreshTokenGuard)
  @Get("refresh")
  async refreshTokens(@Req() req: TokenRequest, @Res({ passthrough: true }) res: Response): Promise<void> {
    const user = await this.usersService.getById(req.user.sub)
    if (user === null) {
      throw new UnauthorizedException("Invalid refreshtoken")
    }

    const tokens = await this.authService.refreshTokens(user)
    this.authService.setTokensAsCookies(tokens, res)
  }
}
