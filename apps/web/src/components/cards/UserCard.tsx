import { Card, Form } from "react-bootstrap"
import DeleteModal from "components/modals/DeleteModal"
import { deleteUser, updateUser } from "api-wrapper"
import { type UpdateUserDto, type User } from "types-custom"
import { useForm } from "react-hook-form"
import UpdateModal from "components/modals/UpdateModal"
import FormInput from "components/forms/FormInput"
import { type ICardProps } from "interfaces/ICardProps"

type UserUpdateForm = UpdateUserDto & { repeatPassword?: string }

type UserCardProps = ICardProps<User>

// eslint-disable-next-line max-lines-per-function
function UserUpdateModal({ data, fetchData }: UserCardProps) {
  const form = useForm<UserUpdateForm>({
    defaultValues: {
      username: data.username
    }
  })

  const {
    register,
    watch,
    formState: { errors }
  } = form

  const handleUpdate = (updateData: UpdateUserDto) => {
    return updateUser(data.id, updateData)
  }

  return (
    <UpdateModal<UpdateUserDto, User>
      form={form}
      handler={handleUpdate}
      fetchData={fetchData}
      typeName="user"
    >
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
      {/* jscpd:ignore-start */}
      <Form.Group className="mb-3">
        <Form.Label>Repeat password</Form.Label>
        <Form.Control
          type="password"
          placeholder="Repeat password"
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
      {/* jscpd:ignore-end */}
    </UpdateModal>
  )
}

function UserDeleteModal({ data, fetchData }: UserCardProps) {
  const handleDelete = () => {
    return deleteUser(data.id)
  }

  return (
    <DeleteModal
      name={data.username}
      handler={handleDelete}
      fetchData={fetchData}
      typeName="user"
    />
  )
}

export default function UserCard({ data, fetchData }: UserCardProps) {
  return (
    <>
      <Card>
        <Card.Body>
          <div className="d-flex flex-row">
            <div>
              <Card.Title>{data.username}</Card.Title>
            </div>
            <div className="ms-auto">
              <UserUpdateModal data={data} fetchData={fetchData} />
              <UserDeleteModal data={data} fetchData={fetchData} />
            </div>
          </div>
        </Card.Body>
      </Card>
      <br />
    </>
  )
}
