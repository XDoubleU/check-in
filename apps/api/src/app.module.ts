/* eslint-disable @typescript-eslint/no-extraneous-class */
import { Module } from "@nestjs/common"
import { CheckInsModule } from "./checkins/checkins.module"
import { LocationsModule } from "./locations/locations.module"
import { SchoolsModule } from "./schools/schools.module"
import { AuthModule } from "./auth/auth.module"
import { APP_GUARD } from "@nestjs/core"
import { AccessTokenGuard } from "./auth/guards/accessToken.guard"
import { UsersModule } from "./users/users.module"
import { RolesGuard } from "./auth/guards/roles.guard"
import config from "mikro-orm-config"
import { MikroOrmModule } from "@mikro-orm/nestjs"
import { SseModule } from "./sse/sse.module"
import { ThrottlerGuard, ThrottlerModule } from "@nestjs/throttler"
import { SentryModule } from "@ntegral/nestjs-sentry"
import { MigrationsModule } from "./migrations/migrations.module"

@Module({
  imports: [
    SentryModule.forRoot({
      dsn: process.env.API_SENTRY_DSN ?? "",
      tracesSampleRate: 1.0,
      logLevels: ["error", "warn"]
    }),
    MikroOrmModule.forRoot({
      ...config,
      autoLoadEntities: true
    }),
    ThrottlerModule.forRoot({
      ttl: 10, // the number of seconds that each request will last in storage
      limit: 30 // the maximum number of requests within the TTL limit
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
