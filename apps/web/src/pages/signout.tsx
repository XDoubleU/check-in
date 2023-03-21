import { useAuth } from "@/contexts"
import LoadingLayout from "@/layouts/LoadingLayout"
import { signOut } from "my-api-wrapper"
import Router from "next/router"
import { useEffect } from "react"

export default function SignOut() {
  const { setUser } = useAuth()

  useEffect(() => {
    void signOut()
      .then(() => setUser(undefined))
      .then(() => Router.push("/signin"))
  }, [setUser])

  return <LoadingLayout />
}
