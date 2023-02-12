import { Module } from "@nestjs/common"
import { CheckInsController } from "./checkins.controller"
import { LocationsModule } from "../locations/locations.module"
import { SchoolsModule } from "../schools/schools.module"
import { CheckInsService } from "./checkins.service"

@Module({
  imports: [SchoolsModule, LocationsModule],
  controllers: [CheckInsController],
  providers: [CheckInsService]
})
export class CheckInsModule {}
