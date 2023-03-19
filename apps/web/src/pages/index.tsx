import { useAuth } from "@/contexts"
import LoadingLayout from "@/layouts/LoadingLayout"
import Router from "next/router"
import { Role } from "types-custom"

export default function Home() {
  const { user } = useAuth()

  if (user?.roles.includes(Role.Admin)) {
    void Router.push("/settings")
  } else {
    void Router.push("/check-in")
  }

  return <LoadingLayout />
}
