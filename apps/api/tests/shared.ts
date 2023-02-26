import { INestApplication } from "@nestjs/common"
import { Test, TestingModule } from "@nestjs/testing"
import { AppModule } from "../src/app.module"
import cookieParser from "cookie-parser"
import { AuthService } from "../src/auth/auth.service"
import { UsersService } from "../src/users/users.service"
import { Role, User } from "types"
import { LocationsService } from "../src/locations/locations.service"
import { SchoolsService } from "../src/schools/schools.service"
import { CheckInsService } from "../src/checkins/checkins.service"
import { hashSync } from "bcrypt"

type GetAccessTokenReturn = {
  accessToken: string,
  user: User
}

export async function getApp(): Promise<INestApplication> {
  const moduleFixture: TestingModule = await Test.createTestingModule({
    imports: [AppModule],
  }).compile()

  const app = moduleFixture.createNestApplication()

  app.use(cookieParser())

  return app
}

export async function getAdminAccessToken(app: INestApplication): Promise<string> {
  const usersService = app.get<UsersService>(UsersService)

  let user = await usersService.getByUserName("AdminUser")
  if (!user) {
    user = await usersService.user.create({
      data: {
        username: "AdminUser",
        passwordHash: hashSync("AdminPassword", 12),
        roles: [Role.Admin]
      }
    })
  }

  const authService = app.get<AuthService>(AuthService)
  return (await authService.getTokens(user)).accessToken
}

export async function getAccessToken(app: INestApplication): Promise<GetAccessTokenReturn> {
  const usersService = app.get<UsersService>(UsersService)
  const locationsService = app.get<LocationsService>(LocationsService)

  let user = await usersService.getByUserName("AuthUser")
  if (!user) {
    user = await usersService.create("AuthUser", "AuthUserPassword") as User
  }
  await locationsService.create("AuthUserLocation", 10, user)
  user = await usersService.getById(user.id) as User

  const authService = app.get<AuthService>(AuthService)

  return {
    accessToken: (await authService.getTokens(user)).accessToken,
    user
  }
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