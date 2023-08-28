/* eslint-disable sonarjs/no-duplicate-string */
import { getAllLocationsPaged, getMyUser, getUser } from "api-wrapper"
import { mocked } from "jest-mock"
import { defaultUserMock, managerUserMock, noUserMock } from "mocks"
import mockRouter from "next-router-mock"
import LocationListView from "pages/settings/locations"
import { screen, render, waitFor } from "test-utils"

// eslint-disable-next-line max-lines-per-function
describe("LocationListView (page)", () => {
  it("Shows overview of locations", async () => {
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

  it("Redirect default", async () => {
    mocked(getMyUser).mockImplementation(defaultUserMock)

    await mockRouter.push("/settings/locations")

    render(<LocationListView />)

    await waitFor(() => expect(document.title).toBe("Loading..."))

    await waitFor(() => expect(mockRouter.isReady))
    expect(mockRouter.asPath).toBe("/settings")
  })

  it("Redirect anonymous", async () => {
    mocked(getMyUser).mockImplementation(noUserMock)

    await mockRouter.push("/settings/locations")

    render(<LocationListView />)

    await waitFor(() => expect(document.title).toBe("Loading..."))

    await waitFor(() => expect(mockRouter.isReady))
    expect(mockRouter.asPath).toBe(
      "/signin?redirect_to=%2Fsettings%2Flocations"
    )
  })
})
