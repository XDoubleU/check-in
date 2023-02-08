import Button from "react-bootstrap/Button"
import Modal from "react-bootstrap/Modal"
import BaseLayout from "./BaseLayout"
import { useState } from "react"
import styles from "./CheckInLayout.module.css"
import { Container } from "react-bootstrap"

export default function CheckInLayout(){
  const [showSchools, setShowSchools] = useState(false)
  const handleClose = () => setShowSchools(false)
  const handleShow = () => setShowSchools(true)

  const count = 1

  return (
    <BaseLayout>
      <Modal show={showSchools} onHide={handleClose} backdrop="static" fullscreen={true}>
        <Modal.Body>
          Bleep
        </Modal.Body>
      </Modal>

      <div className="d-flex align-items-center min-vh-80">
        <Container className="text-center">
        <h1 className="bold" style={{"fontSize": "5rem"}}>Welkom bij INSERT-NAME!</h1>
          <br/>
          
          {
            (count <= 0) ? (
              <Button className={`${styles.btnCheckIn} bold text-white`}>VOLZET</Button>
            ) : (
              <>
                <h2>Nog <span id="count" className="bold">{count}</span> plekken vrij</h2>
                <br/>
                <Button className={`${styles.btnCheckIn} bold text-white`} onClick={handleShow}>
                  CHECK-IN
                </Button>
              </>
            )
          }
        </Container>
      </div>
    </BaseLayout>
  )   
}