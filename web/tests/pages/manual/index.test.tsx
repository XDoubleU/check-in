import { getMyUser } from "api-wrapper"
import { mocked } from "jest-mock"
import ManualHome, { ManualNavigation } from "pages/manual"
import { screen, render, waitFor } from "test-utils"
import {
  adminUserMock,
  defaultUserMock,
  managerUserMock,
  noUserMock
} from "mocks"
import mockRouter from "next-router-mock"

describe("ManualNavigation (component)", () => {
  it("As seen by admin, english", async () => {
    mocked(getMyUser).mockImplementation(adminUserMock)

    await mockRouter.push("/manual/en/manager")

    render(<ManualNavigation />)

    await screen.findByRole("link", { name: "Verander naar Nederlands" })
    await screen.findByRole("link", { name: "Manual Manager" })
    await screen.findByRole("link", { name: "Manual Location" })
  })

  it("As seen by admin, dutch", async () => {
    mocked(getMyUser).mockImplementation(adminUserMock)

    await mockRouter.push("/manual/nl/manager")

    render(<ManualNavigation />)

    await screen.findByRole("link", { name: "Switch to English" })
    await screen.findByRole("link", { name: "Handleiding Beheerder" })
    await screen.findByRole("link", { name: "Handleiding Locatie" })
  })

  it("As seen by manager", async () => {
    mocked(getMyUser).mockImplementation(managerUserMock)

    await mockRouter.push("/manual/en/manager")

    render(<ManualNavigation />)

    await screen.findByRole("link", { name: "Verander naar Nederlands" })
    await screen.findByRole("link", { name: "Manual Manager" })
    await screen.findByRole("link", { name: "Manual Location" })
  })

  it("As seen by default", async () => {
    mocked(getMyUser).mockImplementation(defaultUserMock)

    await mockRouter.push("/manual/en/location")

    render(<ManualNavigation />)

    await screen.findByRole("link", { name: "Verander naar Nederlands" })
  })
})

describe("ManualHome (page)", () => {
  it("Redirect admin", async () => {
    mocked(getMyUser).mockImplementation(adminUserMock)

    await mockRouter.push("/manual")

    render(<ManualHome />)

    await screen.findByText("Loading manual.", { selector: "p" })

    await waitFor(() => {
      expect(mockRouter.asPath).toBe("/manual/en/manager")
    })
  })

  it("Redirect manager", async () => {
    mocked(getMyUser).mockImplementation(managerUserMock)

    await mockRouter.push("/manual")

    render(<ManualHome />)

    await screen.findByText("Loading manual.", { selector: "p" })

    await waitFor(() => {
      expect(mockRouter.asPath).toBe("/manual/en/manager")
    })
  })

  it("Redirect default", async () => {
    mocked(getMyUser).mockImplementation(defaultUserMock)

    await mockRouter.push("/manual")

    render(<ManualHome />)

    await screen.findByText("Loading manual.", { selector: "p" })

    await waitFor(() => {
      expect(mockRouter.asPath).toBe("/manual/en/location")
    })
  })

  it("Redirect anonymous", async () => {
    mocked(getMyUser).mockImplementation(noUserMock)

    await mockRouter.push("/manual")

    render(<ManualHome />)

    await waitFor(() => {
      expect(document.title).toBe("Loading...")
    })

    await waitFor(() => {
      expect(mockRouter.asPath).toBe("/signin?redirect_to=%2Fmanual")
    })
  })
})
