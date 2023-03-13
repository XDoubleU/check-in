import { PassportStrategy } from "@nestjs/passport"
import { ExtractJwt, Strategy } from "passport-jwt"
import { Request } from "express"
import { Injectable } from "@nestjs/common"
import { JwtPayload } from "./accessToken.strategy"
import { UsersService } from "../../users/users.service"
import { Tokens } from "types-custom"
import { UserEntity } from "mikro-orm-config"

@Injectable()
export class RefreshTokenStrategy extends PassportStrategy(
  Strategy,
  "jwt-refresh",
) {
  constructor(private readonly usersService: UsersService) {
    super({
      jwtFromRequest: ExtractJwt.fromExtractors([(request: Request): string | null => {
        const cookies = request.cookies as Tokens
        const refreshToken = cookies.refreshToken
        if (!refreshToken) {
          return null
        }

        return refreshToken
      }]),
      secretOrKey: process.env.JWT_REFRESH_SECRET
    })
  }

  async validate(payload: JwtPayload): Promise<UserEntity | null> {
    return await this.usersService.getById(payload.sub)
  }
}