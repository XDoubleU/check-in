/* eslint-disable @typescript-eslint/naming-convention */
import { useContext } from "react"
import { AuthContext, type AuthContextProps } from "./authContext"

export function useAuth(): AuthContextProps {
  return useContext(AuthContext)
}
