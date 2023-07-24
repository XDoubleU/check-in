import LoadingLayout from "layouts/LoadingLayout"
import { type NextRouter, useRouter } from "next/router"
import React, { useEffect, useState, type ReactNode } from "react"
import { useAuth } from "./authContext"
import { type User } from "api-wrapper/types/apiTypes"

interface Props {
  children: ReactNode
}

// eslint-disable-next-line @typescript-eslint/naming-convention, sonarjs/cognitive-complexity
function checkRedirects(currentUser: User | undefined, router: NextRouter) {
  if (!currentUser && router.pathname !== "/signin") {
    return router.push("/signin")
  }

  if (router.pathname === "/" && currentUser?.role === "manager") {
    return router.push("/settings")
  }

  if (router.pathname === "/" && currentUser?.role === "admin") {
    return router.push("/admin")
  }

  if (router.pathname === "/admin" && currentUser?.role === "admin") {
    return router.push("/admin/migrations")
  }

  // eslint-disable-next-line sonarjs/no-collapsible-if
  if (router.pathname === "/settings") {
    if (currentUser?.role === "manager") {
      return router.push("/settings/locations")
    }
    // eslint-disable-next-line no-warning-comments
    //TODO: fix
    /*if (currentUser?.role === "default" && currentUser.locationId) {
      return router.push(`/settings/locations/${currentUser.locationId}`)
    }*/
  }
  if (
    router.pathname === "/settings/schools" &&
    currentUser?.role !== "manager"
  ) {
    return router.push("/settings")
  }
  if (
    router.pathname === "/settings/locations" &&
    currentUser?.role !== "manager"
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
