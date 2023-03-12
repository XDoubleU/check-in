import request from "supertest"
import { CreateSchoolDto, GetAllPaginatedSchoolDto, School, UpdateSchoolDto } from "types-custom"
import Fixture, { ErrorResponse, TokensAndUser } from "./fixture"
import { CheckInEntity, SchoolEntity } from "mikro-orm-config"
import { expect } from "chai"


describe("SchoolsController (e2e)", () => {
  let fixture: Fixture

  let tokensAndUser: TokensAndUser
  let adminTokensAndUser: TokensAndUser

  let schools: SchoolEntity[]

  const defaultPage = 1
  const defaultPageSize = 4

  before(() => {
    fixture = new Fixture()
    return fixture.init()
      .then(() => fixture.seedDatabase())
      .then(() => fixture.getTokens("User"))
      .then((data) => tokensAndUser = data)
      .then(() => fixture.getTokens("Admin"))
      .then((data) => adminTokensAndUser = data)
      .then(() => fixture.em.find(SchoolEntity, {}))
      .then((data) => {
        schools = data
      })
  })

  after(() => {
    return fixture.clearDatabase()
  })

  describe("/schools/all (GET)", () => {
    it("gets all Schools (200)", async () => {
      const location = tokensAndUser.user.location
      console.log(tokensAndUser.user)
      if (!location) {
        throw new Error("Location is undefined")
      }

      const andere = await fixture.em.findOne(SchoolEntity, 1)
      if (!andere) {
        throw new Error("andere is undefined")
      }

      const school = schools[5]

      for (let i = 0; i < 10; i++) {
        const newCheckIn = new CheckInEntity(location, andere)
        await fixture.em.persistAndFlush(newCheckIn)
      }

      for (let i = 0; i < 15; i++) {
        const newCheckIn = new CheckInEntity(location, school)
        await fixture.em.persistAndFlush(newCheckIn)
      }

      const response = await request(fixture.app.getHttpServer())
        .get("/schools/all")
        .set("Cookie", [`accessToken=${tokensAndUser.tokens.accessToken}`])
        .expect(200)

      const schoolsResponse = response.body as School[]
      expect(schoolsResponse).to.be.a("Array<School>")
      expect(schoolsResponse[schoolsResponse.length - 1]).to.deep.equal(andere)
      expect(schoolsResponse[0]).to.deep.equal(school)
    })

    it("returns Forbidden (403)", async () => {
      return await request(fixture.app.getHttpServer())
        .get("/schools/all")
        .set("Cookie", [`accessToken=${adminTokensAndUser.tokens.accessToken}`])
        .expect(403)
    })
  })

  describe("/schools (GET)", () => {
    it("gets all Schools with default page (200)", async () => {
      const response = await request(fixture.app.getHttpServer())
        .get("/schools")
        .set("Cookie", [`accessToken=${adminTokensAndUser.tokens.accessToken}`])
        .expect(200)

      const paginatedSchoolsResponse = response.body as GetAllPaginatedSchoolDto
      expect(paginatedSchoolsResponse.page).to.be.equal(defaultPage)
      expect(paginatedSchoolsResponse.totalPages).to.be.equal(Math.ceil(schools.length / defaultPageSize))
      expect(paginatedSchoolsResponse.schools.length).to.be.equal(defaultPageSize)
    })

    it("gets certain page of all Schools (200)", async () => {
      const page = 2

      const response = await request(fixture.app.getHttpServer())
        .get("/schools")
        .query({ page })
        .set("Cookie", [`accessToken=${adminTokensAndUser.tokens.accessToken}`])
        .expect(200)
      
      const paginatedSchoolsResponse = response.body as GetAllPaginatedSchoolDto
      expect(paginatedSchoolsResponse.page).to.be.equal(page)
      expect(paginatedSchoolsResponse.totalPages).to.be.equal(Math.ceil(schools.length / defaultPageSize))
      expect(paginatedSchoolsResponse.schools.length).to.be.equal(defaultPageSize)
    })

    it("returns Forbidden (403)", async () => {
      return await request(fixture.app.getHttpServer())
        .get("/schools")
        .set("Cookie", [`accessToken=${tokensAndUser.tokens.accessToken}`])
        .expect(403)
    })
  })

  describe("/schools (POST)", () => {
    it("creates a new School (201)", async () => {
      const data: CreateSchoolDto = {
        name: "NewSchool"
      }
  
      const response = await request(fixture.app.getHttpServer())
        .post("/schools")
        .set("Cookie", [`accessToken=${adminTokensAndUser.tokens.accessToken}`])
        .send(data)
        .expect(201)
      
      const schoolResponse = response.body as School
      expect(schoolResponse.id).to.exist
      expect(schoolResponse.name).to.be.equal(data.name)
    })

    it("returns School with this name already exists (409)", async () => {
      const data: CreateSchoolDto = {
        name: schools[1].name
      }
  
      const response = await request(fixture.app.getHttpServer())
        .post("/schools")
        .set("Cookie", [`accessToken=${adminTokensAndUser.tokens.accessToken}`])
        .send(data)
        .expect(409)
      
      const errorResponse = response.body as ErrorResponse
      expect(errorResponse.message).to.be.equal("School with this name already exists")
    })

    it("returns Forbidden (403)", async () => {
      const data: CreateSchoolDto = {
        name: "NewSchool"
      }
  
      return await request(fixture.app.getHttpServer())
        .post("/schools")
        .set("Cookie", [`accessToken=${tokensAndUser.tokens.accessToken}`])
        .send(data)
        .expect(403)
    })
  })

  describe("/schools/:id (PATCH)", () => {
    it("updates a new School (200)", async () => {
      const id = schools[1].id
      const data: UpdateSchoolDto = {
        name: "NewSchool2"
      }
  
      const response = await request(fixture.app.getHttpServer())
        .patch(`/schools/${id}`)
        .set("Cookie", [`accessToken=${adminTokensAndUser.tokens.accessToken}`])
        .send(data)
        .expect(200)
      
      const schoolResponse = response.body as School
      expect(schoolResponse.id).to.be.equal(id)
      expect(schoolResponse.name).to.be.equal(data.name)
    })

    it("returns School with this name already exists (409)", async () => {
      const id = schools[1].id
      const data: UpdateSchoolDto = {
        name: schools[2].name
      }
  
      const response = await request(fixture.app.getHttpServer())
        .patch(`/schools/${id}`)
        .set("Cookie", [`accessToken=${adminTokensAndUser.tokens.accessToken}`])
        .send(data)
        .expect(409)
      
      const errorResponse = response.body as ErrorResponse
      expect(errorResponse.message).to.be.equal("School with this name already exists")
    })

    it("returns School not found (404)", async () => {
      const id = -1
      const data: UpdateSchoolDto = {
        name: "NewSchool2"
      }
  
      const response = await request(fixture.app.getHttpServer())
        .patch(`/schools/${id}`)
        .set("Cookie", [`accessToken=${adminTokensAndUser.tokens.accessToken}`])
        .send(data)
        .expect(404)
      
      const errorResponse = response.body as ErrorResponse
      expect(errorResponse.message).to.be.equal("School not found")
    })

    it("returns Forbidden (403)", async () => {
      const id = schools[1].id
      const data: UpdateSchoolDto = {
        name: "NewSchool2"
      }
  
      return await request(fixture.app.getHttpServer())
        .patch(`/schools/${id}`)
        .set("Cookie", [`accessToken=${tokensAndUser.tokens.accessToken}`])
        .send(data)
        .expect(403)
    })
  })

  describe("/schools/:id (DELETE)", () => {
    it("deletes a School (200)", async () => {
      const id = schools[1].id
  
      const response = await request(fixture.app.getHttpServer())
        .delete(`/schools/${id}`)
        .set("Cookie", [`accessToken=${adminTokensAndUser.tokens.accessToken}`])
        .expect(200)
      
      const schoolResponse = response.body as School
      expect(schoolResponse.id).to.be.equal(id)
    })

    it("returns School not found (404)", async () => {
      const id = -1
  
      const response = await request(fixture.app.getHttpServer())
        .delete(`/schools/${id}`)
        .set("Cookie", [`accessToken=${adminTokensAndUser.tokens.accessToken}`])
        .expect(404)
      
      const errorResponse = response.body as ErrorResponse
      expect(errorResponse.message).to.be.equal("School not found")
    })

    it("returns Forbidden (403)", async () => {
      const id = schools[1].id
  
      return await request(fixture.app.getHttpServer())
        .delete(`/schools/${id}`)
        .set("Cookie", [`accessToken=${tokensAndUser.tokens.accessToken}`])
        .expect(403)
    })
  })
})
