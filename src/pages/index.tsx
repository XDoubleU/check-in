import CheckInLayout from "@/layouts/CheckInLayout"
import LoadingLayout from "@/layouts/LoadingLayout"
import { useSession } from "next-auth/react"
import Router from "next/router"

export default function CheckIn() {
  const {data, status} = useSession({
    required: true
  })

  if (status == "loading") {
    return <LoadingLayout/>
  }

  if (data.user.isAdmin) {
    Router.push("/settings")
    return <LoadingLayout/>
  }

  return <CheckInLayout/>
}
