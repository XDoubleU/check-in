import { Location as BaseLocation, User } from "database"

export type { User, CheckIn, School } from "database"

export type Location = BaseLocation & { user: User }