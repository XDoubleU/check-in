import { INestApplication } from "@nestjs/common"
import request from "supertest"
import { User } from "types"
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
    user = await usersService.getById(user.id) as User

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
      
      expect(response.body.id).toBe(user.id)
      expect(response.body.username).toBe(user.username)
      expect(response.body.roles).toStrictEqual(user.roles)
      expect(response.body.passwordHash).toBeUndefined()
      expect(response.body.locationId).toBe(user.locationId)
    })
  })
})
