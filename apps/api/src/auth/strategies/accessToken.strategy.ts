import { Injectable } from "@nestjs/common"
import { PassportStrategy } from "@nestjs/passport"
import { type Request } from "express"
import { ExtractJwt, Strategy } from "passport-jwt"
import { UsersService } from "../../users/users.service"
import { Role, type Tokens } from "types-custom"
import { UserEntity } from "../../entities"

export interface JwtPayload {
  admin?: boolean
  sub: string
}

@Injectable()
export class AccessTokenStrategy extends PassportStrategy(Strategy, "jwt") {
  private readonly usersService: UsersService

  public constructor(usersService: UsersService) {
    super({
      jwtFromRequest: ExtractJwt.fromExtractors([
        (request: Request): string | null => {
          const cookies = request.cookies as Tokens
          const accessToken = cookies.accessToken
          if (!accessToken) {
            return null
          }

          return accessToken
        }
      ]),
      secretOrKey: process.env.JWT_ACCESS_SECRET
    })

    this.usersService = usersService
  }

  public async validate(payload: JwtPayload): Promise<UserEntity | null> {
    if (payload.admin) {
      // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
      return new UserEntity(
        // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
        process.env.ADMIN_USERNAME!,
        undefined,
        Role.Admin
      )
    }
    return await this.usersService.getById(payload.sub)
  }
}
