import { useAuth } from "@/contexts"
import LoadingLayout from "@/layouts/LoadingLayout"
import Router from "next/router"
import { Role } from "types-custom"

export default function SettingsHome() {
  const { user } = useAuth()

  if (user?.roles.includes(Role.Admin) || !user?.location?.id) {
    void Router.push("/settings/locations")
  } else {
    void Router.push(`/settings/locations/${user.location.id}`)
  }

  return <LoadingLayout />
}
