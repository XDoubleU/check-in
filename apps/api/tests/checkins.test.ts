import { INestApplication } from "@nestjs/common"
import request from "supertest"
import { CreateCheckInDto, Location, School, User } from "types"
import { LocationsService } from "../src/locations/locations.service"
import { UsersService } from "../src/users/users.service"
import { SchoolsService } from "../src/schools/schools.service"
import { clearDatabase, getAccessToken, getAdminAccessToken, getApp } from "./shared"


describe("CheckInsController (e2e)", () => {
  let app: INestApplication

  let accessToken: string
  let adminAccessToken: string

  let locationsService: LocationsService
  let usersService: UsersService
  let schoolsService: SchoolsService

  let user: User
  let location: Location
  let school: School
  
  beforeAll(async () => {
    app = await getApp()

    locationsService = app.get<LocationsService>(LocationsService)
    usersService = app.get<UsersService>(UsersService)
    schoolsService = app.get<SchoolsService>(SchoolsService)

    await app.init()
  })

  beforeEach(async () => {
    // AccessTokens
    const getAccessTokenObject = await getAccessToken(app)
    accessToken = getAccessTokenObject.accessToken
    adminAccessToken = await getAdminAccessToken(app)

    // UsersService
    user = await usersService.create("TestUser", "testpassword")

    // LocationsService
    location = await locationsService.create("TestLocation", 10, user)

    // SchoolsService
    school = await schoolsService.create("TestSchool")
  })

  afterEach(async () => {
    await clearDatabase(app)
  })

  describe("/checkins (POST)", () => {
    it("creates a new CheckIn (201)", async () => {
      const data: CreateCheckInDto = {
        locationId: location.id,
        schoolId: school.id
      }
  
      const response = await request(app.getHttpServer())
        .post("/checkins")
        .set("Cookie", [`accessToken=${accessToken}`])
        .send(data)
        .expect(201)
      
      expect(response.body.id).toBeDefined()
      expect(response.body.locationId).toBe(location.id)
      expect(response.body.capacity).toBe(location.capacity)
      expect(response.body.datetime).toBeDefined()
      expect(response.body.schoolId).toBe(school.id)
    })

    it("returns Location not found (404)", async () => {
      const data: CreateCheckInDto = {
        locationId: "notfound",
        schoolId: school.id
      }
  
      const response = await request(app.getHttpServer())
        .post("/checkins")
        .set("Cookie", [`accessToken=${accessToken}`])
        .send(data)
        .expect(404)
      
      expect(response.body.message).toBe("Location not found")
    })

    it("returns School not found (404)", async () => {
      const data: CreateCheckInDto = {
        locationId: location.id,
        schoolId: -1
      }
  
      const response = await request(app.getHttpServer())
        .post("/checkins")
        .set("Cookie", [`accessToken=${accessToken}`])
        .send(data)
        .expect(404)
      
      expect(response.body.message).toBe("School not found")
    })

    it("returns Forbidden (403)", async () => {
      const data: CreateCheckInDto = {
        locationId: location.id,
        schoolId: school.id
      }
  
      return await request(app.getHttpServer())
        .post("/checkins")
        .set("Cookie", [`accessToken=${adminAccessToken}`])
        .send(data)
        .expect(403)
    })
  })
})
