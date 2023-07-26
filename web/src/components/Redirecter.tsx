import { useRouter } from "next/router"
import React, { useEffect, useState, type ReactNode } from "react"
import { type Role, type User } from "api-wrapper/types/apiTypes"
import LoadingLayout from "layouts/LoadingLayout"
import { useAuth } from "contexts/authContext"

interface Props {
  children: ReactNode
  redirects?: Map<Role, string>
}

function parseUrlVariables(url: string, user: User): string {
  if (user.location) {
    url = url.replace("{locationId}", user.location?.id)
  }

  return url
}

export function Redirecter({ children, redirects }: Props) {
  const router = useRouter()
  const { user, loadingUser } = useAuth()
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    if (!loadingUser) {
      void (async () => {
        if (!user) {
          if (router.pathname !== "/signin") {
            return router.push("/signin")
          }
          return new Promise((resolve) => resolve(true))
        }
  
        if (router.pathname === "/signin") {
          return router.push("/")
        }
  
        let redirectUrl = redirects?.get(user.role)
  
        if (redirectUrl) {
          redirectUrl = parseUrlVariables(redirectUrl, user)
          return router.push(redirectUrl)
        }
        return new Promise((resolve) => resolve(true))
      })().then(() => setLoading(false))
    }
    
  }, [loadingUser, redirects, router, user])

  return <>{loading || loadingUser ? <LoadingLayout /> : children}</>
}
