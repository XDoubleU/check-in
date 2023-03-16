import { expect } from "chai"
import request from "supertest"
import { type SignInDto } from "types-custom"
import Fixture, { type ErrorResponse, type RequestHeaders, type TokensAndUser } from "./fixture"


describe("AuthController (e2e)", () => {
  let fixture: Fixture

  let tokensAndUser: TokensAndUser
  
  before(() => {
    fixture = new Fixture()
    return fixture.init()
      .then(() => fixture.seedDatabase())
      .then(() => fixture.getTokens("User"))
      .then((data) => tokensAndUser = data)
  })

  after(() => {
    return fixture.clearDatabase()
      .then(() => fixture.app.close())
  })

  describe("/auth/signin (POST)", () => {
    it("signs in user (200)", async () => {
      const data: SignInDto = {
        username: tokensAndUser.user.username,
        password: "testpassword"
      }
      
      const response = await request(fixture.app.getHttpServer())
        .post("/auth/signin")
        .send(data)
        .expect(200)
      
      const responseHeaders = response.headers as RequestHeaders
      expect(responseHeaders["set-cookie"][0]).to.contain("accessToken")
      expect(responseHeaders["set-cookie"][1]).to.contain("refreshToken")
    })

    it("returns Invalid credentials (401)", async () => {
      const data: SignInDto = {
        username: "inexistentuser",
        password: "testpassword"
      }
      
      const response = await request(fixture.app.getHttpServer())
        .post("/auth/signin")
        .send(data)
        .expect(401)
      
      const errorResponse = response.body as ErrorResponse
      expect(errorResponse.message).to.be.equal("Invalid credentials")
    })
  })

  describe("/auth/signout (GET)", () => {
    it("signs out user (200)", async () => {      
      const response = await request(fixture.app.getHttpServer())
        .get("/auth/signout")
        .set("Cookie", [`accessToken=${tokensAndUser.tokens.accessToken}`])
        .expect(200)
      
      const responseHeaders = response.headers as RequestHeaders
      expect(responseHeaders["set-cookie"][0]).to.contain("accessToken=;")
      expect(responseHeaders["set-cookie"][1]).to.contain("refreshToken=;")
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
      expect(responseHeaders["set-cookie"][0]).to.contain("accessToken")
      expect(responseHeaders["set-cookie"][1]).to.contain("refreshToken")
    })

    it("returns unauthorized (401)", async () => {
      return await request(fixture.app.getHttpServer())
        .get("/auth/refresh")
        .expect(401)
    })
  })
})
