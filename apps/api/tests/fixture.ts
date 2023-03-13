import { INestApplication } from "@nestjs/common"
import { Test } from "@nestjs/testing"
import { AppModule } from "../src/app.module"
import cookieParser from "cookie-parser"
import { AuthService } from "../src/auth/auth.service"
import { LocationEntity, SchoolEntity, UserEntity } from "mikro-orm-config"
import { MikroORM } from "@mikro-orm/core"
import { Role, Tokens } from "types-custom"
import { EntityManager, PostgreSqlDriver } from "@mikro-orm/postgresql"

export interface TokensAndUser {
  tokens: Tokens,
  user: UserEntity
}

export interface RequestHeaders {
  "set-cookie": string
}

export interface ErrorResponse {
  message: string
}

export default class Fixture {
  app: INestApplication
  em: EntityManager


  async init(): Promise<void> {
    const module = await Test.createTestingModule({
      imports: [AppModule],
    }).compile()

    this.app = module.createNestApplication()
    this.app.use(cookieParser())
    this.app.enableShutdownHooks()

    await this.app.init()

    const orm = this.app.get<MikroORM<PostgreSqlDriver>>(MikroORM)
    this.em = orm.em.fork()
  }

  async seedDatabase(): Promise<void> {
    const users = [
      new UserEntity("Admin", "testpassword", Role.Admin),
      new UserEntity("User", "testpassword")
    ]

    for (let i = 0; i < 20; i++){
      const newUser = new UserEntity(`TestUser${i}`, "testpassword")
      users.push(newUser)
    }

    await this.em.persistAndFlush(users)

    const locations = [
      new LocationEntity("TestLocation", 20, users[1])
    ]

    for (let i = 0; i < 20; i++){
      const newLocation = new LocationEntity(`TestLocation${i}`, 20, users[i + 2])
      locations.push(newLocation)
    }

    await this.em.persistAndFlush(locations)

    const schools = []

    for (let i = 0; i < 20; i++){
      const newSchool = new SchoolEntity(`TestSchool${i}`)
      schools.push(newSchool)
    }

    await this.em.persistAndFlush(schools)
  }

  async clearDatabase(): Promise<void> {
    await this.em.nativeDelete(SchoolEntity, { id: { $gt: 1 } })
    await this.em.nativeDelete(LocationEntity, {})
    await this.em.nativeDelete(UserEntity, {})
  }

  async getTokens(username: string): Promise<TokensAndUser> {
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
}