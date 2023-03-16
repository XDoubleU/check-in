import { type CustomDecorator, SetMetadata } from "@nestjs/common"
import { type Role } from "types-custom"

export const ROLES_KEY = "roles"
export const Roles = (...roles: Role[]): CustomDecorator => SetMetadata(ROLES_KEY, roles)