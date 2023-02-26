import { INestApplication } from "@nestjs/common"
import request from "supertest"
import { SignInDto, User } from "types"
import { LocationsService } from "../src/locations/locations.service"
import { UsersService } from "../src/users/users.service"
import { clearDatabase, getApp } from "./shared"
import { AuthService } from "../src/auth/auth.service"


describe("AuthController (e2e)", () => {
  let app: INestApplication

  let usersService: UsersService
  let locationsService: LocationsService
  let authService: AuthService

  let user: User
  let accessToken: string
  let refreshToken: string
  
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
    const tokens = (await authService.getTokens(user))
    accessToken = tokens.accessToken
    refreshToken = tokens.refreshToken
  })

  afterEach(async () => {
    await clearDatabase(app)
  })

  describe("/auth/signin (POST)", () => {
    it("signs in user (200)", async () => {
      const data: SignInDto = {
        username: user.username,
        password: "testpassword"
      }
      
      const response = await request(app.getHttpServer())
        .post("/auth/signin")
        .send(data)
        .expect(200)
      
      expect(response.headers["set-cookie"][0]).toContain("accessToken")
      expect(response.headers["set-cookie"][1]).toContain("refreshToken")
    })

    it("returns Invalid credentials (401)", async () => {
      const data: SignInDto = {
        username: "inexistentuser",
        password: "testpassword"
      }
      
      const response = await request(app.getHttpServer())
        .post("/auth/signin")
        .send(data)
        .expect(401)
      
      expect(response.body.message).toBe("Invalid credentials")
    })
  })

  describe("/auth/signout (GET)", () => {
    it("signs out user (200)", async () => {      
      const response = await request(app.getHttpServer())
        .get("/auth/signout")
        .set("Cookie", [`accessToken=${accessToken}`])
        .expect(200)
      
      expect(response.headers["set-cookie"][0]).toContain("accessToken=;")
      expect(response.headers["set-cookie"][1]).toContain("refreshToken=;")
    })

    it("returns unauthorized (401)", async () => {      
      return await request(app.getHttpServer())
        .get("/auth/signout")
        .expect(401)
    })
  })

  describe("/auth/refresh (GET)", () => {
    it("refreshes users tokens (200)", async () => {
      const response = await request(app.getHttpServer())
        .get("/auth/refresh")
        .set("Cookie", [`refreshToken=${refreshToken}`])
        .expect(200)
      
      expect(response.headers["set-cookie"][0]).toContain("accessToken")
      expect(response.headers["set-cookie"][1]).toContain("refreshToken")
    })

    it("returns unauthorized (401)", async () => {
      return await request(app.getHttpServer())
        .get("/auth/refresh")
        .expect(401)
    })
  })
})
