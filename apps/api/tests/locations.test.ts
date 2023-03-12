import request from "supertest"
import { CreateLocationDto, GetAllPaginatedLocationDto, Location, UpdateLocationDto } from "types-custom"
import Fixture, { ErrorResponse, TokensAndUser } from "./fixture"
import { v4 } from "uuid"
import { expect } from "chai"
import { LocationEntity } from "mikro-orm-config"


describe("LocationsController (e2e)", () => {
  let fixture: Fixture

  let tokensAndUser: TokensAndUser
  let adminTokensAndUser: TokensAndUser

  let locations: LocationEntity[]

  const defaultPage = 1
  const defaultPageSize = 3
  
  before(() => {
    fixture = new Fixture()
    return fixture.init()
      .then(() => fixture.seedDatabase())
      .then(() => fixture.getTokens("User"))
      .then((data) => tokensAndUser = data)
      .then(() => fixture.getTokens("Admin"))
      .then((data) => adminTokensAndUser = data)
      .then(() => fixture.em.find(LocationEntity, {}))
      .then((data) => {
        locations = data
      })
  })

  after(() => {
    return fixture.clearDatabase()
  })

  describe("/locations (GET)", () => {
    it("gets all Locations with default page size (200)", async () => {
      const response = await request(fixture.app.getHttpServer())
        .get("/locations")
        .set("Cookie", [`accessToken=${adminTokensAndUser.tokens.accessToken}`])
        .expect(200)
      
      const paginatedLocationsResponse = response.body as GetAllPaginatedLocationDto
      expect(paginatedLocationsResponse.page).to.be.equal(defaultPage)
      expect(paginatedLocationsResponse.totalPages).to.be.equal(Math.ceil(locations.length / defaultPageSize))
      expect(paginatedLocationsResponse.locations.length).to.be.equal(defaultPageSize)
    })

    it("gets certain page of all Locations (200)", async () => {
      const page = 2

      const response = await request(fixture.app.getHttpServer())
        .get("/locations")
        .query({ page })
        .set("Cookie", [`accessToken=${adminTokensAndUser.tokens.accessToken}`])
        .expect(200)
      
      const paginatedLocationsResponse = response.body as GetAllPaginatedLocationDto
      expect(paginatedLocationsResponse.page).to.be.equal(page)
      expect(paginatedLocationsResponse.totalPages).to.be.equal(Math.ceil(locations.length / defaultPageSize))
      expect(paginatedLocationsResponse.locations.length).to.be.equal(defaultPageSize)
    })

    it("returns Forbidden (403)", async () => {
      return await request(fixture.app.getHttpServer())
        .get("/locations")
        .set("Cookie", [`accessToken=${tokensAndUser.tokens.accessToken}`])
        .expect(403)
    })
  })

  describe("/locations/me (GET)", () => {
    it("get my Location (200)", async () => {
      const response = await request(fixture.app.getHttpServer())
        .get("/locations/me")
        .set("Cookie", [`accessToken=${tokensAndUser.tokens.accessToken}`])
        .expect(200)

      const locationResponse = response.body as Location

      expect(locationResponse.id).to.be.equal(tokensAndUser.user.location?.id)
      expect(locationResponse.name).to.exist
      expect(locationResponse.available).to.exist
      expect(locationResponse.capacity).to.exist
      expect(locationResponse.userId).to.be.equal(tokensAndUser.user.id)
    })

    it("returns Forbidden (403)", async () => {
      return await request(fixture.app.getHttpServer())
        .get("/locations/me")
        .set("Cookie", [`accessToken=${adminTokensAndUser.tokens.accessToken}`])
        .expect(403)
    })
  })

  describe("/locations/:id (GET)", () => {
    it("get Location as Admin (200)", async () => {
      const location = locations[0]
      
      const response = await request(fixture.app.getHttpServer())
        .get(`/locations/${location.id}`)
        .set("Cookie", [`accessToken=${tokensAndUser.tokens.accessToken}`])
        .expect(200)

      const locationResponse = response.body as Location
      expect(locationResponse.id).to.be.equal(location.id)
      expect(locationResponse.name).to.be.equal(location.name)
      expect(locationResponse.available).to.be.equal(location.available)
      expect(locationResponse.capacity).to.be.equal(location.capacity)
      expect(locationResponse.userId).to.be.equal(location.user.id)
    })

    it("get Location as User (200)", async () => {
      const response = await request(fixture.app.getHttpServer())
        .get(`/locations/${tokensAndUser.user.location?.id ?? 0}`)
        .set("Cookie", [`accessToken=${tokensAndUser.tokens.accessToken}`])
        .expect(200)

      const locationResponse = response.body as Location
      expect(locationResponse.id).to.be.equal(tokensAndUser.user.location?.id)
      expect(locationResponse.name).to.exist
      expect(locationResponse.available).to.exist
      expect(locationResponse.capacity).to.exist
      expect(locationResponse.userId).to.be.equal(tokensAndUser.user.id)
    })

    it("returns Location not found because Location doesn't exist (404)", async () => {
      const response = await request(fixture.app.getHttpServer())
        .get(`/locations/${v4()}`)
        .set("Cookie", [`accessToken=${adminTokensAndUser.tokens.accessToken}`])
        .expect(404)
      
      const errorResponse = response.body as ErrorResponse
      expect(errorResponse.message).to.be.equal("Location not found")
    })

    it("returns Location not found because User doesn't own Location (404)", async () => {
      const location = locations.find((location) => location.name === "TestLocation0")
      if (!location) {
        throw new Error("Location is undefined")
      }

      const response = await request(fixture.app.getHttpServer())
        .get(`/locations/${location.id}`)
        .set("Cookie", [`accessToken=${tokensAndUser.tokens.accessToken}`])
        .expect(404)
      
      const errorResponse = response.body as ErrorResponse
      expect(errorResponse.message).to.be.equal("Location not found")
    })
  })

  describe("/locations (POST)", () => {
    it("creates a new Location (201)", async () => {
      const data: CreateLocationDto = {
        name: "NewTestLocation",
        capacity: 10,
        username: "NewTestLocationUser",
        password: "testpassword"
      }
  
      const response = await request(fixture.app.getHttpServer())
        .post("/locations")
        .set("Cookie", [`accessToken=${adminTokensAndUser.tokens.accessToken}`])
        .send(data)
        .expect(201)
      
      const locationResponse = response.body as Location
      expect(locationResponse.id).to.exist
      expect(locationResponse.name).to.be.equal(data.name)
      expect(locationResponse.available).to.be.equal(data.capacity)
      expect(locationResponse.capacity).to.be.equal(data.capacity)
      expect(locationResponse.userId).to.exist
    })

    it("returns Location with this name already exists (409)", async () => {
      const data: CreateLocationDto = {
        name: locations[0].name,
        capacity: 10,
        username: "NewTestLocationUser",
        password: "testpassword"
      }
  
      const response = await request(fixture.app.getHttpServer())
        .post("/locations")
        .set("Cookie", [`accessToken=${adminTokensAndUser.tokens.accessToken}`])
        .send(data)
        .expect(409)
      
      const errorResponse = response.body as ErrorResponse
      expect(errorResponse.message).to.be.equal("Location with this name already exists")
    })

    it("returns User with this username already exists (409)", async () => {
      const data: CreateLocationDto = {
        name: "NewTestLocation2",
        capacity: 10,
        username: locations[0].user.username,
        password: "testpassword"
      }
  
      const response = await request(fixture.app.getHttpServer())
        .post("/locations")
        .set("Cookie", [`accessToken=${adminTokensAndUser.tokens.accessToken}`])
        .send(data)
        .expect(409)
      
      const errorResponse = response.body as ErrorResponse
      expect(errorResponse.message).to.be.equal("User with this username already exists")
    })

    it("returns Forbidden (403)", async () => {
      const data: CreateLocationDto = {
        name: "TestLocation",
        capacity: 10,
        username: "TestLocationUser",
        password: "testpassword"
      }
  
      return await request(fixture.app.getHttpServer())
        .post("/locations")
        .set("Cookie", [`accessToken=${tokensAndUser.tokens.accessToken}`])
        .send(data)
        .expect(403)
    })
  })

  describe("/locations/:id (PATCH)", () => {
    it("updates a Location (200)", async () => {
      const id = tokensAndUser.user.location?.id ?? 0
      
      const data: UpdateLocationDto = {
        name: "NewTestLocation2",
        username: "NewTestLocationUser2"
      }
  
      const response = await request(fixture.app.getHttpServer())
        .patch(`/locations/${id}`)
        .set("Cookie", [`accessToken=${tokensAndUser.tokens.accessToken}`])
        .send(data)
        .expect(200)
      
      const locationResponse = response.body as Location
      expect(locationResponse.id).to.exist
      expect(locationResponse.name).to.be.equal(data.name)
      expect(locationResponse.available).to.exist
      expect(locationResponse.capacity).to.exist
      expect(locationResponse.userId).to.be.equal(tokensAndUser.user.id)
    })

    it("returns Location with this name already exists (409)", async () => {
      const id = tokensAndUser.user.location?.id ?? 0 
      
      const data: UpdateLocationDto = {
        name: locations[1].name,
        username: "NewTestLocationUser3"
      }
  
      const response = await request(fixture.app.getHttpServer())
        .patch(`/locations/${id}`)
        .set("Cookie", [`accessToken=${tokensAndUser.tokens.accessToken}`])
        .send(data)
        .expect(409)
      
      const errorResponse = response.body as ErrorResponse
      expect(errorResponse.message).to.be.equal("Location with this name already exists")
    })

    it("returns User with this username already exists (409)", async () => {
      const id = tokensAndUser.user.location?.id ?? 0
      
      const data: UpdateLocationDto = {
        name: "NewTestLocation3",
        username: locations[1].user.username
      }
  
      const response = await request(fixture.app.getHttpServer())
        .patch(`/locations/${id}`)
        .set("Cookie", [`accessToken=${tokensAndUser.tokens.accessToken}`])
        .send(data)
        .expect(409)
      
      const errorResponse = response.body as ErrorResponse
      expect(errorResponse.message).to.be.equal("User with this username already exists")
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
  
      return await request(fixture.app.getHttpServer())
        .patch(`/locations/${location.id}`)
        .set("Cookie", [`accessToken=${tokensAndUser.tokens.accessToken}`])
        .send(data)
        .expect(404)
    })
  })

  describe("/locations/:id (DELETE)", () => {
    it("deletes a Location (200)", async () => {
      const id = locations[0].id
  
      const response = await request(fixture.app.getHttpServer())
        .delete(`/locations/${id}`)
        .set("Cookie", [`accessToken=${adminTokensAndUser.tokens.accessToken}`])
        .expect(200)
      
      const locationResponse = response.body as Location
      expect(locationResponse.id).to.be.equal(id)
    })

    it("returns Location not found (404)", async () => {
      const response = await request(fixture.app.getHttpServer())
        .delete(`/locations/${v4()}`)
        .set("Cookie", [`accessToken=${adminTokensAndUser.tokens.accessToken}`])
        .expect(404)
      
      const errorResponse = response.body as ErrorResponse
      expect(errorResponse.message).to.be.equal("Location not found")
    })

    it("returns Unauthorized (401)", async () => {
      const id = locations[0].id
  
      await request(fixture.app.getHttpServer())
        .delete(`/locations/${id}`)
        .set("Cookie", [`accessToken=${tokensAndUser.tokens.accessToken}`])
        .expect(401)
    })
  })
})
