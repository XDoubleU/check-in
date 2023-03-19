import { useState } from "react"
import { Form, Modal } from "react-bootstrap"
import CustomButton from "@/components/CustomButton"
import { type FieldValues, type SubmitHandler, useForm } from "react-hook-form"

interface DeleteModalProps {
  name: string
  handler: () => Promise<void>
}

export default function DeleteModal<T extends FieldValues>({
  name,
  handler
}: DeleteModalProps) {
  const [showDelete, setShowDelete] = useState(false)
  const handleCloseDelete = () => setShowDelete(false)
  const handleShowDelete = () => setShowDelete(true)

  const { handleSubmit } = useForm<T>()

  const onSubmit: SubmitHandler<T> = async () => {
    await handler()
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
          <Form onSubmit={handleSubmit(onSubmit)}>
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
