import { type CreateUserDto, type User } from "types-custom"
import { createUser, getAllUsersPaged } from "api-wrapper"
import { useState } from "react"
import { useForm } from "react-hook-form"
import CreateModal from "components/modals/CreateModal"
import FormInput from "components/forms/FormInput"
import ListViewLayout, { type List } from "layouts/ListViewLayout"
import { type ICreateModalProps } from "interfaces/ICreateModalProps"
import { Form } from "react-bootstrap"
import UserCard from "components/cards/UserCard"

type CreateUserForm = CreateUserDto & { repeatPassword?: string }

export type CreateUserModalProps = ICreateModalProps<CreateUserForm>

// eslint-disable-next-line max-lines-per-function
function CreateUserModal({ form, fetchData }: CreateUserModalProps) {
  const {
    register,
    watch,
    formState: { errors }
  } = form

  return (
    <CreateModal<CreateUserDto, User>
      form={form}
      handler={createUser}
      refetchData={fetchData}
      typeName="user"
    >
      <FormInput
        label="Username"
        type="text"
        placeholder="Username"
        required
        register={register("username")}
      />
      <FormInput
        label="Password"
        type="password"
        placeholder="Password"
        required
        register={register("password")}
      />
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
  )
}

type UserList = List<User>

export default function UserListView() {
  const [userList, setUserList] = useState<UserList>({
    data: undefined,
    pagination: {
      current: 0,
      total: 0
    }
  })

  const form = useForm<CreateUserDto>()

  return (
    <ListViewLayout
      title="Users"
      form={form}
      list={userList}
      setList={setUserList}
      apiCall={getAllUsersPaged}
      createModal={CreateUserModal}
      card={UserCard}
    />
  )
}
