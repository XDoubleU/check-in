import {
  getAllLocations,
  getDataForDayChart,
  getDataForRangeChart,
  getMyUser
} from "api-wrapper"
import { mocked } from "jest-mock"
import {
  DefaultLocation,
  adminUserMock,
  defaultUserMock,
  noUserMock
} from "mocks"
import mockRouter from "next-router-mock"
import { screen, render, waitFor, fireEvent } from "test-utils"
import Graphs from "pages/settings/graphs"
import userEvent from "@testing-library/user-event"

describe("Graphs (page)", () => {
  it("View range chart", async () => {
    mocked(getMyUser).mockImplementation(adminUserMock)

    mocked(getDataForRangeChart).mockImplementation(() => {
      return Promise.resolve({
        ok: true,
        data: {
          dates: [
            "2023-08-24",
            "2023-08-25",

          ],
          capacitiesPerLocation: {
            locationId: [10, 10]
          },
          valuesPerSchool: {
            Andere: [5, 5]
          }
        }
      })
    })

    mocked(getAllLocations).mockImplementation(() => {
      return Promise.resolve({
        ok: true,
        data: [DefaultLocation]
      })
    })

    await mockRouter.push("/settings/graphs")

    render(<Graphs />)

    await screen.findByRole("heading", { name: "Graphs" })

    const formField = screen.getByRole("listbox")
    const locationOption = screen.getByRole("option", { name: "location" })

    await userEvent.selectOptions(formField, locationOption)

    const startDateField = screen.getByLabelText("Start date")
    fireEvent.change(startDateField, { target: { value: "2001-11-14" } })

    const endDateField = screen.getByLabelText("End date")
    fireEvent.change(endDateField, { target: { value: "2001-11-14" } })

    const downloadCSVButton = screen.getByRole("button", {
      name: "Download as CSV"
    })
    await userEvent.click(downloadCSVButton)
  })

  it("View range chart, no data", async () => {
    mocked(getMyUser).mockImplementation(adminUserMock)

    mocked(getDataForRangeChart).mockImplementation(() => {
      return Promise.resolve({
        ok: true
      })
    })

    mocked(getAllLocations).mockImplementation(() => {
      return Promise.resolve({
        ok: true,
        data: [DefaultLocation]
      })
    })

    await mockRouter.push("/settings/graphs")

    render(<Graphs />)

    await screen.findByRole("heading", { name: "Graphs" })
  })

  it("View day chart", async () => {
    mocked(getMyUser).mockImplementation(adminUserMock)

    mocked(getDataForRangeChart).mockImplementation(() => {
      return Promise.resolve({
        ok: true,
        data: {
          dates: [
            "2023-08-24"
          ],
          capacitiesPerLocation: {
            locationId: [10]
          },
          valuesPerSchool: {
            Andere: [5]
          }
        }
      })
    })

    mocked(getDataForDayChart).mockImplementation(() => {
      return Promise.resolve({
        ok: true,
        data: {
          dates: [
            "2023-08-24"
          ],
          capacitiesPerLocation: {
            locationId: [10]
          },
          valuesPerSchool: {
            Andere: [5]
          }
        }
      })
    })

    mocked(getAllLocations).mockImplementation(() => {
      return Promise.resolve({
        ok: true,
        data: [DefaultLocation]
      })
    })

    await mockRouter.push("/settings/graphs")

    render(<Graphs />)

    await screen.findByRole("heading", { name: "Graphs" })

    const dayTab = screen.getByRole("tab", { name: "Day" })
    await userEvent.click(dayTab)

    const formField = screen.getByRole("listbox")
    const locationOption = screen.getByRole("option", { name: "location" })

    await userEvent.selectOptions(formField, locationOption)

    const dateField = screen.getByLabelText("Date")
    fireEvent.change(dateField, { target: { value: "2001-11-14" } })

    const downloadCSVButton = screen.getByRole("button", {
      name: "Download as CSV"
    })
    await userEvent.click(downloadCSVButton)
  })

  it("View day chart, no data", async () => {
    mocked(getMyUser).mockImplementation(adminUserMock)

    mocked(getDataForRangeChart).mockImplementation(() => {
      return Promise.resolve({
        ok: true
      })
    })

    mocked(getDataForDayChart).mockImplementation(() => {
      return Promise.resolve({
        ok: true
      })
    })

    mocked(getAllLocations).mockImplementation(() => {
      return Promise.resolve({
        ok: true,
        data: [DefaultLocation]
      })
    })

    await mockRouter.push("/settings/graphs")

    render(<Graphs />)

    await screen.findByRole("heading", { name: "Graphs" })

    const dayTab = screen.getByRole("tab", { name: "Day" })
    await userEvent.click(dayTab)
  })

  it("Redirect default", async () => {
    mocked(getMyUser).mockImplementation(defaultUserMock)

    await mockRouter.push("/settings/graphs")

    render(<Graphs />)

    await waitFor(() => {
      expect(document.title).toBe("Loading...")
    })

    await waitFor(() => {
      expect(mockRouter.asPath).toBe("/settings")
    })
  })

  it("Redirect anonymous", async () => {
    mocked(getMyUser).mockImplementation(noUserMock)

    await mockRouter.push("/settings/graphs")

    render(<Graphs />)

    await waitFor(() => {
      expect(document.title).toBe("Loading...")
    })

    await waitFor(() => {
      expect(mockRouter.asPath).toBe("/signin?redirect_to=%2Fsettings%2Fgraphs")
    })
  })
})
