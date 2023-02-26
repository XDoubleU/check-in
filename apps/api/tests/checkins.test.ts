import { INestApplication } from "@nestjs/common"
import request from "supertest"
import { CreateCheckInDto, Location, School, User } from "types"
import { LocationsService } from "../src/locations/locations.service"
import { UsersService } from "../src/users/users.service"
import { SchoolsService } from "../src/schools/schools.service"
import { clearDatabase, getAccessToken, getApp } from "./shared"


describe("CheckInsController (e2e)", () => {
  let app: INestApplication
  let accessToken: string

  let locationsService: LocationsService
  let usersService: UsersService
  let schoolsService: SchoolsService

  let user: User
  let location: Location
  let school: School
  
  beforeAll(async () => {
    app = await getApp()
    accessToken = await getAccessToken(app)

    locationsService = app.get<LocationsService>(LocationsService)
    usersService = app.get<UsersService>(UsersService)
    schoolsService = app.get<SchoolsService>(SchoolsService)

    await app.init()
  })

  beforeEach(async () => {
    // UsersService
    let tempUser = await usersService.getByUserName("TestUser")
    if (!tempUser){
      tempUser = await usersService.create("TestUser", "testpassword")
      if (!tempUser) throw new Error()
    }
    user = tempUser

    // LocationsService
    let tempLocation = await locationsService.getByName("TestLocation")
    if (!tempLocation){
      tempLocation = await locationsService.create("TestLocation", 10, user)
      if (!tempLocation) throw new Error()
    }
    location = tempLocation

    // SchoolsService
    let tempSchool = await schoolsService.getByName("TestSchool")
    if (!tempSchool){
      tempSchool = await schoolsService.create("TestSchool")
      if (!tempSchool) throw new Error()
    }
    school = tempSchool
  })

  afterAll(async () => {
    clearDatabase(app)
  })

  describe("/checkins (POST)", () => {
    // TODO: test role access
    it("creates a new CheckIn (200)", async () => {
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
  })
})
