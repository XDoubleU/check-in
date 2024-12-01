import userEvent from "@testing-library/user-event"
import {
  createUser,
  getAllUsersPaged,
  getMyUser,
  updateUser
} from "api-wrapper"
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
    expect(getAllUsersPaged).toHaveBeenCalledTimes(1)
  })

  it("Creates a user", async () => {
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

    mocked(createUser).mockImplementation(() => {
      return Promise.resolve({
        ok: true
      })
    })

    await mockRouter.push("/settings/users")

    render(<UserListView />)

    const createButton = await screen.findByRole("button", { name: "Create" })
    await userEvent.click(createButton)

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
  })

  it("Creates a user, passwords don't match", async () => {
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

    mocked(createUser).mockImplementation(() => {
      return Promise.resolve({
        ok: true
      })
    })

    await mockRouter.push("/settings/users")

    render(<UserListView />)

    const createButton = await screen.findByRole("button", { name: "Create" })
    await userEvent.click(createButton)

    const userNameField = screen.getByLabelText("Username")
    await userEvent.type(userNameField, "newUserName")

    const passwordField = screen.getByLabelText("Password")
    await userEvent.type(passwordField, "newPassword")

    const repeatPasswordField = screen.getByLabelText("Repeat password")
    await userEvent.type(repeatPasswordField, "newPasswordWrong")

    const createButtons = screen.getAllByRole("button", { name: "Create" })

    const createButtonIndex = createButtons.indexOf(createButton)
    createButtons.splice(createButtonIndex, 1)

    const confirmCreateButton = createButtons[0]
    await userEvent.click(confirmCreateButton)

    await screen.findByText("Your passwords do no match")
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

    const updateButton = await screen.findByRole("button", { name: "Update" })
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

    const deleteButton = await screen.findByRole("button", { name: "Delete" })
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

    await waitFor(() => {
      expect(document.title).toBe("Loading...")
    })

    await waitFor(() => {
      expect(mockRouter.asPath).toBe("/settings")
    })
  })

  it("Redirect default", async () => {
    mocked(getMyUser).mockImplementation(defaultUserMock)

    await mockRouter.push("/settings/users")

    render(<UserListView />)

    await waitFor(() => {
      expect(document.title).toBe("Loading...")
    })

    await waitFor(() => {
      expect(mockRouter.asPath).toBe("/settings")
    })
  })

  it("Redirect anonymous", async () => {
    mocked(getMyUser).mockImplementation(noUserMock)

    await mockRouter.push("/settings/users")

    render(<UserListView />)

    await waitFor(() => {
      expect(document.title).toBe("Loading...")
    })

    await waitFor(() => {
      expect(mockRouter.asPath).toBe("/signin?redirect_to=%2Fsettings%2Fusers")
    })
  })
})
