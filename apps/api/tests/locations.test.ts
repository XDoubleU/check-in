import { INestApplication } from "@nestjs/common"
import request from "supertest"
import { CreateLocationDto, Location, UpdateLocationDto, User } from "types"
import { clearDatabase, getAccessToken, getAdminAccessToken, getApp } from "./shared"
import { LocationsService } from "../src/locations/locations.service"
import { UsersService } from "../src/users/users.service"


describe("LocationsController (e2e)", () => {
  let app: INestApplication

  let accessToken: string
  let authUser: User
  let adminAccessToken: string

  let locationsService: LocationsService
  let usersService: UsersService

  let locations: Location[]

  const defaultPage = 1
  const defaultPageSize = 3
  
  beforeAll(async () => {
    app = await getApp()

    locationsService = app.get<LocationsService>(LocationsService)
    usersService = app.get<UsersService>(UsersService)

    await app.init()
  })

  beforeEach(async () => {
    // AccessTokens
    const getAccessTokenObject = await getAccessToken(app)
    accessToken = getAccessTokenObject.accessToken
    authUser = getAccessTokenObject.user

    adminAccessToken = await getAdminAccessToken(app)

    // LocationsService
    for (let i = 0; i < 20; i++){
      const user = await usersService.create(`TestUser${i}`, "testpassword")
      await locationsService.create(`TestLocation${i}`, 10, user)
    }

    locations = await locationsService.getAll()
  })

  afterEach(async () => {
    await clearDatabase(app)
  })

  describe("/locations (GET)", () => {
    it("gets all Locations with default page (200)", async () => {
      const response = await request(app.getHttpServer())
        .get("/locations")
        .set("Cookie", [`accessToken=${adminAccessToken}`])
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
        .set("Cookie", [`accessToken=${adminAccessToken}`])
        .expect(200)
      
      expect(response.body.page).toBe(page)
      expect(response.body.totalPages).toBe(Math.ceil(locations.length / defaultPageSize))
      expect(response.body.locations).toBeInstanceOf(Array<Location>)
    })

    it("returns Forbidden (403)", async () => {
      return await request(app.getHttpServer())
        .get("/locations")
        .set("Cookie", [`accessToken=${accessToken}`])
        .expect(403)
    })
  })

  describe("/locations/me (GET)", () => {
    it("get my Location (200)", async () => {
      const response = await request(app.getHttpServer())
        .get("/locations/me")
        .set("Cookie", [`accessToken=${accessToken}`])
        .expect(200)

      expect(response.body.id).toBe(authUser.locationId)
      expect(response.body.name).toBeDefined()
      expect(response.body.available).toBeDefined()
      expect(response.body.capacity).toBeDefined()
      expect(response.body.user.id).toBe(authUser.id)
      expect(response.body.user.username).toBe(authUser.username)
      expect(response.body.user.passwordHash).toBeUndefined()
      expect(response.body.user.role).toBeUndefined()
    })
  })

  describe("/locations/:id (GET)", () => {
    it("get Location as Admin (200)", async () => {
      const location = locations[0]
      
      const response = await request(app.getHttpServer())
        .get(`/locations/${location.id}`)
        .set("Cookie", [`accessToken=${adminAccessToken}`])
        .expect(200)

      expect(response.body.id).toBe(location.id)
      expect(response.body.name).toBe(location.name)
      expect(response.body.available).toBe(location.available)
      expect(response.body.capacity).toBe(location.capacity)
      expect(response.body.user.id).toBe(location.user.id)
      expect(response.body.user.username).toBe(location.user.username)
      expect(response.body.user.passwordHash).toBeUndefined()
      expect(response.body.user.role).toBeUndefined()
    })

    it("get Location as User (200)", async () => {
      const response = await request(app.getHttpServer())
        .get(`/locations/${authUser.locationId}`)
        .set("Cookie", [`accessToken=${accessToken}`])
        .expect(200)

      expect(response.body.id).toBe(authUser.locationId)
      expect(response.body.name).toBeDefined()
      expect(response.body.available).toBeDefined()
      expect(response.body.capacity).toBeDefined()
      expect(response.body.user.id).toBe(authUser.id)
      expect(response.body.user.username).toBe(authUser.username)
      expect(response.body.user.passwordHash).toBeUndefined()
      expect(response.body.user.role).toBeUndefined()
    })

    it("returns Location not found because Location doesn't exist (404)", async () => {
      const response = await request(app.getHttpServer())
        .get("/locations/notfound")
        .set("Cookie", [`accessToken=${adminAccessToken}`])
        .expect(404)

      expect(response.body.message).toBe("Location not found")
    })

    it("returns Location not found because User doesn't own Location (404)", async () => {
      const location = locations.find((location) => location.name === "TestLocation0") as Location

      const response = await request(app.getHttpServer())
        .get(`/locations/${location.id}`)
        .set("Cookie", [`accessToken=${accessToken}`])
        .expect(404)

      expect(response.body.message).toBe("Location not found")
    })
  })

  describe("/locations (POST)", () => {
    it("creates a new Location (201)", async () => {
      const data: CreateLocationDto = {
        name: "TestLocation",
        capacity: 10,
        username: "TestLocationUser",
        password: "testpassword"
      }
  
      const response = await request(app.getHttpServer())
        .post("/locations")
        .set("Cookie", [`accessToken=${adminAccessToken}`])
        .send(data)
        .expect(201)
      
      expect(response.body.id).toBeDefined()
      expect(response.body.name).toBe(data.name)
      expect(response.body.available).toBe(data.capacity)
      expect(response.body.capacity).toBe(data.capacity)
      expect(response.body.user.id).toBeDefined()
      expect(response.body.user.username).toBe(data.username)
      expect(response.body.user.passwordHash).toBeUndefined()
      expect(response.body.user.role).toBeUndefined()
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
        .set("Cookie", [`accessToken=${adminAccessToken}`])
        .send(data)
        .expect(409)
      
      expect(response.body.message).toBe("Location with this name already exists")
    })

    it("returns User with this username already exists (409)", async () => {
      const data: CreateLocationDto = {
        name: "TestLocation",
        capacity: 10,
        username: locations[0].user.username,
        password: "testpassword"
      }
  
      const response = await request(app.getHttpServer())
        .post("/locations")
        .set("Cookie", [`accessToken=${adminAccessToken}`])
        .send(data)
        .expect(409)
      
      expect(response.body.message).toBe("User with this username already exists")
    })

    it("returns Forbidden (403)", async () => {
      const data: CreateLocationDto = {
        name: "TestLocation",
        capacity: 10,
        username: "TestLocationUser",
        password: "testpassword"
      }
  
      return await request(app.getHttpServer())
        .post("/locations")
        .set("Cookie", [`accessToken=${accessToken}`])
        .send(data)
        .expect(403)
    })
  })

  describe("/locations/:id (PATCH)", () => {
    it("updates a Location (200)", async () => {
      const id = authUser.locationId
      
      const data: UpdateLocationDto = {
        name: "NewTestLocation",
        username: "NewTestLocationUser"
      }
  
      const response = await request(app.getHttpServer())
        .patch(`/locations/${id}`)
        .set("Cookie", [`accessToken=${accessToken}`])
        .send(data)
        .expect(200)
      
      expect(response.body.id).toBeDefined()
      expect(response.body.name).toBe(data.name)
      expect(response.body.available).toBeDefined()
      expect(response.body.capacity).toBeDefined()
      expect(response.body.user.id).toBe(authUser.id)
      expect(response.body.user.username).toBe(data.username)
      expect(response.body.user.passwordHash).toBeUndefined()
      expect(response.body.user.role).toBeUndefined()
    })

    it("returns Location with this name already exists (409)", async () => {
      const id = authUser.locationId
      
      const data: UpdateLocationDto = {
        name: locations[0].name,
        username: "NewTestLocationUser"
      }
  
      const response = await request(app.getHttpServer())
        .patch(`/locations/${id}`)
        .set("Cookie", [`accessToken=${accessToken}`])
        .send(data)
        .expect(409)
      
      expect(response.body.message).toBe("Location with this name already exists")
    })

    it("returns User with this username already exists (409)", async () => {
      const id = authUser.locationId
      
      const data: UpdateLocationDto = {
        name: "NewTestLocation",
        username: locations[0].user.username
      }
  
      const response = await request(app.getHttpServer())
        .patch(`/locations/${id}`)
        .set("Cookie", [`accessToken=${accessToken}`])
        .send(data)
        .expect(409)
      
      expect(response.body.message).toBe("User with this username already exists")
    })
  })

  describe("/locations/:id (DELETE)", () => {
    it("deletes a Location (200)", async () => {
      const id = locations[0].id
  
      const response = await request(app.getHttpServer())
        .delete(`/locations/${id}`)
        .set("Cookie", [`accessToken=${adminAccessToken}`])
        .expect(200)
      
      expect(response.body.id).toBe(id)
    })

    it("returns Location not found (404)", async () => {
      const id = -1
  
      const response = await request(app.getHttpServer())
        .delete(`/locations/${id}`)
        .set("Cookie", [`accessToken=${adminAccessToken}`])
        .expect(404)
      
      expect(response.body.message).toBe("Location not found")
    })

    it("returns Forbidden (403)", async () => {
      const id = locations[0].id
  
      await request(app.getHttpServer())
        .delete(`/locations/${id}`)
        .set("Cookie", [`accessToken=${accessToken}`])
        .expect(403)
    })
  })
})
