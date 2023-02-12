import Button from "react-bootstrap/Button"
import Modal from "react-bootstrap/Modal"
import { SyntheticEvent, useState } from "react"
import styles from "./check-in.module.css"
import { Container, Form } from "react-bootstrap"
import { GetServerSidePropsContext } from "next"
import { getServerSession } from "next-auth"
import { authOptions } from "@/pages/api/auth/[...nextauth]"
import { Location, School } from "types"
import BaseLayout from "@/layouts/BaseLayout"
import CustomButton from "@/components/CustomButton"
import { createCheckIn, getAllSchools, getLocation } from "api-wrapper"

type CheckInProps = {
  location: Location,
  available: number
}

export default function CheckIn({ location, available }: CheckInProps){
  const [count, setCount] = useState(available)
  const [schools, setSchools] = useState(new Array<School>())
  const [isDisabled, setDisabled] = useState(false)
  const [showSchools, setShowSchools] = useState(false)
  const handleClose = () => setShowSchools(false)
  const handleShow = () => setShowSchools(true)

  const loadSchools = async () => {
    const paginatedSchools = await getAllSchools(undefined, +Infinity)
    setSchools(paginatedSchools.schools)
    handleShow()
  }

  const handleSubmit = async (event: SyntheticEvent<HTMLFormElement>) => {
    event.preventDefault()
  
    const pickedSchool = (event.nativeEvent as SubmitEvent).submitter as HTMLButtonElement
    await createCheckIn(location.id, parseInt(pickedSchool.value))

    setCount(count - 1)
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
            (count <= 0) ? (
              <Button className={`${styles.btnCheckIn} bold text-white`}>VOLZET</Button>
            ) : (
              <>
                <h2>Nog <span id="count" className="bold">{count}</span> plekken vrij</h2>
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

export async function getServerSideProps(context: GetServerSidePropsContext) {
  const session = await getServerSession(context.req, context.res, authOptions)

  if (!session) {
    return {
      redirect: {
        destination: "/",
        permanent: false,
      },
    }
  }

  const location = await getLocation(session.user.locationId as string)
  //TODO: websocket
  const available = 3

  return {
    props: {
      location,
      available
    }
  }
}