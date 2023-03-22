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

export type LocationWithUsername = Omit<Location, "userId"> & {
  username: string
}
type LocationCreateForm = CreateLocationDto & { repeatPassword?: string }

interface LocationList {
  locations: LocationWithUsername[]
  pagination: CustomPaginationProps
}

// eslint-disable-next-line max-lines-per-function
export default function LocationList() {
  const router = useRouter()

  const [locationList, setLocationList] = useState<LocationList>({
    locations: [],
    pagination: {
      current: 0,
      total: 0
    }
  })

  const form = useForm<LocationCreateForm>()

  const fetchData = useCallback(async () => {
    if (!router.isReady) return

    const page = router.query.page
      ? parseInt(router.query.page as string)
      : undefined

    const response = await getAllLocations(page)
    if (!response.data) return

    const locationsWithUsernames = Array<LocationWithUsername>()

    await Promise.all(
      response.data.locations.map(async (location) => {
        const username = (await getUser(location.userId)).data?.username

        locationsWithUsernames.push({
          id: location.id,
          name: location.name,
          normalizedName: location.normalizedName,
          capacity: location.capacity,
          username: username ?? "",
          available: location.available,
          checkIns: location.checkIns
        })
      })
    )

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
            {...form.register("name")}
          ></Form.Control>
        </Form.Group>
        <Form.Group className="mb-3">
          <Form.Label>Capacity</Form.Label>
          <Form.Control
            type="number"
            required
            {...form.register("capacity")}
          ></Form.Control>
        </Form.Group>
        <Form.Group className="mb-3">
          <Form.Label>Username</Form.Label>
          <Form.Control
            type="text"
            placeholder="Username"
            required
            {...form.register("username")}
          ></Form.Control>
        </Form.Group>
        <Form.Group className="mb-3">
          <Form.Label>Password</Form.Label>
          <Form.Control
            type="password"
            placeholder="Password"
            required
            {...form.register("password")}
          ></Form.Control>
        </Form.Group>
        <Form.Group className="mb-3">
          <Form.Label>Repeat password</Form.Label>
          <Form.Control
            type="password"
            placeholder="Repeat password"
            required
            isInvalid={!!form.formState.errors.repeatPassword}
            {...form.register("repeatPassword", {
              validate: (val: string | undefined) => {
                if (form.watch("password") != val) {
                  return "Your passwords do no match"
                }
                return undefined
              }
            })}
          ></Form.Control>
          <Form.Control.Feedback type="invalid">
            {form.formState.errors.repeatPassword?.message}
          </Form.Control.Feedback>
        </Form.Group>
      </CreateModal>

      <br />

      <div className="min-vh-51">
        {locationList.locations.length == 0 ? "Nothing to see here." : ""}

        {locationList.locations.map((location) => {
          return (
            <LocationCard
              id={location.id}
              key={location.id}
              name={location.name}
              normalizedName={location.normalizedName}
              capacity={location.capacity}
              username={location.username}
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
