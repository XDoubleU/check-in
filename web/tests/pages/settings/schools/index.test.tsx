import userEvent from "@testing-library/user-event"
import {
  createSchool,
  getAllSchoolsPaged,
  getMyUser,
  updateSchool
} from "api-wrapper"
import { mocked } from "jest-mock"
import { defaultUserMock, managerUserMock, noUserMock } from "mocks"
import mockRouter from "next-router-mock"
import SchoolListView from "pages/settings/schools"
import { screen, render, waitFor } from "test-utils"

describe("SchoolListView (page)", () => {
  it("Shows overview of schools", async () => {
    mocked(getMyUser).mockImplementation(managerUserMock)

    mocked(getAllSchoolsPaged).mockImplementation(() => {
      return Promise.resolve({
        ok: true,
        data: {
          data: [
            {
              id: 1,
              name: "Andere",
              readOnly: true
            }
          ],
          pagination: {
            current: 1,
            total: 1
          }
        }
      })
    })

    await mockRouter.push("/settings/schools")

    render(<SchoolListView />)

    await screen.findByRole("heading", { name: "Schools" })
    expect(getAllSchoolsPaged).toHaveBeenCalledTimes(1)
  })

  it("Creates a school", async () => {
    mocked(getMyUser).mockImplementation(managerUserMock)

    mocked(getAllSchoolsPaged).mockImplementation(() => {
      return Promise.resolve({
        ok: true,
        data: {
          data: [
            {
              id: 2,
              name: "School",
              readOnly: false
            }
          ],
          pagination: {
            current: 1,
            total: 1
          }
        }
      })
    })

    mocked(createSchool).mockImplementation(() => {
      return Promise.resolve({
        ok: true
      })
    })

    await mockRouter.push("/settings/schools")

    render(<SchoolListView />)

    const createButton = await screen.findByRole("button", { name: "Create" })
    await userEvent.click(createButton)

    const nameField = screen.getByLabelText("Name")
    await userEvent.type(nameField, "newName")

    const createButtons = screen.getAllByRole("button", { name: "Create" })

    const createButtonIndex = createButtons.indexOf(createButton)
    createButtons.splice(createButtonIndex, 1)

    const confirmCreateButton = createButtons[0]
    await userEvent.click(confirmCreateButton)
  })

  it("Creates a school, school with name already exists", async () => {
    mocked(getMyUser).mockImplementation(managerUserMock)

    mocked(getAllSchoolsPaged).mockImplementation(() => {
      return Promise.resolve({
        ok: true,
        data: {
          data: [
            {
              id: 2,
              name: "School",
              readOnly: false
            }
          ],
          pagination: {
            current: 1,
            total: 1
          }
        }
      })
    })

    mocked(createSchool).mockImplementation(() => {
      return Promise.resolve({
        ok: false,
        message: {
          name: "school with this name already exists"
        }
      })
    })

    await mockRouter.push("/settings/schools")

    render(<SchoolListView />)

    const createButton = await screen.findByRole("button", { name: "Create" })
    await userEvent.click(createButton)

    const nameField = screen.getByLabelText("Name")
    await userEvent.type(nameField, "newName")

    const createButtons = screen.getAllByRole("button", { name: "Create" })

    const createButtonIndex = createButtons.indexOf(createButton)
    createButtons.splice(createButtonIndex, 1)

    const confirmCreateButton = createButtons[0]
    await userEvent.click(confirmCreateButton)

    await screen.findByText("school with this name already exists")
  })

  it("Cancel creating a school", async () => {
    mocked(getMyUser).mockImplementation(managerUserMock)

    mocked(getAllSchoolsPaged).mockImplementation(() => {
      return Promise.resolve({
        ok: true,
        data: {
          data: [
            {
              id: 2,
              name: "School",
              readOnly: false
            }
          ],
          pagination: {
            current: 1,
            total: 1
          }
        }
      })
    })

    await mockRouter.push("/settings/schools")

    render(<SchoolListView />)

    const createButton = await screen.findByRole("button", { name: "Create" })
    await userEvent.click(createButton)

    const cancelButton = screen.getByRole("button", { name: "Cancel" })
    await userEvent.click(cancelButton)
  })

  it("Updates a school", async () => {
    mocked(getMyUser).mockImplementation(managerUserMock)

    mocked(getAllSchoolsPaged).mockImplementation(() => {
      return Promise.resolve({
        ok: true,
        data: {
          data: [
            {
              id: 2,
              name: "School",
              readOnly: false
            }
          ],
          pagination: {
            current: 1,
            total: 1
          }
        }
      })
    })

    mocked(updateSchool).mockImplementation(() => {
      return Promise.resolve({
        ok: true
      })
    })

    await mockRouter.push("/settings/schools")

    render(<SchoolListView />)

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

  it("Updates a school, school with name already exists", async () => {
    mocked(getMyUser).mockImplementation(managerUserMock)

    mocked(getAllSchoolsPaged).mockImplementation(() => {
      return Promise.resolve({
        ok: true,
        data: {
          data: [
            {
              id: 2,
              name: "School",
              readOnly: false
            }
          ],
          pagination: {
            current: 1,
            total: 1
          }
        }
      })
    })

    mocked(updateSchool).mockImplementation(() => {
      return Promise.resolve({
        ok: false,
        message: {
          name: "school with this name already exists"
        }
      })
    })

    await mockRouter.push("/settings/schools")

    render(<SchoolListView />)

    const updateButton = await screen.findByRole("button", { name: "Update" })
    await userEvent.click(updateButton)

    const nameField = screen.getByLabelText("Name")
    await userEvent.type(nameField, "newName")

    const updateButtons = screen.getAllByRole("button", { name: "Update" })

    const updateButtonIndex = updateButtons.indexOf(updateButton)
    updateButtons.splice(updateButtonIndex, 1)

    const confirmUpdateButton = updateButtons[0]
    await userEvent.click(confirmUpdateButton)

    await screen.findByText("school with this name already exists")
  })

  it("Updates a school, something went wrong", async () => {
    mocked(getMyUser).mockImplementation(managerUserMock)

    mocked(getAllSchoolsPaged).mockImplementation(() => {
      return Promise.resolve({
        ok: true,
        data: {
          data: [
            {
              id: 2,
              name: "School",
              readOnly: false
            }
          ],
          pagination: {
            current: 1,
            total: 1
          }
        }
      })
    })

    mocked(updateSchool).mockImplementation(() => {
      return Promise.resolve({
        ok: false,
        message: "Something went wrong"
      })
    })

    await mockRouter.push("/settings/schools")

    render(<SchoolListView />)

    const updateButton = await screen.findByRole("button", { name: "Update" })
    await userEvent.click(updateButton)

    const nameField = screen.getByLabelText("Name")
    await userEvent.type(nameField, "newName")

    const updateButtons = screen.getAllByRole("button", { name: "Update" })

    const updateButtonIndex = updateButtons.indexOf(updateButton)
    updateButtons.splice(updateButtonIndex, 1)

    const confirmUpdateButton = updateButtons[0]
    await userEvent.click(confirmUpdateButton)

    await screen.findByText("Something went wrong")
  })

  it("Cancel updating a school", async () => {
    mocked(getMyUser).mockImplementation(managerUserMock)

    mocked(getAllSchoolsPaged).mockImplementation(() => {
      return Promise.resolve({
        ok: true,
        data: {
          data: [
            {
              id: 2,
              name: "School",
              readOnly: false
            }
          ],
          pagination: {
            current: 1,
            total: 1
          }
        }
      })
    })

    await mockRouter.push("/settings/schools")

    render(<SchoolListView />)

    const updateButton = await screen.findByRole("button", { name: "Update" })
    await userEvent.click(updateButton)

    const cancelButton = screen.getByRole("button", { name: "Cancel" })
    await userEvent.click(cancelButton)
  })

  it("Deletes a school", async () => {
    mocked(getMyUser).mockImplementation(managerUserMock)

    mocked(getAllSchoolsPaged).mockImplementation(() => {
      return Promise.resolve({
        ok: true,
        data: {
          data: [
            {
              id: 2,
              name: "School",
              readOnly: false
            }
          ],
          pagination: {
            current: 1,
            total: 1
          }
        }
      })
    })

    await mockRouter.push("/settings/schools")

    render(<SchoolListView />)

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

    await mockRouter.push("/settings/schools")

    render(<SchoolListView />)

    await waitFor(() => {
      expect(document.title).toBe("Loading...")
    })

    await waitFor(() => {
      expect(mockRouter.asPath).toBe("/settings")
    })
  })

  it("Redirect anonymous", async () => {
    mocked(getMyUser).mockImplementation(noUserMock)

    await mockRouter.push("/settings/schools")

    render(<SchoolListView />)

    await waitFor(() => {
      expect(document.title).toBe("Loading...")
    })

    await waitFor(() => {
      expect(mockRouter.asPath).toBe("/signin?redirect_to=%2Fsettings%2Fschools")
    })
  })
})
