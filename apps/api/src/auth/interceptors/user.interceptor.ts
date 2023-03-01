import { CallHandler, ExecutionContext, Injectable, NestInterceptor } from "@nestjs/common"
import { UsersService } from "../../users/users.service"
import { Observable } from "rxjs"

@Injectable()
export class UserInterceptor implements NestInterceptor {
    constructor(private readonly usersService: UsersService) {}

    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    async intercept(context: ExecutionContext, handler: CallHandler): Promise<Observable<any>> {
        const request = context.switchToHttp().getRequest()

        if (request.user) {
            request.user = await this.usersService.getById(request.user.sub)
        }

        return handler.handle()
    }
}