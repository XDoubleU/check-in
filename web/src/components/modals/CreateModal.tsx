import { useState } from "react"
import { Modal } from "react-bootstrap"
import CustomButton from "components/CustomButton"
import { type FieldValues, type SubmitHandler } from "react-hook-form"
import BaseForm from "components/forms/BaseForm"
import { type IModalProps } from "interfaces/IModalProps"
import { setErrors } from "./helpers"

type CreateModalProps<T extends FieldValues, Y> = IModalProps<T, Y>

// eslint-disable-next-line max-lines-per-function
export default function CreateModal<T extends FieldValues, Y>({
  children,
  form,
  handler,
  fetchData,
  typeName
}: CreateModalProps<T, Y>) {
  const [showCreate, setShowCreate] = useState(false)
  const handleCloseCreate = () => setShowCreate(false)
  const handleShowCreate = () => setShowCreate(true)

  const {
    handleSubmit,
    formState: { errors },
    setError,
    reset
  } = form

  const onSubmit: SubmitHandler<T> = async (data) => {
    const response = await handler(data)
    if (!response.ok) {
      setErrors(response, setError)
    } else {
      handleCloseCreate()
      reset()
      await fetchData()
    }
  }

  const onCancel = () => {
    handleCloseCreate()
    reset()
  }

  return (
    <>
      <Modal show={showCreate} onHide={onCancel} animation={false}>
        <Modal.Body>
          <Modal.Title>Create {typeName.toLowerCase()}</Modal.Title>
          <br />
          <BaseForm
            onSubmit={handleSubmit(onSubmit)}
            errors={errors}
            submitBtnText="Create"
            submitBtnDisabled={Object.keys(errors).length !== 0}
            onCancelCallback={onCancel}
          >
            {children}
          </BaseForm>
        </Modal.Body>
      </Modal>

      <CustomButton onClick={handleShowCreate}>Create</CustomButton>
      <br />
    </>
  )
}
