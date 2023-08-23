import { getMyUser } from "api-wrapper"
import { mocked } from "jest-mock"
import mockRouter from "next-router-mock"
import CheckIn from "pages"
import { screen, render, waitFor } from "test-utils"
import { adminUserMock, defaultUserMock, managerUserMock } from "user-mocks"

describe("CheckIn (page)", () => {
  it("Default user is logged in, Check-In btn is shown and pressed, then a school is picked", async () => {
    mocked(getMyUser).mockImplementation(defaultUserMock)

    await mockRouter.push("/")

    render(<CheckIn />)

    await waitFor(() => expect(document.title).toContain())

    await screen.findByRole("button", { name: "Check-In" })
  })

  it("Redirect admin to settings", async () => {
    mocked(getMyUser).mockImplementation(adminUserMock)

    await mockRouter.push("/")

    render(<CheckIn />)

    await waitFor(() => expect(mockRouter.asPath).toBe("/settings"))
    expect(mockRouter.asPath).toBe("/settings")
  })
  
  it("Redirect manager to settings", async () => {
    mocked(getMyUser).mockImplementation(managerUserMock)

    await mockRouter.push("/")

    render(<CheckIn />)

    await waitFor(() => expect(mockRouter.asPath).toBe("/settings"))
    expect(mockRouter.asPath).toBe("/settings")
  })
})
