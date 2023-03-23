/* eslint-disable sonarjs/no-duplicate-string */
/* eslint-disable max-lines-per-function */
import request from "supertest"
import {
  type CreateLocationDto,
  type GetAllPaginatedLocationDto,
  type Location,
  type UpdateLocationDto
} from "types-custom"
import Fixture, { type ErrorResponse, type TokensAndUser } from "./fixture"
import { v4 } from "uuid"
import { LocationEntity } from "mikro-orm-config"

describe("LocationsController (e2e)", () => {
  let fixture: Fixture

  let tokensAndUser: TokensAndUser
  let adminTokensAndUser: TokensAndUser

  let locations: LocationEntity[]

  const defaultPage = 1
  const defaultPageSize = 3

  beforeEach(() => {
    fixture = new Fixture()
    return fixture
      .init()
      .then(() => fixture.seedDatabase())
      .then(() => fixture.getTokens("User"))
      .then((data) => (tokensAndUser = data))
      .then(() => fixture.getTokens("Admin"))
      .then((data) => (adminTokensAndUser = data))
      .then(() => fixture.em.find(LocationEntity, {}))
      .then((data) => {
        locations = data
      })
  })

  afterEach(() => {
    return fixture.clearDatabase().then(() => fixture.app.close())
  })

  describe("/locations (GET)", () => {
    it("gets all Locations with default page size (200)", async () => {
      const response = await request(fixture.app.getHttpServer())
        .get("/locations")
        .set("Cookie", [`accessToken=${adminTokensAndUser.tokens.accessToken}`])
        .expect(200)

      const paginatedLocationsResponse =
        response.body as GetAllPaginatedLocationDto
      expect(paginatedLocationsResponse.page).toBe(defaultPage)
      expect(paginatedLocationsResponse.totalPages).toBe(
        Math.ceil(locations.length / defaultPageSize)
      )
      expect(paginatedLocationsResponse.locations.length).toBe(defaultPageSize)
    })

    it("gets certain page of all Locations (200)", async () => {
      const page = 2

      const response = await request(fixture.app.getHttpServer())
        .get("/locations")
        .query({ page })
        .set("Cookie", [`accessToken=${adminTokensAndUser.tokens.accessToken}`])
        .expect(200)

      const paginatedLocationsResponse =
        response.body as GetAllPaginatedLocationDto
      expect(paginatedLocationsResponse.page).toBe(page)
      expect(paginatedLocationsResponse.totalPages).toBe(
        Math.ceil(locations.length / defaultPageSize)
      )
      expect(paginatedLocationsResponse.locations.length).toBe(defaultPageSize)
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

      expect(locationResponse.id).toBe(tokensAndUser.user.location?.id)
      expect(locationResponse.name).toBeDefined()
      expect(locationResponse.available).toBeDefined()
      expect(locationResponse.capacity).toBeDefined()
      expect(locationResponse.userId).toBe(tokensAndUser.user.id)
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
      const location = await fixture.em.refresh(locations[0])
      if (!location) throw new Error("Location is undefined")

      const response = await request(fixture.app.getHttpServer())
        .get(`/locations/${location.id}`)
        .set("Cookie", [`accessToken=${tokensAndUser.tokens.accessToken}`])
        .expect(200)

      const locationResponse = response.body as Location
      expect(locationResponse.id).toBe(location.id)
      expect(locationResponse.name).toBe(location.name)
      expect(locationResponse.available).toBe(location.available)
      expect(locationResponse.capacity).toBe(location.capacity)
      expect(locationResponse.userId).toBe(location.user.id)
    })

    it("get Location as User (200)", async () => {
      const response = await request(fixture.app.getHttpServer())
        .get(`/locations/${tokensAndUser.user.location?.id ?? 0}`)
        .set("Cookie", [`accessToken=${tokensAndUser.tokens.accessToken}`])
        .expect(200)

      const locationResponse = response.body as Location
      expect(locationResponse.id).toBe(tokensAndUser.user.location?.id)
      expect(locationResponse.name).toBeDefined()
      expect(locationResponse.available).toBeDefined()
      expect(locationResponse.capacity).toBeDefined()
      expect(locationResponse.userId).toBe(tokensAndUser.user.id)
    })

    it("returns Location not found because Location doesn't exist (404)", async () => {
      const response = await request(fixture.app.getHttpServer())
        .get(`/locations/${v4()}`)
        .set("Cookie", [`accessToken=${adminTokensAndUser.tokens.accessToken}`])
        .expect(404)

      const errorResponse = response.body as ErrorResponse
      expect(errorResponse.message).toBe("Location not found")
    })

    it("returns Location not found because User doesn't own Location (404)", async () => {
      const location = locations.find(
        (location) => location.name === "TestLocation0"
      )
      if (!location) {
        throw new Error("Location is undefined")
      }

      const response = await request(fixture.app.getHttpServer())
        .get(`/locations/${location.id}`)
        .set("Cookie", [`accessToken=${tokensAndUser.tokens.accessToken}`])
        .expect(404)

      const errorResponse = response.body as ErrorResponse
      expect(errorResponse.message).toBe("Location not found")
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
      expect(locationResponse.id).toBeDefined()
      expect(locationResponse.name).toBe(data.name)
      expect(locationResponse.available).toBe(data.capacity)
      expect(locationResponse.capacity).toBe(data.capacity)
      expect(locationResponse.userId).toBeDefined()
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
      expect(errorResponse.message).toBe(
        "Location with this name already exists"
      )
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
      expect(errorResponse.message).toBe(
        "User with this username already exists"
      )
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
      expect(locationResponse.id).toBeDefined()
      expect(locationResponse.name).toBe(data.name)
      expect(locationResponse.available).toBeDefined()
      expect(locationResponse.capacity).toBeDefined()
      expect(locationResponse.userId).toBe(tokensAndUser.user.id)
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
      expect(errorResponse.message).toBe(
        "Location with this name already exists"
      )
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
      expect(errorResponse.message).toBe(
        "User with this username already exists"
      )
    })

    it("returns Location not found because User doesn't own Location (404)", async () => {
      const location = locations.find(
        (location) => location.name === "TestLocation0"
      )
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
      expect(locationResponse.id).toBe(id)
    })

    it("returns Location not found (404)", async () => {
      const response = await request(fixture.app.getHttpServer())
        .delete(`/locations/${v4()}`)
        .set("Cookie", [`accessToken=${adminTokensAndUser.tokens.accessToken}`])
        .expect(404)

      const errorResponse = response.body as ErrorResponse
      expect(errorResponse.message).toBe("Location not found")
    })

    it("returns Forbidden (403)", async () => {
      const id = locations[0].id

      await request(fixture.app.getHttpServer())
        .delete(`/locations/${id}`)
        .set("Cookie", [`accessToken=${tokensAndUser.tokens.accessToken}`])
        .expect(403)
    })
  })
})
