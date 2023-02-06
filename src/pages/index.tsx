import Spinner from "@/components/Spinner"
import CheckInLayout from "@/layouts/CheckInLayout"
import { useSession } from "next-auth/react"

export default function CheckIn() {
  const {data, status} = useSession({
    required: true
  })

  if (status == "loading") {
    return <Spinner/>
  }

  console.log(data.user)
  return <CheckInLayout/>
}
