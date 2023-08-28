/* eslint-disable sonarjs/no-duplicate-string */
import userEvent from "@testing-library/user-event"
import { getAllSchoolsPaged, getMyUser, updateSchool } from "api-wrapper"
import { mocked } from "jest-mock"
import { defaultUserMock, managerUserMock, noUserMock } from "mocks"
import mockRouter from "next-router-mock"
import SchoolListView from "pages/settings/schools"
import { screen, render, waitFor } from "test-utils"

// eslint-disable-next-line max-lines-per-function
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

    await screen.findByRole("heading", { name: "Schools" })

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

    await screen.findByRole("heading", { name: "Schools" })

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

    await mockRouter.push("/settings/schools")

    render(<SchoolListView />)

    await waitFor(() => expect(document.title).toBe("Loading..."))

    await waitFor(() => expect(mockRouter.isReady))
    expect(mockRouter.asPath).toBe("/settings")
  })

  it("Redirect anonymous", async () => {
    mocked(getMyUser).mockImplementation(noUserMock)

    await mockRouter.push("/settings/schools")

    render(<SchoolListView />)

    await waitFor(() => expect(document.title).toBe("Loading..."))

    await waitFor(() => expect(mockRouter.isReady))
    expect(mockRouter.asPath).toBe("/signin?redirect_to=%2Fsettings%2Fschools")
  })
})
