/* eslint-disable @typescript-eslint/no-extraneous-class */
import { Module } from "@nestjs/common"
import { SchoolsController } from "./schools.controller"
import { SchoolsService } from "./schools.service"

@Module({
  controllers: [SchoolsController],
  providers: [SchoolsService],
  exports: [SchoolsService]
})
export class SchoolsModule {}
