import { Location as dbLocation, User as dbUser } from "prisma-database"

export type { CheckIn, School, User as UserWithPasswordHash } from "prisma-database"

export { Role } from "prisma-database"

export type BaseLocation = Omit<dbLocation, "userId"> & { user: Omit<dbUser, "passwordHash"|"roles"> }

export type BaseUser = Omit<dbUser, "passwordHash"> & { location: { id: string } | null }

export type Location = BaseLocation & { normalizedName: string, available: number }

export type User = Omit<dbUser, "passwordHash"> & { locationId?: string }