import { useState } from "react"
import { Modal } from "react-bootstrap"
import CustomButton from "components/CustomButton"
import { type FieldValues, type SubmitHandler, useForm } from "react-hook-form"
import BaseForm from "components/forms/BaseForm"
import { type APIResponse } from "api-wrapper"

interface DeleteModalProps<T> {
  name: string
  handler: () => Promise<APIResponse<T>>
  refetchData: () => Promise<void>
  typeName: string
}

export default function DeleteModal<T extends FieldValues>({
  name,
  handler,
  refetchData,
  typeName
}: DeleteModalProps<T>) {
  const [showDelete, setShowDelete] = useState(false)
  const handleCloseDelete = () => setShowDelete(false)
  const handleShowDelete = () => setShowDelete(true)

  const {
    handleSubmit,
    formState: { errors }
  } = useForm<T>()

  const onSubmit: SubmitHandler<T> = async () => {
    await handler()
    handleCloseDelete()
    await refetchData()
  }

  return (
    <>
      <Modal show={showDelete} onHide={handleCloseDelete}>
        <Modal.Body>
          <Modal.Title>Delete {typeName.toLowerCase()}</Modal.Title>
          <br />
          Are you sure you want to delete &quot;{name}&quot;?
          <br />
          <br />
          <BaseForm
            onSubmit={handleSubmit(onSubmit)}
            onCancelCallback={handleCloseDelete}
            submitBtnText="Delete"
            errors={errors}
          ></BaseForm>
        </Modal.Body>
      </Modal>
      <CustomButton onClick={handleShowDelete}>Delete</CustomButton>
    </>
  )
}
