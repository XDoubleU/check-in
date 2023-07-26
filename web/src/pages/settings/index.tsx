import { type Role } from "api-wrapper/types/apiTypes"
import Loader from "components/Loader"
import { Redirecter } from "components/Redirecter"
import ManagerLayout from "layouts/AdminLayout"

export default function SettingsHome() {
  const redirects = new Map<Role, string>([
    ["admin", "/settings/locations"],
    ["manager", "/settings/locations"],
    ["default", "/settings/locations/{locationId}"]
  ])

  return (
    <Redirecter redirects={redirects}>
      <ManagerLayout title="">
        <Loader message="Loading home page." />
      </ManagerLayout>
    </Redirecter>
  )
}
