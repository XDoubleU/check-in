import { expect } from "chai"
import request from "supertest"
import { type User } from "types-custom"
import Fixture, { type TokensAndUser } from "./fixture"

describe("UsersController (e2e)", () => {
  let fixture: Fixture

  let tokensAndUser: TokensAndUser

  before(() => {
    fixture = new Fixture()
    return fixture
      .init()
      .then(() => fixture.seedDatabase())
      .then(() => fixture.getTokens("User"))
      .then((data) => (tokensAndUser = data))
  })

  after(() => {
    return fixture.clearDatabase().then(() => fixture.app.close())
  })

  describe("/users/me (GET)", () => {
    it("gets info about logged in User (200)", async () => {
      const response = await request(fixture.app.getHttpServer())
        .get("/users/me")
        .set("Cookie", [`accessToken=${tokensAndUser.tokens.accessToken}`])
        .expect(200)

      const userResponse = response.body as User
      expect(userResponse.id).to.be.equal(tokensAndUser.user.id)
      expect(userResponse.username).to.be.equal(tokensAndUser.user.username)
      expect(userResponse.roles).to.deep.equal(tokensAndUser.user.roles)
      expect(userResponse.location?.id).to.be.equal(
        tokensAndUser.user.location?.id
      )
    })
  })
})
