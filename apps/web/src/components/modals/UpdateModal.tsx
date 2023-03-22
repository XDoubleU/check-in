import { type ReactElement, useState } from "react"
import { Modal } from "react-bootstrap"
import CustomButton from "@/components/CustomButton"
import {
  type FieldValues,
  type SubmitHandler,
  type UseFormReturn
} from "react-hook-form"
import type APIResponse from "my-api-wrapper/dist/src/types/apiResponse"
import BaseForm from "../forms/BaseForm"

interface UpdateModalProps<T extends FieldValues> {
  children: ReactElement | ReactElement[]
  form: UseFormReturn<T>
  handler: (data: T) => Promise<APIResponse<T>>
  refetchData: () => Promise<void>
  typeName: string
}

// eslint-disable-next-line max-lines-per-function
export default function UpdateModal<T extends FieldValues>({
  children,
  form,
  handler,
  refetchData,
  typeName
}: UpdateModalProps<T>) {
  const [showUpdate, setShowUpdate] = useState(false)
  const handleCloseUpdate = () => setShowUpdate(false)
  const handleShowUpdate = () => setShowUpdate(true)

  const { dirtyFields } = form.formState
  const onSubmit: SubmitHandler<T> = async (data) => {
    const dataToSubmit = Object.fromEntries(
      Object.keys(dirtyFields).map((key) => [key, data[key]])
    )

    const response = await handler(dataToSubmit as T)
    if (!response.ok) {
      form.setError("root", {
        message: response.message ?? "Something went wrong"
      })
    } else {
      handleCloseUpdate()
      form.reset(data)
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
            onSubmit={form.handleSubmit(onSubmit)}
            errors={form.formState.errors}
            submitBtnText="Update"
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
