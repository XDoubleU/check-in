import Fixture, { TokensAndUser } from "./fixture"
import EventSource from "eventsource"
import { Server } from "http"
import { AddressInfo } from "net"
import { LocationUpdateEventData } from "../src/sse/sse.service"
import { SseService } from "../src/sse/sse.service"
import { LocationEntity } from "mikro-orm-config"
import { expect } from "chai"


describe("SseController (e2e)", () => {
  let fixture: Fixture

  let port: number
  let eventSource: EventSource
  
  let tokensAndUser: TokensAndUser

  let sseService: SseService

  let location: LocationEntity
  
  before(() => {
    fixture = new Fixture()
    return fixture.init()
      .then(() => fixture.seedDatabase())
      .then(() => fixture.getTokens("User"))
      .then((data) => tokensAndUser = data)
      .then(() => fixture.em.find(LocationEntity, {}))
      .then((data) => {
        location = data[0]
      })
  })

  beforeEach(() => {
    const address = (fixture.app.getHttpServer() as Server).listen().address() as AddressInfo
    port = address.port

    sseService = fixture.app.get<SseService>(SseService)
  })

  afterEach(() => {
    (fixture.app.getHttpServer() as Server).close()
    eventSource.close()
  })

  after(() => {
    return fixture.clearDatabase()
      .then(() => fixture.app.close())
  })

  describe("/sse (SSE)", () => {
    it("receives event after updating Location capacity (200)", (done: Mocha.Done) => {
      eventSource = new EventSource(`http://localhost:${port}/sse`)
      
      eventSource.onerror = (error): void => {
        done(error)
      }

      eventSource.onopen = (): void => {
        sseService.addLocationUpdate(location)
      }

      eventSource.onmessage = (event): void => {
        const locationUpdateEventData = JSON.parse(event.data as string) as LocationUpdateEventData

        try {
          expect(locationUpdateEventData.normalizedName).to.be.equal(location.normalizedName)
          expect(locationUpdateEventData.available).to.be.equal(location.available)
          expect(locationUpdateEventData.capacity).to.be.equal(location.capacity)
          done()
        } catch (error) {
          done(error)
        }   
      }
    })
  })

  describe("/sse/:normalizedName (SSE)", () => {
    it("receives event after updating Location capacity (200)", (done: Mocha.Done) => {
        const eventSourceInitDict = {
          headers: {
            "Cookie": `accessToken=${tokensAndUser.tokens.accessToken}`
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
          sseService.addLocationUpdate(location)
        }
  
        eventSource.onmessage = (event): void => {
          const locationUpdateEventData = JSON.parse(event.data as string) as LocationUpdateEventData
  
          try {
            expect(locationUpdateEventData.normalizedName).to.be.equal(location.normalizedName)
            expect(locationUpdateEventData.available).to.be.equal(location.available)
            expect(locationUpdateEventData.capacity).to.be.equal(location.capacity)
            done()
          } catch (error) {
            done(error)
          }   
        }
    })

    it("returns Unauthorized (401)", (done: Mocha.Done) => {
      eventSource = new EventSource(`http://localhost:${port}/sse/${location.normalizedName}`)
      
      eventSource.onerror = (event): void => {
        try{
          expect(event.status).to.be.equal(401)
          done()
        } catch (error) {
          done(error)
        }
      }
    })
  })
})
