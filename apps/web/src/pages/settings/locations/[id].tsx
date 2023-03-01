import { Col } from "react-bootstrap"
import AdminLayout from "@/layouts/AdminLayout"
import { LocationUpdateModal } from "@/components/cards/LocationCard"
import CustomButton from "@/components/CustomButton"
import { Location } from "types-custom"
import { getLocation } from "my-api-wrapper"
import { useEffect, useState } from "react"
import Router from "next/router"
import LoadingLayout from "@/layouts/LoadingLayout"

export default function LocationDetail() {
  const [location, updateLocation] = useState<Location>()

  useEffect(() => {
    void getLocation(Router.query.id as string)
      .then((data) => {
        if (data) {
          updateLocation(data)
        } else {
          console.log("ERROR")
        }
      })
  }, [])

  if (location === undefined) {
    return <LoadingLayout/>
  }

  return (
    <AdminLayout title={location.name}>
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