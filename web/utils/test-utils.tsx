import React, { type ReactElement } from "react"
import { render, type RenderOptions } from "@testing-library/react"
import { AuthProvider } from "contexts/authContext"

const allTheProviders = ({ children }: { children: React.ReactNode }) => {
  return <AuthProvider>{children}</AuthProvider>
}

const CustomRender = (
  ui: ReactElement,
  options?: Omit<RenderOptions, "wrapper">
) => render(ui, { wrapper: allTheProviders, ...options })

export * from "@testing-library/react"
export { CustomRender as render }
