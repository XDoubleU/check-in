import { PassportStrategy } from "@nestjs/passport"
import { ExtractJwt, Strategy } from "passport-jwt"
import { Request } from "express"
import { Injectable } from "@nestjs/common"
import { JwtPayload } from "./accessToken.strategy"

export type JwtPayloadWithRefresh = {
  sub: string
  username: string
  refreshToken: string
};

@Injectable()
export class RefreshTokenStrategy extends PassportStrategy(
  Strategy,
  "jwt-refresh",
) {
  constructor() {
    super({
      jwtFromRequest: ExtractJwt.fromAuthHeaderAsBearerToken(),
      secretOrKey: process.env.JWT_REFRESH_SECRET,
      passReqToCallback: true,
    })
  }

  validate(req: Request, payload: JwtPayload): JwtPayloadWithRefresh {
    const authorization = req.get("Authorization")
    if (authorization === undefined){
      throw new Error("Authorization Error")
    }

    const refreshToken = authorization.replace("Bearer", "").trim()
    return { ...payload, refreshToken }
  }
}