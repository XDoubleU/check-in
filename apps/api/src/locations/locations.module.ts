import { Module } from "@nestjs/common"
import { LocationsController } from "./locations.controller"
import { LocationsService } from "./locations.service"
import { UsersService } from "src/users/users.service"

@Module({
  imports: [UsersService],
  controllers: [LocationsController],
  providers: [LocationsService],
  exports: [LocationsService]
})
export class LocationsModule {}
