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

@Module({
  imports: [
    MikroOrmModule.forRoot({
      ...config,
      autoLoadEntities: true
    }),
    CheckInsModule,
    LocationsModule,
    SchoolsModule,
    AuthModule,
    UsersModule
  ],
  providers: [
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
