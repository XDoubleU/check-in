import { type CanActivate, type ExecutionContext, Injectable } from "@nestjs/common"
import { Reflector } from "@nestjs/core"
import { type Role, type User } from "types-custom"
import { ROLES_KEY } from "../decorators/roles.decorator"
import { type Request } from "express"

@Injectable()
export class RolesGuard implements CanActivate {
  private readonly reflector: Reflector

  public constructor(reflector: Reflector) {
    this.reflector = reflector
  }

  public canActivate(context: ExecutionContext): boolean {
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