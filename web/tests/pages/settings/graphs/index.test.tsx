/* eslint-disable sonarjs/no-duplicate-string */
import { getMyUser } from "api-wrapper"
import { mocked } from "jest-mock"
import { adminUserMock, defaultUserMock } from "mocks"
import mockRouter from "next-router-mock"
import { screen, render, waitFor } from "test-utils"
import Graphs from "pages/settings/graphs"
import userEvent from "@testing-library/user-event"

// eslint-disable-next-line max-lines-per-function
describe("Graphs (page)", () => {
  it("Graphs load", async () => {
    mocked(getMyUser).mockImplementation(adminUserMock)

    await mockRouter.push("settings/graphs")

    render(<Graphs />)

    await screen.findByRole("heading", { name: "Graphs" })

    const formField = screen.getByRole("listbox")
    const locationOption = screen.getByRole("option", { name: "location" })

    await userEvent.selectOptions(formField, locationOption)

    // eslint-disable-next-line no-warning-comments
    //TODO: tests graphs more extensively here
  })

  it("Redirect default", async () => {
    mocked(getMyUser).mockImplementation(defaultUserMock)

    await mockRouter.push("settings/graphs")

    render(<Graphs />)

    await waitFor(() => expect(mockRouter.isReady))
    expect(mockRouter.asPath).toBe("/settings")
  })
})
