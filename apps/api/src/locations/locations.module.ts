/* eslint-disable @typescript-eslint/no-extraneous-class */
import { Module } from "@nestjs/common"
import { LocationsController } from "./locations.controller"
import { LocationsService } from "./locations.service"
import { UsersModule } from "../users/users.module"
import { SseModule } from "../sse/sse.module"
import { MikroOrmModule } from "@mikro-orm/nestjs"
import { LocationEntity } from "mikro-orm-config"

@Module({
  imports: [
    MikroOrmModule.forFeature([LocationEntity]),
    UsersModule,
    SseModule
  ],
  controllers: [LocationsController],
  providers: [LocationsService],
  exports: [LocationsService]
})
export class LocationsModule {}
