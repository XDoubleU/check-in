import { getMyUser } from "api-wrapper"
import { type User } from "api-wrapper/types/apiTypes"
import { useRouter } from "next/router"
import React, {
  useState,
  type SetStateAction,
  type Dispatch,
  type ReactNode,
  useEffect,
  useContext
} from "react"

interface AuthContextProps {
  user: User | undefined
  setUser: Dispatch<SetStateAction<User | undefined>>
  loadingUser: boolean
}

interface Props {
  children: ReactNode
}

// eslint-disable-next-line @typescript-eslint/naming-convention
const AuthContext = React.createContext<AuthContextProps>({
  user: undefined,
  // eslint-disable-next-line @typescript-eslint/no-empty-function
  setUser: () => {},
  loadingUser: true
})

// eslint-disable-next-line @typescript-eslint/naming-convention
export function useAuth(): AuthContextProps {
  return useContext(AuthContext)
}

export const AuthProvider = ({ children }: Props) => {
  const router = useRouter()
  const [currentUser, setCurrentUser] = useState<User | undefined>()
  const [loading, setLoading] = useState(true)

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
