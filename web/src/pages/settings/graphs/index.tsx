import { type Role } from "api-wrapper/types/apiTypes"
import Charts from "components/charts/Charts"
import { AuthRedirecter } from "contexts/authContext"
import { useState } from "react"

export default function Graphs() {
  const [locationIds] = useState<string[]>([])
  //const [locationIds, setLocationIds] = useState<string[]>([])

  const redirects = new Map<Role, string>([["default", "/settings"]])

  return (
    <AuthRedirecter redirects={redirects}>
      <Charts locationIds={locationIds} />
    </AuthRedirecter>
  )
}