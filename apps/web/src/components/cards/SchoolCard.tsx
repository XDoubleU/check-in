import { Card, Form } from "react-bootstrap"
import DeleteModal from "@/components/modals/DeleteModal"
import { deleteSchool, updateSchool } from "my-api-wrapper"
import { type School, type UpdateSchoolDto } from "types-custom"
import { useForm } from "react-hook-form"
import UpdateModal from "../modals/UpdateModal"

interface SchoolCardProps {
  school: School
  refetchData: () => Promise<void>
}

function SchoolUpdateModal({ school, refetchData }: SchoolCardProps) {
  const form = useForm<UpdateSchoolDto>({
    defaultValues: {
      name: school.name
    }
  })

  const { register } = form

  const handleUpdate = (data: UpdateSchoolDto) => {
    return updateSchool(school.id, data)
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

function SchoolDeleteModal({ school, refetchData }: SchoolCardProps) {
  const handleDelete = () => {
    return deleteSchool(school.id)
  }

  return (
    <DeleteModal
      name={school.name}
      handler={handleDelete}
      refetchData={refetchData}
      typeName="school"
    />
  )
}

export default function SchoolCard({ school, refetchData }: SchoolCardProps) {
  return (
    <>
      <Card>
        <Card.Body>
          <div className="d-flex flex-row">
            <div>
              <Card.Title>{school.name}</Card.Title>
            </div>
            <div className="ms-auto">
              <SchoolUpdateModal school={school} refetchData={refetchData} />
              <SchoolDeleteModal school={school} refetchData={refetchData} />
            </div>
          </div>
        </Card.Body>
      </Card>
      <br />
    </>
  )
}
