/* eslint-disable max-lines-per-function */
/* eslint-disable sonarjs/no-duplicate-string */
import userEvent from "@testing-library/user-event"
import { getAllUsersPaged, getMyUser, updateUser } from "api-wrapper"
import { mocked } from "jest-mock"
import {
  DefaultLocation,
  adminUserMock,
  defaultUserMock,
  managerUserMock,
  noUserMock
} from "mocks"
import mockRouter from "next-router-mock"
import UserListView from "pages/settings/users"
import { screen, render, waitFor } from "test-utils"

// eslint-disable-next-line max-lines-per-function
describe("UserListView (page)", () => {
  it("Shows overview of users", async () => {
    mocked(getMyUser).mockImplementation(adminUserMock)

    mocked(getAllUsersPaged).mockImplementation(() => {
      return Promise.resolve({
        ok: true,
        data: {
          data: [
            {
              id: "userId",
              username: "default",
              location: DefaultLocation,
              role: "default"
            }
          ],
          pagination: {
            current: 1,
            total: 1
          }
        }
      })
    })

    await mockRouter.push("/settings/users")

    render(<UserListView />)

    await screen.findByRole("heading", { name: "Users" })
  })

  it("Updates a user", async () => {
    mocked(getMyUser).mockImplementation(adminUserMock)

    mocked(getAllUsersPaged).mockImplementation(() => {
      return Promise.resolve({
        ok: true,
        data: {
          data: [
            {
              id: "userId",
              username: "default",
              location: DefaultLocation,
              role: "default"
            }
          ],
          pagination: {
            current: 1,
            total: 1
          }
        }
      })
    })

    mocked(updateUser).mockImplementation(() => {
      return Promise.resolve({
        ok: true
      })
    })

    await mockRouter.push("/settings/users")

    render(<UserListView />)

    await screen.findByRole("heading", { name: "Users" })

    const updateButton = screen.getByRole("button", { name: "Update" })
    await userEvent.click(updateButton)

    const nameField = screen.getByLabelText("Username")
    await userEvent.type(nameField, "newUserName")

    const updateButtons = screen.getAllByRole("button", { name: "Update" })

    const updateButtonIndex = updateButtons.indexOf(updateButton)
    updateButtons.splice(updateButtonIndex, 1)

    const confirmUpdateButton = updateButtons[0]
    await userEvent.click(confirmUpdateButton)
  })

  it("Deletes a user", async () => {
    mocked(getMyUser).mockImplementation(adminUserMock)

    mocked(getAllUsersPaged).mockImplementation(() => {
      return Promise.resolve({
        ok: true,
        data: {
          data: [
            {
              id: "userId",
              username: "default",
              location: DefaultLocation,
              role: "default"
            }
          ],
          pagination: {
            current: 1,
            total: 1
          }
        }
      })
    })

    await mockRouter.push("/settings/users")

    render(<UserListView />)

    await screen.findByRole("heading", { name: "Users" })

    const deleteButton = screen.getByRole("button", { name: "Delete" })
    await userEvent.click(deleteButton)

    const deleteButtons = screen.getAllByRole("button", { name: "Delete" })

    const deleteButtonIndex = deleteButtons.indexOf(deleteButton)
    deleteButtons.splice(deleteButtonIndex, 1)

    const confirmDeleteButton = deleteButtons[0]
    await userEvent.click(confirmDeleteButton)
  })

  it("Redirect Manager", async () => {
    mocked(getMyUser).mockImplementation(managerUserMock)

    await mockRouter.push("/settings/users")

    render(<UserListView />)

    await waitFor(() => expect(document.title).toBe("Loading..."))

    await waitFor(() => expect(mockRouter.isReady))
    expect(mockRouter.asPath).toBe("/settings")
  })

  it("Redirect default", async () => {
    mocked(getMyUser).mockImplementation(defaultUserMock)

    await mockRouter.push("/settings/users")

    render(<UserListView />)

    await waitFor(() => expect(document.title).toBe("Loading..."))

    await waitFor(() => expect(mockRouter.isReady))
    expect(mockRouter.asPath).toBe("/settings")
  })

  it("Redirect anonymous", async () => {
    mocked(getMyUser).mockImplementation(noUserMock)

    await mockRouter.push("/settings/users")

    render(<UserListView />)

    await waitFor(() => expect(document.title).toBe("Loading..."))

    await waitFor(() => expect(mockRouter.isReady))
    expect(mockRouter.asPath).toBe("/signin?redirect_to=%2Fsettings%2Fusers")
  })
})
