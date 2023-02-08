import { prisma } from "@/common/prisma"
import CustomButton from "@/components/CustomButton"
import CustomPagination from "@/components/CustomPagination"
import SchoolCard from "@/components/SchoolCard"
import AdminLayout from "@/layouts/AdminLayout"
import LoadingLayout from "@/layouts/LoadingLayout"
import { School } from "@prisma/client"
import { useSession } from "next-auth/react"
import { useRouter } from "next/router"
import { FormEvent, MouseEvent, useEffect, useState } from "react"
import { Col, Form, Modal } from "react-bootstrap"

const PAGE_SIZE = 4


type SchoolListProps = {
  totalSchools: number
}

export default function SchoolList({totalSchools}: SchoolListProps) {
  const router = useRouter()

  const [schools, setSchools] = useState(new Array<School>())
  const [currentPage, setCurrentPage] = useState(1)
  const [addInfo, setAddInfo] = useState({ name: "" })
  const [showAdd, setShowAdd] = useState(false)
  const handleCloseAdd = () => setShowAdd(false)
  const handleShowAdd = () => setShowAdd(true)

  const {data, status} = useSession({
    required: true
  })
  
  useEffect(() => {
    fetch("/api/schools")
      .then((res) => res.json())
      .then((data) => setSchools(data))
  }, [])

  const changePage = async (event: MouseEvent<HTMLElement>) => {
    const clickedPage = event.target as HTMLElement
    const value = parseInt(clickedPage.innerHTML.split("<")[0])

    const response = await fetch(`/api/schools?page=${value}&pageSize=${PAGE_SIZE}`)

    setSchools(await response.json())
    setCurrentPage(value)
  }

  const handleAdd = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault()

    const data = {
      ...addInfo
    }

    const response = await fetch("/api/schools", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(data)
    })

    if (response.status < 300) {
      router.replace(router.asPath)
      addInfo.name = ""
      handleCloseAdd()
    }
  }

  if (status == "loading") {
    return <LoadingLayout/>
  }

  return (
    <AdminLayout title="Schools" isAdmin={data.user.isAdmin}>
      <h1>Schools</h1>
      <br/>
      
      <Modal show={showAdd} onHide={handleCloseAdd}>
        <Modal.Body>
        <Modal.Title>Add school</Modal.Title>
        <br/>
          <Form onSubmit={handleAdd}>
            <Form.Group className="mb-3">
              <Form.Label>Name</Form.Label>
              <Form.Control type="text" placeholder="Name" value={addInfo.name} onChange={({ target}) => setAddInfo({ ...addInfo, name: target.value })}></Form.Control>
            </Form.Group>
            <CustomButton type="button" style={{"float": "left"}}>Cancel</CustomButton>
            <CustomButton type="submit" style={{"float": "right"}}>Add</CustomButton>
          </Form>
        </Modal.Body>
      </Modal>

      <Col size={2}>
        <CustomButton onClick={handleShowAdd}>
          Add
        </CustomButton>
      </Col>

      <br/>

      {
        (schools === undefined || schools.length == 0) ? "Nothing to see here." : ""
      }

      {
        schools.map((school) => {
          return <SchoolCard id={school.id} key={school.id} title={school.name} />
        })
      }

      <CustomPagination current={currentPage} total={totalSchools} pageSize={PAGE_SIZE} onClick={changePage} />
      
    </AdminLayout>
  )  
}

export async function getServerSideProps() {
  const totalSchools = await prisma.school.count({
    where: {
      id: {
        not: 1
      }
    }
  })

  return {
    props: {
      totalSchools
    }
  }
}