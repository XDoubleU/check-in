import { type ReactElement, useState } from "react"
import { Modal } from "react-bootstrap"
import CustomButton from "components/CustomButton"
import {
  type FieldValues,
  type SubmitHandler,
  type UseFormReturn
} from "react-hook-form"
import { type APIResponse } from "api-wrapper"
import BaseForm from "components/forms/BaseForm"

interface UpdateModalProps<T extends FieldValues, Y> {
  children: ReactElement | ReactElement[]
  form: UseFormReturn<T>
  handler: (data: T) => Promise<APIResponse<Y>>
  refetchData: () => Promise<void>
  typeName: string
}

// eslint-disable-next-line max-lines-per-function
export default function UpdateModal<T extends FieldValues, Y>({
  children,
  form,
  handler,
  refetchData,
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
      setError("root", {
        message: response.message ?? "Something went wrong"
      })
    } else {
      handleCloseUpdate()
      reset(data)
      await refetchData()
    }
  }

  return (
    <>
      <Modal show={showUpdate} onHide={handleCloseUpdate}>
        <Modal.Body>
          <Modal.Title>Update {typeName.toLowerCase()}</Modal.Title>
          <br />
          <BaseForm
            onSubmit={handleSubmit(onSubmit)}
            errors={errors}
            submitBtnText="Update"
            submitBtnDisabled={Object.keys(dirtyFields).length === 0}
            onCancelCallback={handleCloseUpdate}
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
