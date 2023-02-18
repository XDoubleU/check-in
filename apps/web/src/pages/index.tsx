import LoadingLayout from "@/layouts/LoadingLayout"
import { getUserInfo } from "api-wrapper"
import Router from "next/router"
import { useEffect, useState } from "react"
import { User } from "types"

export default function Home() {
  const [userInfo, setUserInfo] = useState<User>()

  useEffect(() => {
    getUserInfo()
      .then(data => {
        if (data === null) {
          Router.push("/signin")
        } else {
          setUserInfo(data)
        }
      })
  }, [])

  if (userInfo === undefined) {
    return <LoadingLayout/>
  }

  if (userInfo.isAdmin) {
    Router.push("/settings")
  } else {
    Router.push("/check-in")
  }

  return <LoadingLayout/>
}

