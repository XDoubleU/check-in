import { INestApplication } from "@nestjs/common"
import request from "supertest"
import { CreateSchoolDto, School, UpdateSchoolDto } from "types"
import { SchoolsService } from "../src/schools/schools.service"
import { clearDatabase, getAccessToken, getApp } from "./shared"


describe("SchoolsController (e2e)", () => {
  let app: INestApplication
  let accessToken: string

  let schoolsService: SchoolsService

  let schools: School[]

  const defaultPage = 1
  const defaultPageSize = 4
  
  beforeAll(async () => {
    app = await getApp()
    accessToken = await getAccessToken(app)

    schoolsService = app.get<SchoolsService>(SchoolsService)

    await app.init()
  })

  beforeEach(async () => {
    // SchoolsService
    for (let i = 0; i < 20; i++){
      let tempSchool = await schoolsService.getByName(`TestSchool${i}`)
      if (!tempSchool){
        tempSchool = await schoolsService.create(`TestSchool${i}`)
        if (!tempSchool) throw new Error()
      }
    }

    schools = await schoolsService.getAll(undefined)
  })

  afterAll(async () => {
    clearDatabase(app)
  })

  describe("/schools/all (GET)", () => {
    // TODO: test ordering
    it("gets all Schools (200)", async () => {
      const response = await request(app.getHttpServer())
        .get("/schools/all")
        .set("Cookie", [`accessToken=${accessToken}`])
        .expect(200)

      expect(response.body).toBeInstanceOf(Array<School>)
    })

    it("gets certain page of all Schools (200)", async () => {
      const page = 2

      const response = await request(app.getHttpServer())
        .get("/schools")
        .query({ page })
        .set("Cookie", [`accessToken=${accessToken}`])
        .expect(200)
      
      expect(response.body.page).toBe(page)
      expect(response.body.totalPages).toBe(Math.ceil(schools.length / defaultPageSize))
      expect(response.body.schools).toBeInstanceOf(Array<School>)
    })
  })

  describe("/schools (GET)", () => {
    it("gets all Schools with default page (200)", async () => {
      const response = await request(app.getHttpServer())
        .get("/schools")
        .set("Cookie", [`accessToken=${accessToken}`])
        .expect(200)

      expect(response.body.page).toBe(defaultPage)
      expect(response.body.totalPages).toBe(Math.ceil(schools.length / defaultPageSize))
      expect(response.body.schools).toBeInstanceOf(Array<School>)
    })

    it("gets certain page of all Schools (200)", async () => {
      const page = 2

      const response = await request(app.getHttpServer())
        .get("/schools")
        .query({ page })
        .set("Cookie", [`accessToken=${accessToken}`])
        .expect(200)
      
      expect(response.body.page).toBe(page)
      expect(response.body.totalPages).toBe(Math.ceil(schools.length / defaultPageSize))
      expect(response.body.schools).toBeInstanceOf(Array<School>)
    })
  })

  describe("/schools (POST)", () => {
    it("creates a new School (200)", async () => {
      const data: CreateSchoolDto = {
        name: "NewSchool"
      }
  
      const response = await request(app.getHttpServer())
        .post("/schools")
        .set("Cookie", [`accessToken=${accessToken}`])
        .send(data)
        .expect(201)
      
      expect(response.body.id).toBeDefined()
      expect(response.body.name).toBe(data.name)
    })

    it("returns School with this name already exists (409)", async () => {
      const data: CreateSchoolDto = {
        name: schools[1].name
      }
  
      const response = await request(app.getHttpServer())
        .post("/schools")
        .set("Cookie", [`accessToken=${accessToken}`])
        .send(data)
        .expect(409)
      
      expect(response.body.message).toBe("School with this name already exists")
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
        .set("Cookie", [`accessToken=${accessToken}`])
        .send(data)
        .expect(200)
      
      expect(response.body.id).toBe(id)
      expect(response.body.name).toBe(data.name)
    })

    it("returns School with this name already exists (409)", async () => {
      const id = schools[1].id
      const data: UpdateSchoolDto = {
        name: schools[2].name
      }
  
      const response = await request(app.getHttpServer())
        .patch(`/schools/${id}`)
        .set("Cookie", [`accessToken=${accessToken}`])
        .send(data)
        .expect(409)
      
      expect(response.body.message).toBe("School with this name already exists")
    })

    it("returns School not found (404)", async () => {
      const id = -1
      const data: UpdateSchoolDto = {
        name: "NewSchool2"
      }
  
      const response = await request(app.getHttpServer())
        .patch(`/schools/${id}`)
        .set("Cookie", [`accessToken=${accessToken}`])
        .send(data)
        .expect(404)
      
      expect(response.body.message).toBe("School not found")
    })
  })

  describe("/schools/:id (DELETE)", () => {
    it("deletes a School (200)", async () => {
      const id = schools[1].id
  
      const response = await request(app.getHttpServer())
        .delete(`/schools/${id}`)
        .set("Cookie", [`accessToken=${accessToken}`])
        .expect(200)
      
      expect(response.body.id).toBe(id)
    })

    it("returns School not found (404)", async () => {
      const id = -1
  
      const response = await request(app.getHttpServer())
        .delete(`/schools/${id}`)
        .set("Cookie", [`accessToken=${accessToken}`])
        .expect(404)
      
      expect(response.body.message).toBe("School not found")
    })
  })
})
