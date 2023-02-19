import LoadingLayout from "@/layouts/LoadingLayout"
import { getMyUser } from "api-wrapper"
import Router from "next/router"
import { useEffect, useState } from "react"
import { User } from "types"

export default function SettingsHome() {
  const [user, setUser] = useState<User>()

  useEffect(() => {
    getMyUser()
      .then(data => {
        if (data === null) {
          Router.push("/signin")
        } else {
          setUser(data)
        }
      })
  })

  if (user === undefined) {
    return <LoadingLayout/>
  }

  if (user.isAdmin) {
    Router.push("/settings/locations")
  } else {
    Router.push(`/settings/locations/${user.locationId}`)
  }

  return <LoadingLayout/>
}