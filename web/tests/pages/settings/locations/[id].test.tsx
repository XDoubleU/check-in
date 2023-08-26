import { getCheckInsToday, getLocation, getUser } from "api-wrapper"
import { mocked } from "jest-mock"
import mockRouter from "next-router-mock"
import LocationDetail from "pages/settings/locations/[id]"
import { screen, render } from "test-utils"
import { DefaultLocation, defaultUserMock } from "mocks"

// eslint-disable-next-line max-lines-per-function
describe("LocationDetail (page)", () => {
  it("Show detailed information of location", async () => {
    mocked(getUser).mockImplementation(defaultUserMock)

    mocked(getLocation).mockImplementation(() => {
      return Promise.resolve({
        ok: true,
        data: DefaultLocation
      })
    })

    mocked(getCheckInsToday).mockImplementation(() => {
      return Promise.resolve({
        ok: true,
        data: [{
          id: 1,
          capacity: 10,
          createdAt: new Date().toISOString(),
          locationId: "locationId",
          schoolName: "Andere"
        }]
      })
    })

    await mockRouter.push("/locations/locationId")

    render(<LocationDetail />) 

    await screen.findByRole("heading", { name: "location" })
  })

  it("Can't fetch detailed information of location", async () => {
    mocked(getUser).mockImplementation(defaultUserMock)

    mocked(getLocation).mockImplementation(() => {
      return Promise.resolve({
        ok: false
      })
    })

    await mockRouter.push("/locations/notALocationId")

    render(<LocationDetail />)

    await screen.findByText("User has no location")
  })
})