import { INestApplication } from "@nestjs/common"
import request from "supertest"
import { CreateSchoolDto, GetAllPaginatedSchoolDto, School, UpdateSchoolDto, User } from "types"
import { SchoolsService } from "../src/schools/schools.service"
import { clearDatabase, ErrorResponse, getAccessToken, getAdminAccessToken, getApp } from "./shared"
import { CheckInsService } from "../src/checkins/checkins.service"
import { LocationsService } from "../src/locations/locations.service"


describe("SchoolsController (e2e)", () => {
  let app: INestApplication

  let accessToken: string
  let authUser: Omit<User, "locationId"> & { locationId: string }
  let adminAccessToken: string

  let schoolsService: SchoolsService
  let checkinsService: CheckInsService
  let locationsService: LocationsService

  let schools: School[]

  const defaultPage = 1
  const defaultPageSize = 4
  
  beforeAll(async () => {
    app = await getApp()

    schoolsService = app.get<SchoolsService>(SchoolsService)
    checkinsService = app.get<CheckInsService>(CheckInsService)
    locationsService = app.get<LocationsService>(LocationsService)

    await app.init()
  })

  beforeEach(async () => {
    // AccessTokens
    const getAccessTokenObject = await getAccessToken(app)
    accessToken = getAccessTokenObject.accessToken
    authUser = getAccessTokenObject.user
    adminAccessToken = await getAdminAccessToken(app)

    // SchoolsService
    for (let i = 0; i < 20; i++){
      await schoolsService.create(`TestSchool${i}`)
    }
    schools = await schoolsService.getAll(undefined)
  })

  afterEach(async () => {
    await clearDatabase(app)
  })

  describe("/schools/all (GET)", () => {
    it("gets all Schools (200)", async () => {
      const location = await locationsService.getById(authUser.locationId)
      if (!location) {
        throw new Error("Location is undefined")
      }

      const andere = await schoolsService.getById(1)
      if (!andere) {
        throw new Error("andere is undefined")
      }

      const school = schools[5]

      for (let i = 0; i < 10; i++) {
        await checkinsService.create(location, andere)
      }

      for (let i = 0; i < 15; i++) {
        await checkinsService.create(location, school)
      }

      const response = await request(app.getHttpServer())
        .get("/schools/all")
        .set("Cookie", [`accessToken=${accessToken}`])
        .expect(200)

      const schoolsResponse = response.body as School[]
      expect(schoolsResponse).toBeInstanceOf(Array<School>)
      expect(schoolsResponse[schoolsResponse.length - 1]).toStrictEqual(andere)
      expect(schoolsResponse[0]).toStrictEqual(school)
    })

    it("returns Forbidden (403)", async () => {
      return await request(app.getHttpServer())
        .get("/schools/all")
        .set("Cookie", [`accessToken=${adminAccessToken}`])
        .expect(403)
    })
  })

  describe("/schools (GET)", () => {
    it("gets all Schools with default page (200)", async () => {
      const response = await request(app.getHttpServer())
        .get("/schools")
        .set("Cookie", [`accessToken=${adminAccessToken}`])
        .expect(200)

      const paginatedSchoolsResponse = response.body as GetAllPaginatedSchoolDto
      expect(paginatedSchoolsResponse.page).toBe(defaultPage)
      expect(paginatedSchoolsResponse.totalPages).toBe(Math.ceil(schools.length / defaultPageSize))
      expect(paginatedSchoolsResponse.schools).toBeInstanceOf(Array<School>)
    })

    it("gets certain page of all Schools (200)", async () => {
      const page = 2

      const response = await request(app.getHttpServer())
        .get("/schools")
        .query({ page })
        .set("Cookie", [`accessToken=${adminAccessToken}`])
        .expect(200)
      
      const paginatedSchoolsResponse = response.body as GetAllPaginatedSchoolDto
      expect(paginatedSchoolsResponse.page).toBe(page)
      expect(paginatedSchoolsResponse.totalPages).toBe(Math.ceil(schools.length / defaultPageSize))
      expect(paginatedSchoolsResponse.schools).toBeInstanceOf(Array<School>)
    })

    it("returns Forbidden (403)", async () => {
      return await request(app.getHttpServer())
        .get("/schools")
        .set("Cookie", [`accessToken=${accessToken}`])
        .expect(403)
    })
  })

  describe("/schools (POST)", () => {
    it("creates a new School (201)", async () => {
      const data: CreateSchoolDto = {
        name: "NewSchool"
      }
  
      const response = await request(app.getHttpServer())
        .post("/schools")
        .set("Cookie", [`accessToken=${adminAccessToken}`])
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
  
      const response = await request(app.getHttpServer())
        .post("/schools")
        .set("Cookie", [`accessToken=${adminAccessToken}`])
        .send(data)
        .expect(409)
      
      const errorResponse = response.body as ErrorResponse
      expect(errorResponse.message).toBe("School with this name already exists")
    })

    it("returns Forbidden (403)", async () => {
      const data: CreateSchoolDto = {
        name: "NewSchool"
      }
  
      return await request(app.getHttpServer())
        .post("/schools")
        .set("Cookie", [`accessToken=${accessToken}`])
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
  
      const response = await request(app.getHttpServer())
        .patch(`/schools/${id}`)
        .set("Cookie", [`accessToken=${adminAccessToken}`])
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
  
      const response = await request(app.getHttpServer())
        .patch(`/schools/${id}`)
        .set("Cookie", [`accessToken=${adminAccessToken}`])
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
  
      const response = await request(app.getHttpServer())
        .patch(`/schools/${id}`)
        .set("Cookie", [`accessToken=${adminAccessToken}`])
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
  
      return await request(app.getHttpServer())
        .patch(`/schools/${id}`)
        .set("Cookie", [`accessToken=${accessToken}`])
        .send(data)
        .expect(403)
    })
  })

  describe("/schools/:id (DELETE)", () => {
    it("deletes a School (200)", async () => {
      const id = schools[1].id
  
      const response = await request(app.getHttpServer())
        .delete(`/schools/${id}`)
        .set("Cookie", [`accessToken=${adminAccessToken}`])
        .expect(200)
      
      const schoolResponse = response.body as School
      expect(schoolResponse.id).toBe(id)
    })

    it("returns School not found (404)", async () => {
      const id = -1
  
      const response = await request(app.getHttpServer())
        .delete(`/schools/${id}`)
        .set("Cookie", [`accessToken=${adminAccessToken}`])
        .expect(404)
      
      const errorResponse = response.body as ErrorResponse
      expect(errorResponse.message).toBe("School not found")
    })

    it("returns Forbidden (403)", async () => {
      const id = schools[1].id
  
      return await request(app.getHttpServer())
        .delete(`/schools/${id}`)
        .set("Cookie", [`accessToken=${accessToken}`])
        .expect(403)
    })
  })
})
