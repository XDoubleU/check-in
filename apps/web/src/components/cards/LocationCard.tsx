import { Alert, Card, Form } from "react-bootstrap"
import Link from "next/link"
import UpdateModal from "@/components/modals/UpdateModal"
import DeleteModal from "@/components/modals/DeleteModal"
import { deleteLocation, updateLocation } from "my-api-wrapper"
import { type UpdateLocationDto } from "types-custom"
import { useForm } from "react-hook-form"

type LocationUpdateProps = Omit<LocationCardProps, "normalizedName">
type LocationUpdateForm = UpdateLocationDto & { repeatPassword?: string }

interface LocationCardProps {
  id: string
  name: string
  normalizedName: string
  capacity: number
  username: string
  refetchData: () => Promise<void>
}

// eslint-disable-next-line max-lines-per-function
export function LocationUpdateModal({
  id,
  name,
  capacity,
  username,
  refetchData
}: LocationUpdateProps) {
  const form = useForm<LocationUpdateForm>({
    defaultValues: {
      name: name,
      capacity: capacity,
      username: username
    }
  })

  const handleUpdate = (data: UpdateLocationDto) => {
    return updateLocation(id, data)
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
          {...form.register("name")}
        ></Form.Control>
      </Form.Group>
      <Form.Group className="mb-3">
        <Form.Label>Capacity</Form.Label>
        <Form.Control
          type="number"
          placeholder="Capacity"
          {...form.register("capacity")}
        ></Form.Control>
      </Form.Group>
      <Form.Group className="mb-3">
        <Form.Label>Username</Form.Label>
        <Form.Control
          type="text"
          placeholder="Username"
          {...form.register("username")}
        ></Form.Control>
      </Form.Group>
      <Form.Group className="mb-3">
        <Form.Label>Password</Form.Label>
        <Form.Control
          type="password"
          placeholder="Password"
          {...form.register("password")}
        ></Form.Control>
      </Form.Group>
      <Form.Group className="mb-3">
        <Form.Label>Repeat password</Form.Label>
        <Form.Control
          type="password"
          placeholder="Repeat password"
          {...form.register("repeatPassword", {
            validate: (val: string | undefined) => {
              if (form.watch("password") != val) {
                return "Your passwords do no match"
              }
              return undefined
            }
          })}
        ></Form.Control>
        {form.formState.errors.repeatPassword && (
          <Alert key="danger">
            {form.formState.errors.repeatPassword.message}
          </Alert>
        )}
      </Form.Group>
    </UpdateModal>
  )
}

function LocationDeleteModal({ id, name, refetchData }: LocationUpdateProps) {
  const handleDelete = () => {
    return deleteLocation(id)
  }

  return (
    <DeleteModal
      name={name}
      handler={handleDelete}
      refetchData={refetchData}
      typeName="location"
    />
  )
}

export default function LocationCard({
  id,
  name,
  normalizedName,
  capacity,
  username,
  refetchData
}: LocationCardProps) {
  return (
    <>
      <Card>
        <Card.Body>
          <div className="d-flex flex-row">
            <div>
              <Card.Title>
                <Link href={`/settings/locations/${id}`}>{name}</Link> (
                {normalizedName})
              </Card.Title>
              <Card.Subtitle className="mb-2 text-muted">
                {capacity}
              </Card.Subtitle>
            </div>
            <div className="ms-auto">
              <LocationUpdateModal
                id={id}
                name={name}
                username={username}
                capacity={capacity}
                refetchData={refetchData}
              />
              <LocationDeleteModal
                id={id}
                name={name}
                username={username}
                capacity={capacity}
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
