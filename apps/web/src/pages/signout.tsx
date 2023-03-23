import { useAuth } from "@/contexts"
import LoadingLayout from "@/layouts/LoadingLayout"
import { signOut } from "my-api-wrapper"
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

  return <LoadingLayout />
}
