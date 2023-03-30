import {
  Body,
  Controller,
  Get,
  Post,
  Res,
  UnauthorizedException,
  UseGuards
} from "@nestjs/common"
import { AuthService } from "./auth.service"
import { RefreshTokenGuard } from "./guards/refreshToken.guard"
import { Response } from "express"
import { Public } from "./decorators/public.decorator"
import { ReqUser } from "./decorators/user.decorator"
import { SignInDto } from "types-custom"
import { UserEntity } from "mikro-orm-config"

@Controller("auth")
export class AuthController {
  private readonly authService: AuthService

  public constructor(authService: AuthService) {
    this.authService = authService
  }

  @Public()
  @Post("signin")
  public async signin(
    @Body() signinDto: SignInDto,
    @Res({ passthrough: true }) res: Response
  ): Promise<UserEntity> {
    const userAndTokens = await this.authService.signin(
      signinDto.username,
      signinDto.password
    )
    if (!userAndTokens) {
      throw new UnauthorizedException("Invalid credentials")
    }

    this.authService.setTokensAsCookies(
      userAndTokens.tokens,
      res,
      signinDto.rememberMe
    )

    res.status(200)

    return userAndTokens.user
  }

  @Get("signout")
  public signout(@Res({ passthrough: true }) res: Response): void {
    res.clearCookie("accessToken")
    res.clearCookie("refreshToken", {
      path: "/auth/refresh"
    })
  }

  @Public()
  @UseGuards(RefreshTokenGuard)
  @Get("refresh")
  public async refreshTokens(
    @ReqUser() user: UserEntity,
    @Res({ passthrough: true }) res: Response
  ): Promise<void> {
    const tokens = await this.authService.refreshTokens(user)
    this.authService.setTokensAsCookies(tokens, res, true)
  }
}
