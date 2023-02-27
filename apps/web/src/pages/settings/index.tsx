import LoadingLayout from "@/layouts/LoadingLayout"
import { getMyUser } from "api-wrapper"
import Router from "next/router"
import { useEffect, useState } from "react"
import { Role, User } from "types"

export default function SettingsHome() {
  const [user, setUser] = useState<User>()

  useEffect(() => {
    void getMyUser()
      .then(async (data) => {
        if (data === null) {
          await Router.push("/signin")
        } else {
          setUser(data)
        }
      })
  })

  if (user === undefined) {
    return <LoadingLayout/>
  }

  if (user.roles.includes(Role.Admin) || !user.locationId) {
    void Router.push("/settings/locations")
  } else {
    void Router.push(`/settings/locations/${user.locationId}`)
  }

  return <LoadingLayout/>
}