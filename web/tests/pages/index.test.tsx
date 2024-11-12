import {
  checkinsWebsocket,
  getAllSchoolsSortedForLocation,
  getMyUser
} from "api-wrapper"
import { mocked } from "jest-mock"
import mockRouter from "next-router-mock"
import CheckIn from "pages"
import { screen, render, waitFor, fireEvent } from "test-utils"
import {
  adminUserMock,
  defaultUserMock,
  managerUserMock,
  noUserMock
} from "mocks"
import WS from "jest-websocket-mock"
import {
  type LocationUpdateEvent,
  type State
} from "api-wrapper/types/apiTypes"

mocked(checkinsWebsocket).mockImplementation(() => {
  return new WebSocket("ws://localhost:8000")
})

mocked(getAllSchoolsSortedForLocation).mockImplementation(() => {
  return Promise.resolve({
    ok: true,
    data: [
      {
        id: 1,
        name: "Andere",
        readOnly: true
      }
    ]
  })
})

describe("CheckIn (page)", () => {
  const server = new WS("ws://localhost:8000")

  it("Default user is logged in, Check-In btn is shown, 2 people check in and location is full", async () => {
    mocked(getMyUser).mockImplementation(defaultUserMock)

    await mockRouter.push("/")

    render(<CheckIn />)

    await waitFor(() => {
      expect(document.title).toBe("Check-In")
    })

    // First check-in
    const button = screen.getByRole("button", { name: "CHECK-IN" })
    fireEvent.click(button)

    await screen.findByRole("heading", { name: "KIES JE SCHOOL:" })
    let school = screen.getByRole("button", { name: "ANDERE" })
    fireEvent.click(school)

    await waitFor(() => {
      expect(school).not.toBeVisible()
    })

    // Check if button is disabled and becomes enabled again
    expect(button).toBeDisabled()
    await waitFor(
      () => {
        expect(button).toBeEnabled()
      },
      { timeout: 1500 }
    )

    // Second check-in
    fireEvent.click(button)

    await screen.findByRole("heading", { name: "KIES JE SCHOOL:" })
    school = screen.getByRole("button", { name: "ANDERE" })
    fireEvent.click(school)

    // Check that location is full
    await screen.findByRole("button", { name: "VOLZET" })
  })

  it("Default user is logged in, Check-In btn is shown, server sends an update", async () => {
    mocked(getMyUser).mockImplementation(defaultUserMock)

    await mockRouter.push("/")

    render(<CheckIn />)

    await waitFor(() => {
      expect(document.title).toBe("Check-In")
    })

    await screen.findByText("2", { selector: "span" })

    const update: LocationUpdateEvent = {
      available: 1,
      capacity: 10,
      normalizedName: "location",
      yesterdayFullAt: "",
      availableYesterday: 0,
      capacityYesterday: 0
    }

    server.send(JSON.stringify(update))

    await screen.findByText("1", { selector: "span" })
  })

  it("Default user is logged in, Check-In btn is shown, server sends a state update", async () => {
    mocked(getMyUser).mockImplementation(defaultUserMock)

    await mockRouter.push("/")

    render(<CheckIn />)

    await waitFor(() => {
      expect(document.title).toBe("Check-In")
    })

    await screen.findByText("2", { selector: "span" })

    const update: State = {
      isMaintenance: true,
      isDatabaseActive: false
    }

    server.send(JSON.stringify(update))

    const alert = await screen.findByRole("alert")
    expect(alert.innerHTML).toContain(
      "The Check-In is currently under maintenance."
    )
  })

  it("Redirect admin to settings", async () => {
    mocked(getMyUser).mockImplementation(adminUserMock)

    await mockRouter.push("/")

    render(<CheckIn />)

    await waitFor(() => {
      expect(document.title).toBe("Loading...")
    })

    await waitFor(() => expect(mockRouter.isReady))
    expect(mockRouter.asPath).toBe("/settings")
  })

  it("Redirect manager to settings", async () => {
    mocked(getMyUser).mockImplementation(managerUserMock)

    await mockRouter.push("/")

    render(<CheckIn />)

    await waitFor(() => {
      expect(document.title).toBe("Loading...")
    })

    await waitFor(() => expect(mockRouter.isReady))
    expect(mockRouter.asPath).toBe("/settings")
  })

  it("Redirect anonymous", async () => {
    mocked(getMyUser).mockImplementation(noUserMock)

    await mockRouter.push("/")

    render(<CheckIn />)

    await waitFor(() => {
      expect(document.title).toBe("Loading...")
    })

    await waitFor(() => expect(mockRouter.isReady))
    expect(mockRouter.asPath).toBe("/signin")
  })
})
