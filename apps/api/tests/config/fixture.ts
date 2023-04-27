import { ValidationPipe, type INestApplication } from "@nestjs/common"
import { Test } from "@nestjs/testing"
import { AppModule } from "../../src/app.module"
import cookieParser from "cookie-parser"
import { AuthService, type UserAndTokens } from "../../src/auth/auth.service"
import {
  CheckInEntity,
  LocationEntity,
  SchoolEntity,
  UserEntity
} from "../../src/entities"
import { MikroORM, type Transaction } from "@mikro-orm/core"
import { Role } from "types-custom"
import {
  type Knex,
  type EntityManager,
  type PostgreSqlDriver
} from "@mikro-orm/postgresql"
import { TestModule } from "./test.module"
import { ContextManager } from "./test.middleware"
import helmet from "helmet"
import { WsAdapter } from "@nestjs/platform-ws"

export interface ErrorResponse {
  message: string
}

export default class Fixture {
  public app!: INestApplication
  public em!: EntityManager
  public contextManager!: ContextManager
  public mainTransaction!: Transaction<Knex.Transaction>

  public async beforeAll(): Promise<void> {
    const module = await Test.createTestingModule({
      imports: [
        // Import the AppModule without any change to config
        AppModule,
        // Add the test module to register the TransactionContextMiddleware
        TestModule
      ]
    }).compile()

    this.app = module.createNestApplication()
    this.app.use(helmet())
    this.app.use(cookieParser())
    this.app.useGlobalPipes(new ValidationPipe({ transform: true }))

    this.app.useWebSocketAdapter(new WsAdapter(this.app))

    this.app.enableShutdownHooks()

    await this.app.init()

    const orm = this.app.get<MikroORM<PostgreSqlDriver>>(MikroORM)
    this.contextManager = this.app.get(ContextManager)

    this.em = orm.em.fork()

    this.mainTransaction = await this.em.getConnection().begin()
    this.em.setTransactionContext(this.mainTransaction)

    await this.seedDatabase()
  }

  public async beforeEach(): Promise<void> {
    const testTransaction = await this.em
      .getConnection()
      .begin({ ctx: this.mainTransaction })
    this.em.setTransactionContext(testTransaction)
    this.contextManager.setContext(testTransaction)
  }

  public async afterEach(): Promise<void> {
    const testTransaction = this.contextManager.resetContext()
    if (!testTransaction) {
      throw new Error("testTransaction is undefined")
    }

    await this.em.getConnection().rollback(testTransaction)
    this.em.clear()
  }

  public async afterAll(): Promise<void> {
    await this.em.getConnection().rollback(this.mainTransaction)
    await this.app.close()
  }

  public async getTokens(username: string): Promise<UserAndTokens> {
    const authService = this.app.get<AuthService>(AuthService)

    await this.em.find(LocationEntity, {})
    const user = await this.em.findOneOrFail(UserEntity, {
      username: username
    })

    return {
      tokens: await authService.getTokens(user),
      user
    }
  }

  public async getAdminTokens(): Promise<UserAndTokens> {
    const authService = this.app.get<AuthService>(AuthService)

    return {
      user: new UserEntity(
        // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
        process.env.ADMIN_USERNAME!,
        // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
        process.env.ADMIN_PASSWORD!,
        Role.Admin
      ),
      tokens: await authService.getAdminAccessToken()
    }
  }

  // eslint-disable-next-line max-lines-per-function
  private async seedDatabase(): Promise<void> {
    const users = [
      new UserEntity("Admin", "testpassword", Role.Admin),
      new UserEntity("Manager", "testpassword", Role.Manager),
      new UserEntity("User", "testpassword")
    ]
    for (let i = 0; i < 20; i++) {
      const newUser = new UserEntity(`TestUser${i}`, "testpassword")
      users.push(newUser)
    }

    await this.em.persistAndFlush(users)

    const managerUsers: UserEntity[] = []
    for (let i = 0; i < 10; i++) {
      const newManagerUser = new UserEntity(
        `ManagerUser${i}`,
        "testpassword",
        Role.Manager
      )
      managerUsers.push(newManagerUser)
    }

    await this.em.persistAndFlush(managerUsers)

    const locations = [new LocationEntity("TestLocation", 20, users[2])]
    for (let i = 0; i < 20; i++) {
      const newLocation = new LocationEntity(
        `TestLocation${i}`,
        20,
        users[i + 3]
      )
      locations.push(newLocation)
    }

    await this.em.persistAndFlush(locations)

    const schools: SchoolEntity[] = []
    for (let i = 0; i < 20; i++) {
      const newSchool = new SchoolEntity(`TestSchool${i}`)
      schools.push(newSchool)
    }

    await this.em.persistAndFlush(schools)

    const checkIns: CheckInEntity[] = []
    for (let i = 0; i < 5; i++) {
      checkIns.push(new CheckInEntity(locations[0], schools[0]))
    }

    await this.em.persistAndFlush(checkIns)

    locations[0].capacity++

    await this.em.flush()

    for (let i = 0; i < 5; i++) {
      checkIns.push(new CheckInEntity(locations[0], schools[0]))
    }

    await this.em.persistAndFlush(checkIns)

    this.em.clear()
  }
}
