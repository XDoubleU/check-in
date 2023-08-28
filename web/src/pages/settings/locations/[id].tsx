import ManagerLayout from "layouts/ManagerLayout"
import { LocationUpdateModal } from "components/cards/LocationCard"
import { getLocation, getUser, getCheckInsToday } from "api-wrapper"
import { useCallback, useEffect, useState } from "react"
import { useRouter } from "next/router"
import LoadingLayout from "layouts/LoadingLayout"
import { type LocationWithUsername } from "."
import { AuthRedirecter, useAuth } from "contexts/authContext"
import Charts from "components/charts/Charts"
import { type CheckIn, type User } from "api-wrapper/types/apiTypes"
import CheckInCard from "components/cards/CheckInCard"
import Loader from "components/Loader"

// eslint-disable-next-line max-lines-per-function
export default function LocationDetail() {
  const { user } = useAuth()
  const router = useRouter()
  const [location, updateLocation] = useState<LocationWithUsername>()
  const [checkInsList, setCheckInsList] = useState<CheckIn[]>([])

  const fetchCheckInData = useCallback(async () => {
    const locationId = router.query.id as string

    // eslint-disable-next-line @typescript-eslint/no-unsafe-argument
    const response = await getCheckInsToday(locationId)
    setCheckInsList(response.data as CheckIn[])
  }, [router])

  const fetchData = useCallback(async () => {
    const locationId = router.query.id as string

    const responseLocation = await getLocation(locationId)
    if (!responseLocation.data) return

    let responseUser: User = user as User
    if (user?.role !== "default") {
      const response = await getUser(responseLocation.data.userId)
      responseUser = response.data as User
    }

    const locationWithUsername = {
      id: responseLocation.data.id,
      name: responseLocation.data.name,
      normalizedName: responseLocation.data.normalizedName,
      capacity: responseLocation.data.capacity,
      username: responseUser.username,
      available: responseLocation.data.available,
      yesterdayFullAt: responseLocation.data.yesterdayFullAt,
      timeZone: responseLocation.data.timeZone
    }

    updateLocation(locationWithUsername)
  }, [router, user])

  useEffect(() => {
    void fetchData()
    void fetchCheckInData()
  }, [fetchCheckInData, fetchData])

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
          <Charts locationIds={[location.id]} />

          <h2>Todays Check-Ins</h2>
          <br />

          {!checkInsList && <Loader message="Fetching data." />}

          {checkInsList.length == 0 ? "Nothing to see here." : ""}

          {checkInsList.map((item) => {
            return (
              <div key={item.id}>
                <CheckInCard
                  data={item}
                  readonly={user?.role !== "admin" && user?.role !== "manager"}
                  fetchData={fetchCheckInData}
                />
              </div>
            )
          })}
        </ManagerLayout>
      )}
    </AuthRedirecter>
  )
}
