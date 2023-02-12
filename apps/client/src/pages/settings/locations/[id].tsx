import { GetServerSidePropsContext } from "next"
import { useSession } from "next-auth/react"
import { Col } from "react-bootstrap"
import LoadingLayout from "../../../layouts/LoadingLayout"
import AdminLayout from "../../../layouts/AdminLayout"
import { LocationUpdateModal } from "../../../components/cards/LocationCard"
import CustomButton from "../../../components/CustomButton"
import { Location } from "database"
import { LocationWithUser } from "../../../types/prisma"

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
          username={(location as LocationWithUser).user.username}
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
  const response = await fetch(`${process.env.NEXTAUTH_URL}/api/locations/${context.query.id}`)
  
  if (response.status == 404) {
    return {
      notFound: true
    }
  }
  
  const location = await response.json()

  return {
    props: {
      location
    }
  }
}