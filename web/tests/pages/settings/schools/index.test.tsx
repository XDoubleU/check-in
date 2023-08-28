/* eslint-disable sonarjs/no-duplicate-string */
import { getAllSchoolsPaged, getMyUser } from "api-wrapper"
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
