import { CanActivate, ExecutionContext, Injectable } from "@nestjs/common"
import { Reflector } from "@nestjs/core"
import { Role, User } from "types-custom"
import { ROLES_KEY } from "../decorators/roles.decorator"
import { Request } from "express"

@Injectable()
export class RolesGuard implements CanActivate {
  constructor(private reflector: Reflector) {}

  canActivate(context: ExecutionContext): boolean {
    const requiredRoles = this.reflector.getAllAndOverride<Role[] | undefined>(ROLES_KEY, [
      context.getHandler(),
      context.getClass(),
    ])
    if (!requiredRoles) {
      return true
    }
    const { user }: Request = context.switchToHttp().getRequest()
    return requiredRoles.some((role) => (user as User).roles.includes(role))
  }
}