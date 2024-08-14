import { getMyUser } from "api-wrapper"
import { type Role, type User } from "api-wrapper/types/apiTypes"
import LoadingLayout from "layouts/LoadingLayout"
import Router from "next/router";
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
  }, [setLoading])

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
  redirects: Map<Role, string> | undefined,
  user: User | undefined
) {
  if (!user) {
    if (!Router.asPath.includes("/signin")) {
      let query: ParsedUrlQueryInput | undefined

      if (!Router.asPath.includes("/signout") && Router.asPath !== "/") {
        query = { redirect_to: Router.asPath }
      }

      if (Router.asPath.includes("/settings/locations")) {
        query = { redirect_to: "/settings/locations" }
      }

      return Router.push({ pathname: `/signin`, query })
    }

    return new Promise((resolve) => resolve(true))
  }

  if (Router.asPath.includes("/signin")) {
    const redirectPath = (Router.query.redirect_to as string | undefined) ?? "/"
    return Router.push(redirectPath)
  }

  const redirectUrl = redirects?.get(user.role)
  if (redirectUrl) {
    let query: ParsedUrlQueryInput | undefined
    if (redirectUrl.includes("[id]")) {
      query = { id: user.location?.id }
    }

    return Router.push({
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
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    if (!loadingUser) {
      void redirect(redirects, user).then(() => setLoading(false))
    }
  }, [loadingUser, redirects, user])

  return <>{loading ? <LoadingLayout /> : children}</>
}
