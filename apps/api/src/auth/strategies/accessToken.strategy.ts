import { Injectable } from "@nestjs/common"
import { PassportStrategy } from "@nestjs/passport"
import { Request } from "express"
import { ExtractJwt, Strategy } from "passport-jwt"
import { UsersService } from "../../users/users.service"
import { Tokens } from "types-custom"
import { UserEntity } from "mikro-orm-config"

export interface JwtPayload {
  sub: string
}

@Injectable()
export class AccessTokenStrategy extends PassportStrategy(Strategy, "jwt") {
  constructor(private readonly usersService: UsersService) {
    super({
      jwtFromRequest: ExtractJwt.fromExtractors([(request: Request): string | null => {
        const cookies = request.cookies as Tokens
        const accessToken = cookies.accessToken
        if (!accessToken) {
          return null
        }

        return accessToken
      }]),
      secretOrKey: process.env.JWT_ACCESS_SECRET
    })
  }

  async validate(payload: JwtPayload): Promise<UserEntity | null> {
    return await this.usersService.getById(payload.sub)
  }
}