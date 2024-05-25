import { getAllLocationsPaged, getMyUser, getUser } from "api-wrapper"
import { mocked } from "jest-mock"
import { defaultUserMock, managerUserMock } from "mocks"
import mockRouter from "next-router-mock"
import LocationListView from "pages/settings/locations"
import { screen, render } from "test-utils"

describe("LocationListView (page)", () => {
  it("Redirect when requesting out of bounds", async () => {
    mocked(getMyUser).mockImplementation(managerUserMock)
    mocked(getUser).mockImplementation(defaultUserMock)

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
              yesterdayFullAt: "",
              availableYesterday: 0,
              capacityYesterday: 0
            },
            {
              id: "locationId2",
              name: "location2",
              normalizedName: "location2",
              available: 10,
              capacity: 10,
              timeZone: "Europe/Brussels",
              userId: "userId2",
              yesterdayFullAt: new Date().toISOString(),
              availableYesterday: 0,
              capacityYesterday: 0
            }
          ],
          pagination: {
            current: 2,
            total: 1
          }
        }
      })
    })

    await mockRouter.push("/settings/locations")

    render(<LocationListView />)

    await screen.findByRole("heading", { name: "Locations" })
    expect(getAllLocationsPaged).toHaveBeenCalledTimes(2)
  })
})
