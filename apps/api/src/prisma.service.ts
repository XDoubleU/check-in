import { INestApplication, Injectable, OnModuleInit } from "@nestjs/common"
import { PrismaClient } from "prisma-database"

@Injectable()
export class PrismaService extends PrismaClient implements OnModuleInit {
  async onModuleInit(): Promise<void> {
    await this.$connect()
  }

  enableShutdownHooks(app: INestApplication): void {
    this.$on("beforeExit", () => {
      async (): Promise<void> => {
        await app.close()
      }
    })
  }
}
