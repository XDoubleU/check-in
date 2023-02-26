import { INestApplication } from "@nestjs/common"
import request from "supertest"
import { CreateLocationDto, Location, Role, User } from "types"
import { clearDatabase, getAccessToken, getApp } from "./shared"
import { LocationsService } from "../src/locations/locations.service"
import { UsersService } from "../src/users/users.service"


describe("LocationsController (e2e)", () => {
  let app: INestApplication
  let accessToken: string

  let locationsService: LocationsService
  let usersService: UsersService

  let locations: Location[]

  const defaultPage = 1
  const defaultPageSize = 3
  
  beforeAll(async () => {
    app = await getApp()
    accessToken = await getAccessToken(app)

    locationsService = app.get<LocationsService>(LocationsService)
    usersService = app.get<UsersService>(UsersService)

    await app.init()
  })

  beforeEach(async () => {
    // LocationsService
    for (let i = 0; i < 20; i++){
      let tempLocation = await locationsService.getByName(`TestLocation${i}`)
      if (!tempLocation){
        const user = await usersService.create(`TestUser${i}`, "testpassword") as User
        tempLocation = await locationsService.create(`TestLocation${i}`, 10, user)
        if (!tempLocation) throw new Error()
      }
    }

    locations = await locationsService.getAll()
  })

  afterAll(async () => {
    clearDatabase(app)
  })

  describe("/locations (GET)", () => {
    it("gets all Locations with default page (200)", async () => {
      const response = await request(app.getHttpServer())
        .get("/locations")
        .set("Cookie", [`accessToken=${accessToken}`])
        .expect(200)

      expect(response.body.page).toBe(defaultPage)
      expect(response.body.totalPages).toBe(Math.ceil(locations.length / defaultPageSize))
      expect(response.body.locations).toBeInstanceOf(Array<Location>)
    })

    it("gets certain page of all Locations (200)", async () => {
      const page = 2

      const response = await request(app.getHttpServer())
        .get("/locations")
        .query({ page })
        .set("Cookie", [`accessToken=${accessToken}`])
        .expect(200)
      
      expect(response.body.page).toBe(page)
      expect(response.body.totalPages).toBe(Math.ceil(locations.length / defaultPageSize))
      expect(response.body.locations).toBeInstanceOf(Array<Location>)
    })
  })

  describe("/locations (POST)", () => {
    it("creates a new Location (200)", async () => {
      const data: CreateLocationDto = {
        name: "TestLocation",
        capacity: 10,
        username: "TestLocationUser",
        password: "testpassword"
      }
  
      const response = await request(app.getHttpServer())
        .post("/locations")
        .set("Cookie", [`accessToken=${accessToken}`])
        .send(data)
        .expect(201)
      
      expect(response.body.id).toBeDefined()
      expect(response.body.name).toBe(data.name)
      expect(response.body.available).toBe(data.capacity)
      expect(response.body.capacity).toBe(data.capacity)
      expect(response.body.user.id).toBeDefined()
      expect(response.body.user.username).toBe(data.username)
      expect(response.body.user.passwordHash).toBeUndefined()
      expect(response.body.user.role).toBe(Role.User)
    })

    it("returns Location with this name already exists (409)", async () => {
      const data: CreateLocationDto = {
        name: locations[0].name,
        capacity: 10,
        username: "TestLocationUser",
        password: "testpassword"
      }
  
      const response = await request(app.getHttpServer())
        .post("/locations")
        .set("Cookie", [`accessToken=${accessToken}`])
        .send(data)
        .expect(409)
      
      expect(response.body.message).toBe("Location with this name already exists")
    })
  })

  describe("/locations/:id (DELETE)", () => {
    it("deletes a Location (200)", async () => {
      const id = locations[0].id
  
      const response = await request(app.getHttpServer())
        .delete(`/locations/${id}`)
        .set("Cookie", [`accessToken=${accessToken}`])
        .expect(200)
      
      expect(response.body.id).toBe(id)
    })

    it("returns Location not found (404)", async () => {
      const id = -1
  
      const response = await request(app.getHttpServer())
        .delete(`/locations/${id}`)
        .set("Cookie", [`accessToken=${accessToken}`])
        .expect(404)
      
      expect(response.body.message).toBe("Location not found")
    })
  })
})
