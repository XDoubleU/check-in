import { type State } from "api-wrapper/types/apiTypes"
import { Alert } from "react-bootstrap"

interface StateAlertProps {
  state: State | undefined
}

// eslint-disable-next-line sonarjs/function-return-type
export default function StateAlert({ state }: Readonly<StateAlertProps>) {
  return (
    (state && state.isMaintenance && (
      <Alert variant="danger">
        The Check-In is currently under maintenance. Changes you make might not
        be saved.
      </Alert>
    )) ??
    (state && !state.isDatabaseActive && (
      <Alert variant="danger">
        The Check-In is currently experiencing some issues. We&apos;re looking
        into this and will be back up soon.
      </Alert>
    ))
  )
}
