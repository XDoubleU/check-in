import { FormEvent, useState } from "react"
import { Card, Form, Modal } from "react-bootstrap"
import CustomButton from "./CustomButton"
import { useRouter } from "next/router"

type SchoolCardProps = {
  title: string,
  subtitle?: string,
  id: number
}

type ModalProps = {
  id: number,
  name: string
}

function UpdateModal({id, name}: ModalProps) {
  const router = useRouter()
  const [updateInfo, setUpdateInfo] = useState({name: name})
  const [showUpdate, setShowUpdate] = useState(false)
  const handleCloseUpdate = () => setShowUpdate(false)
  const handleShowUpdate = () => setShowUpdate(true)

  const handleUpdate = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault()

    const data = {
      id: id,
      name: updateInfo.name
    }

    const response = await fetch("/api/schools", {
      method: "PUT",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(data)
    })

    if (response.status < 300) {
      router.replace(router.asPath)
      handleCloseUpdate()
    }
  }

  return (
    <>
      <Modal show={showUpdate} onHide={handleCloseUpdate}>
        <Modal.Body>
          <Modal.Title>Update school</Modal.Title>
          <br/>
          <Form onSubmit={handleUpdate}>
            <Form.Group className="mb-3">
              <Form.Label>Name</Form.Label>
              <Form.Control type="text" placeholder="Name" value={updateInfo.name} onChange={({ target}) => setUpdateInfo({ ...updateInfo, name: target.value })}></Form.Control>
            </Form.Group>
            <CustomButton type="button" style={{"float": "left"}}>Cancel</CustomButton>
            <CustomButton type="submit" style={{"float": "right"}}>Update</CustomButton>
          </Form>
        </Modal.Body>
      </Modal>
      <CustomButton onClick={handleShowUpdate} style={{"marginRight":"0.25em"}}>Update</CustomButton>
    </>
  )
}

function DeleteModal({id, name}: ModalProps) {
  const router = useRouter()
  const [showDelete, setShowDelete] = useState(false)
  const handleCloseDelete = () => setShowDelete(false)
  const handleShowDelete = () => setShowDelete(true)

  const handleDelete = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault()

    const data = {
      id: id
    }

    const response = await fetch("/api/schools", {
      method: "DELETE",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(data)
    })

    if (response.status < 300) {
      router.replace(router.asPath)
      handleCloseDelete()
    }
  }

  return (
    <>
      <Modal show={showDelete} onHide={handleCloseDelete}>
        <Modal.Body>
          <Modal.Title>Delete school</Modal.Title>
          <br/>
          Are you sure you want to delete &quot;{name}&quot;?
          <br/>
          <br/>
          <Form onSubmit={handleDelete}>
            <CustomButton type="button" style={{"float": "left"}}>Cancel</CustomButton>
            <CustomButton type="submit" style={{"float": "right"}}>Delete</CustomButton>
          </Form>
        </Modal.Body>
      </Modal>
      <CustomButton onClick={handleShowDelete}>Delete</CustomButton>
    </>
  )
}

export default function SchoolCard({title, subtitle, id}: SchoolCardProps) {
  return (
    <>
      <Card>
        <Card.Body>
          <div className="d-flex flex-row">
            <div>
              <Card.Title>{title}</Card.Title>
              <Card.Subtitle className="mb-2 text-muted">{subtitle}</Card.Subtitle>
            </div>
            <div className="ms-auto">
              <UpdateModal id={id} name={title} />
              <DeleteModal id={id} name={title} />
            </div>
          </div>
        </Card.Body>
      </Card>
      <br/>
    </>
  )
}