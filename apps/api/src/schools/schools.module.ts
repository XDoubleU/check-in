/* eslint-disable @typescript-eslint/no-extraneous-class */
import { MikroOrmModule } from "@mikro-orm/nestjs"
import { Module } from "@nestjs/common"
import { SchoolEntity } from "../entities"
import { SchoolsController } from "./schools.controller"
import { SchoolsService } from "./schools.service"

@Module({
  imports: [MikroOrmModule.forFeature([SchoolEntity])],
  controllers: [SchoolsController],
  providers: [SchoolsService],
  exports: [SchoolsService]
})
export class SchoolsModule {}
