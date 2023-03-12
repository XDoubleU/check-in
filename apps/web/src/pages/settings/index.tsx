import LoadingLayout from "@/layouts/LoadingLayout"
import { getMyUser } from "my-api-wrapper"
import Router from "next/router"
import { useEffect, useState } from "react"
import { Role, User } from "types-custom"

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

  if (user.roles.includes(Role.Admin) || !user.location?.id) {
    void Router.push("/settings/locations")
  } else {
    void Router.push(`/settings/locations/${user.location.id}`)
  }

  return <LoadingLayout/>
}