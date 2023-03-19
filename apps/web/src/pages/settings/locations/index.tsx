import CustomButton from "@/components/CustomButton"
import CustomPagination, {
  type CustomPaginationProps
} from "@/components/CustomPagination"
import LocationCard from "@/components/cards/LocationCard"
import AdminLayout from "@/layouts/AdminLayout"
import { type CreateLocationDto, type Location } from "types-custom"
import { useRouter } from "next/router"
import { useCallback, useEffect, useState } from "react"
import { Alert, Col, Form, Modal } from "react-bootstrap"
import { createLocation, getAllLocations } from "my-api-wrapper"
import { type SubmitHandler, useForm } from "react-hook-form"

interface LocationList {
  locations: Location[]
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

  const {
    register,
    watch,
    handleSubmit,
    setError,
    reset,
    formState: { errors }
  } = useForm<CreateLocationDto & { repeatPassword: string }>()

  const [showCreate, setShowCreate] = useState(false)
  const handleCloseCreate = () => setShowCreate(false)
  const handleShowCreate = () => setShowCreate(true)
  const onCloseCreate = useCallback(() => {
    return !showCreate
  }, [showCreate])

  useEffect(() => {
    if (!router.isReady) return
    const page = router.query.page
      ? parseInt(router.query.page as string)
      : undefined
    void getAllLocations(page).then(async (response) => {
      if (!response.ok) {
        await router.push("/signin")
        return
      }

      setLocationList({
        locations: response.data?.locations ?? Array<Location>(),
        pagination: {
          current: response.data?.page ?? 1,
          total: response.data?.totalPages ?? 1
        }
      })
    })
  }, [onCloseCreate, router])

  const onSubmit: SubmitHandler<CreateLocationDto> = async (data) => {
    const response = await createLocation(data)

    if (!response.ok) {
      setError("root", {
        message: response.message ?? "Something went wrong"
      })
    } else {
      handleCloseCreate()
      reset()
    }
  }

  return (
    <AdminLayout title="Locations">
      <Modal show={showCreate} onHide={handleCloseCreate}>
        <Modal.Body>
          <Modal.Title>Create location</Modal.Title>
          <br />
          <Form onSubmit={handleSubmit(onSubmit)}>
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
                {...register("repeatPassword", {
                  validate: (val: string) => {
                    if (watch("password") != val) {
                      return "Your passwords do no match"
                    }
                    return ""
                  }
                })}
              ></Form.Control>
            </Form.Group>
            {errors.root && <Alert key="danger">{errors.root.message}</Alert>}
            <br />
            <CustomButton
              type="button"
              style={{ float: "left" }}
              onClick={handleCloseCreate}
            >
              Cancel
            </CustomButton>
            <CustomButton type="submit" style={{ float: "right" }}>
              Create
            </CustomButton>
          </Form>
        </Modal.Body>
      </Modal>

      <Col size={2}>
        <CustomButton onClick={handleShowCreate}>Create</CustomButton>
      </Col>

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
              username={"TODO"}
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
