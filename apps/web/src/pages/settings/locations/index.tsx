import CustomPagination, {
  type CustomPaginationProps
} from "@/components/CustomPagination"
import LocationCard from "@/components/cards/LocationCard"
import AdminLayout from "@/layouts/AdminLayout"
import { type CreateLocationDto, type Location } from "types-custom"
import { useRouter } from "next/router"
import { useCallback, useEffect, useState } from "react"
import { Form } from "react-bootstrap"
import { createLocation, getAllLocations, getUser } from "my-api-wrapper"
import { useForm } from "react-hook-form"
import CreateModal from "@/components/modals/CreateModal"
import Loader from "@/components/Loader"

export type LocationWithUsername = Omit<Location, "userId"> & {
  username: string
}
type LocationCreateForm = CreateLocationDto & { repeatPassword?: string }

interface LocationList {
  locations: LocationWithUsername[] | undefined
  pagination: CustomPaginationProps
}

// eslint-disable-next-line max-lines-per-function
export default function LocationList() {
  const router = useRouter()

  const [locationList, setLocationList] = useState<LocationList>({
    locations: undefined,
    pagination: {
      current: 0,
      total: 0
    }
  })

  const form = useForm<LocationCreateForm>()

  const {
    register,
    watch,
    formState: { errors }
  } = form

  const fetchData = useCallback(async () => {
    if (!router.isReady) return

    const page = router.query.page
      ? parseInt(router.query.page as string)
      : undefined

    const response = await getAllLocations(page)
    if (!response.data) return

    if (
      response.data.totalPages !== 0 &&
      response.data.page > response.data.totalPages
    ) {
      await router.push(`locations?page=${response.data.totalPages}`)
    }

    const locationsWithUsernames = Array<LocationWithUsername>()

    for (const location of response.data.locations) {
      const username = (await getUser(location.userId)).data?.username

      locationsWithUsernames.push({
        id: location.id,
        name: location.name,
        normalizedName: location.normalizedName,
        capacity: location.capacity,
        username: username ?? "",
        available: location.available,
        checkIns: location.checkIns,
        yesterdayFullAt: location.yesterdayFullAt
      })
    }

    setLocationList({
      locations: locationsWithUsernames,
      pagination: {
        current: response.data.page,
        total: response.data.totalPages
      }
    })
  }, [router])

  useEffect(() => {
    void fetchData()
  }, [fetchData])

  const handleCreate = (data: CreateLocationDto) => {
    return createLocation(data)
  }

  return (
    <AdminLayout title="Locations">
      <CreateModal<CreateLocationDto, Location>
        form={form}
        handler={handleCreate}
        refetchData={fetchData}
        typeName="location"
      >
        <Form.Group className="mb-3">
          <Form.Label>Name</Form.Label>
          <Form.Control
            type="text"
            placeholder="Name"
            required
            {...register("name")}
          ></Form.Control>
        </Form.Group>
        <Form.Group className="mb-3">
          <Form.Label>Capacity</Form.Label>
          <Form.Control
            type="number"
            required
            {...register("capacity")}
          ></Form.Control>
        </Form.Group>
        <Form.Group className="mb-3">
          <Form.Label>Username</Form.Label>
          <Form.Control
            type="text"
            placeholder="Username"
            required
            {...register("username")}
          ></Form.Control>
        </Form.Group>
        <Form.Group className="mb-3">
          <Form.Label>Password</Form.Label>
          <Form.Control
            type="password"
            placeholder="Password"
            required
            {...register("password")}
          ></Form.Control>
        </Form.Group>
        <Form.Group className="mb-3">
          <Form.Label>Repeat password</Form.Label>
          <Form.Control
            type="password"
            placeholder="Repeat password"
            required
            isInvalid={!!errors.repeatPassword}
            {...register("repeatPassword", {
              validate: (val: string | undefined) => {
                if (watch("password") != val) {
                  return "Your passwords do no match"
                }
                return undefined
              }
            })}
          ></Form.Control>
          <Form.Control.Feedback type="invalid">
            {errors.repeatPassword?.message}
          </Form.Control.Feedback>
        </Form.Group>
      </CreateModal>

      <br />

      <div className="min-vh-51">
        {!locationList.locations && <Loader />}

        {locationList.locations && locationList.locations.length == 0
          ? "Nothing to see here."
          : ""}

        {locationList.locations?.map((location) => {
          return (
            <LocationCard
              key={location.id}
              location={location}
              refetchData={fetchData}
            />
          )
        })}
      </div>

      <CustomPagination
        current={locationList.pagination.current}
        total={locationList.pagination.total}
      />
    </AdminLayout>
  )
}
