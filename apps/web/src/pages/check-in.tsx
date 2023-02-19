import Button from "react-bootstrap/Button"
import Modal from "react-bootstrap/Modal"
import { SyntheticEvent, useEffect, useState } from "react"
import styles from "./check-in.module.css"
import { Container, Form } from "react-bootstrap"
import { Location, School } from "types"
import BaseLayout from "@/layouts/BaseLayout"
import CustomButton from "@/components/CustomButton"
import { createCheckIn, getAllSchools, getMyLocation } from "api-wrapper"
import LoadingLayout from "@/layouts/LoadingLayout"

// TODO

export default function CheckIn(){
  const [location, setLocation] = useState<Location>()
  const [schools, setSchools] = useState(new Array<School>())
  const [isDisabled, setDisabled] = useState(false)
  const [showSchools, setShowSchools] = useState(false)
  const handleClose = () => setShowSchools(false)
  const handleShow = () => setShowSchools(true)

  useEffect(() => {
    getMyLocation()
      .then(data => {
        if(data) {
          setLocation(data)
        } else {
          console.log("ERROR")
        }
      })
    
    
  }, [])

  if (!location) {
    return <LoadingLayout/>
  }

  const loadSchools = async () => {
    const paginatedSchools = await getAllSchools(undefined, +Infinity)
    if (paginatedSchools === null) {
      throw new Error()
    }

    setSchools(paginatedSchools.schools)
    handleShow()
  }

  const handleSubmit = async (event: SyntheticEvent<HTMLFormElement>) => {
    event.preventDefault()
  
    const pickedSchool = (event.nativeEvent as SubmitEvent).submitter as HTMLButtonElement
    await createCheckIn(location.id, parseInt(pickedSchool.value))

    handleClose()

    setTimeout( () => {
      setDisabled(true)
    })

    setTimeout(function(){
      setDisabled(false)
    }, 1500)
  }

  return (
    <BaseLayout>
      <Modal show={showSchools} onHide={handleClose} backdrop="static" fullscreen={true} scrollable={true}>
        <div className={styles.modalContent}>
          <Modal.Body>
            <h1 className="bold" style={{"fontSize": "4rem"}}>KIES JE SCHOOL:</h1>
            <h2 style={{"fontSize": "3rem"}}>(scroll voor meer opties)</h2>
            <br/>
            <Form onSubmit={handleSubmit}>
              {
                schools.map((school) => {
                  return (
                    <CustomButton key={school.id} value={school.id} type="submit" className={`${styles.btnSchool} bold`}>
                      {school.name.toUpperCase()}
                    </CustomButton>
                  )
                })
              }
            </Form>
          </Modal.Body>
        </div>
      </Modal>

      <div className="d-flex align-items-center min-vh-80">
        <Container className="text-center">
        <h1 className="bold" style={{"fontSize": "5rem"}}>Welkom bij {location.name}!</h1>
          <br/>
          {
            (location.available <= 0) ? (
              <Button className={`${styles.btnCheckIn} bold text-white`}>VOLZET</Button>
            ) : (
              <>
                <h2>Nog <span id="count" className="bold">{location.available}</span> plekken vrij</h2>
                <br/>
                <Button className={`${styles.btnCheckIn} bold text-white`} onClick={loadSchools} disabled={isDisabled}>
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