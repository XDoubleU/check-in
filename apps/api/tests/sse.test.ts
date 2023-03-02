import { INestApplication } from "@nestjs/common"
import { Location, School, User } from "types-custom"
import { LocationsService } from "../src/locations/locations.service"
import { UsersService } from "../src/users/users.service"
import { clearDatabase, getAccessToken, getApp } from "./shared"
import EventSource from "eventsource"
import { Server } from "http"
import { AddressInfo } from "net"
import { LocationUpdateEventData } from "../src/sse/sse.service"
import { CheckInsService } from "../src/checkins/checkins.service"
import { SchoolsService } from "../src/schools/schools.service"


describe("SseController (e2e)", () => {
  let app: INestApplication
  let port: number
  let eventSource: EventSource
  let accessToken: string

  let locationsService: LocationsService
  let usersService: UsersService
  let checkinsService: CheckInsService
  let schoolsService: SchoolsService

  let user: User
  let location: Location
  let school: School
  
  beforeAll(async () => {
    app = await getApp()

    locationsService = app.get<LocationsService>(LocationsService)
    usersService = app.get<UsersService>(UsersService)
    checkinsService = app.get<CheckInsService>(CheckInsService)
    schoolsService = app.get<SchoolsService>(SchoolsService)

    await app.init()
  })

  beforeEach(async () => {
    // AccessTokens
    const getAccessTokenObject = await getAccessToken(app)
    accessToken = getAccessTokenObject.accessToken

    // UsersService
    user = await usersService.create("TestUser", "testpassword")

    // LocationsService
    location = await locationsService.create("TestLocation", 10, user)

    // SchoolsService
    school = await schoolsService.create("TestSchool")

    const address = (app.getHttpServer() as Server).listen().address() as AddressInfo
    port = address.port
  })

  afterEach(async () => {
    (app.getHttpServer() as Server).close()
    eventSource.close()
    
    await clearDatabase(app)
  })

  describe("/sse (SSE)", () => {
    it("receives event after updating Location capacity (200)", (done: jest.DoneCallback) => {
      eventSource = new EventSource(`http://localhost:${port}/sse`)
      
      eventSource.onerror = (error): void => {
        done(error)
      }

      eventSource.onopen = async (): Promise<void> => {
        await locationsService.update(location, undefined, location.capacity)
      }

      eventSource.onmessage = (event): void => {
        const locationUpdateEventData = JSON.parse(event.data as string) as LocationUpdateEventData

        try {
          expect(locationUpdateEventData.normalizedName).toBe(location.normalizedName)
          expect(locationUpdateEventData.available).toBe(location.available)
          expect(locationUpdateEventData.capacity).toBe(location.capacity)
          done()
        } catch (error) {
          done(error)
        }   
      }
    })

    it("receives event after creating CheckIn (200)", (done: jest.DoneCallback) => {
      eventSource = new EventSource(`http://localhost:${port}/sse`)
      
      eventSource.onerror = (error): void => {
        done(error)
      }

      eventSource.onopen = async (): Promise<void> => {
        await checkinsService.create(location, school)
      }

      eventSource.onmessage = (event): void => {
        const locationUpdateEventData = JSON.parse(event.data as string) as LocationUpdateEventData

        try {
          expect(locationUpdateEventData.normalizedName).toBe(location.normalizedName)
          expect(locationUpdateEventData.available).toBe(location.available - 1)
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
            "Cookie": `accessToken=${accessToken}`
          }
        }

        eventSource = new EventSource(
          `http://localhost:${port}/sse/${location.normalizedName}`,
          eventSourceInitDict
        )
          
        eventSource.onerror = (error): void => {
          done(error)
        }
  
        eventSource.onopen = async (): Promise<void> => {
          await locationsService.update(location, undefined, location.capacity)
        }
  
        eventSource.onmessage = (event): void => {
          const locationUpdateEventData = JSON.parse(event.data as string) as LocationUpdateEventData
  
          try {
            expect(locationUpdateEventData.normalizedName).toBe(location.normalizedName)
            expect(locationUpdateEventData.available).toBe(location.available)
            expect(locationUpdateEventData.capacity).toBe(location.capacity)
            done()
          } catch (error) {
            done(error)
          }   
        }
    })

    it("receives event after creating CheckIn (200)", (done: jest.DoneCallback) => {
      const eventSourceInitDict = {
        headers: {
          "Cookie": `accessToken=${accessToken}`
        }
      }

      eventSource = new EventSource(
        `http://localhost:${port}/sse/${location.normalizedName}`,
        eventSourceInitDict
      )
      
      eventSource.onerror = (error): void => {
        done(error)
      }

      eventSource.onopen = async (): Promise<void> => {
        await checkinsService.create(location, school)
      }

      eventSource.onmessage = (event): void => {
        const locationUpdateEventData = JSON.parse(event.data as string) as LocationUpdateEventData

        try {
          expect(locationUpdateEventData.normalizedName).toBe(location.normalizedName)
          expect(locationUpdateEventData.available).toBe(location.available - 1)
          expect(locationUpdateEventData.capacity).toBe(location.capacity)
          done()
        } catch (error) {
          done(error)
        }   
      }
    })

    it("returns Unauthorized (401)", (done: jest.DoneCallback) => {
      eventSource = new EventSource(`http://localhost:${port}/sse/${location.normalizedName}`)
      
      eventSource.onerror = (event): void => {
        try{
          expect(event.status).toBe(401)
          done()
        } catch (error) {
          done(error)
        }
      }
    })
  })
})
