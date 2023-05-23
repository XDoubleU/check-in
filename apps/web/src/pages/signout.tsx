import { useAuth } from "contexts/authContext"
import LoadingLayout from "layouts/LoadingLayout"
import { signOut } from "api-wrapper"
import { useRouter } from "next/router"
import { useEffect } from "react"

export default function SignOut() {
  const router = useRouter()
  const { setUser } = useAuth()

  useEffect(() => {
    void signOut()
      .then(() => setUser(undefined))
      .then(() => router.push("/signin"))
  }, [router, setUser])

  return <LoadingLayout message="Signing out." />
}
