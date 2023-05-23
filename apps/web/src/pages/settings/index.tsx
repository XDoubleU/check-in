import Loader from "components/Loader"
import ManagerLayout from "layouts/AdminLayout"

export default function SettingsHome() {
  return (
    <ManagerLayout title="">
      <Loader message="Loading home page." />
    </ManagerLayout>
  )
}
