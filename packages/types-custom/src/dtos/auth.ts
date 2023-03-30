export interface SignInDto {
  username: string
  password: string
  rememberMe: boolean
}

export interface Tokens {
  accessToken: string
  refreshToken?: string
}
