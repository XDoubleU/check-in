import { Module } from "@nestjs/common"
import { CheckInsModule } from "./checkins/checkins.module"
import { LocationsModule } from "./locations/locations.module"
import { SchoolsModule } from "./schools/schools.module"

@Module({
  imports: [CheckInsModule, LocationsModule, SchoolsModule]
})
export class AppModule {}
