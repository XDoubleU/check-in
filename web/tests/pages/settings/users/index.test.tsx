/* eslint-disable sonarjs/no-duplicate-string */
import { getAllUsersPaged, getMyUser } from "api-wrapper"
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
