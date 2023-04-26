/* eslint-disable sonarjs/no-duplicate-string */
/* eslint-disable max-lines-per-function */
import request from "supertest"
import {
  Role,
  type GetAllPaginatedUserDto,
  type User,
  type CreateUserDto,
  type UpdateUserDto
} from "types-custom"
import { type UserAndTokens } from "../src/auth/auth.service"
import Fixture, { type ErrorResponse } from "./config/fixture"
import { UserEntity } from "../src/entities"

describe("UsersController (e2e)", () => {
  const fixture: Fixture = new Fixture()

  let userAndTokens: UserAndTokens
  let managerUserAndTokens: UserAndTokens
  let adminUserAndTokens: UserAndTokens

  let managerUsers: User[]

  const defaultPage = 1
  const defaultPageSize = 4

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
      .then(() => fixture.getTokens("Manager"))
      .then((data) => (managerUserAndTokens = data))
      .then(() => fixture.getAdminTokens())
      .then((data) => (adminUserAndTokens = data))
      .then(() =>
        fixture.em.find(UserEntity, {
          roles: {
            $contains: [Role.Manager]
          }
        })
      )
      .then((data) => {
        managerUsers = data
      })
  })

  afterEach(() => {
    return fixture.afterEach()
  })

  describe("/users/me (GET)", () => {
    it("gets info about logged in User (200)", async () => {
      const response = await request(fixture.app.getHttpServer())
        .get("/users/me")
        .set("Cookie", [`accessToken=${userAndTokens.tokens.accessToken}`])
        .expect(200)

      const userResponse = response.body as User
      expect(userResponse.id).toBe(userAndTokens.user.id)
      expect(userResponse.username).toBe(userAndTokens.user.username)
      expect(userResponse.roles).toStrictEqual(userAndTokens.user.roles)
      expect(userResponse.locationId).toBe(userAndTokens.user.location?.id)
    })
  })

  describe("/users/:id (GET)", () => {
    it("gets User as manager (200)", async () => {
      const response = await request(fixture.app.getHttpServer())
        .get(`/users/${userAndTokens.user.id}`)
        .set("Cookie", [
          `accessToken=${managerUserAndTokens.tokens.accessToken}`
        ])
        .expect(200)

      const userResponse = response.body as User
      expect(userResponse.id).toBe(userAndTokens.user.id)
      expect(userResponse.username).toBe(userAndTokens.user.username)
      expect(userResponse.roles).toStrictEqual(userAndTokens.user.roles)
      expect(userResponse.locationId).toBe(userAndTokens.user.location?.id)
    })

    it("returns User not found (404)", async () => {
      const response = await request(fixture.app.getHttpServer())
        .get("/users/random")
        .set("Cookie", [
          `accessToken=${managerUserAndTokens.tokens.accessToken}`
        ])
        .expect(404)

      const errorResponse = response.body as ErrorResponse
      expect(errorResponse.message).toBe("User not found")
    })
  })

  describe("/users (GET)", () => {
    it("gets all manager Users with default page size (200)", async () => {
      const response = await request(fixture.app.getHttpServer())
        .get("/users")
        .set("Cookie", [`accessToken=${adminUserAndTokens.tokens.accessToken}`])
        .expect(200)

      const paginatedManagerUsersResponse =
        response.body as GetAllPaginatedUserDto
      expect(paginatedManagerUsersResponse.pagination.current).toBe(defaultPage)
      expect(paginatedManagerUsersResponse.pagination.total).toBe(
        Math.ceil(managerUsers.length / defaultPageSize)
      )
      expect(paginatedManagerUsersResponse.data.length).toBe(defaultPageSize)
    })

    it("gets certain page of all manager Users (200)", async () => {
      const page = 2

      const response = await request(fixture.app.getHttpServer())
        .get("/users")
        .query({ page })
        .set("Cookie", [`accessToken=${adminUserAndTokens.tokens.accessToken}`])
        .expect(200)

      const paginatedManagerUsersResponse =
        response.body as GetAllPaginatedUserDto
      expect(paginatedManagerUsersResponse.pagination.current).toBe(page)
      expect(paginatedManagerUsersResponse.pagination.total).toBe(
        Math.ceil(managerUsers.length / defaultPageSize)
      )
      expect(paginatedManagerUsersResponse.data.length).toBe(defaultPageSize)
    })

    it("returns Page should be greater than 0 (400)", async () => {
      const page = 0

      const response = await request(fixture.app.getHttpServer())
        .get("/users")
        .query({ page })
        .set("Cookie", [
          `accessToken=${adminUserAndTokens.tokens.accessToken}`
        ])
        .expect(400)

        const errorResponse = response.body as ErrorResponse
        expect(errorResponse.message).toBe("Page should be greater than 0")
    })
  })

  describe("/users (POST)", () => {
    it("creates a new User (201)", async () => {
      const data: CreateUserDto = {
        username: "ManagerUser",
        password: "testpassword"
      }

      const response = await request(fixture.app.getHttpServer())
        .post("/users")
        .set("Cookie", [`accessToken=${adminUserAndTokens.tokens.accessToken}`])
        .send(data)
        .expect(201)

      const userResponse = response.body as User
      expect(userResponse.id).toBeDefined()
      expect(userResponse.username).toBe(data.username)
      expect(userResponse.passwordHash).toBeUndefined()
    })

    it("returns User with this username already exists (409)", async () => {
      const data: CreateUserDto = {
        username: "Manager",
        password: "testpassword"
      }

      const response = await request(fixture.app.getHttpServer())
        .post("/users")
        .set("Cookie", [`accessToken=${adminUserAndTokens.tokens.accessToken}`])
        .send(data)
        .expect(409)

      const errorResponse = response.body as ErrorResponse
      expect(errorResponse.message).toBe(
        "User with this username already exists"
      )
    })
  })

  describe("/users/:id (PATCH)", () => {
    it("updates a User (200)", async () => {
      const id = managerUserAndTokens.user.id

      const data: UpdateUserDto = {
        username: "Manager2",
        password: "newPassword"
      }

      const response = await request(fixture.app.getHttpServer())
        .patch(`/users/${id}`)
        .set("Cookie", [`accessToken=${adminUserAndTokens.tokens.accessToken}`])
        .send(data)
        .expect(200)

      const userResponse = response.body as User
      expect(userResponse.id).toBe(id)
      expect(userResponse.username).toBe(data.username)
    })

    it("returns User with this username already exists (409)", async () => {
      const id = managerUserAndTokens.user.id

      const data: UpdateUserDto = {
        username: managerUserAndTokens.user.username
      }

      const response = await request(fixture.app.getHttpServer())
        .patch(`/users/${id}`)
        .set("Cookie", [`accessToken=${adminUserAndTokens.tokens.accessToken}`])
        .send(data)
        .expect(409)

      const errorResponse = response.body as ErrorResponse
      expect(errorResponse.message).toBe(
        "User with this username already exists"
      )
    })

    it("returns User not found (404)", async () => {
      const data: UpdateUserDto = {
        username: "Manager2",
        password: "newPassword"
      }

      const response = await request(fixture.app.getHttpServer())
        .patch("/users/random")
        .set("Cookie", [`accessToken=${adminUserAndTokens.tokens.accessToken}`])
        .send(data)
        .expect(404)

      const errorResponse = response.body as ErrorResponse
      expect(errorResponse.message).toBe("User not found")
    })
  })

  describe("/users/:id (DELETE)", () => {
    it("deletes a User (200)", async () => {
      const id = managerUserAndTokens.user.id

      const response = await request(fixture.app.getHttpServer())
        .delete(`/users/${id}`)
        .set("Cookie", [`accessToken=${adminUserAndTokens.tokens.accessToken}`])
        .expect(200)

      const userResponse = response.body as User
      expect(userResponse.id).toBe(id)
    })

    it("returns User not found (404)", async () => {
      const response = await request(fixture.app.getHttpServer())
        .delete("/users/random")
        .set("Cookie", [`accessToken=${adminUserAndTokens.tokens.accessToken}`])
        .expect(404)

      const errorResponse = response.body as ErrorResponse
      expect(errorResponse.message).toBe("User not found")
    })
  })
})
