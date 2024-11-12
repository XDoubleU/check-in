import { getAllLocations } from "api-wrapper"
import { type Location, type Role } from "api-wrapper/types/apiTypes"
import Charts from "components/charts/Charts"
import { AuthRedirecter } from "contexts/authContext"
import ManagerLayout from "layouts/ManagerLayout"
import { useEffect, useState } from "react"
import { Form } from "react-bootstrap"

export default function Graphs() {
  const redirects = new Map<Role, string>([["default", "/settings"]])

  const [locationIds, setLocationIds] = useState<string[]>([])
  const [locations, setLocations] = useState<Location[]>([])

  useEffect(() => {
    void getAllLocations().then((response) => {
      if (response.data) {
        setLocations(response.data)
      }
    })
  }, [])

  return (
    <AuthRedirecter redirects={redirects}>
      <ManagerLayout title="Graphs">
        <Charts locationIds={locationIds} />
        <Form.Group className="mb-3">
          <Form.Label>Locations</Form.Label>
          <Form.Select
            multiple
            onChange={(e) => {
              setLocationIds(
                Array<HTMLOptionElement>()
                  .slice.call(e.target.selectedOptions)
                  .map((item) => item.value)
              )
            }}
          >
            {locations.map((location) => {
              return (
                <option key={location.id} value={location.id}>
                  {location.name}
                </option>
              )
            })}
          </Form.Select>
        </Form.Group>
      </ManagerLayout>
    </AuthRedirecter>
  )
}
