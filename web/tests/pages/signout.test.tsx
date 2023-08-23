import { render, waitFor } from "test-utils"
import { signOut, getMyUser } from "api-wrapper"
import SignOut from "pages/signout"
import React from "react"
import mockRouter from "next-router-mock"

(getMyUser as jest.Mock).mockImplementation(() => Promise.resolve({
  ok: false
}))

describe("SignOut (page)", () => {
  it("Performs signout on API, sets User on undefined and redirects to /signin", async () => {
    await mockRouter.push("/signout")

    render(<SignOut />)

    await waitFor(() => expect(document.title).toBe("Loading..."))

    await waitFor(() => expect(signOut).toBeCalled())
    await waitFor(() => mockRouter.isReady)

    expect(mockRouter.asPath).toBe("/signin")
  })
})
