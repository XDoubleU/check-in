import { useAuth } from "contexts/authContext"
import LoadingLayout from "layouts/LoadingLayout"
import { signOut } from "api-wrapper"
import { useRouter } from "next/router"
import { useEffect } from "react"
import { Redirecter } from "components/Redirecter"

export default function SignOut() {
  const { setUser } = useAuth()
  const router = useRouter()

  useEffect(() => {
    void signOut()
      .then(() => setUser(undefined))
      .then(() => router.push("/signin"))
  }, [router, setUser])

  return (
    <Redirecter>
      <LoadingLayout message="Signing out." />
    </Redirecter>
  )
}
