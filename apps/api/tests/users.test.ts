import request from "supertest"
import { type User } from "types-custom"
import { type UserAndTokens } from "../src/auth/auth.service"
import Fixture from "./config/fixture"

describe("UsersController (e2e)", () => {
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

  describe("/users/me (GET)", () => {
    it("gets info about logged in User (200)", async () => {
      const response = await request(fixture.app.getHttpServer())
        .get("/users/me")
        .set("Cookie", [`accessToken=${userAndTokens.tokens.accessToken}`])
        .expect(200)

      const userResponse = response.body as User
      expect(userResponse.id).toBe(userAndTokens.user.id)
      expect(userResponse.username).toBe(userAndTokens.user.username)
      expect(userResponse.roles).toStrictEqual(userAndTokens.user.roles)
      expect(userResponse.location?.id).toBe(userAndTokens.user.location?.id)
    })
  })
})
