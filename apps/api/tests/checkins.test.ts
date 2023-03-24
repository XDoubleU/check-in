/* eslint-disable max-lines-per-function */
import request from "supertest"
import { type CheckIn, type CreateCheckInDto } from "types-custom"
import Fixture, {
  type ErrorResponse,
  type TokensAndUser
} from "./config/fixture"
import { LocationEntity, SchoolEntity } from "mikro-orm-config"

describe("CheckInsController (e2e)", () => {
  const fixture: Fixture = new Fixture()

  let tokensAndUser: TokensAndUser
  let adminTokensAndUser: TokensAndUser

  let location: LocationEntity
  let school: SchoolEntity

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
      .then(() => fixture.getTokens("Admin"))
      .then((data) => (adminTokensAndUser = data))
      .then(() => fixture.em.findOne(SchoolEntity, { id: 1 }))
      .then((data) => {
        if (!data) {
          throw new Error("School is null")
        }

        school = data
      })
      .then(() => fixture.em.find(LocationEntity, {}))
      .then((data) => {
        location = data[0]
      })
  })

  afterEach(() => {
    return fixture.afterEach()
  })

  describe("/checkins (POST)", () => {
    it("creates a new CheckIn (201)", async () => {
      const data: CreateCheckInDto = {
        schoolId: school.id
      }

      const response = await request(fixture.app.getHttpServer())
        .post("/checkins")
        .set("Cookie", [`accessToken=${tokensAndUser.tokens.accessToken}`])
        .send(data)
        .expect(201)

      const responseCheckIn = response.body as CheckIn
      expect(responseCheckIn.id).toBeDefined()
      expect(responseCheckIn.location.id).toBe(location.id)
      expect(responseCheckIn.capacity).toBe(location.capacity)
      expect(responseCheckIn.createdAt).toBeDefined()
      expect(responseCheckIn.school.id).toBe(school.id)
    })

    it("returns School not found (404)", async () => {
      const data: CreateCheckInDto = {
        schoolId: -1
      }

      const response = await request(fixture.app.getHttpServer())
        .post("/checkins")
        .set("Cookie", [`accessToken=${tokensAndUser.tokens.accessToken}`])
        .send(data)
        .expect(404)

      const errorResponse = response.body as ErrorResponse
      expect(errorResponse.message).toBe("School not found")
    })

    it("returns Forbidden (403)", async () => {
      const data: CreateCheckInDto = {
        schoolId: school.id
      }

      return await request(fixture.app.getHttpServer())
        .post("/checkins")
        .set("Cookie", [`accessToken=${adminTokensAndUser.tokens.accessToken}`])
        .send(data)
        .expect(403)
    })
  })
})
