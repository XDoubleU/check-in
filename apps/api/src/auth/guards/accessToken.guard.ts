import { type ExecutionContext, Injectable } from "@nestjs/common"
import { Reflector } from "@nestjs/core"
import { AuthGuard } from "@nestjs/passport"
import { IS_PUBLIC_KEY } from "../decorators/public.decorator"
import { type Observable } from "rxjs"

@Injectable()
export class AccessTokenGuard extends AuthGuard("jwt") {
  private readonly reflector: Reflector

  public constructor(reflector: Reflector) {
    super()

    this.reflector = reflector
  }

  public override canActivate(
    context: ExecutionContext
  ): boolean | Promise<boolean> | Observable<boolean> {
    const isPublic = this.reflector.getAllAndOverride<boolean>(IS_PUBLIC_KEY, [
      context.getHandler(),
      context.getClass()
    ])

    if (isPublic) {
      return true
    }
    return super.canActivate(context)
  }
}
