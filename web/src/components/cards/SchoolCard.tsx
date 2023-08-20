import { Card } from "react-bootstrap"
import DeleteModal from "components/modals/DeleteModal"
import { deleteSchool, updateSchool } from "api-wrapper"
import { useForm } from "react-hook-form"
import UpdateModal from "components/modals/UpdateModal"
import FormInput from "components/forms/FormInput"
import { type ICardProps } from "interfaces/ICardProps"
import { type SchoolDto, type School } from "api-wrapper/types/apiTypes"

type SchoolCardProps = ICardProps<School>

function SchoolUpdateModal({ data, fetchData }: SchoolCardProps) {
  const form = useForm<SchoolDto>({
    defaultValues: {
      name: data.name
    }
  })

  const { register, formState: { errors } } = form

  const handleUpdate = (updateData: SchoolDto) => {
    return updateSchool(data.id, updateData)
  }

  return (
    <UpdateModal<SchoolDto, School>
      form={form}
      handler={handleUpdate}
      fetchData={fetchData}
      typeName="school"
    >
      <FormInput
        label="Name"
        type="text"
        placeholder="Name"
        register={register("name")}
        errors={errors.name}
      />
    </UpdateModal>
  )
}

function SchoolDeleteModal({ data, fetchData }: SchoolCardProps) {
  const handleDelete = () => {
    return deleteSchool(data.id)
  }

  return (
    <DeleteModal
      name={data.name}
      handler={handleDelete}
      fetchData={fetchData}
      typeName="school"
    />
  )
}

export default function SchoolCard({ data, fetchData }: SchoolCardProps) {
  return (
    <>
      <Card>
        <Card.Body>
          <div className="d-flex flex-row">
            <div>
              <Card.Title>{data.name}</Card.Title>
            </div>
            {!data.readOnly && (
              <div className="ms-auto">
                <SchoolUpdateModal data={data} fetchData={fetchData} />
                <SchoolDeleteModal data={data} fetchData={fetchData} />
              </div>
            )}
          </div>
        </Card.Body>
      </Card>
      <br />
    </>
  )
}
