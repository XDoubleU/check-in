/* eslint-disable max-lines-per-function */
/* eslint-disable sonarjs/no-duplicate-string */
import userEvent from "@testing-library/user-event"
import {
  getAllLocationsPaged,
  getMyUser,
  getUser,
  updateLocation
} from "api-wrapper"
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
            },
            {
              id: "locationId2",
              name: "location2",
              normalizedName: "location2",
              available: 10,
              capacity: 10,
              timeZone: "Europe/Brussels",
              userId: "userId2",
              yesterdayFullAt: new Date().toISOString()
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

  it("Updates a location", async () => {
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

    mocked(updateLocation).mockImplementation(() => {
      return Promise.resolve({
        ok: true
      })
    })

    await mockRouter.push("/settings/locations")

    render(<LocationListView />)

    await screen.findByRole("heading", { name: "Locations" })

    const updateButton = screen.getByRole("button", { name: "Update" })
    await userEvent.click(updateButton)

    const nameField = screen.getByLabelText("Name")
    await userEvent.type(nameField, "newName")

    const updateButtons = screen.getAllByRole("button", { name: "Update" })

    const updateButtonIndex = updateButtons.indexOf(updateButton)
    updateButtons.splice(updateButtonIndex, 1)

    const confirmUpdateButton = updateButtons[0]
    await userEvent.click(confirmUpdateButton)
  })

  it("Deletes a location", async () => {
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

    const deleteButton = screen.getByRole("button", { name: "Delete" })
    await userEvent.click(deleteButton)

    const deleteButtons = screen.getAllByRole("button", { name: "Delete" })

    const deleteButtonIndex = deleteButtons.indexOf(deleteButton)
    deleteButtons.splice(deleteButtonIndex, 1)

    const confirmDeleteButton = deleteButtons[0]
    await userEvent.click(confirmDeleteButton)
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
