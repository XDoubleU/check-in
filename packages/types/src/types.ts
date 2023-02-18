import { Location as BaseLocation, User as BaseUser } from "database"

export type { CheckIn, School, User as UserWithPasswordHash } from "database"

export type Location = BaseLocation & { user: User }

export type User = Omit<BaseUser, "passwordHash"> & {locationId?: string}