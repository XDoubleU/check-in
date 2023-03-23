/* eslint-disable @typescript-eslint/naming-convention */
import { useContext } from "react"
import { AuthContext, type AuthContextProps } from "./authContext"
import { LoadingContext, type LoadingContextProps } from "./loadingContext"

export function useAuth(): AuthContextProps {
  return useContext(AuthContext)
}

export function useLoading(): LoadingContextProps {
  return useContext(LoadingContext)
}
