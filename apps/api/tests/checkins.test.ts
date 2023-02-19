import { Test, TestingModule } from "@nestjs/testing"
import { INestApplication } from "@nestjs/common"
import request from "supertest"
import { AppModule } from "../src/app.module"
import { CreateCheckInDto, User, Location, School } from "types"
import { LocationsService } from "../src/locations/locations.service"
import { SseService } from "../src/sse/sse.service"
import { UsersService } from "../src/users/users.service"
import { SchoolsService } from "../src/schools/schools.service"

describe("CheckInsController (e2e)", () => {
  let app: INestApplication
  let locationsService: LocationsService
  let usersService: UsersService
  let schoolsService: SchoolsService
  
  beforeAll(() => {
    const sseService = new SseService()
    locationsService = new LocationsService(sseService)
    usersService = new UsersService()
    schoolsService = new SchoolsService()
  })

  beforeEach(async () => {
    const moduleFixture: TestingModule = await Test.createTestingModule({
      imports: [AppModule],
    }).compile()

    app = moduleFixture.createNestApplication()
    await app.init()
  })

  test("/ (POST)", async () => {
    it("creates a new checkin", async () => {
      const user = await usersService.create("TestUser", "testpassword") as User
      const location = await locationsService.create("TestLocation", 10, user) as Location
      const school = await schoolsService.create("TestSchool") as School

      const data: CreateCheckInDto = {
        locationId: location.id,
        schoolId: school.id
      }
  
      return request(app.getHttpServer())
        .post("/")
        .send(data)
        .expect(200)
        .expect("Hello World!")
    })
  })
})
