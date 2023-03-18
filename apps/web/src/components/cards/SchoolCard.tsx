import { useState } from "react"
import { Card, Form } from "react-bootstrap"
import UpdateModal from "@/components/modals/UpdateModal"
import DeleteModal from "@/components/modals/DeleteModal"
import { deleteSchool, updateSchool } from "my-api-wrapper"

interface SchoolCardProps {
  id: number
  name: string
}

function SchoolUpdateModal({ id, name }: SchoolCardProps) {
  const [updateInfo, setUpdateInfo] = useState({
    id: id,
    name: name
  })

  const handleUpdate = async () => {
    await updateSchool(id, updateInfo.name)
  }

  return (
    <UpdateModal handler={handleUpdate}>
      <Form.Group className="mb-3">
        <Form.Label>Name</Form.Label>
        <Form.Control
          type="text"
          placeholder="Name"
          value={updateInfo.name}
          onChange={({ target }) =>
            setUpdateInfo({ ...updateInfo, name: target.value })
          }
        ></Form.Control>
      </Form.Group>
    </UpdateModal>
  )
}

export default function SchoolCard({ id, name }: SchoolCardProps) {
  const handleDelete = async () => {
    await deleteSchool(id)
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
              <SchoolUpdateModal id={id} name={name} />
              <DeleteModal name={name} handler={handleDelete} />
            </div>
          </div>
        </Card.Body>
      </Card>
      <br />
    </>
  )
}
