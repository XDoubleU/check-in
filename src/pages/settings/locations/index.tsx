import { prisma } from "@/common/prisma"
import CustomButton from "@/components/CustomButton"
import AdminLayout from "@/layouts/AdminLayout"
import LoadingLayout from "@/layouts/LoadingLayout"
import { Location } from "@prisma/client"
import { useSession } from "next-auth/react"
import { useState } from "react"
import { Card, Modal } from "react-bootstrap"

type LocationListProps = {
  locations: Location[]
}

export default function LocationList({locations}: LocationListProps) {
  const [showAdd, setShowAdd] = useState(false)
  const handleCloseAdd = () => setShowAdd(false)
  const handleShowAdd = () => setShowAdd(true)

  const {data, status} = useSession({
    required: true
  })

  if (status == "loading") {
    return <LoadingLayout/>
  }

  return (
    <AdminLayout title="Locations" isAdmin={data.user.isAdmin}>
      <h1>Locations</h1>
      <br/>
      
      <Modal show={showAdd} onHide={handleCloseAdd}>
        <Modal.Body>
          Bleep
        </Modal.Body>
      </Modal>

      <div className="col-2">
        <CustomButton onClick={handleShowAdd}>
          Add
        </CustomButton>
      </div>

      <br/>

      {
        locations.length == 0 ? "Nothing to see here." : ""
      }

      {
        locations.map((location) => {
          return <Card id={location.id} key={location.id} title={location.name} subtitle={location.capacity.toString()} />
        })
      }

      
    </AdminLayout>
  )  
}

export async function getServerSideProps() {
  const locations = await prisma.location.findMany()
  return {
    props: {
      locations
    }
  }
}