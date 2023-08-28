import { getAllSchoolsPaged, getUser } from "api-wrapper"
import { mocked } from "jest-mock"
import { managerUserMock } from "mocks"
import mockRouter from "next-router-mock"
import SchoolListView from "pages/settings/schools"
import { screen, render } from "test-utils"

describe("SchoolListView (page)", () => {
  it("Shows overview of schools", async () => {
    mocked(getUser).mockImplementation(managerUserMock)

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
})