/* eslint-disable max-lines-per-function */
import Fixture from "./config/fixture"
import EventSource from "eventsource"
import { type Server } from "http"
import { type AddressInfo } from "net"
import { SseService } from "../src/sse/sse.service"
import { CheckInEntity, LocationEntity, SchoolEntity } from "mikro-orm-config"
import { type LocationUpdateEventDto } from "types-custom"
import { type UserAndTokens } from "../src/auth/auth.service"

describe("SseController (e2e)", () => {
  const fixture: Fixture = new Fixture()

  let port: number
  let eventSource: EventSource

  let userAndTokens: UserAndTokens

  let sseService: SseService

  let location: LocationEntity
  let school: SchoolEntity

  beforeAll(() => {
    return fixture.beforeAll().then(() => {
      const address = (fixture.app.getHttpServer() as Server)
        .listen()
        .address() as AddressInfo
      port = address.port

      sseService = fixture.app.get<SseService>(SseService)
    })
  })

  afterAll(() => {
    const server = fixture.app.getHttpServer() as Server
    server.close()
    return fixture.afterAll()
  })

  beforeEach(() => {
    return fixture
      .beforeEach()
      .then(() => fixture.getTokens("User"))
      .then((data) => (userAndTokens = data))
      .then(() => fixture.em.find(LocationEntity, {}))
      .then((data) => {
        location = data[0]
      })
      .then(() => fixture.em.find(SchoolEntity, {}))
      .then((data) => {
        school = data[0]
      })
  })

  afterEach(() => {
    return fixture.afterEach().then(() => eventSource.close())
  })

  describe("/sse (SSE)", () => {
    it("receives event after updating Location capacity (200)", (done: jest.DoneCallback) => {
      eventSource = new EventSource(`http://localhost:${port}/sse`)

      eventSource.onerror = (error): void => {
        done(error)
      }

      eventSource.onopen = (): void => {
        const newCheckIn = new CheckInEntity(location, school)

        void fixture.em
          .persistAndFlush(newCheckIn)
          .then(() => fixture.em.findOneOrFail(LocationEntity, location.id))
          .then((data) => {
            location = data
          })
          .then(() => sseService.addLocationUpdate(location))
      }

      eventSource.onmessage = (event): void => {
        const locationUpdateEventData = JSON.parse(
          event.data as string
        ) as LocationUpdateEventDto

        try {
          expect(locationUpdateEventData.normalizedName).toBe(
            location.normalizedName
          )
          expect(locationUpdateEventData.available).toBe(location.available)
          expect(locationUpdateEventData.capacity).toBe(location.capacity)
          done()
        } catch (error) {
          done(error)
        }
      }
    })
  })

  describe("/sse/:normalizedName (SSE)", () => {
    it("receives event after updating Location capacity (200)", (done: jest.DoneCallback) => {
      const eventSourceInitDict = {
        headers: {
          Cookie: `accessToken=${userAndTokens.tokens.accessToken}`
        }
      }

      eventSource = new EventSource(
        `http://localhost:${port}/sse/${location.normalizedName}`,
        eventSourceInitDict
      )

      eventSource.onerror = (error): void => {
        done(error)
      }

      eventSource.onopen = (): void => {
        location.capacity--

        void fixture.em
          .flush()
          .then(() => fixture.em.findOneOrFail(LocationEntity, location.id))
          .then((data) => {
            location = data
          })
          .then(() => sseService.addLocationUpdate(location))
      }

      eventSource.onmessage = (event): void => {
        const locationUpdateEventData = JSON.parse(
          event.data as string
        ) as LocationUpdateEventDto

        try {
          expect(locationUpdateEventData.normalizedName).toBe(
            location.normalizedName
          )
          expect(locationUpdateEventData.available).toBe(location.available)
          expect(locationUpdateEventData.capacity).toBe(location.capacity)
          done()
        } catch (error) {
          done(error)
        }
      }
    })

    it("returns Unauthorized (401)", (done: jest.DoneCallback) => {
      eventSource = new EventSource(
        `http://localhost:${port}/sse/${location.normalizedName}`
      )

      eventSource.onerror = (event): void => {
        try {
          expect(event.status).toBe(401)
          done()
        } catch (error) {
          done(error)
        }
      }
    })
  })
})
