import { render, waitFor } from "test-utils"
import { signOut } from "api-wrapper"
import SignOut from "pages/signout"
import React from "react"
import mockRouter from "next-router-mock"
import { mocked } from "jest-mock"

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
})
