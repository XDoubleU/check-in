import LoadingLayout from "@/layouts/LoadingLayout"
import { getMyUser } from "api-wrapper"
import Router from "next/router"
import { useEffect, useState } from "react"
import { Role, User } from "types"

export default function Home() {
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
  }, [])

  if (user === undefined) {
    return <LoadingLayout/>
  }

  if (user.roles.includes(Role.Admin)) {
    void Router.push("/settings")
  } else {
    void Router.push("/check-in")
  }

  return <LoadingLayout/>
}

