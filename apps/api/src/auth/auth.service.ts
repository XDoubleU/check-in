import { Injectable } from "@nestjs/common"
import { UsersService } from "../users/users.service"
import { User } from "types"
import { compareSync } from "bcrypt"
import { JwtService } from "@nestjs/jwt"

export type Tokens = {
  accessToken: string,
  refreshToken: string
}

@Injectable()
export class AuthService {
  constructor(
    private usersService: UsersService,
    private jwtService: JwtService
  ) {}

  async login(username: string, pass: string): Promise<Tokens | null> {
    const user = await this.usersService.getByUserName(username)
    if (user && compareSync(pass, user.passwordHash)) {
      return await this.getTokens(user)
    }
    return null
  }

  async refreshTokens(user: User): Promise<Tokens> {
    return this.getTokens(user)
  }

  private async getTokens(user: User): Promise<Tokens> {
    const [accessToken, refreshToken] = await Promise.all([
      this.jwtService.signAsync(
        {
          sub: user.id,
          username: user.username
        },
        {
          secret: process.env.JWT_ACCESS_SECRET,
          expiresIn: process.env.JWT_ACCESS_EXPIRATION
        }
      ),
      this.jwtService.signAsync(
        {
          sub: user.id,
          username: user.username,
        },
        {
          secret: process.env.JWT_REFRESH_SECRET,
          expiresIn: process.env.JWT_REFRESH_EXPIRATION
        }
      ),
    ])

    return {
      accessToken,
      refreshToken,
    }
  }
}
