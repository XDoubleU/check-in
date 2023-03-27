import SchoolCard from "../../../components/cards/SchoolCard"
import { type CreateSchoolDto, type School } from "types-custom"
import { createSchool, getAllSchoolsPaged } from "my-api-wrapper"
import { useState } from "react"
import { useForm } from "react-hook-form"
import CreateModal from "../../../components/modals/CreateModal"
import FormInput from "../../../components/forms/FormInput"
import ListViewLayout, { type List } from "../../../layouts/ListViewLayout"
import { type ICreateModalProps } from "../../../interfaces/ICreateModalProps"

export type CreateSchoolModalProps = ICreateModalProps<CreateSchoolDto>

function CreateSchoolModal({ form, fetchData }: CreateSchoolModalProps) {
  const { register } = form

  return (
    <CreateModal<CreateSchoolDto, School>
      form={form}
      handler={createSchool}
      refetchData={fetchData}
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
  const [schoolList, setSchoolList] = useState<SchoolList>({
    data: undefined,
    pagination: {
      current: 0,
      total: 0
    }
  })

  const form = useForm<CreateSchoolDto>()

  return (
    <ListViewLayout
      title="Schools"
      form={form}
      list={schoolList}
      setList={setSchoolList}
      apiCall={getAllSchoolsPaged}
      createModal={CreateSchoolModal}
      card={SchoolCard}
    />
  )
}
