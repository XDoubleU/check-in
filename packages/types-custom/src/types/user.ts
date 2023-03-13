import { Location } from "./location"


export interface User {
  id: string
  username: string
  passwordHash: string
  roles: Role[]
  location?: Location
}

export enum Role {
  User = "user",
  Admin = "admin"
}