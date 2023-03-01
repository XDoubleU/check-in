import { INestApplication } from "@nestjs/common"
import request from "supertest"
import { User } from "types-custom"
import { LocationsService } from "../src/locations/locations.service"
import { UsersService } from "../src/users/users.service"
import { clearDatabase, getApp } from "./shared"
import { AuthService } from "../src/auth/auth.service"


describe("UsersController (e2e)", () => {
  let app: INestApplication
  let accessToken: string

  let usersService: UsersService
  let locationsService: LocationsService
  let authService: AuthService

  let user: User
  
  beforeAll(async () => {
    app = await getApp()

    usersService = app.get<UsersService>(UsersService)
    locationsService = app.get<LocationsService>(LocationsService)
    authService = app.get<AuthService>(AuthService)

    await app.init()
  })

  beforeEach(async () => {
    // UsersService
    user = await usersService.create("TestUser", "testpassword")

    // LocationsService
    await locationsService.create("TestLocation", 10, user)
    const tempUser = await usersService.getById(user.id)
    if (tempUser) {
      user = tempUser
    }

    // AuthService
    accessToken = (await authService.getTokens(user)).accessToken
  })

  afterEach(async () => {
    await clearDatabase(app)
  })

  describe("/users/me (GET)", () => {
    it("gets info about logged in User (200)", async () => {  
      const response = await request(app.getHttpServer())
        .get("/users/me")
        .set("Cookie", [`accessToken=${accessToken}`])
        .expect(200)
      
      const userResponse = response.body as User
      expect(userResponse.id).toBe(user.id)
      expect(userResponse.username).toBe(user.username)
      expect(userResponse.roles).toStrictEqual(user.roles)
      expect(userResponse.locationId).toBe(user.locationId)
    })
  })
})
