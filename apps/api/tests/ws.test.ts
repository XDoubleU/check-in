/* eslint-disable max-lines-per-function */
import Fixture from "./config/fixture"
import { type Server } from "http"
import { type AddressInfo } from "net"
import { WsService } from "../src/ws/ws.service"
import { CheckInEntity, LocationEntity, SchoolEntity } from "mikro-orm-config"
import { type LocationUpdateEventDto } from "types-custom"
import WebSocket from "ws"

describe("WsGateway (e2e)", () => {
  const fixture: Fixture = new Fixture()

  let port: number
  let webSocket: WebSocket

  let wsService: WsService

  let location: LocationEntity
  let school: SchoolEntity

  beforeAll(() => {
    return fixture.beforeAll().then(() => {
      const address = (fixture.app.getHttpServer() as Server)
        .listen()
        .address() as AddressInfo
      port = address.port

      wsService = fixture.app.get<WsService>(WsService)
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
    return fixture.afterEach().then(() => webSocket.close())
  })

  describe("all-locations (WS)", () => {
    it("receives event after updating Location capacity (200)", (done: jest.DoneCallback) => {
      webSocket = new WebSocket(`ws://localhost:${port}`)

      webSocket.onerror = (error): void => {
        done(error)
      }

      webSocket.onopen = (): void => {
        webSocket.send(
          JSON.stringify({
            event: "all-locations"
          })
        )

        const newCheckIn = new CheckInEntity(location, school)

        void fixture.em
          .persistAndFlush(newCheckIn)
          .then(() => fixture.em.findOneOrFail(LocationEntity, location.id))
          .then((data) => {
            location = data
          })
          .then(() => wsService.addLocationUpdate(location))
      }

      webSocket.onmessage = (event): void => {
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

  describe("single-location (WS)", () => {
    it("receives event after updating Location capacity (200)", (done: jest.DoneCallback) => {
      webSocket = new WebSocket(`ws://localhost:${port}`)

      webSocket.onerror = (error): void => {
        done(error)
      }

      webSocket.onopen = (): void => {
        webSocket.send(
          JSON.stringify({
            event: "single-location",
            data: {
              normalizedName: location.normalizedName
            }
          })
        )

        location.capacity--

        void fixture.em
          .flush()
          .then(() => fixture.em.findOneOrFail(LocationEntity, location.id))
          .then((data) => {
            location = data
          })
          .then(() => wsService.addLocationUpdate(location))
      }

      webSocket.onmessage = (event): void => {
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
})
