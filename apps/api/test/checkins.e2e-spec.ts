import { Test, TestingModule } from "@nestjs/testing"
import { INestApplication } from "@nestjs/common"
import request from "supertest"
import { AppModule } from "../src/app.module"

describe("CheckInsController (e2e)", () => {
  let app: INestApplication

  beforeEach(async () => {
    const moduleFixture: TestingModule = await Test.createTestingModule({
      imports: [AppModule],
    }).compile()

    app = moduleFixture.createNestApplication()
    await app.init()
  })

  it("/ (POST)", () => {
    /*const data: CreateCheckInDto = {
      locationId: "1",
      schoolId: 1
    }*/

    return request(app.getHttpServer())
      .post("/")
      //.send(data)
      .expect(200)
      .expect("Hello World!")
  })
})
