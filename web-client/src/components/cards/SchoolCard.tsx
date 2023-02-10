import { useState } from "react"
import { Card, Form } from "react-bootstrap"
import DeleteModal from "../modals/DeleteModal"
import UpdateModal from "../modals/UpdateModal"

type SchoolCardProps = {
  id: number,
  name: string
}

function SchoolUpdateModal({id, name}: SchoolCardProps) {
  const [updateInfo, setUpdateInfo] = useState({
    id: id,
    name: name
  })

  return (
    <UpdateModal<SchoolCardProps> updateInfo={updateInfo} endpoint="/api/schools">
      <Form.Group className="mb-3">
        <Form.Label>Name</Form.Label>
        <Form.Control type="text" placeholder="Name" value={updateInfo.name} onChange={({ target}) => setUpdateInfo({ ...updateInfo, name: target.value })}></Form.Control>
      </Form.Group>
    </UpdateModal>
  )
}

export default function SchoolCard({id, name}: SchoolCardProps) {
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
              <DeleteModal id={id} name={name} endpoint="/api/schools" />
            </div>
          </div>
        </Card.Body>
      </Card>
      <br/>
    </>
  )
}