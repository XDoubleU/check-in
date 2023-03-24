import { Col } from "react-bootstrap"
import AdminLayout from "@/layouts/AdminLayout"
import { LocationUpdateModal } from "@/components/cards/LocationCard"
import CustomButton from "@/components/CustomButton"
import { getLocation, getUser } from "my-api-wrapper"
import { useCallback, useEffect, useState } from "react"
import { useRouter } from "next/router"
import LoadingLayout from "@/layouts/LoadingLayout"
import { type LocationWithUsername } from "."
import { useAuth } from "@/contexts"
import { Role, type User } from "types-custom"
import type APIResponse from "my-api-wrapper/dist/src/types/apiResponse"

// eslint-disable-next-line max-lines-per-function
export default function LocationDetail() {
  const router = useRouter()
  const { user } = useAuth()
  const [location, updateLocation] = useState<LocationWithUsername>()

  const fetchData = useCallback(async () => {
    if (!router.isReady) return

    const locationId = router.query.id as string

    const responseLocation = await getLocation(locationId)
    if (!responseLocation.data) {
      await router.push("locations")
      return
    }

    let responseUser: APIResponse<User> | undefined = undefined
    if (user?.roles.includes(Role.Admin)) {
      responseUser = await getUser(responseLocation.data.userId)
      if (!responseUser.data) return
    }

    const locationWithUsername = {
      id: responseLocation.data.id,
      name: responseLocation.data.name,
      normalizedName: responseLocation.data.normalizedName,
      capacity: responseLocation.data.capacity,
      username: responseUser?.data?.username ?? user?.username ?? "",
      available: responseLocation.data.available,
      checkIns: responseLocation.data.checkIns
    }

    updateLocation(locationWithUsername)
  }, [router, user?.roles, user?.username])

  useEffect(() => {
    void fetchData()
  }, [fetchData])

  if (!location) {
    return <LoadingLayout />
  }

  return (
    <AdminLayout title={location.name}>
      <Col size={2}>
        <LocationUpdateModal
          id={location.id}
          name={location.name}
          username={location.username}
          capacity={location.capacity}
          refetchData={fetchData}
        />
      </Col>
      <br />
      <Col size={2}>
        <CustomButton>Download CSV (TODO)</CustomButton>
      </Col>
      <br />
      Still needs a chart :)
    </AdminLayout>
  )
}
