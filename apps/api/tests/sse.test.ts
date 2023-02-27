import { INestApplication } from "@nestjs/common"
import { Location, User } from "types"
import { LocationsService } from "../src/locations/locations.service"
import { UsersService } from "../src/users/users.service"
import { clearDatabase, getApp } from "./shared"


describe("SseController (e2e)", () => {
  let app: INestApplication

  //let accessToken: string
  //let adminAccessToken: string

  let locationsService: LocationsService
  let usersService: UsersService
  //let schoolsService: SchoolsService

  let user: User
  let location: Location
  //let school: School
  
  beforeAll(async () => {
    app = await getApp()

    locationsService = app.get<LocationsService>(LocationsService)
    usersService = app.get<UsersService>(UsersService)
    //schoolsService = app.get<SchoolsService>(SchoolsService)

    await app.init()
  })

  beforeEach(async () => {
    // AccessTokens
    //const getAccessTokenObject = await getAccessToken(app)
    //accessToken = getAccessTokenObject.accessToken
    //adminAccessToken = await getAdminAccessToken(app)

    // UsersService
    user = await usersService.create("TestUser", "testpassword")

    // LocationsService
    location = await locationsService.create("TestLocation", 10, user)

    // SchoolsService
    //school = await schoolsService.create("TestSchool")
  })

  afterEach(async () => {
    await clearDatabase(app)
  })

  describe("/sse (SSE)", () => {
    it("receives event after updating Location capacity (200)", async () => {
      const eventSource = new EventSource("/sse")
      
      eventSource.onmessage = (event: MessageEvent<Location>): void => {
        console.log(event)
      }

      await locationsService.update(location, undefined, location.capacity + 10)
    })

    it("receives event after creating CheckIn (200)", () => {
      return 
    })
  })

  describe("/sse/:id (SSE)", () => {
    it("receives event after updating Location capacity (200)", () => {
      return 
    })

    it("receives event after creating CheckIn (200)", () => {
      return 
    })

    it("returns Forbidden (403)", () => {
      return
    })
  })
})
