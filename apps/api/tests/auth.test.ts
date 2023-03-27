/* eslint-disable sonarjs/no-duplicate-string */
/* eslint-disable max-lines-per-function */
import request from "supertest"
import { type SignInDto } from "types-custom"
import { type UserAndTokens } from "../src/auth/auth.service"
import Fixture, { type ErrorResponse } from "./config/fixture"

describe("AuthController (e2e)", () => {
  const fixture: Fixture = new Fixture()

  let userAndTokens: UserAndTokens

  beforeAll(() => {
    return fixture.beforeAll()
  })

  afterAll(() => {
    return fixture.afterAll()
  })

  beforeEach(() => {
    return fixture
      .beforeEach()
      .then(() => fixture.getTokens("User"))
      .then((data) => (userAndTokens = data))
  })

  afterEach(() => {
    return fixture.afterEach()
  })

  describe("/auth/signin (POST)", () => {
    it("signs in user (200)", async () => {
      const data: SignInDto = {
        username: userAndTokens.user.username,
        password: "testpassword",
        rememberMe: true
      }

      const response = await request(fixture.app.getHttpServer())
        .post("/auth/signin")
        .send(data)
        .expect(200)

      // eslint-disable-next-line @typescript-eslint/no-unsafe-member-access
      expect(response.headers["set-cookie"][0]).toContain("accessToken")
      // eslint-disable-next-line @typescript-eslint/no-unsafe-member-access
      expect(response.headers["set-cookie"][1]).toContain("refreshToken")
    })

    it("returns Invalid credentials because of inexistent user (401)", async () => {
      const data: SignInDto = {
        username: "inexistentuser",
        password: "testpassword",
        rememberMe: true
      }

      const response = await request(fixture.app.getHttpServer())
        .post("/auth/signin")
        .send(data)
        .expect(401)

      const errorResponse = response.body as ErrorResponse
      expect(errorResponse.message).toBe("Invalid credentials")
    })

    it("returns Invalid credentials because of wrong password (401)", async () => {
      const data: SignInDto = {
        username: userAndTokens.user.username,
        password: "wrongpassword",
        rememberMe: true
      }

      const response = await request(fixture.app.getHttpServer())
        .post("/auth/signin")
        .send(data)
        .expect(401)

      const errorResponse = response.body as ErrorResponse
      expect(errorResponse.message).toBe("Invalid credentials")
    })

    it("returns Internal server error exception because of missing JWT config (500)", async () => {
      const temp = process.env.JWT_ACCESS_SECRET
      process.env.JWT_ACCESS_SECRET = ""

      const data: SignInDto = {
        username: userAndTokens.user.username,
        password: "testpassword",
        rememberMe: true
      }

      const response = await request(fixture.app.getHttpServer())
        .post("/auth/signin")
        .send(data)
        .expect(500)

      const errorResponse = response.body as ErrorResponse
      expect(errorResponse.message).toBe(
        "JWT secrets or expirations missing in environment"
      )

      process.env.JWT_ACCESS_SECRET = temp
    })
  })

  describe("/auth/signout (GET)", () => {
    it("signs out user (200)", async () => {
      const response = await request(fixture.app.getHttpServer())
        .get("/auth/signout")
        .set("Cookie", [`accessToken=${userAndTokens.tokens.accessToken}`])
        .expect(200)

      // eslint-disable-next-line @typescript-eslint/no-unsafe-member-access
      expect(response.headers["set-cookie"][0]).toContain("accessToken=;")
      // eslint-disable-next-line @typescript-eslint/no-unsafe-member-access
      expect(response.headers["set-cookie"][1]).toContain("refreshToken=;")
    })

    it("returns unauthorized (401)", async () => {
      return await request(fixture.app.getHttpServer())
        .get("/auth/signout")
        .expect(401)
    })
  })

  describe("/auth/refresh (GET)", () => {
    it("refreshes users tokens (200)", async () => {
      const response = await request(fixture.app.getHttpServer())
        .get("/auth/refresh")
        .set("Cookie", [`refreshToken=${userAndTokens.tokens.refreshToken}`])
        .expect(200)

      // eslint-disable-next-line @typescript-eslint/no-unsafe-member-access
      expect(response.headers["set-cookie"][0]).toContain("accessToken")
      // eslint-disable-next-line @typescript-eslint/no-unsafe-member-access
      expect(response.headers["set-cookie"][1]).toContain("refreshToken")
    })

    it("returns unauthorized (401)", async () => {
      return await request(fixture.app.getHttpServer())
        .get("/auth/refresh")
        .expect(401)
    })
  })
})
