/* eslint-disable @typescript-eslint/no-extraneous-class */
import { Module } from "@nestjs/common"
import { CheckInsController } from "./checkins.controller"
import { LocationsModule } from "../locations/locations.module"
import { SchoolsModule } from "../schools/schools.module"
import { CheckInsService } from "./checkins.service"
import { SseModule } from "../sse/sse.module"
import { MikroOrmModule } from "@mikro-orm/nestjs"
import { CheckInEntity } from "mikro-orm-config"

@Module({
  imports: [
    MikroOrmModule.forFeature([CheckInEntity]),
    SchoolsModule,
    LocationsModule,
    SseModule
  ],
  controllers: [CheckInsController],
  providers: [CheckInsService]
})
export class CheckInsModule {}
