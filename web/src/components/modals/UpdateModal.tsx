import { useState } from "react"
import { Modal } from "react-bootstrap"
import CustomButton from "components/CustomButton"
import { type FieldValues, type SubmitHandler } from "react-hook-form"
import BaseForm from "components/forms/BaseForm"
import { type IModalProps } from "interfaces/IModalProps"
import { setErrors } from "./helpers"

type UpdateModalProps<T extends FieldValues, Y> = IModalProps<T, Y>

// eslint-disable-next-line max-lines-per-function
export default function UpdateModal<T extends FieldValues, Y>({
  children,
  form,
  handler,
  fetchData,
  typeName
}: UpdateModalProps<T, Y>) {
  const [showUpdate, setShowUpdate] = useState(false)
  const handleCloseUpdate = () => setShowUpdate(false)
  const handleShowUpdate = () => setShowUpdate(true)

  const {
    handleSubmit,
    formState: { dirtyFields, errors },
    setError,
    reset
  } = form

  const onSubmit: SubmitHandler<T> = async (data) => {
    const dataToSubmit = Object.fromEntries(
      Object.keys(dirtyFields).map((key) => [key, data[key]])
    )

    const response = await handler(dataToSubmit as T)
    if (!response.ok) {
      setErrors(response, setError)
    } else {
      handleCloseUpdate()
      reset(data)
      await fetchData()
    }
  }

  const onCancel = () => {
    handleCloseUpdate()
    reset()
  }

  return (
    <>
      <Modal show={showUpdate} onHide={onCancel}>
        <Modal.Body>
          <Modal.Title>Update {typeName.toLowerCase()}</Modal.Title>
          <br />
          <BaseForm
            onSubmit={handleSubmit(onSubmit)}
            errors={errors}
            submitBtnText="Update"
            submitBtnDisabled={
              Object.keys(dirtyFields).length === 0 ||
              Object.keys(errors).length !== 0
            }
            onCancelCallback={onCancel}
          >
            {children}
          </BaseForm>
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
