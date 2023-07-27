import SchoolCard from "components/cards/SchoolCard"
import { createSchool, getAllSchoolsPaged } from "api-wrapper"
import { useState } from "react"
import { useForm } from "react-hook-form"
import CreateModal from "components/modals/CreateModal"
import FormInput from "components/forms/FormInput"
import ListViewLayout, { type List } from "layouts/ListViewLayout"
import { type ICreateModalProps } from "interfaces/ICreateModalProps"
import {
  type Role,
  type School,
  type SchoolDto
} from "api-wrapper/types/apiTypes"
import { AuthRedirecter } from "contexts/authContext"

export type CreateSchoolModalProps = ICreateModalProps<SchoolDto>

function CreateSchoolModal({ form, fetchData }: CreateSchoolModalProps) {
  const { register } = form

  return (
    <CreateModal<SchoolDto, School>
      form={form}
      handler={createSchool}
      fetchData={fetchData}
      typeName="school"
    >
      <FormInput
        label="Name"
        type="text"
        placeholder="Name"
        required
        register={register("name")}
      />
    </CreateModal>
  )
}

type SchoolList = List<School>

export default function SchoolListView() {
  const redirects = new Map<Role, string>([["default", "/settings"]])

  const [schoolList, setSchoolList] = useState<SchoolList>({
    data: undefined,
    pagination: {
      current: 0,
      total: 0
    }
  })

  const form = useForm<SchoolDto>()

  return (
    <AuthRedirecter redirects={redirects}>
      <ListViewLayout
        title="Schools"
        form={form}
        list={schoolList}
        setList={setSchoolList}
        apiCall={getAllSchoolsPaged}
        createModal={CreateSchoolModal}
        card={SchoolCard}
      />
    </AuthRedirecter>
  )
}
