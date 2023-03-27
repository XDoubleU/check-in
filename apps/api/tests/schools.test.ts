/* eslint-disable max-lines-per-function */
/* eslint-disable sonarjs/no-duplicate-string */
import request from "supertest"
import {
  type CreateSchoolDto,
  type GetAllPaginatedSchoolDto,
  type School,
  type UpdateSchoolDto
} from "types-custom"
import Fixture, { type ErrorResponse } from "./config/fixture"
import { CheckInEntity, SchoolEntity } from "mikro-orm-config"
import { type UserAndTokens } from "../src/auth/auth.service"

describe("SchoolsController (e2e)", () => {
  const fixture: Fixture = new Fixture()

  let userAndTokens: UserAndTokens
  let adminUserAndTokens: UserAndTokens

  let schools: SchoolEntity[]

  const defaultPage = 1
  const defaultPageSize = 4

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
      .then(() => fixture.getTokens("Admin"))
      .then((data) => (adminUserAndTokens = data))
      .then(() => fixture.em.find(SchoolEntity, {}))
      .then((data) => {
        schools = data
      })
  })

  afterEach(() => {
    return fixture.afterEach()
  })

  describe("/schools/location (GET)", () => {
    it("gets all Schools sorted by checkins at location (200)", async () => {
      const location = userAndTokens.user.location

      if (!location) {
        throw new Error("Location is undefined")
      }

      const andere = await fixture.em.findOneOrFail(SchoolEntity, 1)
      const school = schools[5]

      for (let i = 0; i < 5; i++) {
        const newCheckIn = new CheckInEntity(location, andere)
        await fixture.em.persistAndFlush(newCheckIn)
      }

      for (let i = 0; i < 15; i++) {
        const newCheckIn = new CheckInEntity(location, school)
        await fixture.em.persistAndFlush(newCheckIn)
      }

      const response = await request(fixture.app.getHttpServer())
        .get("/schools/location")
        .set("Cookie", [`accessToken=${userAndTokens.tokens.accessToken}`])
        .expect(200)

      const schoolsResponse = response.body as School[]
      expect(schoolsResponse.length).toBe(schools.length)
      expect(schoolsResponse[schoolsResponse.length - 1].name).toBe(andere.name)
      expect(schoolsResponse[0].name).toBe(school.name)
    })

    it("returns Forbidden (403)", async () => {
      return await request(fixture.app.getHttpServer())
        .get("/schools/location")
        .set("Cookie", [`accessToken=${adminUserAndTokens.tokens.accessToken}`])
        .expect(403)
    })
  })

  describe("/schools (GET)", () => {
    it("gets all Schools with default page (200)", async () => {
      const response = await request(fixture.app.getHttpServer())
        .get("/schools")
        .set("Cookie", [`accessToken=${adminUserAndTokens.tokens.accessToken}`])
        .expect(200)

      const paginatedSchoolsResponse = response.body as GetAllPaginatedSchoolDto
      expect(paginatedSchoolsResponse.pagination.current).toBe(defaultPage)
      expect(paginatedSchoolsResponse.pagination.total).toBe(
        Math.ceil((schools.length - 1) / defaultPageSize)
      )
      expect(paginatedSchoolsResponse.data.length).toBe(defaultPageSize)
    })

    it("gets certain page of all Schools (200)", async () => {
      const page = 2

      const response = await request(fixture.app.getHttpServer())
        .get("/schools")
        .query({ page })
        .set("Cookie", [`accessToken=${adminUserAndTokens.tokens.accessToken}`])
        .expect(200)

      const paginatedSchoolsResponse = response.body as GetAllPaginatedSchoolDto
      expect(paginatedSchoolsResponse.pagination.current).toBe(page)
      expect(paginatedSchoolsResponse.pagination.total).toBe(
        Math.ceil((schools.length - 1) / defaultPageSize)
      )
      expect(paginatedSchoolsResponse.data.length).toBe(defaultPageSize)
    })

    it("returns Forbidden (403)", async () => {
      return await request(fixture.app.getHttpServer())
        .get("/schools")
        .set("Cookie", [`accessToken=${userAndTokens.tokens.accessToken}`])
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
        .set("Cookie", [`accessToken=${adminUserAndTokens.tokens.accessToken}`])
        .send(data)
        .expect(201)

      const schoolResponse = response.body as School
      expect(schoolResponse.id).toBeDefined()
      expect(schoolResponse.name).toBe(data.name)
    })

    it("returns School with this name already exists (409)", async () => {
      const data: CreateSchoolDto = {
        name: schools[1].name
      }

      const response = await request(fixture.app.getHttpServer())
        .post("/schools")
        .set("Cookie", [`accessToken=${adminUserAndTokens.tokens.accessToken}`])
        .send(data)
        .expect(409)

      const errorResponse = response.body as ErrorResponse
      expect(errorResponse.message).toBe("School with this name already exists")
    })

    it("returns Forbidden (403)", async () => {
      const data: CreateSchoolDto = {
        name: "NewSchool"
      }

      return await request(fixture.app.getHttpServer())
        .post("/schools")
        .set("Cookie", [`accessToken=${userAndTokens.tokens.accessToken}`])
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
        .set("Cookie", [`accessToken=${adminUserAndTokens.tokens.accessToken}`])
        .send(data)
        .expect(200)

      const schoolResponse = response.body as School
      expect(schoolResponse.id).toBe(id)
      expect(schoolResponse.name).toBe(data.name)
    })

    it("returns School with this name already exists (409)", async () => {
      const id = schools[1].id
      const data: UpdateSchoolDto = {
        name: schools[2].name
      }

      const response = await request(fixture.app.getHttpServer())
        .patch(`/schools/${id}`)
        .set("Cookie", [`accessToken=${adminUserAndTokens.tokens.accessToken}`])
        .send(data)
        .expect(409)

      const errorResponse = response.body as ErrorResponse
      expect(errorResponse.message).toBe("School with this name already exists")
    })

    it("returns School not found (404)", async () => {
      const id = -1
      const data: UpdateSchoolDto = {
        name: "NewSchool2"
      }

      const response = await request(fixture.app.getHttpServer())
        .patch(`/schools/${id}`)
        .set("Cookie", [`accessToken=${adminUserAndTokens.tokens.accessToken}`])
        .send(data)
        .expect(404)

      const errorResponse = response.body as ErrorResponse
      expect(errorResponse.message).toBe("School not found")
    })

    it("returns Forbidden (403)", async () => {
      const id = schools[1].id
      const data: UpdateSchoolDto = {
        name: "NewSchool2"
      }

      return await request(fixture.app.getHttpServer())
        .patch(`/schools/${id}`)
        .set("Cookie", [`accessToken=${userAndTokens.tokens.accessToken}`])
        .send(data)
        .expect(403)
    })
  })

  describe("/schools/:id (DELETE)", () => {
    it("deletes a School (200)", async () => {
      const id = schools[1].id

      const response = await request(fixture.app.getHttpServer())
        .delete(`/schools/${id}`)
        .set("Cookie", [`accessToken=${adminUserAndTokens.tokens.accessToken}`])
        .expect(200)

      const schoolResponse = response.body as School
      expect(schoolResponse.id).toBe(id)
    })

    it("returns School not found (404)", async () => {
      const id = -1

      const response = await request(fixture.app.getHttpServer())
        .delete(`/schools/${id}`)
        .set("Cookie", [`accessToken=${adminUserAndTokens.tokens.accessToken}`])
        .expect(404)

      const errorResponse = response.body as ErrorResponse
      expect(errorResponse.message).toBe("School not found")
    })

    it("returns Forbidden (403)", async () => {
      const id = schools[1].id

      return await request(fixture.app.getHttpServer())
        .delete(`/schools/${id}`)
        .set("Cookie", [`accessToken=${userAndTokens.tokens.accessToken}`])
        .expect(403)
    })
  })
})
