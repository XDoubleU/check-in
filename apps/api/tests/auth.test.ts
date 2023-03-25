/* eslint-disable sonarjs/no-duplicate-string */
/* eslint-disable max-lines-per-function */
import request from "supertest"
import { type SignInDto } from "types-custom"
import Fixture, {
  type ErrorResponse,
  type RequestHeaders,
  type TokensAndUser
} from "./config/fixture"

describe("AuthController (e2e)", () => {
  const fixture: Fixture = new Fixture()

  let tokensAndUser: TokensAndUser

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
      .then((data) => (tokensAndUser = data))
  })

  afterEach(() => {
    return fixture.afterEach()
  })

  describe("/auth/signin (POST)", () => {
    it("signs in user (200)", async () => {
      const data: SignInDto = {
        username: tokensAndUser.user.username,
        password: "testpassword",
        rememberMe: true
      }

      const response = await request(fixture.app.getHttpServer())
        .post("/auth/signin")
        .send(data)
        .expect(200)

      const responseHeaders = response.headers as RequestHeaders
      expect(responseHeaders["set-cookie"][0]).toContain("accessToken")
      expect(responseHeaders["set-cookie"][1]).toContain("refreshToken")
    })

    it("returns Invalid credentials (401)", async () => {
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
  })

  describe("/auth/signout (GET)", () => {
    it("signs out user (200)", async () => {
      const response = await request(fixture.app.getHttpServer())
        .get("/auth/signout")
        .set("Cookie", [`accessToken=${tokensAndUser.tokens.accessToken}`])
        .expect(200)

      const responseHeaders = response.headers as RequestHeaders
      expect(responseHeaders["set-cookie"][0]).toContain("accessToken=;")
      expect(responseHeaders["set-cookie"][1]).toContain("refreshToken=;")
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
        .set("Cookie", [`refreshToken=${tokensAndUser.tokens.refreshToken}`])
        .expect(200)

      const responseHeaders = response.headers as RequestHeaders
      expect(responseHeaders["set-cookie"][0]).toContain("accessToken")
      expect(responseHeaders["set-cookie"][1]).toContain("refreshToken")
    })

    it("returns unauthorized (401)", async () => {
      return await request(fixture.app.getHttpServer())
        .get("/auth/refresh")
        .expect(401)
    })
  })
})
