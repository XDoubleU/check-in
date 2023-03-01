import { Module } from "@nestjs/common"
import { CheckInsModule } from "./checkins/checkins.module"
import { LocationsModule } from "./locations/locations.module"
import { SchoolsModule } from "./schools/schools.module"
import { AuthModule } from "./auth/auth.module"
import { APP_GUARD, APP_INTERCEPTOR } from "@nestjs/core"
import { AccessTokenGuard } from "./auth/guards/accessToken.guard"
import { UserInterceptor } from "./auth/interceptors/user.interceptor"
import { UsersModule } from "./users/users.module"

@Module({
  imports: [CheckInsModule, LocationsModule, SchoolsModule, AuthModule, UsersModule],
  providers: [
    {
      provide: APP_GUARD,
      useClass: AccessTokenGuard,
    },
    {
      provide: APP_INTERCEPTOR,
      useClass: UserInterceptor
    }
  ],
})
export class AppModule {}
