import { GetServerSidePropsContext } from "next"
import { useSession } from "next-auth/react"
import { Col } from "react-bootstrap"
import LoadingLayout from "@/layouts/LoadingLayout"
import AdminLayout from "@/layouts/AdminLayout"
import { LocationUpdateModal } from "@/components/cards/LocationCard"
import CustomButton from "@/components/CustomButton"
import { Location } from "types"
import { getLocation } from "api-wrapper"

type LocationDetailProps = {
  location: Location
}

export default function LocationDetail({location}: LocationDetailProps) {
  const {data, status} = useSession({
    required: true
  })

  if (status == "loading") {
    return <LoadingLayout/>
  }

  return (
    <AdminLayout title={location.name} user={data.user}>
      <Col size={2}>
        <LocationUpdateModal
          id={location.id}
          name={location.name}
          username={location.user.username}
          capacity={location.capacity}
        />
      </Col>
      <br/>

      <Col size={2}>
        <CustomButton>Download CSV (TODO)</CustomButton>
      </Col>
      <br/>

      Still needs a chart :)

    </AdminLayout>
  )  
}

export async function getServerSideProps(context: GetServerSidePropsContext) {
  const location = getLocation(context.query.id as string)

  return {
    props: {
      location
    }
  }
}