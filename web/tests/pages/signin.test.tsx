import userEvent from "@testing-library/user-event"
import { render, waitFor, screen, fireEvent } from "test-utils"
import SignIn from "pages/signin"
import React from "react"
import mockRouter from "next-router-mock"
import { signIn } from "api-wrapper"

describe("SignIn (page)", () => {
  it("Performs a successful signin", async () => {
    await mockRouter.push("/signin")

    render(<SignIn />)

    await waitFor(() => expect(document.title).toContain("Sign In"))

    const usernameInput = screen.getByLabelText("Username")
    const passwordInput = screen.getByLabelText("Password")
    const signInButton = screen.getByRole("button", { name: "Sign In" })

    await userEvent.type(usernameInput, "validusername")
    await userEvent.type(passwordInput, "validpassword")
    fireEvent.click(signInButton)

    await waitFor(() => expect(signIn).toBeCalled())
    expect(mockRouter.asPath).toBe("/")
  })

  it("Performs a non-successful signin", async () => {
    await mockRouter.push("/signin")

    render(<SignIn />)

    await waitFor(() => expect(document.title).toContain("Sign In"))

    const usernameInput = screen.getByLabelText("Username")
    const passwordInput = screen.getByLabelText("Password")
    const signInButton = screen.getByRole("button", { name: "Sign In" })

    await userEvent.type(usernameInput, "invalidusername")
    await userEvent.type(passwordInput, "invalidpassword")
    fireEvent.click(signInButton)

    await waitFor(() => expect(signIn).toBeCalled())

    await waitFor(() => expect(screen.getByRole("alert")).toHaveTextContent("Invalid credentials"))
  })
})
