export interface SignInDto {
  username: string
  password: string
}

export interface Tokens {
  accessToken: string
  refreshToken: string
}