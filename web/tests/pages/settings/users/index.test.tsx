import { getAllUsersPaged, getUser } from "api-wrapper"
import { mocked } from "jest-mock"
import { DefaultLocation, adminUserMock } from "mocks"
import mockRouter from "next-router-mock"
import UserListView from "pages/settings/users"
import { screen, render } from "test-utils"

describe("UserListView (page)", () => {
  it("Shows overview of users", async () => {
    mocked(getUser).mockImplementation(adminUserMock)

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
})