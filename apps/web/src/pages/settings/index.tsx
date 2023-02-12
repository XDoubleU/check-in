import { useSession } from "next-auth/react"
import Router from "next/router"
import LoadingLayout from "@/layouts/LoadingLayout"

export default function SettingsHome() {
  const {data, status} = useSession({
    required: true
  })

  if (status == "loading") {
    return <LoadingLayout/>
  }

  if (data.user.isAdmin) {
    Router.push("/settings/locations")
  } else {
    Router.push(`/settings/locations/${data.user.locationId}`)
  }
  
  return <LoadingLayout/>
}