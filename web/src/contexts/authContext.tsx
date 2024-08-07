import { getMyUser } from "api-wrapper"
import { type Role, type User } from "api-wrapper/types/apiTypes"
import LoadingLayout from "layouts/LoadingLayout"
import { type NextRouter, useRouter } from "next/router"
import { type ParsedUrlQueryInput } from "querystring"
import React, {
  useState,
  type SetStateAction,
  type Dispatch,
  type ReactNode,
  useEffect,
  useContext
} from "react"
import * as Sentry from "@sentry/nextjs"

interface AuthContextProps {
  user: User | undefined
  setUser: Dispatch<SetStateAction<User | undefined>>
  loadingUser: boolean
}

interface AuthProviderProps {
  children: ReactNode
}

interface AuthRedirecterProps {
  children: ReactNode
  redirects?: Map<Role, string>
}

// eslint-disable-next-line @typescript-eslint/naming-convention
const AuthContext = React.createContext<AuthContextProps>({
  user: undefined,
  // eslint-disable-next-line @typescript-eslint/no-empty-function
  setUser: () => {},
  loadingUser: true
})

export function useAuth(): AuthContextProps {
  return useContext(AuthContext)
}

export const AuthProvider = ({ children }: AuthProviderProps) => {
  const router = useRouter()
  const [currentUser, setCurrentUser] = useState<User | undefined>()
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    if (!currentUser) {
      Sentry.setUser({})
      return
    }

    Sentry.setUser({
      id: currentUser.id,
      username: currentUser.username
    })
  }, [currentUser])

  useEffect(() => {
    void getMyUser()
      .then((response) => {
        if (!response.ok) {
          return
        }
        setCurrentUser(response.data)
        return response.data
      })
      .then(() => setLoading(false))
  }, [router, setLoading])

  return (
    <AuthContext.Provider
      value={{
        user: currentUser,
        setUser: setCurrentUser,
        loadingUser: loading
      }}
    >
      {children}
    </AuthContext.Provider>
  )
}

function redirect(
  router: NextRouter,
  redirects: Map<Role, string> | undefined,
  user: User | undefined
) {
  if (!user) {
    if (!router.asPath.includes("/signin")) {
      let query: ParsedUrlQueryInput | undefined

      if (!router.asPath.includes("/signout") && router.asPath !== "/") {
        query = { redirect_to: router.asPath }
      }

      if (router.asPath.includes("/settings/locations")) {
        query = { redirect_to: "/settings/locations" }
      }

      return router.push({ pathname: `/signin`, query })
    }

    return new Promise((resolve) => resolve(true))
  }

  if (router.asPath.includes("/signin")) {
    const redirectPath = (router.query.redirect_to as string | undefined) ?? "/"
    return router.push(redirectPath)
  }

  const redirectUrl = redirects?.get(user.role)
  if (redirectUrl) {
    let query: ParsedUrlQueryInput | undefined
    if (redirectUrl.includes("[id]")) {
      query = { id: user.location?.id }
    }

    return router.push({
      pathname: redirectUrl,
      query
    })
  }

  return new Promise((resolve) => resolve(true))
}

export const AuthRedirecter = ({
  children,
  redirects
}: AuthRedirecterProps) => {
  const { user, loadingUser } = useContext(AuthContext)
  const router = useRouter()
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    if (!loadingUser) {
      void redirect(router, redirects, user).then(() => setLoading(false))
    }
  }, [loadingUser, redirects, router, user])

  return <>{loading ? <LoadingLayout /> : children}</>
}
