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
  setUser: Dispatch<SetStateAction<User | undefined>>
}

interface Props {
  children: ReactNode
}

export const AuthContext = React.createContext<AuthContextProps>({
  user: undefined,
  // eslint-disable-next-line @typescript-eslint/no-empty-function
  setUser: () => {}
})

export const AuthProvider = ({ children }: Props) => {
  const [currentUser, setCurrentUser] = useState<User | undefined>()

  useEffect(() => {
    void getMyUser().then(async (response) => {
      if (!response.ok) {
        await Router.push("/signin")
      } else {
        setCurrentUser(response.data)
      }
    })
  }, [])

  if (!currentUser) {
    return <LoadingLayout />
  }

  return (
    <AuthContext.Provider
      value={{
        user: currentUser,
        setUser: setCurrentUser
      }}
    >
      {children}
    </AuthContext.Provider>
  )
}
