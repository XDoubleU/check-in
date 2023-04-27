/* eslint-disable sonarjs/no-duplicate-string */
/* eslint-disable max-lines-per-function */
import request from "supertest"
import { DATE_FORMAT, type CheckIn, type CreateCheckInDto } from "types-custom"
import Fixture, { type ErrorResponse } from "./config/fixture"
import { SchoolEntity } from "../src/entities"
import { type UserAndTokens } from "../src/auth/auth.service"
import { add, format } from "date-fns"
import { v4 } from "uuid"

describe("CheckInsController (e2e)", () => {
  const fixture: Fixture = new Fixture()

  let userAndTokens: UserAndTokens
  let managerUserAndTokens: UserAndTokens

  let school: SchoolEntity

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
      .then(() => fixture.em.findOne(SchoolEntity, { id: 1 }))
      .then((data) => {
        if (!data) {
          throw new Error("School is null")
        }

        school = data
      })
  })

  afterEach(() => {
    return fixture.afterEach()
  })

  describe("/checkins/range/:locationId (GET)", () => {
    it("fetches data for range chart as owner (200)", async () => {
      const startDate = format(new Date(), DATE_FORMAT)
      const endDate = format(add(new Date(), { days: 1 }), DATE_FORMAT)

      const response = await request(fixture.app.getHttpServer())
        .get(
          `/checkins/range/${
            userAndTokens.user.location?.id ?? ""
          }?startDate=${startDate}&endDate=${endDate}`
        )
        .set("Cookie", [`accessToken=${userAndTokens.tokens.accessToken}`])
        .expect(200)

      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      const data = response.body as any[]

      // eslint-disable-next-line @typescript-eslint/no-unsafe-member-access
      expect(data[0].datetime).toBeDefined()
      // eslint-disable-next-line @typescript-eslint/no-unsafe-member-access
      expect(data[0].capacity).toBeDefined()
      // eslint-disable-next-line @typescript-eslint/no-unsafe-member-access
      expect(data[0].Andere).toBeDefined()
    })

    it("fetches data for range chart as manager (200)", async () => {
      const startDate = format(new Date(), DATE_FORMAT)
      const endDate = format(add(new Date(), { days: 1 }), DATE_FORMAT)

      const response = await request(fixture.app.getHttpServer())
        .get(
          `/checkins/range/${
            userAndTokens.user.location?.id ?? ""
          }?startDate=${startDate}&endDate=${endDate}`
        )
        .set("Cookie", [
          `accessToken=${managerUserAndTokens.tokens.accessToken}`
        ])
        .expect(200)

      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      const data = response.body as any[]

      // eslint-disable-next-line @typescript-eslint/no-unsafe-member-access
      expect(data[0].datetime).toBeDefined()
      // eslint-disable-next-line @typescript-eslint/no-unsafe-member-access
      expect(data[0].capacity).toBeDefined()
      // eslint-disable-next-line @typescript-eslint/no-unsafe-member-access
      expect(data[0].Andere).toBeDefined()
    })

    it("returns startDate and endDate need to be a valid date - missing value (400)", async () => {
      const response = await request(fixture.app.getHttpServer())
        .get(`/checkins/range/${userAndTokens.user.location?.id ?? ""}`)
        .set("Cookie", [`accessToken=${userAndTokens.tokens.accessToken}`])
        .expect(400)

      const errorResponse = response.body as ErrorResponse
      expect(errorResponse.message).toBe("Validation failed (Date is expected)")
    })

    it("returns startDate and endDate need to be a valid date - invalid value (400)", async () => {
      const response = await request(fixture.app.getHttpServer())
        .get(
          `/checkins/range/${
            userAndTokens.user.location?.id ?? ""
          }?startDate=random&endDate=random`
        )
        .set("Cookie", [
          `accessToken=${managerUserAndTokens.tokens.accessToken}`
        ])
        .expect(400)

      const errorResponse = response.body as ErrorResponse
      expect(errorResponse.message).toBe("Validation failed (Date is expected)")
    })

    it("returns Location not found (404)", async () => {
      const startDate = format(new Date(), DATE_FORMAT)
      const endDate = format(add(new Date(), { days: 1 }), DATE_FORMAT)

      const response = await request(fixture.app.getHttpServer())
        .get(
          `/checkins/range/${v4()}?startDate=${startDate}&endDate=${endDate}`
        )
        .set("Cookie", [`accessToken=${userAndTokens.tokens.accessToken}`])
        .expect(404)

      const errorResponse = response.body as ErrorResponse
      expect(errorResponse.message).toBe("Location not found")
    })

    it("returns Bad request, id is not uuid (400)", async () => {
      const startDate = format(new Date(), DATE_FORMAT)
      const endDate = format(add(new Date(), { days: 1 }), DATE_FORMAT)

      const response = await request(fixture.app.getHttpServer())
        .get(`/checkins/range/random?startDate=${startDate}&endDate=${endDate}`)
        .set("Cookie", [`accessToken=${userAndTokens.tokens.accessToken}`])
        .expect(400)

      const errorResponse = response.body as ErrorResponse
      expect(errorResponse.message).toBe("Validation failed (uuid is expected)")
    })
  })

  describe("/checkins/csv/range/:locationId (GET)", () => {
    it("fetches csv with data from range chart as owner (200)", async () => {
      const startDate = format(new Date(), DATE_FORMAT)
      const endDate = format(add(new Date(), { days: 1 }), DATE_FORMAT)

      const response = await request(fixture.app.getHttpServer())
        .get(
          `/checkins/csv/range/${
            userAndTokens.user.location?.id ?? ""
          }?startDate=${startDate}&endDate=${endDate}`
        )
        .set("Cookie", [`accessToken=${userAndTokens.tokens.accessToken}`])
        .expect(200)

      // eslint-disable-next-line @typescript-eslint/no-unsafe-member-access
      expect(response.headers["content-type"]).toContain("text/csv")
    })

    it("fetches csv with data from range chart as manager (200)", async () => {
      const startDate = format(new Date(), DATE_FORMAT)
      const endDate = format(add(new Date(), { days: 1 }), DATE_FORMAT)

      const response = await request(fixture.app.getHttpServer())
        .get(
          `/checkins/csv/range/${
            userAndTokens.user.location?.id ?? ""
          }?startDate=${startDate}&endDate=${endDate}`
        )
        .set("Cookie", [
          `accessToken=${managerUserAndTokens.tokens.accessToken}`
        ])
        .expect(200)

      // eslint-disable-next-line @typescript-eslint/no-unsafe-member-access
      expect(response.headers["content-type"]).toContain("text/csv")
    })

    it("returns startDate and endDate need to be a valid date - missing value (400)", async () => {
      const response = await request(fixture.app.getHttpServer())
        .get(`/checkins/csv/range/${userAndTokens.user.location?.id ?? ""}`)
        .set("Cookie", [`accessToken=${userAndTokens.tokens.accessToken}`])
        .expect(400)

      const errorResponse = response.body as ErrorResponse
      expect(errorResponse.message).toBe("Validation failed (Date is expected)")
    })

    it("returns startDate and endDate need to be a valid date - invalid value (400)", async () => {
      const response = await request(fixture.app.getHttpServer())
        .get(
          `/checkins/csv/range/${
            userAndTokens.user.location?.id ?? ""
          }?startDate=random&endDate=random`
        )
        .set("Cookie", [
          `accessToken=${managerUserAndTokens.tokens.accessToken}`
        ])
        .expect(400)

      const errorResponse = response.body as ErrorResponse
      expect(errorResponse.message).toBe("Validation failed (Date is expected)")
    })

    it("returns Location not found (404)", async () => {
      const startDate = format(new Date(), DATE_FORMAT)
      const endDate = format(add(new Date(), { days: 1 }), DATE_FORMAT)

      const response = await request(fixture.app.getHttpServer())
        .get(
          `/checkins/csv/range/${v4()}?startDate=${startDate}&endDate=${endDate}`
        )
        .set("Cookie", [`accessToken=${userAndTokens.tokens.accessToken}`])
        .expect(404)

      const errorResponse = response.body as ErrorResponse
      expect(errorResponse.message).toBe("Location not found")
    })

    it("returns Bad request, id is not uuid (400)", async () => {
      const startDate = format(new Date(), DATE_FORMAT)
      const endDate = format(add(new Date(), { days: 1 }), DATE_FORMAT)

      const response = await request(fixture.app.getHttpServer())
        .get(
          `/checkins/csv/range/random?startDate=${startDate}&endDate=${endDate}`
        )
        .set("Cookie", [`accessToken=${userAndTokens.tokens.accessToken}`])
        .expect(400)

      const errorResponse = response.body as ErrorResponse
      expect(errorResponse.message).toBe("Validation failed (uuid is expected)")
    })
  })

  describe("/checkins/day/:locationId (GET)", () => {
    it("fetches data for day chart as owner (200)", async () => {
      const date = format(new Date(), DATE_FORMAT)

      const response = await request(fixture.app.getHttpServer())
        .get(
          `/checkins/day/${userAndTokens.user.location?.id ?? ""}?date=${date}`
        )
        .set("Cookie", [`accessToken=${userAndTokens.tokens.accessToken}`])
        .expect(200)

      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      const data = response.body as any[]

      // eslint-disable-next-line @typescript-eslint/no-unsafe-member-access
      expect(data[0].datetime).toBeDefined()
      // eslint-disable-next-line @typescript-eslint/no-unsafe-member-access
      expect(data[0].capacity).toBeDefined()
      // eslint-disable-next-line @typescript-eslint/no-unsafe-member-access
      expect(data[0].Andere).toBeDefined()
    })

    it("fetches data for day chart as manager (200)", async () => {
      const date = format(new Date(), DATE_FORMAT)

      const response = await request(fixture.app.getHttpServer())
        .get(
          `/checkins/day/${userAndTokens.user.location?.id ?? ""}?date=${date}`
        )
        .set("Cookie", [
          `accessToken=${managerUserAndTokens.tokens.accessToken}`
        ])
        .expect(200)

      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      const data = response.body as any[]

      // eslint-disable-next-line @typescript-eslint/no-unsafe-member-access
      expect(data[0].datetime).toBeDefined()
      // eslint-disable-next-line @typescript-eslint/no-unsafe-member-access
      expect(data[0].capacity).toBeDefined()
      // eslint-disable-next-line @typescript-eslint/no-unsafe-member-access
      expect(data[0].Andere).toBeDefined()
    })

    it("returns date needs to be a valid date - missing value (400)", async () => {
      const response = await request(fixture.app.getHttpServer())
        .get(`/checkins/day/${userAndTokens.user.location?.id ?? ""}`)
        .set("Cookie", [`accessToken=${userAndTokens.tokens.accessToken}`])
        .expect(400)

      const errorResponse = response.body as ErrorResponse
      expect(errorResponse.message).toBe("Validation failed (Date is expected)")
    })

    it("returns date needs to be a valid date - invalid value (400)", async () => {
      const response = await request(fixture.app.getHttpServer())
        .get(
          `/checkins/day/${userAndTokens.user.location?.id ?? ""}?date=random`
        )
        .set("Cookie", [
          `accessToken=${managerUserAndTokens.tokens.accessToken}`
        ])
        .expect(400)

      const errorResponse = response.body as ErrorResponse
      expect(errorResponse.message).toBe("Validation failed (Date is expected)")
    })

    it("returns Location not found (404)", async () => {
      const date = format(new Date(), DATE_FORMAT)

      const response = await request(fixture.app.getHttpServer())
        .get(`/checkins/day/${v4()}?date=${date}`)
        .set("Cookie", [`accessToken=${userAndTokens.tokens.accessToken}`])
        .expect(404)

      const errorResponse = response.body as ErrorResponse
      expect(errorResponse.message).toBe("Location not found")
    })

    it("returns Bad request, id is not uuid (400)", async () => {
      const date = format(new Date(), DATE_FORMAT)

      const response = await request(fixture.app.getHttpServer())
        .get(`/checkins/day/random?date=${date}`)
        .set("Cookie", [`accessToken=${userAndTokens.tokens.accessToken}`])
        .expect(400)

      const errorResponse = response.body as ErrorResponse
      expect(errorResponse.message).toBe("Validation failed (uuid is expected)")
    })
  })

  describe("/checkins/csv/day/:locationId (GET)", () => {
    it("fetches csv with data from day chart as owner (200)", async () => {
      const date = format(new Date(), DATE_FORMAT)

      const response = await request(fixture.app.getHttpServer())
        .get(
          `/checkins/csv/day/${
            userAndTokens.user.location?.id ?? ""
          }?date=${date}`
        )
        .set("Cookie", [`accessToken=${userAndTokens.tokens.accessToken}`])
        .expect(200)

      // eslint-disable-next-line @typescript-eslint/no-unsafe-member-access
      expect(response.headers["content-type"]).toContain("text/csv")
    })

    it("fetches csv with data from day chart as manager (200)", async () => {
      const date = format(new Date(), DATE_FORMAT)

      const response = await request(fixture.app.getHttpServer())
        .get(
          `/checkins/csv/day/${
            userAndTokens.user.location?.id ?? ""
          }?date=${date}`
        )
        .set("Cookie", [
          `accessToken=${managerUserAndTokens.tokens.accessToken}`
        ])
        .expect(200)

      // eslint-disable-next-line @typescript-eslint/no-unsafe-member-access
      expect(response.headers["content-type"]).toContain("text/csv")
    })

    it("returns date needs to be a valid date - missing value (400)", async () => {
      const response = await request(fixture.app.getHttpServer())
        .get(`/checkins/csv/day/${userAndTokens.user.location?.id ?? ""}`)
        .set("Cookie", [`accessToken=${userAndTokens.tokens.accessToken}`])
        .expect(400)

      const errorResponse = response.body as ErrorResponse
      expect(errorResponse.message).toBe("Validation failed (Date is expected)")
    })

    it("returns date needs to be a valid date - invalid value (400)", async () => {
      const response = await request(fixture.app.getHttpServer())
        .get(
          `/checkins/csv/day/${
            userAndTokens.user.location?.id ?? ""
          }?date=random`
        )
        .set("Cookie", [
          `accessToken=${managerUserAndTokens.tokens.accessToken}`
        ])
        .expect(400)

      const errorResponse = response.body as ErrorResponse
      expect(errorResponse.message).toBe("Validation failed (Date is expected)")
    })

    it("returns Location not found (404)", async () => {
      const date = format(new Date(), DATE_FORMAT)

      const response = await request(fixture.app.getHttpServer())
        .get(`/checkins/csv/day/${v4()}?date=${date}`)
        .set("Cookie", [`accessToken=${userAndTokens.tokens.accessToken}`])
        .expect(404)

      const errorResponse = response.body as ErrorResponse
      expect(errorResponse.message).toBe("Location not found")
    })

    it("returns Bad request, id is not uuid (400)", async () => {
      const date = format(new Date(), DATE_FORMAT)

      const response = await request(fixture.app.getHttpServer())
        .get(`/checkins/csv/day/random?date=${date}`)
        .set("Cookie", [`accessToken=${userAndTokens.tokens.accessToken}`])
        .expect(400)

      const errorResponse = response.body as ErrorResponse
      expect(errorResponse.message).toBe("Validation failed (uuid is expected)")
    })
  })

  describe("/checkins (POST)", () => {
    it("creates a new CheckIn (201)", async () => {
      const data: CreateCheckInDto = {
        schoolId: school.id
      }

      const response = await request(fixture.app.getHttpServer())
        .post("/checkins")
        .set("Cookie", [`accessToken=${userAndTokens.tokens.accessToken}`])
        .send(data)
        .expect(201)

      const responseCheckIn = response.body as CheckIn
      expect(responseCheckIn.id).toBeDefined()
      expect(responseCheckIn.location.id).toBe(userAndTokens.user.location?.id)
      expect(responseCheckIn.capacity).toBe(
        userAndTokens.user.location?.capacity
      )
      expect(responseCheckIn.createdAt).toBeDefined()
      expect(responseCheckIn.school.id).toBe(school.id)
    })

    it("returns School not found (404)", async () => {
      const data: CreateCheckInDto = {
        schoolId: -1
      }

      const response = await request(fixture.app.getHttpServer())
        .post("/checkins")
        .set("Cookie", [`accessToken=${userAndTokens.tokens.accessToken}`])
        .send(data)
        .expect(404)

      const errorResponse = response.body as ErrorResponse
      expect(errorResponse.message).toBe("School not found")
    })

    it("returns Forbidden (403)", async () => {
      const data: CreateCheckInDto = {
        schoolId: school.id
      }

      return await request(fixture.app.getHttpServer())
        .post("/checkins")
        .set("Cookie", [
          `accessToken=${managerUserAndTokens.tokens.accessToken}`
        ])
        .send(data)
        .expect(403)
    })
  })
})
