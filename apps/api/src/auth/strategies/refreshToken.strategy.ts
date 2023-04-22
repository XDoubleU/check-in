import { PassportStrategy } from "@nestjs/passport"
import { ExtractJwt, Strategy } from "passport-jwt"
import { type Request } from "express"
import { Injectable } from "@nestjs/common"
import { type JwtPayload } from "./accessToken.strategy"
import { UsersService } from "../../users/users.service"
import { type Tokens } from "types-custom"
import { type UserEntity } from "../../entities"

@Injectable()
export class RefreshTokenStrategy extends PassportStrategy(
  Strategy,
  "jwt-refresh"
) {
  private readonly usersService: UsersService

  public constructor(usersService: UsersService) {
    super({
      jwtFromRequest: ExtractJwt.fromExtractors([
        (request: Request): string | null => {
          const cookies = request.cookies as Tokens
          const refreshToken = cookies.refreshToken
          if (!refreshToken) {
            return null
          }

          return refreshToken
        }
      ]),
      secretOrKey: process.env.JWT_REFRESH_SECRET
    })

    this.usersService = usersService
  }

  public async validate(payload: JwtPayload): Promise<UserEntity | null> {
    return await this.usersService.getById(payload.sub)
  }
}
