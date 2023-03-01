import { INestApplication } from "@nestjs/common"
import request from "supertest"
import { CreateLocationDto, GetAllPaginatedLocationDto, Location, UpdateLocationDto, User } from "types"
import { clearDatabase, ErrorResponse, getAccessToken, getAdminAccessToken, getApp } from "./shared"
import { LocationsService } from "../src/locations/locations.service"
import { UsersService } from "../src/users/users.service"


describe("LocationsController (e2e)", () => {
  let app: INestApplication

  let accessToken: string
  let authUser: Omit<User, "locationId"> & { locationId: string }
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
      
      const paginatedLocationsResponse = response.body as GetAllPaginatedLocationDto
      expect(paginatedLocationsResponse.page).toBe(defaultPage)
      expect(paginatedLocationsResponse.totalPages).toBe(Math.ceil(locations.length / defaultPageSize))
      expect(paginatedLocationsResponse.locations).toBeInstanceOf(Array<Location>)
    })

    it("gets certain page of all Locations (200)", async () => {
      const page = 2

      const response = await request(app.getHttpServer())
        .get("/locations")
        .query({ page })
        .set("Cookie", [`accessToken=${adminAccessToken}`])
        .expect(200)
      
      const paginatedLocationsResponse = response.body as GetAllPaginatedLocationDto
      expect(paginatedLocationsResponse.page).toBe(page)
      expect(paginatedLocationsResponse.totalPages).toBe(Math.ceil(locations.length / defaultPageSize))
      expect(paginatedLocationsResponse.locations).toBeInstanceOf(Array<Location>)
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

      const locationResponse = response.body as Location
      expect(locationResponse.id).toBe(authUser.locationId)
      expect(locationResponse.name).toBeDefined()
      expect(locationResponse.available).toBeDefined()
      expect(locationResponse.capacity).toBeDefined()
      expect(locationResponse.user.id).toBe(authUser.id)
      expect(locationResponse.user.username).toBe(authUser.username)
    })

    it("returns Forbidden (403)", async () => {
      return await request(app.getHttpServer())
        .get("/locations/me")
        .set("Cookie", [`accessToken=${adminAccessToken}`])
        .expect(403)
    })
  })

  describe("/locations/:id (GET)", () => {
    it("get Location as Admin (200)", async () => {
      const location = locations[0]
      
      const response = await request(app.getHttpServer())
        .get(`/locations/${location.id}`)
        .set("Cookie", [`accessToken=${adminAccessToken}`])
        .expect(200)

      const locationResponse = response.body as Location
      expect(locationResponse.id).toBe(location.id)
      expect(locationResponse.name).toBe(location.name)
      expect(locationResponse.available).toBe(location.available)
      expect(locationResponse.capacity).toBe(location.capacity)
      expect(locationResponse.user.id).toBe(location.user.id)
      expect(locationResponse.user.username).toBe(location.user.username)
    })

    it("get Location as User (200)", async () => {
      const response = await request(app.getHttpServer())
        .get(`/locations/${authUser.locationId}`)
        .set("Cookie", [`accessToken=${accessToken}`])
        .expect(200)

      const locationResponse = response.body as Location
      expect(locationResponse.id).toBe(authUser.locationId)
      expect(locationResponse.name).toBeDefined()
      expect(locationResponse.available).toBeDefined()
      expect(locationResponse.capacity).toBeDefined()
      expect(locationResponse.user.id).toBe(authUser.id)
      expect(locationResponse.user.username).toBe(authUser.username)
    })

    it("returns Location not found because Location doesn't exist (404)", async () => {
      const response = await request(app.getHttpServer())
        .get("/locations/notfound")
        .set("Cookie", [`accessToken=${adminAccessToken}`])
        .expect(404)
      
      const errorResponse = response.body as ErrorResponse
      expect(errorResponse.message).toBe("Location not found")
    })

    it("returns Location not found because User doesn't own Location (404)", async () => {
      const location = locations.find((location) => location.name === "TestLocation0")
      if (!location) {
        throw new Error("Location is undefined")
      }

      const response = await request(app.getHttpServer())
        .get(`/locations/${location.id}`)
        .set("Cookie", [`accessToken=${accessToken}`])
        .expect(404)
      
      const errorResponse = response.body as ErrorResponse
      expect(errorResponse.message).toBe("Location not found")
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
      
      const locationResponse = response.body as Location
      expect(locationResponse.id).toBeDefined()
      expect(locationResponse.name).toBe(data.name)
      expect(locationResponse.available).toBe(data.capacity)
      expect(locationResponse.capacity).toBe(data.capacity)
      expect(locationResponse.user.id).toBeDefined()
      expect(locationResponse.user.username).toBe(data.username)
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
      
      const errorResponse = response.body as ErrorResponse
      expect(errorResponse.message).toBe("Location with this name already exists")
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
      
      const errorResponse = response.body as ErrorResponse
      expect(errorResponse.message).toBe("User with this username already exists")
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
      
      const locationResponse = response.body as Location
      expect(locationResponse.id).toBeDefined()
      expect(locationResponse.name).toBe(data.name)
      expect(locationResponse.available).toBeDefined()
      expect(locationResponse.capacity).toBeDefined()
      expect(locationResponse.user.id).toBe(authUser.id)
      expect(locationResponse.user.username).toBe(data.username)
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
      
      const errorResponse = response.body as ErrorResponse
      expect(errorResponse.message).toBe("Location with this name already exists")
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
      
      const errorResponse = response.body as ErrorResponse
      expect(errorResponse.message).toBe("User with this username already exists")
    })

    it("returns Location not found because User doesn't own Location (404)", async () => {
      const location = locations.find((location) => location.name === "TestLocation0")
      if (!location) {
        throw new Error("Location is undefined")
      }

      const data: UpdateLocationDto = {
        name: "NewTestLocation",
        username: "NewTestLocationUser"
      }
  
      return await request(app.getHttpServer())
        .patch(`/locations/${location.id}`)
        .set("Cookie", [`accessToken=${accessToken}`])
        .send(data)
        .expect(404)
    })
  })

  describe("/locations/:id (DELETE)", () => {
    it("deletes a Location (200)", async () => {
      const id = locations[0].id
  
      const response = await request(app.getHttpServer())
        .delete(`/locations/${id}`)
        .set("Cookie", [`accessToken=${adminAccessToken}`])
        .expect(200)
      
      const locationResponse = response.body as Location
      expect(locationResponse.id).toBe(id)
    })

    it("returns Location not found (404)", async () => {
      const id = -1
  
      const response = await request(app.getHttpServer())
        .delete(`/locations/${id}`)
        .set("Cookie", [`accessToken=${adminAccessToken}`])
        .expect(404)
      
      const errorResponse = response.body as ErrorResponse
      expect(errorResponse.message).toBe("Location not found")
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
