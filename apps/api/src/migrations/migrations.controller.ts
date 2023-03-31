import { Controller, Get } from "@nestjs/common"
import { Role } from "types-custom"
import { Roles } from "../auth/decorators/roles.decorator"
import { MikroORM } from "@mikro-orm/core"
import { DatabaseSeeder } from "mikro-orm-config"

@Controller("migrations")
export class MigrationsController {
  private readonly orm: MikroORM

  public constructor(orm: MikroORM) {
    this.orm = orm
  }

  @Roles(Role.Admin)
  @Get("up")
  public async applyMigrationsUp(): Promise<string> {
    // eslint-disable-next-line @typescript-eslint/no-non-null-assertion

    const response = await this.orm.getMigrator().up()
    return response[0].name
  }

  @Roles(Role.Admin)
  @Get("down")
  public async applyMigrationsDown(): Promise<string> {
    // eslint-disable-next-line @typescript-eslint/no-non-null-assertion

    const response = await this.orm.getMigrator().down()
    return response[0].name
  }

  @Roles(Role.Admin)
  @Get("seed")
  public async applySeeder(): Promise<void> {
    // eslint-disable-next-line @typescript-eslint/no-non-null-assertion

    await this.orm.getSeeder().seed(DatabaseSeeder)
  }
}
