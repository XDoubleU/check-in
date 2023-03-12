import request from "supertest"
import { CheckIn, CreateCheckInDto } from "types-custom"
import Fixture, { ErrorResponse, TokensAndUser } from "./fixture"
import { LocationEntity, SchoolEntity } from "mikro-orm-config"
import { expect } from "chai"
import { v4 } from "uuid"


describe("CheckInsController (e2e)", () => {
  let fixture: Fixture

  let tokensAndUser: TokensAndUser
  let adminTokensAndUser: TokensAndUser

  let location: LocationEntity
  let school: SchoolEntity
  
  before(() => {
    fixture = new Fixture()
    return fixture.init()
      .then(() => fixture.seedDatabase())
      .then(() => fixture.getTokens("User"))
      .then((data) => tokensAndUser = data)
      .then(() => fixture.getTokens("Admin"))
      .then((data) => adminTokensAndUser = data)
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

  after(() => {
    return fixture.clearDatabase()
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
      expect(responseCheckIn.id).to.exist
      expect(responseCheckIn.location.id).to.be.equal(location.id)
      expect(responseCheckIn.capacity).to.be.equal(location.capacity)
      expect(responseCheckIn.createdAt).to.exist
      expect(responseCheckIn.school.id).to.be.equal(school.id)
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
      expect(errorResponse.message).to.be.equal("Location not found")
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
      expect(errorResponse.message).to.be.equal("School not found")
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
