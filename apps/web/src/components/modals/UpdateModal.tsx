import { type ReactElement, useState } from "react"
import { Alert, Form, Modal } from "react-bootstrap"
import CustomButton from "@/components/CustomButton"
import { type FieldValues, type SubmitHandler, useForm } from "react-hook-form"
import type APIResponse from "my-api-wrapper/dist/src/types/apiResponse"

interface UpdateModalProps<T> {
  children: ReactElement | ReactElement[]
  handler: (data: T) => Promise<APIResponse<T>>
}

// eslint-disable-next-line max-lines-per-function
export default function UpdateModal<T extends FieldValues>({
  children,
  handler
}: UpdateModalProps<T>) {
  const [showUpdate, setShowUpdate] = useState(false)
  const handleCloseUpdate = () => setShowUpdate(false)
  const handleShowUpdate = () => setShowUpdate(true)

  const {
    //register,
    handleSubmit,
    setError,
    formState: { errors }
  } = useForm<T>()

  const onSubmit: SubmitHandler<T> = async (data) => {
    const response = await handler(data)
    if (!response.ok) {
      setError("root", {
        message: response.message ?? "Something went wrong"
      })
    } else {
      handleCloseUpdate()
    }
  }

  return (
    <>
      <Modal show={showUpdate} onHide={handleCloseUpdate}>
        <Modal.Body>
          <Modal.Title>Update school</Modal.Title>
          <br />
          <Form onSubmit={handleSubmit(onSubmit)}>
            {children}
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
