import { Location as dbLocation, User as dbUser } from "database"

export type { CheckIn, School, User as UserWithPasswordHash } from "database"

export { Role } from "database"

export type BaseLocation = Omit<dbLocation, "userId"> & { user: User }

export type BaseUser = Omit<dbUser, "passwordHash"> & { location: { id: string } | null }

export type Location = BaseLocation & { normalizedName: string, available: number }

export type User = Omit<dbUser, "passwordHash"> & { locationId?: string }