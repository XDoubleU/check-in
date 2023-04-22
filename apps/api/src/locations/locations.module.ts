/* eslint-disable @typescript-eslint/no-extraneous-class */
import { Module } from "@nestjs/common"
import { LocationsController } from "./locations.controller"
import { LocationsService } from "./locations.service"
import { UsersModule } from "../users/users.module"
import { MikroOrmModule } from "@mikro-orm/nestjs"
import { LocationEntity } from "../entities"
import { WsModule } from "../ws/ws.module"

@Module({
  imports: [MikroOrmModule.forFeature([LocationEntity]), UsersModule, WsModule],
  controllers: [LocationsController],
  providers: [LocationsService],
  exports: [LocationsService]
})
export class LocationsModule {}
