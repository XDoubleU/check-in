import { IsNotEmpty } from "class-validator"

export class SignInDto {
  @IsNotEmpty()
  username: string

  @IsNotEmpty()
  password: string
}

export class Tokens {
  @IsNotEmpty()
  accessToken: string

  @IsNotEmpty()
  refreshToken: string
}