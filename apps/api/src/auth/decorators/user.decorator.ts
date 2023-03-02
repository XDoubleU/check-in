import { createParamDecorator, ExecutionContext } from "@nestjs/common"
import { Request } from "express"

export const ReqUser = createParamDecorator(
  (_data: unknown, ctx: ExecutionContext) => {
    const request: Request = ctx.switchToHttp().getRequest()
    return request.user
  },
)