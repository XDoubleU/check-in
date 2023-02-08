import { useRouter } from "next/router"
import { FormEvent, useState } from "react"
import { Form, Modal } from "react-bootstrap"
import CustomButton from "../CustomButton"

type DeleteModalProps = {
  id: string | number,
  name: string,
  endpoint: string
}

export default function DeleteModal({id, name, endpoint}: DeleteModalProps) {
  const router = useRouter()
  const [showDelete, setShowDelete] = useState(false)
  const handleCloseDelete = () => setShowDelete(false)
  const handleShowDelete = () => setShowDelete(true)

  const handleDelete = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault()

    const data = {
      id: id
    }

    const response = await fetch(endpoint, {
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
            <CustomButton type="button" style={{"float": "left"}} onClick={handleCloseDelete}>Cancel</CustomButton>
            <CustomButton type="submit" style={{"float": "right"}}>Delete</CustomButton>
          </Form>
        </Modal.Body>
      </Modal>
      <CustomButton onClick={handleShowDelete}>Delete</CustomButton>
    </>
  )
}