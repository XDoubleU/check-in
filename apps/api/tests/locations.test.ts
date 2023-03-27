/* eslint-disable max-lines-per-function */
/* eslint-disable sonarjs/no-duplicate-string */
import request from "supertest"
import {
  type CreateLocationDto,
  type GetAllPaginatedLocationDto,
  type Location,
  type UpdateLocationDto
} from "types-custom"
import Fixture, { type ErrorResponse } from "./config/fixture"
import { v4 } from "uuid"
import { LocationEntity } from "mikro-orm-config"
import { type UserAndTokens } from "../src/auth/auth.service"

describe("LocationsController (e2e)", () => {
  const fixture: Fixture = new Fixture()

  let userAndTokens: UserAndTokens
  let adminUserAndTokens: UserAndTokens

  let locations: LocationEntity[]

  const defaultPage = 1
  const defaultPageSize = 3

  beforeAll(() => {
    return fixture.beforeAll()
  })

  afterAll(() => {
    return fixture.afterAll()
  })

  beforeEach(() => {
    return fixture
      .beforeEach()
      .then(() => fixture.getTokens("User"))
      .then((data) => (userAndTokens = data))
      .then(() => fixture.getTokens("Admin"))
      .then((data) => (adminUserAndTokens = data))
      .then(() => fixture.em.find(LocationEntity, {}))
      .then((data) => {
        locations = data
      })
  })

  afterEach(() => {
    return fixture.afterEach()
  })

  describe("/locations/sse (GET)", () => {
    it("gets all locations as LocationUpdateEventDto (200)", async () => {
      const response = await request(fixture.app.getHttpServer())
        .get("/locations/sse")
        .expect(200)
      
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      const data = response.body as any[]

      // eslint-disable-next-line @typescript-eslint/no-unsafe-member-access
      expect(data[0].normalizedName).toBeDefined()
      // eslint-disable-next-line @typescript-eslint/no-unsafe-member-access
      expect(data[0].available).toBeDefined()
      // eslint-disable-next-line @typescript-eslint/no-unsafe-member-access
      expect(data[0].capacity).toBeDefined()
      // eslint-disable-next-line @typescript-eslint/no-unsafe-member-access
      expect(data[0].yesterdayFullAt).toBeDefined()
      // eslint-disable-next-line @typescript-eslint/no-unsafe-member-access
      expect(data[0].id).toBeUndefined()
    })
  })

  describe("/locations (GET)", () => {
    it("gets all Locations with default page size (200)", async () => {
      const response = await request(fixture.app.getHttpServer())
        .get("/locations")
        .set("Cookie", [`accessToken=${adminUserAndTokens.tokens.accessToken}`])
        .expect(200)

      const paginatedLocationsResponse =
        response.body as GetAllPaginatedLocationDto
      expect(paginatedLocationsResponse.pagination.current).toBe(defaultPage)
      expect(paginatedLocationsResponse.pagination.total).toBe(
        Math.ceil(locations.length / defaultPageSize)
      )
      expect(paginatedLocationsResponse.data.length).toBe(defaultPageSize)
    })

    it("gets certain page of all Locations (200)", async () => {
      const page = 2

      const response = await request(fixture.app.getHttpServer())
        .get("/locations")
        .query({ page })
        .set("Cookie", [`accessToken=${adminUserAndTokens.tokens.accessToken}`])
        .expect(200)

      const paginatedLocationsResponse =
        response.body as GetAllPaginatedLocationDto
      expect(paginatedLocationsResponse.pagination.current).toBe(page)
      expect(paginatedLocationsResponse.pagination.total).toBe(
        Math.ceil(locations.length / defaultPageSize)
      )
      expect(paginatedLocationsResponse.data.length).toBe(defaultPageSize)
    })

    it("returns Forbidden (403)", async () => {
      return await request(fixture.app.getHttpServer())
        .get("/locations")
        .set("Cookie", [`accessToken=${userAndTokens.tokens.accessToken}`])
        .expect(403)
    })
  })

  describe("/locations/:id (GET)", () => {
    it("get Location as Admin (200)", async () => {
      const location = await fixture.em.findOneOrFail(
        LocationEntity,
        locations[0].id
      )

      const response = await request(fixture.app.getHttpServer())
        .get(`/locations/${location.id}`)
        .set("Cookie", [`accessToken=${adminUserAndTokens.tokens.accessToken}`])
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
        .get(`/locations/${userAndTokens.user.location?.id ?? 0}`)
        .set("Cookie", [`accessToken=${userAndTokens.tokens.accessToken}`])
        .expect(200)

      const locationResponse = response.body as Location
      expect(locationResponse.id).toBe(userAndTokens.user.location?.id)
      expect(locationResponse.name).toBeDefined()
      expect(locationResponse.available).toBeDefined()
      expect(locationResponse.capacity).toBeDefined()
      expect(locationResponse.userId).toBe(userAndTokens.user.id)
    })

    it("returns Location not found because Location doesn't exist (404)", async () => {
      const response = await request(fixture.app.getHttpServer())
        .get(`/locations/${v4()}`)
        .set("Cookie", [`accessToken=${adminUserAndTokens.tokens.accessToken}`])
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
        .set("Cookie", [`accessToken=${userAndTokens.tokens.accessToken}`])
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
        .set("Cookie", [`accessToken=${adminUserAndTokens.tokens.accessToken}`])
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
        .set("Cookie", [`accessToken=${adminUserAndTokens.tokens.accessToken}`])
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
        username: userAndTokens.user.username,
        password: "testpassword"
      }

      const response = await request(fixture.app.getHttpServer())
        .post("/locations")
        .set("Cookie", [`accessToken=${adminUserAndTokens.tokens.accessToken}`])
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
        .set("Cookie", [`accessToken=${userAndTokens.tokens.accessToken}`])
        .send(data)
        .expect(403)
    })
  })

  describe("/locations/:id (PATCH)", () => {
    it("updates a Location (200)", async () => {
      const id = userAndTokens.user.location?.id ?? 0

      const data: UpdateLocationDto = {
        username: "NewTestLocationUser2"
      }

      const response = await request(fixture.app.getHttpServer())
        .patch(`/locations/${id}`)
        .set("Cookie", [`accessToken=${userAndTokens.tokens.accessToken}`])
        .send(data)
        .expect(200)

      const locationResponse = response.body as Location
      expect(locationResponse.id).toBeDefined()
      expect(locationResponse.name).toBeDefined()
      expect(locationResponse.available).toBeDefined()
      expect(locationResponse.capacity).toBeDefined()
      expect(locationResponse.userId).toBe(userAndTokens.user.id)
    })

    it("updates a Location as admin (200)", async () => {
      const id = userAndTokens.user.location?.id ?? 0

      const data: UpdateLocationDto = {
        name: "NewTestLocation2",
        capacity: 100
      }

      const response = await request(fixture.app.getHttpServer())
        .patch(`/locations/${id}`)
        .set("Cookie", [`accessToken=${adminUserAndTokens.tokens.accessToken}`])
        .send(data)
        .expect(200)

      const locationResponse = response.body as Location
      expect(locationResponse.id).toBeDefined()
      expect(locationResponse.name).toBe(data.name)
      expect(locationResponse.available).toBeDefined()
      expect(locationResponse.capacity).toBe(100)
      expect(locationResponse.userId).toBe(userAndTokens.user.id)
    })

    it("returns Location with this name already exists (409)", async () => {
      const id = userAndTokens.user.location?.id ?? 0

      const data: UpdateLocationDto = {
        name: locations[1].name,
        username: "NewTestLocationUser3"
      }

      const response = await request(fixture.app.getHttpServer())
        .patch(`/locations/${id}`)
        .set("Cookie", [`accessToken=${userAndTokens.tokens.accessToken}`])
        .send(data)
        .expect(409)

      const errorResponse = response.body as ErrorResponse
      expect(errorResponse.message).toBe(
        "Location with this name already exists"
      )
    })

    it("returns User with this username already exists (409)", async () => {
      const id = userAndTokens.user.location?.id ?? 0

      const data: UpdateLocationDto = {
        name: "NewTestLocation3",
        username: adminUserAndTokens.user.username
      }

      const response = await request(fixture.app.getHttpServer())
        .patch(`/locations/${id}`)
        .set("Cookie", [`accessToken=${userAndTokens.tokens.accessToken}`])
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

      const response = await request(fixture.app.getHttpServer())
        .patch(`/locations/${location.id}`)
        .set("Cookie", [`accessToken=${userAndTokens.tokens.accessToken}`])
        .send(data)
        .expect(404)
      
      const errorResponse = response.body as ErrorResponse
      expect(errorResponse.message).toBe(
        "Location not found"
      )
    })
  })

  describe("/locations/:id (DELETE)", () => {
    it("deletes a Location (200)", async () => {
      const id = locations[0].id

      const response = await request(fixture.app.getHttpServer())
        .delete(`/locations/${id}`)
        .set("Cookie", [`accessToken=${adminUserAndTokens.tokens.accessToken}`])
        .expect(200)

      const locationResponse = response.body as Location
      expect(locationResponse.id).toBe(id)
    })

    it("returns Location not found (404)", async () => {
      const response = await request(fixture.app.getHttpServer())
        .delete(`/locations/${v4()}`)
        .set("Cookie", [`accessToken=${adminUserAndTokens.tokens.accessToken}`])
        .expect(404)

      const errorResponse = response.body as ErrorResponse
      expect(errorResponse.message).toBe("Location not found")
    })

    it("returns Forbidden (403)", async () => {
      const id = locations[0].id

      await request(fixture.app.getHttpServer())
        .delete(`/locations/${id}`)
        .set("Cookie", [`accessToken=${userAndTokens.tokens.accessToken}`])
        .expect(403)
    })
  })
})
