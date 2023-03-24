import { type ReactElement, useState } from "react"
import { Col, Modal } from "react-bootstrap"
import CustomButton from "@/components/CustomButton"
import {
  type FieldValues,
  type SubmitHandler,
  type UseFormReturn
} from "react-hook-form"
import type APIResponse from "my-api-wrapper/dist/src/types/apiResponse"
import BaseForm from "../forms/BaseForm"

interface CreateModalProps<T extends FieldValues, Y> {
  children: ReactElement | ReactElement[]
  form: UseFormReturn<T>
  handler: (data: T) => Promise<APIResponse<Y>>
  refetchData: () => Promise<void>
  typeName: string
}

// eslint-disable-next-line max-lines-per-function
export default function CreateModal<
  T extends FieldValues,
  Y extends FieldValues
>({ children, form, handler, refetchData, typeName }: CreateModalProps<T, Y>) {
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
      setError("root", {
        message: response.message ?? "Something went wrong"
      })
    } else {
      handleCloseCreate()
      reset()
      await refetchData()
    }
  }

  return (
    <>
      <Modal show={showCreate} onHide={handleCloseCreate}>
        <Modal.Body>
          <Modal.Title>Create {typeName.toLowerCase()}</Modal.Title>
          <br />
          <BaseForm
            onSubmit={handleSubmit(onSubmit)}
            errors={errors}
            submitBtnText="Create"
            onCancelCallback={handleCloseCreate}
          >
            {children}
          </BaseForm>
        </Modal.Body>
      </Modal>

      <Col size={2}>
        <CustomButton onClick={handleShowCreate}>Create</CustomButton>
      </Col>
    </>
  )
}
