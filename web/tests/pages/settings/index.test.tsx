/* eslint-disable sonarjs/no-duplicate-string */
import { getMyUser } from "api-wrapper"
import { mocked } from "jest-mock"
import SettingsHome from "pages/settings"
import { adminUserMock, defaultUserMock, managerUserMock } from "user-mocks"
import mockRouter from "next-router-mock"
import { screen, render, waitFor } from "test-utils"

describe("SettingsHome (page)", () => {
  it("Redirect admin", async () => {
    mocked(getMyUser).mockImplementation(adminUserMock)

    await mockRouter.push("settings")

    render(<SettingsHome />)

    await screen.findByText("Loading home page.", { selector: "p" })

    await waitFor(() => expect(mockRouter.isReady))
    expect(mockRouter.asPath).toBe("/settings/locations")
  })

  it("Redirect manager", async () => {
    mocked(getMyUser).mockImplementation(managerUserMock)

    await mockRouter.push("settings")

    render(<SettingsHome />)

    await screen.findByText("Loading home page.", { selector: "p" })

    await waitFor(() => expect(mockRouter.isReady))
    expect(mockRouter.asPath).toBe("/settings/locations")
  })

  it("Redirect default", async () => {
    mocked(getMyUser).mockImplementation(defaultUserMock)

    await mockRouter.push("settings")

    render(<SettingsHome />)

    await screen.findByText("Loading home page.", { selector: "p" })

    await waitFor(() => expect(mockRouter.isReady))
    expect(mockRouter.asPath).toBe("/settings/locations/locationId")
  })
})
