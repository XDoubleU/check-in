import { AuthRedirecter, useAuth } from "contexts/authContext"
import LoadingLayout from "layouts/LoadingLayout"
import { signOut } from "api-wrapper"
import { useEffect } from "react"

export default function SignOut() {
  const { setUser } = useAuth()

  useEffect(() => {
    void signOut().then(() => setUser(undefined))
  })

  return (
    <AuthRedirecter>
      <LoadingLayout message="Signing out." />
    </AuthRedirecter>
  )
}
