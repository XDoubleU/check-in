import userEvent from "@testing-library/user-event"
import { getMyUser, getState, updateState } from "api-wrapper"
import { mocked } from "jest-mock"
import {
  adminUserMock,
  defaultUserMock,
  managerUserMock,
  noUserMock
} from "mocks"
import mockRouter from "next-router-mock"
import StateView from "pages/settings/state"
import { screen, render, waitFor } from "test-utils"

describe("StateView (page)", () => {
  it("Updates the state", async () => {
    mocked(getMyUser).mockImplementation(adminUserMock)

    mocked(getState).mockImplementation(() => {
      return Promise.resolve({
        ok: true,
        data: {
          isMaintenance: false,
          isDatabaseActive: true
        }
      })
    })

    mocked(updateState).mockImplementation(() => {
      return Promise.resolve({
        ok: true
      })
    })

    await mockRouter.push("/settings/state")

    render(<StateView />)

    const isMaintenanceField = await screen.findByLabelText(
      "Is maintenance enabled"
    )
    await userEvent.click(isMaintenanceField)

    const updateButton = await screen.findByRole("button", { name: "Update" })
    await userEvent.click(updateButton)
  })

  it("Updates the state, something went wrong", async () => {
    mocked(getMyUser).mockImplementation(adminUserMock)

    mocked(getState).mockImplementation(() => {
      return Promise.resolve({
        ok: true,
        data: {
          isMaintenance: false,
          isDatabaseActive: true
        }
      })
    })

    mocked(updateState).mockImplementation(() => {
      return Promise.resolve({
        ok: false
      })
    })

    await mockRouter.push("/settings/state")

    render(<StateView />)

    const isMaintenanceField = await screen.findByLabelText(
      "Is maintenance enabled"
    )
    await userEvent.click(isMaintenanceField)

    const updateButton = await screen.findByRole("button", { name: "Update" })
    await userEvent.click(updateButton)
  })

  it("Redirect manager", async () => {
    mocked(getMyUser).mockImplementation(managerUserMock)

    await mockRouter.push("/settings/state")

    render(<StateView />)

    await waitFor(() => expect(document.title).toBe("Loading..."))

    await waitFor(() => expect(mockRouter.isReady))
    expect(mockRouter.asPath).toBe("/settings")
  })

  it("Redirect default", async () => {
    mocked(getMyUser).mockImplementation(defaultUserMock)

    await mockRouter.push("/settings/state")

    render(<StateView />)

    await waitFor(() => expect(document.title).toBe("Loading..."))

    await waitFor(() => expect(mockRouter.isReady))
    expect(mockRouter.asPath).toBe("/settings")
  })

  it("Redirect anonymous", async () => {
    mocked(getMyUser).mockImplementation(noUserMock)

    await mockRouter.push("/settings/state")

    render(<StateView />)

    await waitFor(() => expect(document.title).toBe("Loading..."))

    await waitFor(() => expect(mockRouter.isReady))
    expect(mockRouter.asPath).toBe("/signin?redirect_to=%2Fsettings%2Fstate")
  })
})
