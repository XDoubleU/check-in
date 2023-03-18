import { useRouter } from "next/router"
import { type FormEvent, useState } from "react"
import { Form, Modal } from "react-bootstrap"
import CustomButton from "@/components/CustomButton"

interface DeleteModalProps {
  name: string
  handler: () => Promise<void>
}

export default function DeleteModal({ name, handler }: DeleteModalProps) {
  const router = useRouter()
  const [showDelete, setShowDelete] = useState(false)
  const handleCloseDelete = () => setShowDelete(false)
  const handleShowDelete = () => setShowDelete(true)

  const handleDelete = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault()

    await handler()

    await router.replace(router.asPath)
    handleCloseDelete()
  }

  return (
    <>
      <Modal show={showDelete} onHide={handleCloseDelete}>
        <Modal.Body>
          <Modal.Title>Delete school</Modal.Title>
          <br />
          Are you sure you want to delete &quot;{name}&quot;?
          <br />
          <br />
          <Form onSubmit={() => handleDelete}>
            <CustomButton
              type="button"
              style={{ float: "left" }}
              onClick={handleCloseDelete}
            >
              Cancel
            </CustomButton>
            <CustomButton type="submit" style={{ float: "right" }}>
              Delete
            </CustomButton>
          </Form>
        </Modal.Body>
      </Modal>
      <CustomButton onClick={handleShowDelete}>Delete</CustomButton>
    </>
  )
}
