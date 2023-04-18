import { type Role } from "./role"

export interface User {
  id: string
  username: string
  passwordHash: string
  roles: Role[]
  locationId?: string
}
