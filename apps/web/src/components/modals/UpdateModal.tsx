import { useRouter } from "next/router"
import { FormEvent, ReactElement, useState } from "react"
import { Form, Modal } from "react-bootstrap"
import CustomButton from "@/components/CustomButton"

interface UpdateModalProps {
  children: ReactElement | ReactElement[],
  handler: () => Promise<void>
}

export default function UpdateModal({children, handler}: UpdateModalProps) {
  const router = useRouter()
  const [showUpdate, setShowUpdate] = useState(false)
  const handleCloseUpdate = () => setShowUpdate(false)
  const handleShowUpdate = () => setShowUpdate(true)

  const handleUpdate = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault()

    await handler()
    
    await router.replace(router.asPath)
    handleCloseUpdate()
  }

  return (
    <>
      <Modal show={showUpdate} onHide={handleCloseUpdate}>
        <Modal.Body>
          <Modal.Title>Update school</Modal.Title>
          <br/>
          <Form onSubmit={() => handleUpdate}>
            {children}
            <br/>
            <CustomButton type="button" style={{"float": "left"}} onClick={handleCloseUpdate}>Cancel</CustomButton>
            <CustomButton type="submit" style={{"float": "right"}}>Update</CustomButton>
          </Form>
        </Modal.Body>
      </Modal>
      <CustomButton onClick={handleShowUpdate} style={{"marginRight":"0.25em"}}>Update</CustomButton>
    </>
  )
}