import { Injectable } from "@nestjs/common"
import { UsersService } from "../users/users.service"
import { Tokens, User } from "types"
import { JwtService } from "@nestjs/jwt"
import { Response } from "express"

@Injectable()
export class AuthService {
  constructor(
    private usersService: UsersService,
    private jwtService: JwtService
  ) {}

  async signin(username: string, password: string): Promise<Tokens | null> {
    if (await this.usersService.checkPassword(username, password)) {
      const user = await this.usersService.getByUserName(username)
      if (!user) {
        return null
      }

      return await this.getTokens(user)
    }
    return null
  }

  async refreshTokens(user: User): Promise<Tokens> {
    return this.getTokens(user)
  }

  setTokensAsCookies(tokens: Tokens, res: Response): void {
    const accessTokenExpires = parseInt((this.jwtService.decode(tokens.accessToken) as Record<string, string>).exp)
    const refreshTokenExpires = parseInt((this.jwtService.decode(tokens.refreshToken) as Record<string, string>).exp)

    res.cookie("accessToken", tokens.accessToken, {
      expires: new Date(accessTokenExpires * 1000),
      sameSite: "strict",
      httpOnly: true,
      secure: process.env.NODE_ENV === "production"
    })
    res.cookie("refreshToken", tokens.refreshToken, {
      expires: new Date(refreshTokenExpires * 1000),
      sameSite: "strict",
      httpOnly: true,
      path: "/auth/refresh",
      secure: process.env.NODE_ENV === "production" 
    })
  }

  async getTokens(user: User): Promise<Tokens> {
    const [accessToken, refreshToken] = await Promise.all([
      this.jwtService.signAsync(
        {
          sub: user.id
        },
        {
          secret: process.env.JWT_ACCESS_SECRET,
          expiresIn: process.env.JWT_ACCESS_EXPIRATION
        }
      ),
      this.jwtService.signAsync(
        {
          sub: user.id
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
