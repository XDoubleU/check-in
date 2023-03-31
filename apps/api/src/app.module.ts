/* eslint-disable @typescript-eslint/no-extraneous-class */
import { HttpException, Module } from "@nestjs/common"
import { CheckInsModule } from "./checkins/checkins.module"
import { LocationsModule } from "./locations/locations.module"
import { SchoolsModule } from "./schools/schools.module"
import { AuthModule } from "./auth/auth.module"
import { APP_GUARD, APP_INTERCEPTOR } from "@nestjs/core"
import { AccessTokenGuard } from "./auth/guards/accessToken.guard"
import { UsersModule } from "./users/users.module"
import { RolesGuard } from "./auth/guards/roles.guard"
import config from "mikro-orm-config"
import { MikroOrmModule } from "@mikro-orm/nestjs"
import { SseModule } from "./sse/sse.module"
import { ThrottlerGuard, ThrottlerModule } from "@nestjs/throttler"
import { MigrationsModule } from "./migrations/migrations.module"
import { RavenInterceptor, RavenModule } from "nest-raven"

@Module({
  imports: [
    RavenModule,
    MikroOrmModule.forRoot({
      ...config,
      autoLoadEntities: true
    }),
    ThrottlerModule.forRoot({
      ttl: 10, // the number of seconds that each request will last in storage
      limit: 60 // the maximum number of requests within the TTL limit
    }),
    CheckInsModule,
    LocationsModule,
    SchoolsModule,
    AuthModule,
    UsersModule,
    SseModule,
    MigrationsModule
  ],
  providers: [
    {
      provide: APP_INTERCEPTOR,
      useValue: new RavenInterceptor({
        filters: [
          // Filter exceptions of type HttpException. Ignore those that
          // have status code of less than 500
          {
            type: HttpException,
            filter: (exception: HttpException): boolean =>
              500 > exception.getStatus()
          }
        ]
      })
    },
    {
      provide: APP_GUARD,
      useClass: ThrottlerGuard
    },
    {
      provide: APP_GUARD,
      useClass: AccessTokenGuard
    },
    {
      provide: APP_GUARD,
      useClass: RolesGuard
    }
  ]
})
export class AppModule {}
