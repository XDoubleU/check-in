import { Card, Form } from "react-bootstrap"
import DeleteModal from "@/components/modals/DeleteModal"
import { deleteSchool, updateSchool } from "my-api-wrapper"
import { type UpdateSchoolDto } from "types-custom"
import { useForm } from "react-hook-form"
import UpdateModal from "../modals/UpdateModal"

interface SchoolCardProps {
  id: number
  name: string
  refetchData: () => Promise<void>
}

function SchoolUpdateModal({ id, name, refetchData }: SchoolCardProps) {
  const form = useForm<UpdateSchoolDto>({
    defaultValues: {
      name: name
    }
  })

  const { register } = form

  const handleUpdate = (data: UpdateSchoolDto) => {
    return updateSchool(id, data)
  }

  return (
    <UpdateModal<UpdateSchoolDto>
      form={form}
      handler={handleUpdate}
      refetchData={refetchData}
      typeName="school"
    >
      <Form.Group className="mb-3">
        <Form.Label>Name</Form.Label>
        <Form.Control
          type="text"
          placeholder="Name"
          {...register("name")}
        ></Form.Control>
      </Form.Group>
    </UpdateModal>
  )
}

function SchoolDeleteModal({ id, name, refetchData }: SchoolCardProps) {
  const handleDelete = () => {
    return deleteSchool(id)
  }

  return (
    <DeleteModal
      name={name}
      handler={handleDelete}
      refetchData={refetchData}
      typeName="school"
    />
  )
}

export default function SchoolCard({ id, name, refetchData }: SchoolCardProps) {
  return (
    <>
      <Card>
        <Card.Body>
          <div className="d-flex flex-row">
            <div>
              <Card.Title>{name}</Card.Title>
            </div>
            <div className="ms-auto">
              <SchoolUpdateModal
                id={id}
                name={name}
                refetchData={refetchData}
              />
              <SchoolDeleteModal
                id={id}
                name={name}
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
