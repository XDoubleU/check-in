import LoadingLayout from "@/layouts/LoadingLayout"
import { getMyUser } from "my-api-wrapper"
import Router from "next/router"
import React, {
  useState,
  type SetStateAction,
  type Dispatch,
  type ReactNode,
  useEffect
} from "react"
import { type User } from "types-custom"

export interface AuthContextProps {
  user: User | undefined
  loading: boolean
  setUser: Dispatch<SetStateAction<User | undefined>>
  setLoading: Dispatch<SetStateAction<boolean>>
}

interface Props {
  children: ReactNode
}

export const AuthContext = React.createContext<AuthContextProps>({
  user: undefined,
  loading: true,
  // eslint-disable-next-line @typescript-eslint/no-empty-function
  setUser: () => {},
  // eslint-disable-next-line @typescript-eslint/no-empty-function
  setLoading: () => {}
})

export const AuthProvider = ({ children }: Props) => {
  const [currentUser, setCurrentUser] = useState<User | undefined>()
  const [loading, setLoading] = useState<boolean>(true)

  useEffect(() => {
    void getMyUser().then(async (response) => {
      if (!response.ok) {
        await Router.push("/signin")
      } else {
        setCurrentUser(response.data)
      }
      setLoading(false)
    })
  }, [])

  if (loading) {
    return <LoadingLayout />
  }

  return (
    <AuthContext.Provider
      value={{
        user: currentUser,
        loading: loading,
        setUser: setCurrentUser,
        setLoading: setLoading
      }}
    >
      {children}
    </AuthContext.Provider>
  )
}
