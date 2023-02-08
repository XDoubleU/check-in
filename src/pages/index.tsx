import LoadingLayout from "@/layouts/LoadingLayout"
import { useSession } from "next-auth/react"
import Router from "next/router"

export default function Home() {
  const {data, status} = useSession({
    required: true
  })

  if (status == "loading") {
    return <LoadingLayout/>
  }

  if (data.user.isAdmin) {
    Router.push("/settings")
  } else {
    Router.push("/check-in")
  }

  return <LoadingLayout/>
}
