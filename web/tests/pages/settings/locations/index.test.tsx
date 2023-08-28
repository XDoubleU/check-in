import { getAllLocationsPaged, getUser } from "api-wrapper"
import { mocked } from "jest-mock"
import { managerUserMock } from "mocks"
import mockRouter from "next-router-mock"
import LocationListView from "pages/settings/locations"
import { screen, render } from "test-utils"

describe("LocationListView (page)", () => {
  it("Shows overview of locations", async () => {
    mocked(getUser).mockImplementation(managerUserMock)

    mocked(getAllLocationsPaged).mockImplementation(() => {
      return Promise.resolve({
        ok: true,
        data: {
          data: [
            {
              id: "locationId",
              name: "location",
              normalizedName: "location",
              available: 10,
              capacity: 10,
              timeZone: "Europe/Brussels",
              userId: "userId",
              yesterdayFullAt: ""
            }
          ],
          pagination: {
            current: 1,
            total: 1
          }
        }
      })
    })

    await mockRouter.push("/settings/locations")

    render(<LocationListView />)

    await screen.findByRole("heading", { name: "Locations" })
  })
})