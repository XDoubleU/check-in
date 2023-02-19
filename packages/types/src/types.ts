import { Location as dbLocation, User as dbUser } from "database"

export type { CheckIn, School, User as UserWithPasswordHash } from "database"

export type BaseLocation = Omit<dbLocation, "userId"> & { user: User }

export type Location = BaseLocation & { normalizedName: string, available: number }

export type User = Omit<dbUser, "passwordHash"> & { locationId?: string }