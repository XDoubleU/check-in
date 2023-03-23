/* eslint-disable max-lines-per-function */
import request from "supertest"
import { type CheckIn, type CreateCheckInDto } from "types-custom"
import Fixture, { type ErrorResponse, type TokensAndUser } from "./fixture"
import { LocationEntity, SchoolEntity } from "mikro-orm-config"
import { v4 } from "uuid"

describe("CheckInsController (e2e)", () => {
  let fixture: Fixture

  let tokensAndUser: TokensAndUser
  let adminTokensAndUser: TokensAndUser

  let location: LocationEntity
  let school: SchoolEntity

  beforeEach(() => {
    fixture = new Fixture()
    return fixture
      .init()
      .then(() => fixture.seedDatabase())
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
    return fixture.clearDatabase().then(() => fixture.app.close())
  })

  describe("/checkins (POST)", () => {
    it("creates a new CheckIn (201)", async () => {
      const data: CreateCheckInDto = {
        locationId: location.id,
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

    it("returns Location not found (404)", async () => {
      const data: CreateCheckInDto = {
        locationId: v4(),
        schoolId: school.id
      }

      const response = await request(fixture.app.getHttpServer())
        .post("/checkins")
        .set("Cookie", [`accessToken=${tokensAndUser.tokens.accessToken}`])
        .send(data)
        .expect(404)

      const errorResponse = response.body as ErrorResponse
      expect(errorResponse.message).toBe("Location not found")
    })

    it("returns School not found (404)", async () => {
      const data: CreateCheckInDto = {
        locationId: location.id,
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

    // TODO: test user should own location where checkin is made

    it("returns Forbidden (403)", async () => {
      const data: CreateCheckInDto = {
        locationId: location.id,
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
