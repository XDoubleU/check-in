import { Injectable, InternalServerErrorException } from "@nestjs/common"
import { UsersService } from "../users/users.service"
import { Role, type Tokens } from "types-custom"
import { JwtService } from "@nestjs/jwt"
import { type Response } from "express"
import { compare } from "bcrypt"
import { UserEntity } from "mikro-orm-config"

export interface UserAndTokens {
  user: UserEntity
  tokens: Tokens
}

@Injectable()
export class AuthService {
  private readonly usersService: UsersService
  private readonly jwtService: JwtService

  public constructor(usersService: UsersService, jwtService: JwtService) {
    this.usersService = usersService
    this.jwtService = jwtService
  }

  public async signin(
    username: string,
    password: string
  ): Promise<UserAndTokens | null> {
    if (
      username === process.env.ADMIN_USERNAME &&
      password === process.env.ADMIN_PASSWORD
    ) {
      return {
        user: new UserEntity(username, password, Role.Admin),
        tokens: await this.getAdminAccessToken()
      }
    }

    const user = await this.usersService.getByUserName(username)
    if (!user) {
      return null
    }

    if (!(await compare(password, user.passwordHash))) {
      return null
    }

    return {
      user,
      tokens: await this.getTokens(user)
    }
  }

  public async refreshTokens(user: UserEntity): Promise<Tokens> {
    return this.getTokens(user)
  }

  public setTokensAsCookies(
    tokens: Tokens,
    res: Response,
    rememberMe: boolean
  ): void {
    const accessTokenExpires = parseInt(
      (this.jwtService.decode(tokens.accessToken) as Record<string, string>).exp
    )

    res.cookie("accessToken", tokens.accessToken, {
      expires: new Date(accessTokenExpires * 1000),
      sameSite: "strict",
      httpOnly: true,
      secure: process.env.NODE_ENV === "production"
    })

    if (tokens.refreshToken) {
      const refreshTokenExpires = parseInt(
        (this.jwtService.decode(tokens.refreshToken) as Record<string, string>)
          .exp
      )

      if (rememberMe) {
        res.cookie("refreshToken", tokens.refreshToken, {
          expires: new Date(refreshTokenExpires * 1000),
          sameSite: "strict",
          httpOnly: true,
          secure: process.env.NODE_ENV === "production"
        })
      }
    }
  }

  public async getTokens(user: UserEntity): Promise<Tokens> {
    if (
      !process.env.JWT_ACCESS_SECRET ||
      !process.env.JWT_ACCESS_EXPIRATION ||
      !process.env.JWT_REFRESH_SECRET ||
      !process.env.JWT_REFRESH_EXPIRATION
    ) {
      throw new InternalServerErrorException(
        "JWT secrets or expirations missing in environment"
      )
    }

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
      )
    ])

    return {
      accessToken,
      refreshToken
    }
  }

  public async getAdminAccessToken(): Promise<Tokens> {
    if (!process.env.JWT_ACCESS_SECRET || !process.env.JWT_ACCESS_EXPIRATION) {
      throw new InternalServerErrorException(
        "JWT secrets or expirations missing in environment"
      )
    }

    const accessToken = await this.jwtService.signAsync(
      {
        admin: true
      },
      {
        secret: process.env.JWT_ACCESS_SECRET,
        expiresIn: process.env.JWT_ACCESS_EXPIRATION
      }
    )

    return { accessToken }
  }
}
