import { render, waitFor } from "test-utils"
import { getMyUser, signOut } from "api-wrapper"
import SignOut from "pages/signout"
import React from "react"
import mockRouter from "next-router-mock"
import { mocked } from "jest-mock"
import { noUserMock } from "mocks"

describe("SignOut (page)", () => {
  it("Performs signout on API, sets User on undefined and redirects to /signin", async () => {
    mocked(signOut).mockImplementation(() => Promise.resolve(undefined))

    await mockRouter.push("/signout")

    render(<SignOut />)

    await waitFor(() => expect(document.title).toBe("Loading..."))

    await waitFor(() => expect(signOut).toBeCalled())
    await waitFor(() => mockRouter.isReady)

    expect(mockRouter.asPath).toBe("/signin")
  })

  it("Redirect anonymous", async () => {
    mocked(getMyUser).mockImplementation(noUserMock)

    await mockRouter.push("/signout")

    render(<SignOut />)

    await waitFor(() => expect(document.title).toBe("Loading..."))

    await waitFor(() => expect(mockRouter.isReady))
    expect(mockRouter.asPath).toBe("/signin")
  })
})
