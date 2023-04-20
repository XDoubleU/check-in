import { Card } from "react-bootstrap"
import DeleteModal from "components/modals/DeleteModal"
import { deleteSchool, updateSchool } from "api-wrapper"
import { type School, type UpdateSchoolDto } from "types-custom"
import { useForm } from "react-hook-form"
import UpdateModal from "components/modals/UpdateModal"
import FormInput from "components/forms/FormInput"
import { type ICardProps } from "interfaces/ICardProps"

type SchoolCardProps = ICardProps<School>

function SchoolUpdateModal({ data, refetchData }: SchoolCardProps) {
  const form = useForm<UpdateSchoolDto>({
    defaultValues: {
      name: data.name
    }
  })

  const { register } = form

  const handleUpdate = (updateData: UpdateSchoolDto) => {
    return updateSchool(data.id, updateData)
  }

  return (
    <UpdateModal<UpdateSchoolDto, School>
      form={form}
      handler={handleUpdate}
      refetchData={refetchData}
      typeName="school"
    >
      <FormInput
        label="Name"
        type="text"
        placeholder="Name"
        register={register("name")}
      />
    </UpdateModal>
  )
}

function SchoolDeleteModal({ data, refetchData }: SchoolCardProps) {
  const handleDelete = () => {
    return deleteSchool(data.id)
  }

  return (
    <DeleteModal
      name={data.name}
      handler={handleDelete}
      refetchData={refetchData}
      typeName="school"
    />
  )
}

export default function SchoolCard({ data, refetchData }: SchoolCardProps) {
  return (
    <>
      <Card>
        <Card.Body>
          <div className="d-flex flex-row">
            <div>
              <Card.Title>{data.name}</Card.Title>
            </div>
            <div className="ms-auto">
              <SchoolUpdateModal data={data} refetchData={refetchData} />
              <SchoolDeleteModal data={data} refetchData={refetchData} />
            </div>
          </div>
        </Card.Body>
      </Card>
      <br />
    </>
  )
}
