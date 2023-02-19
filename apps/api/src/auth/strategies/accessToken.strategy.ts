import { Injectable } from "@nestjs/common"
import { PassportStrategy } from "@nestjs/passport"
import { Request } from "express"
import { ExtractJwt, Strategy } from "passport-jwt"

export type JwtPayload = {
  sub: string
}

@Injectable()
export class AccessTokenStrategy extends PassportStrategy(Strategy, "jwt") {
  constructor() {
    super({
      jwtFromRequest: ExtractJwt.fromExtractors([(request: Request): string | null => {
        const accessToken = request.cookies["accessToken"]
        if (!accessToken) {
          return null
        }

        return accessToken
      }]),
      secretOrKey: process.env.JWT_ACCESS_SECRET
    })
  }

  validate(payload: JwtPayload): JwtPayload {
    return payload
  }
}