import { createParamDecorator, ExecutionContext } from "@nestjs/common"

export const ReqUser = createParamDecorator(
  (_data: unknown, ctx: ExecutionContext) => {
    const request = ctx.switchToHttp().getRequest()
    return request.user
  },
)