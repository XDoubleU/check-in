import Button from "react-bootstrap/Button"
import Modal from "react-bootstrap/Modal"
import { SyntheticEvent, useState } from "react"
import styles from "./check-in.module.css"
import { Container, Form } from "react-bootstrap"
import { GetServerSidePropsContext } from "next"
import BaseLayout from "@/layouts/BaseLayout"
import { getServerSession } from "next-auth"
import { authOptions } from "@/pages/api/auth/[...nextauth]"
import { Location, School } from "@prisma/client"
import CustomButton from "@/components/CustomButton"

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

  const checkIn = async () => {
    const response = await fetch("/api/schools")
    setSchools(await response.json())
    handleShow()
  }

  const handleSubmit = async (event: SyntheticEvent<HTMLFormElement>) => {
    event.preventDefault()
  
    const pickedSchool = (event.nativeEvent as SubmitEvent).submitter as HTMLButtonElement
    await fetch("/api/check-ins", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        locationId: location.id,
        schoolId: pickedSchool.value
      })
    })

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

                  return 
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
                <Button className={`${styles.btnCheckIn} bold text-white`} onClick={checkIn} disabled={isDisabled}>
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

  const response = await fetch(`${process.env.NEXTAUTH_URL}/api/locations/${session.user.locationId}`)
  const location = await response.json()

  const response2 = await fetch(`${process.env.NEXTAUTH_URL}/api/check-ins`)
  const checkInsToday = await response2.json()
  const available = location.capacity - checkInsToday

  return {
    props: {
      location,
      available
    }
  }
}