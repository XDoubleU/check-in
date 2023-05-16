import Button from "react-bootstrap/Button"
import Modal from "react-bootstrap/Modal"
import { useCallback, useEffect, useState } from "react"
import styles from "./index.module.css"
import { Container, Form } from "react-bootstrap"
import {
  type LocationUpdateEventDto,
  type CreateCheckInDto,
  type School,
  type Location
} from "types-custom"
import BaseLayout from "layouts/BaseLayout"
import CustomButton from "components/CustomButton"
import {
  checkinsWebsocket,
  createCheckIn,
  getAllSchoolsSortedForLocation,
  getMyLocation
} from "api-wrapper"
import { type SubmitHandler, useForm } from "react-hook-form"
import LoadingLayout from "layouts/LoadingLayout"
import * as Sentry from "@sentry/nextjs"

// eslint-disable-next-line max-lines-per-function
export default function CheckIn() {
  const [location, setLocation] = useState<Location | undefined>()
  const [available, setAvailable] = useState(0)
  const [schools, setSchools] = useState(new Array<School>())
  const [isDisabled, setDisabled] = useState(false)

  const [showSchools, setShowSchools] = useState(false)
  const handleClose = () => setShowSchools(false)
  const handleShow = () => setShowSchools(true)

  const { handleSubmit } = useForm<CreateCheckInDto>()

  const connectWebSocket = useCallback((apiLocation: Location): WebSocket => {
    let webSocket = checkinsWebsocket(apiLocation)

    webSocket.onmessage = (event): void => {
      const locationUpdateEvent = JSON.parse(
        event.data as string
      ) as LocationUpdateEventDto

      setAvailable(locationUpdateEvent.available)
    }

    webSocket.onclose = (): void => {
      setTimeout(() => {
        webSocket = connectWebSocket(apiLocation)
      })
    }

    return webSocket
  }, [])

  useEffect(() => {
    void getMyLocation()
      .then((response) => response.data)
      .then((apiLocation) => {
        if (!apiLocation) return

        setLocation(apiLocation)
        setAvailable(apiLocation.available)

        const webSocket = connectWebSocket(apiLocation)

        return () => {
          if (webSocket.readyState === 1) {
            webSocket.close()
          }
        }
      })
  }, [connectWebSocket])

  const loadSchools = async () => {
    const response = await getAllSchoolsSortedForLocation()

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

  if (!location) {
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
          <Modal.Body style={{ maxHeight: "100vh" }}>
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
            Welkom bij {location.name}!
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
