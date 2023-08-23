import React, { type ReactElement } from "react"
import { render, type RenderOptions } from "@testing-library/react"
import { AuthProvider } from "contexts/authContext"

// eslint-disable-next-line @typescript-eslint/naming-convention
const AllTheProviders = ({ children }: { children: React.ReactNode }) => {
  return <AuthProvider>{children}</AuthProvider>
}

// eslint-disable-next-line @typescript-eslint/naming-convention
const customRender = (
  ui: ReactElement,
  options?: Omit<RenderOptions, "wrapper">
) => render(ui, { wrapper: AllTheProviders, ...options })

export * from "@testing-library/react"
export { customRender as render }
