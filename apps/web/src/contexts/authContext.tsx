import { getMyUser } from "my-api-wrapper"
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
  loadingUser: boolean
  setUser: Dispatch<SetStateAction<User | undefined>>
}

interface Props {
  children: ReactNode
}

export const AuthContext = React.createContext<AuthContextProps>({
  user: undefined,
  loadingUser: true,
  // eslint-disable-next-line @typescript-eslint/no-empty-function
  setUser: () => {}
})

// eslint-disable-next-line max-lines-per-function
export const AuthProvider = ({ children }: Props) => {
  const [currentUser, setCurrentUser] = useState<User | undefined>()
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    void getMyUser()
      .then((response) => {
        if (!response.ok) {
          return
        }
        return setCurrentUser(response.data)
      })
      .then(() => setLoading(false))
  }, [loading, setLoading])

  return (
    <AuthContext.Provider
      value={{
        user: currentUser,
        loadingUser: loading,
        setUser: setCurrentUser
      }}
    >
      {children}
    </AuthContext.Provider>
  )
}
