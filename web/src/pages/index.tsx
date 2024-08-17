import Button from "react-bootstrap/Button"
import Modal from "react-bootstrap/Modal"
import { type MouseEventHandler, useCallback, useEffect, useState } from "react"
import styles from "./index.module.css"
import { Container } from "react-bootstrap"
import BaseLayout from "layouts/BaseLayout"
import CustomButton from "components/CustomButton"
import {
  checkinsWebsocket,
  createCheckIn,
  getAllSchoolsSortedForLocation
} from "api-wrapper"
import LoadingLayout from "layouts/LoadingLayout"
import {
  type Role,
  type Location,
  type LocationUpdateEvent,
  type School,
  type State
} from "api-wrapper/types/apiTypes"
import { AuthRedirecter, useAuth } from "contexts/authContext"
import { generateIntegrationScripts } from "../../utils/integration-script"
import StateAlert from "components/StateAlert"

// is executed on compile time
export function getStaticProps() {
  generateIntegrationScripts()

  // otherwise can't compile
  return {props: {}}
}


// eslint-disable-next-line max-lines-per-function
export default function CheckIn() {
  const redirects = new Map<Role, string>([
    ["admin", "/settings"],
    ["manager", "/settings"]
  ])

  const { user } = useAuth()
  const [apiState, setApiState] = useState<State>()
  const [available, setAvailable] = useState(0)
  const [schools, setSchools] = useState(new Array<School>())
  const [isDisabled, setDisabled] = useState(false)

  const [showSchools, setShowSchools] = useState(false)
  const handleClose = () => setShowSchools(false)
  const handleShow = () => setShowSchools(true)

  const connectWebSocket = useCallback((apiLocation: Location): WebSocket => {
    let webSocket = checkinsWebsocket(apiLocation)

    webSocket.onmessage = (event): void => {
      const updateEvent = JSON.parse(
        event.data as string
      ) as LocationUpdateEvent | State

      if ((updateEvent as LocationUpdateEvent).available) {
        setAvailable((updateEvent as LocationUpdateEvent).available)
      }
      else if ((updateEvent as State).isDatabaseActive != undefined) {
        setApiState((updateEvent as State))
        setDisabled((updateEvent as State).isMaintenance || !(updateEvent as State).isDatabaseActive)
      }
    }

    webSocket.onclose = (): void => {
      setTimeout(() => {
        webSocket = connectWebSocket(apiLocation)
      })
    }

    return webSocket
  }, [])

  useEffect(() => {
    let webSocket: WebSocket | undefined
    if (user?.location) {
      setAvailable(user.location.available)
      webSocket = connectWebSocket(user.location)
    }

    return () => {
      if (webSocket && webSocket.readyState === 1) {
        webSocket.close()
      }
    }
  }, [connectWebSocket, user?.location])

  const loadSchools = async () => {
    const response = await getAllSchoolsSortedForLocation()

    setSchools(response.data as School[])
    handleShow()
  }

  const onClick: MouseEventHandler<HTMLButtonElement> = (event) => {
    setDisabled(true)

    void createCheckIn({
      schoolId: parseInt((event.target as HTMLButtonElement).value)
    })

    setAvailable(available - 1)

    handleClose()

    setTimeout(function () {
      setDisabled(false)
    }, 1500)
  }

  return (
    <AuthRedirecter redirects={redirects}>
      {!user?.location ? (
        <LoadingLayout message="User has no location." />
      ) : (
        <BaseLayout>
          <Modal
            show={showSchools}
            onHide={handleClose}
            backdrop="static"
            fullscreen={true}
            scrollable={true}
            animation={false}
          >
            <div className={styles.modalContent}>
              <Modal.Body style={{ maxHeight: "100vh" }}>
                <h1 className="bold" style={{ fontSize: "4rem" }}>
                  KIES JE SCHOOL:
                </h1>
                <h2 style={{ fontSize: "3rem" }}>(scroll voor meer opties)</h2>
                <br />
                {schools.map((school) => {
                  return (
                    <CustomButton
                      key={school.id}
                      value={school.id}
                      onClick={onClick}
                      className={`${styles.btnSchool} bold`}
                    >
                      {school.name.toUpperCase()}
                    </CustomButton>
                  )
                })}
              </Modal.Body>
            </div>
          </Modal>

          <div className="d-flex align-items-center min-vh-80">
            <Container className="text-center">
              <StateAlert state={apiState} />
              <h1 className="bold" style={{ fontSize: "5rem" }}>
                Welkom bij {user?.location?.name}!
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
                    onClick={loadSchools}
                    disabled={isDisabled}
                  >
                    CHECK-IN
                  </Button>
                </>
              )}
            </Container>
          </div>
        </BaseLayout>
      )}
    </AuthRedirecter>
  )
}
