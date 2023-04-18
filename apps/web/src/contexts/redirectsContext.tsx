import LoadingLayout from "../layouts/LoadingLayout"
import { type NextRouter, useRouter } from "next/router"
import React, { useEffect, useState, type ReactNode } from "react"
import { Role, type User } from "types-custom"
import { useAuth } from "./authContext"

interface Props {
  children: ReactNode
}

// eslint-disable-next-line @typescript-eslint/naming-convention, sonarjs/cognitive-complexity
function checkRedirects(currentUser: User | undefined, router: NextRouter) {
  if (!currentUser && router.pathname !== "/signin") {
    return router.push("/signin")
  }

  if (router.pathname === "/" && currentUser?.roles.includes(Role.Manager)) {
    return router.push("/settings")
  }

  if (router.pathname === "/" && currentUser?.roles.includes(Role.Admin)) {
    return router.push("/admin")
  }

  if (router.pathname === "/admin" && currentUser?.roles.includes(Role.Admin)) {
    return router.push("/admin/migrations")
  }

  if (router.pathname === "/settings") {
    if (currentUser?.roles.includes(Role.Manager)) {
      return router.push("/settings/locations")
    }
    if (currentUser?.roles.includes(Role.User) && currentUser.locationId) {
      return router.push(`/settings/locations/${currentUser.locationId}`)
    }
  }
  if (
    router.pathname === "/settings/schools" &&
    !currentUser?.roles.includes(Role.Manager)
  ) {
    return router.push("/settings")
  }
  if (
    router.pathname === "/settings/locations" &&
    !currentUser?.roles.includes(Role.Manager)
  ) {
    return router.push("/settings")
  }
  return new Promise((resolve) => resolve(true))
}

export const RedirectsProvider = ({ children }: Props) => {
  const router = useRouter()
  const { user, loadingUser } = useAuth()
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    if (!loadingUser) {
      void checkRedirects(user, router).then(() => setLoading(false))
    }
  }, [loading, loadingUser, router, setLoading, user])

  if (loading || loadingUser) {
    return <LoadingLayout />
  }

  return <>{children}</>
}
