import request from "supertest"
import { type User } from "types-custom"
import Fixture, { type TokensAndUser } from "./fixture"

describe("UsersController (e2e)", () => {
  let fixture: Fixture

  let tokensAndUser: TokensAndUser

  beforeEach(() => {
    fixture = new Fixture()
    return fixture
      .init()
      .then(() => fixture.seedDatabase())
      .then(() => fixture.getTokens("User"))
      .then((data) => (tokensAndUser = data))
  })

  afterEach(() => {
    return fixture.clearDatabase().then(() => fixture.app.close())
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
