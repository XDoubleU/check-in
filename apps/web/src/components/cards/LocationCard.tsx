import { Alert, Card, Form } from "react-bootstrap"
import Link from "next/link"
import UpdateModal from "@/components/modals/UpdateModal"
import DeleteModal from "@/components/modals/DeleteModal"
import { deleteLocation, updateLocation } from "my-api-wrapper"
import { type UpdateLocationDto } from "types-custom"
import { useForm } from "react-hook-form"
import { format } from "date-fns"
import { type LocationWithUsername } from "@/pages/settings/locations"

type LocationUpdateForm = UpdateLocationDto & { repeatPassword?: string }

interface LocationCardProps {
  location: LocationWithUsername
  refetchData: () => Promise<void>
}

// eslint-disable-next-line max-lines-per-function
export function LocationUpdateModal({
  location,
  refetchData
}: LocationCardProps) {
  const form = useForm<LocationUpdateForm>({
    defaultValues: {
      name: location.name,
      capacity: location.capacity,
      username: location.username
    }
  })

  const {
    register,
    watch,
    formState: { errors }
  } = form

  const handleUpdate = (data: UpdateLocationDto) => {
    return updateLocation(location.id, data)
  }

  return (
    <UpdateModal<UpdateLocationDto>
      form={form}
      handler={handleUpdate}
      refetchData={refetchData}
      typeName="location"
    >
      <Form.Group className="mb-3">
        <Form.Label>Name</Form.Label>
        <Form.Control
          type="text"
          placeholder="Name"
          {...register("name")}
        ></Form.Control>
      </Form.Group>
      <Form.Group className="mb-3">
        <Form.Label>Capacity</Form.Label>
        <Form.Control
          type="number"
          placeholder="Capacity"
          {...register("capacity")}
        ></Form.Control>
      </Form.Group>
      <Form.Group className="mb-3">
        <Form.Label>Username</Form.Label>
        <Form.Control
          type="text"
          placeholder="Username"
          {...register("username")}
        ></Form.Control>
      </Form.Group>
      <Form.Group className="mb-3">
        <Form.Label>Password</Form.Label>
        <Form.Control
          type="password"
          placeholder="Password"
          {...register("password")}
        ></Form.Control>
      </Form.Group>
      <Form.Group className="mb-3">
        <Form.Label>Repeat password</Form.Label>
        <Form.Control
          type="password"
          placeholder="Repeat password"
          {...register("repeatPassword", {
            validate: (val: string | undefined) => {
              if (watch("password") != val) {
                return "Your passwords do no match"
              }
              return undefined
            }
          })}
        ></Form.Control>
        {errors.repeatPassword && (
          <Alert key="danger">{errors.repeatPassword.message}</Alert>
        )}
      </Form.Group>
    </UpdateModal>
  )
}

function LocationDeleteModal({ location, refetchData }: LocationCardProps) {
  const handleDelete = () => {
    return deleteLocation(location.id)
  }

  return (
    <DeleteModal
      name={location.name}
      handler={handleDelete}
      refetchData={refetchData}
      typeName="location"
    />
  )
}

// eslint-disable-next-line max-lines-per-function
export default function LocationCard({
  location,
  refetchData
}: LocationCardProps) {
  return (
    <>
      <Card>
        <Card.Body>
          <div className="d-flex flex-row">
            <div>
              <Card.Title>
                <Link href={`/settings/locations/${location.id}`}>
                  {location.name}
                </Link>{" "}
                ({location.normalizedName})
              </Card.Title>
              <Card.Subtitle className="mb-2 text-muted">
                {location.available} / {location.capacity}
              </Card.Subtitle>
              <Card.Subtitle className="mb-2 text-muted">
                {location.yesterdayFullAt
                  ? `Yesterday full at ${format(
                      new Date(location.yesterdayFullAt),
                      "HH:mm"
                    )}`
                  : "Yesterday not full"}
              </Card.Subtitle>
            </div>
            <div className="ms-auto">
              <LocationUpdateModal
                location={location}
                refetchData={refetchData}
              />
              <LocationDeleteModal
                location={location}
                refetchData={refetchData}
              />
            </div>
          </div>
        </Card.Body>
      </Card>
      <br />
    </>
  )
}
