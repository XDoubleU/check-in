import userEvent from "@testing-library/user-event"
import {
  createLocation,
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
            current: 1,
            total: 1
          }
        }
      })
    })

    await mockRouter.push("/settings/locations")

    render(<LocationListView />)

    await screen.findByRole("heading", { name: "Locations" })
    expect(getAllLocationsPaged).toHaveBeenCalledTimes(1)
  })

  it("Creates a location", async () => {
    mocked(getMyUser).mockImplementation(managerUserMock)

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
            }
          ],
          pagination: {
            current: 1,
            total: 1
          }
        }
      })
    })

    mocked(createLocation).mockImplementation(() => {
      return Promise.resolve({
        ok: true
      })
    })

    await mockRouter.push("/settings/locations")

    render(<LocationListView />)

    const createButton = await screen.findByRole("button", { name: "Create" })
    await userEvent.click(createButton)

    const nameField = screen.getByLabelText("Name")
    await userEvent.type(nameField, "newName")

    const capacityField = screen.getByLabelText("Capacity")
    await userEvent.type(capacityField, "10")

    const userNameField = screen.getByLabelText("Username")
    await userEvent.type(userNameField, "newUserName")

    const passwordField = screen.getByLabelText("Password")
    await userEvent.type(passwordField, "newPassword")

    const repeatPasswordField = screen.getByLabelText("Repeat password")
    await userEvent.type(repeatPasswordField, "newPassword")

    const createButtons = screen.getAllByRole("button", { name: "Create" })

    const createButtonIndex = createButtons.indexOf(createButton)
    createButtons.splice(createButtonIndex, 1)

    const confirmCreateButton = createButtons[0]
    await userEvent.click(confirmCreateButton)
  }, 10000)

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
              yesterdayFullAt: "",
              availableYesterday: 0,
              capacityYesterday: 0
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

    const updateButton = await screen.findByRole("button", { name: "Update" })
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
              yesterdayFullAt: "",
              availableYesterday: 0,
              capacityYesterday: 0
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

    const deleteButton = await screen.findByRole("button", { name: "Delete" })
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

    await waitFor(() => {
      expect(document.title).toBe("Loading...")
    })

    await waitFor(() => expect(mockRouter.isReady))
    expect(mockRouter.asPath).toBe("/settings")
  })

  it("Redirect anonymous", async () => {
    mocked(getMyUser).mockImplementation(noUserMock)

    await mockRouter.push("/settings/locations")

    render(<LocationListView />)

    await waitFor(() => {
      expect(document.title).toBe("Loading...")
    })

    await waitFor(() => expect(mockRouter.isReady))
    expect(mockRouter.asPath).toBe(
      "/signin?redirect_to=%2Fsettings%2Flocations"
    )
  })
})
