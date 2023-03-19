import { useState } from "react"
import { Alert, Card, Form, Modal } from "react-bootstrap"
import DeleteModal from "@/components/modals/DeleteModal"
import { deleteSchool, updateSchool } from "my-api-wrapper"
import { type UpdateSchoolDto } from "types-custom"
import { type SubmitHandler, useForm } from "react-hook-form"
import CustomButton from "../CustomButton"

interface SchoolCardProps {
  id: number
  name: string
  fetchData: () => Promise<void>
}

// eslint-disable-next-line max-lines-per-function
function SchoolUpdateModal({ id, name, fetchData }: SchoolCardProps) {
  const [showUpdate, setShowUpdate] = useState(false)
  const handleCloseUpdate = () => setShowUpdate(false)
  const handleShowUpdate = () => setShowUpdate(true)

  const {
    register,
    handleSubmit,
    setError,
    formState: { errors }
  } = useForm<UpdateSchoolDto>({
    defaultValues: {
      name: name
    }
  })

  const onSubmit: SubmitHandler<UpdateSchoolDto> = async (data) => {
    const response = await updateSchool(id, data)
    if (!response.ok) {
      setError("root", {
        message: response.message ?? "Something went wrong"
      })
    } else {
      handleCloseUpdate()
      await fetchData()
    }
  }

  return (
    <>
      <Modal show={showUpdate} onHide={handleCloseUpdate}>
        <Modal.Body>
          <Modal.Title>Update school</Modal.Title>
          <br />
          <Form onSubmit={handleSubmit(onSubmit)}>
            <Form.Group className="mb-3">
              <Form.Label>Name</Form.Label>
              <Form.Control
                type="text"
                placeholder="Name"
                {...register("name")}
              ></Form.Control>
            </Form.Group>
            {errors.root && <Alert key="danger">{errors.root.message}</Alert>}
            <br />
            <CustomButton
              type="button"
              style={{ float: "left" }}
              onClick={handleCloseUpdate}
            >
              Cancel
            </CustomButton>
            <CustomButton type="submit" style={{ float: "right" }}>
              Update
            </CustomButton>
          </Form>
        </Modal.Body>
      </Modal>
      <CustomButton
        onClick={handleShowUpdate}
        style={{ marginRight: "0.25em" }}
      >
        Update
      </CustomButton>
    </>
  )
}

export default function SchoolCard({ id, name, fetchData }: SchoolCardProps) {
  const handleDelete = async () => {
    await deleteSchool(id)
    await fetchData()
  }

  return (
    <>
      <Card>
        <Card.Body>
          <div className="d-flex flex-row">
            <div>
              <Card.Title>{name}</Card.Title>
            </div>
            <div className="ms-auto">
              <SchoolUpdateModal id={id} name={name} fetchData={fetchData} />
              <DeleteModal name={name} handler={handleDelete} />
            </div>
          </div>
        </Card.Body>
      </Card>
      <br />
    </>
  )
}
