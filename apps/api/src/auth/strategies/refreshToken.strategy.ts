import { PassportStrategy } from "@nestjs/passport"
import { ExtractJwt, Strategy } from "passport-jwt"
import { Request } from "express"
import { Injectable } from "@nestjs/common"
import { JwtPayload } from "./accessToken.strategy"
import { UsersService } from "../../users/users.service"
import { User } from "types"

@Injectable()
export class RefreshTokenStrategy extends PassportStrategy(
  Strategy,
  "jwt-refresh",
) {
  constructor(private readonly usersService: UsersService) {
    super({
      jwtFromRequest: ExtractJwt.fromExtractors([(request: Request): string | null => {
        const refreshToken = request.cookies["refreshToken"]
        if (!refreshToken) {
          return null
        }

        return refreshToken
      }]),
      secretOrKey: process.env.JWT_REFRESH_SECRET,
      passReqToCallback: true,
    })
  }

  async validate(payload: JwtPayload): Promise<User> {
    return await this.usersService.getById(payload.sub) as User
  }
}