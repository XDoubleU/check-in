import { INestApplication } from "@nestjs/common"
import { Test, TestingModule } from "@nestjs/testing"
import { AppModule } from "../src/app.module"
import cookieParser from "cookie-parser"
import { AuthService } from "../src/auth/auth.service"
import { UsersService } from "../src/users/users.service"
import { User } from "types"
import { LocationsService } from "../src/locations/locations.service"
import { SchoolsService } from "../src/schools/schools.service"
import { CheckInsService } from "../src/checkins/checkins.service"

export async function getApp(): Promise<INestApplication> {
  const moduleFixture: TestingModule = await Test.createTestingModule({
    imports: [AppModule],
  }).compile()

  const app = moduleFixture.createNestApplication()

  app.use(cookieParser())

  return app
}

export async function getAccessToken(app: INestApplication): Promise<string> {
  const usersService = app.get<UsersService>(UsersService)

  let user = await usersService.getByUserName("AuthUser")
  if (!user) {
    user = await usersService.create("AuthUser", "AuthUserPassword") as User
  }

  const authService = app.get<AuthService>(AuthService)
  return (await authService.getTokens(user)).accessToken
}

export async function clearDatabase(app: INestApplication): Promise<void> {
  const usersService = app.get<UsersService>(UsersService)
  const checkinsService = app.get<CheckInsService>(CheckInsService)
  const locationsService = app.get<LocationsService>(LocationsService)
  const schoolsService = app.get<SchoolsService>(SchoolsService)

  await usersService.user.deleteMany()
  await checkinsService.checkIn.deleteMany()
  await locationsService.location.deleteMany()
  await schoolsService.school.deleteMany({
    where: {
      id: {
        not: 1
      }
    }
  })
}