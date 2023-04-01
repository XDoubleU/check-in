/* eslint-disable @typescript-eslint/no-extraneous-class */
import { MikroOrmModule } from "@mikro-orm/nestjs"
import { Module } from "@nestjs/common"
import { SchoolEntity } from "mikro-orm-config"
import { MigrationsController } from "./migrations.controller"

@Module({
  imports: [MikroOrmModule.forFeature([SchoolEntity])],
  controllers: [MigrationsController]
})
export class MigrationsModule {}
