import { Module } from "@nestjs/common"
import { LocationsController } from "./locations.controller"
import { LocationsService } from "./locations.service"
import { UsersModule } from "../users/users.module"

@Module({
  imports: [UsersModule],
  controllers: [LocationsController],
  providers: [LocationsService],
  exports: [LocationsService]
})
export class LocationsModule {}
