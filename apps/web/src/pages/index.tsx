import Button from "react-bootstrap/Button"
import Modal from "react-bootstrap/Modal"
import { useEffect, useState } from "react"
import styles from "./index.module.css"
import { Container, Form } from "react-bootstrap"
import {
  type LocationUpdateEventDto,
  type CreateCheckInDto,
  type School
} from "types-custom"
import BaseLayout from "@/layouts/BaseLayout"
import CustomButton from "@/components/CustomButton"
import {
  checkinsEventSource,
  createCheckIn,
  getAllSchools
} from "my-api-wrapper"
import { useAuth } from "@/contexts"
import { type SubmitHandler, useForm } from "react-hook-form"
import LoadingLayout from "@/layouts/LoadingLayout"

// eslint-disable-next-line max-lines-per-function
export default function CheckIn() {
  const { user } = useAuth()
  const [available, setAvailable] = useState(user?.location?.available ?? 0)
  const [schools, setSchools] = useState(new Array<School>())
  const [isDisabled, setDisabled] = useState(false)
  const [showSchools, setShowSchools] = useState(false)
  const handleClose = () => setShowSchools(false)
  const handleShow = () => setShowSchools(true)
  const { handleSubmit } = useForm<CreateCheckInDto>()

  useEffect(() => {
    if (!user?.location) return

    const eventSource = checkinsEventSource(user.location)

    eventSource.onmessage = (event): void => {
      const locationUpdateEvent = JSON.parse(
        event.data as string
      ) as LocationUpdateEventDto

      setAvailable(locationUpdateEvent.available)
    }

    return () => {
      eventSource.close()
    }
  }, [user?.location])

  const loadSchools = async () => {
    const response = await getAllSchools()
    setSchools(response.data ?? Array<School>())
    handleShow()
  }

  const onSubmit: SubmitHandler<CreateCheckInDto> = async (_, event) => {
    const pickedSchool = (event?.nativeEvent as SubmitEvent)
      .submitter as HTMLButtonElement

    await createCheckIn({
      schoolId: parseInt(pickedSchool.value)
    })

    handleClose()

    setTimeout(() => {
      setDisabled(true)
    })

    setTimeout(function () {
      setDisabled(false)
    }, 1500)
  }

  if (!user?.location) {
    return <LoadingLayout />
  }

  return (
    <BaseLayout>
      <Modal
        show={showSchools}
        onHide={handleClose}
        backdrop="static"
        fullscreen={true}
        scrollable={true}
      >
        <div className={styles.modalContent}>
          <Modal.Body>
            <h1 className="bold" style={{ fontSize: "4rem" }}>
              KIES JE SCHOOL:
            </h1>
            <h2 style={{ fontSize: "3rem" }}>(scroll voor meer opties)</h2>
            <br />
            <Form onSubmit={handleSubmit(onSubmit)}>
              {schools.map((school) => {
                return (
                  <CustomButton
                    key={school.id}
                    value={school.id}
                    type="submit"
                    className={`${styles.btnSchool} bold`}
                  >
                    {school.name.toUpperCase()}
                  </CustomButton>
                )
              })}
            </Form>
          </Modal.Body>
        </div>
      </Modal>

      <div className="d-flex align-items-center min-vh-80">
        <Container className="text-center">
          <h1 className="bold" style={{ fontSize: "5rem" }}>
            Welkom bij {user.location.name}!
          </h1>
          <br />
          {available <= 0 ? (
            <Button className={`${styles.btnCheckIn} bold text-white`}>
              VOLZET
            </Button>
          ) : (
            <>
              <h2>
                Nog{" "}
                <span id="count" className="bold">
                  {available}
                </span>{" "}
                plekken vrij
              </h2>
              <br />
              <Button
                className={`${styles.btnCheckIn} bold text-white`}
                onClick={() => loadSchools()}
                disabled={isDisabled}
              >
                CHECK-IN
              </Button>
            </>
          )}
        </Container>
      </div>
    </BaseLayout>
  )
}
