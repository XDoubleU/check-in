/* eslint-disable sonarjs/no-duplicate-string */
import { type Role } from "api-wrapper/types/apiTypes"
import Loader from "components/Loader"
import { AuthRedirecter } from "contexts/authContext"
import ManagerLayout from "layouts/ManagerLayout"

export const ManagerRedirects = new Map<Role, string>([
  ["default", "/manual/location"]
])

export default function SettingsHome() {
  const redirects = new Map<Role, string>([
    ["admin", "/manual/manager"],
    ["manager", "/manual/manager"],
    ["default", "/manual/location"]
  ])

  return (
    <AuthRedirecter redirects={redirects}>
      <ManagerLayout title="">
        <Loader message="Loading manual." />
      </ManagerLayout>
    </AuthRedirecter>
  )
}
