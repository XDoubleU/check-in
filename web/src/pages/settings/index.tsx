import { type Role } from "api-wrapper/types/apiTypes"
import Loader from "components/Loader"
import { AuthRedirecter } from "contexts/authContext"
import ManagerLayout from "layouts/AdminLayout"

export default function SettingsHome() {
  const redirects = new Map<Role, string>([
    ["admin", "/settings/locations"],
    ["manager", "/settings/locations"],
    ["default", "/settings/locations/[id]"]
  ])

  return (
    <AuthRedirecter redirects={redirects}>
      <ManagerLayout title="">
        <Loader message="Loading home page." />
      </ManagerLayout>
    </AuthRedirecter>
  )
}
