import { Card } from "react-bootstrap"
import DeleteModal from "components/modals/DeleteModal"
import { deleteUser, updateUser } from "api-wrapper"
import { useForm } from "react-hook-form"
import UpdateModal from "components/modals/UpdateModal"
import { type ICardProps } from "interfaces/ICardProps"
import { type User, type UpdateUserDto } from "api-wrapper/types/apiTypes"
import UserInputs from "components/forms/UserInputs"

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
      <UserInputs
        required={false}
        register={register}
        watch={watch}
        errors={errors}
      />
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
