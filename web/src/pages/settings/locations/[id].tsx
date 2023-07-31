import ManagerLayout from "layouts/AdminLayout"
import { LocationUpdateModal } from "components/cards/LocationCard"
import { getLocation, getUser, type APIResponse } from "api-wrapper"
import { useCallback, useEffect, useState } from "react"
import { useRouter } from "next/router"
import LoadingLayout from "layouts/LoadingLayout"
import { type LocationWithUsername } from "."
import { AuthRedirecter, useAuth } from "contexts/authContext"
import Charts from "components/charts/Charts"
import { type User } from "api-wrapper/types/apiTypes"

// eslint-disable-next-line max-lines-per-function
export default function LocationDetail() {
  const { user } = useAuth()
  const router = useRouter()
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
    if (user?.role !== "default") {
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
      yesterdayFullAt: responseLocation.data.yesterdayFullAt,
      timeZone: responseLocation.data.timeZone
    }

    updateLocation(locationWithUsername)
  }, [router, user?.role, user?.username])

  useEffect(() => {
    void fetchData()
  }, [fetchData])

  return (
    <AuthRedirecter>
      {!location ? (
        <LoadingLayout message="User has no location" />
      ) : (
        <ManagerLayout
          title={location.name}
          titleButton={
            <LocationUpdateModal data={location} fetchData={fetchData} />
          }
        >
          <Charts locationId={location.id} />
        </ManagerLayout>
      )}
    </AuthRedirecter>
  )
}
