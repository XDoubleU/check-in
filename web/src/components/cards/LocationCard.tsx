import { Card, Form } from "react-bootstrap"
import Link from "next/link"
import UpdateModal from "components/modals/UpdateModal"
import DeleteModal from "components/modals/DeleteModal"
import { deleteLocation, updateLocation } from "api-wrapper"
import { useForm } from "react-hook-form"
import { type LocationWithUsername } from "pages/settings/locations/index"
import FormInput from "components/forms/FormInput"
import { type ICardProps } from "interfaces/ICardProps"
import {
  type UpdateLocationDto,
  type Location,
  TIME_FORMAT
} from "api-wrapper/types/apiTypes"
import moment from "moment"

type LocationUpdateForm = UpdateLocationDto & { repeatPassword?: string }

type LocationCardProps = ICardProps<LocationWithUsername>

// eslint-disable-next-line max-lines-per-function
export function LocationUpdateModal({ data, fetchData }: LocationCardProps) {
  const form = useForm<LocationUpdateForm>({
    defaultValues: {
      name: data.name,
      capacity: data.capacity,
      username: data.username,
      timeZone: data.timeZone
    }
  })

  const {
    register,
    watch,
    formState: { errors }
  } = form

  const handleUpdate = (updateData: UpdateLocationDto) => {
    return updateLocation(data.id, updateData)
  }

  return (
    <UpdateModal<UpdateLocationDto, Location>
      form={form}
      handler={handleUpdate}
      fetchData={fetchData}
      typeName="location"
    >
      <FormInput
        label="Name"
        type="text"
        placeholder="Name"
        register={register("name")}
      />
      <FormInput
        label="Capacity"
        type="number"
        placeholder={10}
        register={register("capacity")}
      />
      <Form.Group
        className="mb-3"
        hidden={process.env.NEXT_PUBLIC_EDIT_TIME_ZONE !== "true"}
      >
        <Form.Label>Time zone</Form.Label>
        <Form.Select {...register("timeZone")}>
          {Intl.supportedValuesOf("timeZone").map((timeZone) => {
            return (
              <option key={timeZone} value={timeZone}>
                {timeZone}
              </option>
            )
          })}
        </Form.Select>
      </Form.Group>
      <FormInput
        label="Username"
        type="text"
        placeholder="Username"
        register={register("username")}
      />
      <FormInput
        label="Password"
        type="password"
        placeholder="Password"
        register={register("password")}
      />
      <FormInput
        label="Repeat password"
        type="password"
        placeholder="Repeat password"
        register={register("repeatPassword", {
          validate: (val: string | undefined) => {
            if (watch("password") != val) {
              return "Your passwords do no match"
            }
            return undefined
          }
        })}
        errors={errors.repeatPassword}
      />
    </UpdateModal>
  )
}

function LocationDeleteModal({ data, fetchData }: LocationCardProps) {
  const handleDelete = () => {
    return deleteLocation(data.id)
  }

  return (
    <DeleteModal
      name={data.name}
      handler={handleDelete}
      fetchData={fetchData}
      typeName="location"
    />
  )
}

export default function LocationCard({ data, fetchData }: LocationCardProps) {
  return (
    <>
      <Card>
        <Card.Body>
          <div className="d-flex flex-row">
            <div>
              <Card.Title>
                <Link
                  href={{
                    pathname: "/settings/locations/[id]",
                    query: { id: data.id }
                  }}
                >
                  {data.name}
                </Link>{" "}
                ({data.normalizedName})
              </Card.Title>
              <Card.Subtitle className="mb-2 text-muted">
                {data.available} / {data.capacity}
              </Card.Subtitle>
              <Card.Subtitle className="mb-2 text-muted">
                {data.yesterdayFullAt
                  ? `Yesterday full at ${moment
                      .utc(data.yesterdayFullAt)
                      .format(TIME_FORMAT)}`
                  : "Yesterday not full"}
              </Card.Subtitle>
            </div>
            <div className="ms-auto">
              <LocationUpdateModal data={data} fetchData={fetchData} />
              <LocationDeleteModal data={data} fetchData={fetchData} />
            </div>
          </div>
        </Card.Body>
      </Card>
      <br />
    </>
  )
}
