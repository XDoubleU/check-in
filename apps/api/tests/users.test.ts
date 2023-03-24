import request from "supertest"
import { type User } from "types-custom"
import Fixture, { type TokensAndUser } from "./config/fixture"

describe("UsersController (e2e)", () => {
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

  describe("/users/me (GET)", () => {
    it("gets info about logged in User (200)", async () => {
      const response = await request(fixture.app.getHttpServer())
        .get("/users/me")
        .set("Cookie", [`accessToken=${tokensAndUser.tokens.accessToken}`])
        .expect(200)

      const userResponse = response.body as User
      expect(userResponse.id).toBe(tokensAndUser.user.id)
      expect(userResponse.username).toBe(tokensAndUser.user.username)
      expect(userResponse.roles).toStrictEqual(tokensAndUser.user.roles)
      expect(userResponse.location?.id).toBe(tokensAndUser.user.location?.id)
    })
  })
})
