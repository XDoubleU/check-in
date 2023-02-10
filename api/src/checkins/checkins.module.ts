import { Module } from "@nestjs/common"
import { CheckInsService } from "src/checkins/checkins.service"
import { CheckInsController } from "./checkins.controller"
import { SchoolsModule } from "src/schools/schools.module"
import { LocationsModule } from "src/locations/locations.module"

@Module({
  imports: [SchoolsModule, LocationsModule],
  controllers: [CheckInsController],
  providers: [CheckInsService]
})
export class CheckInsModule {}
