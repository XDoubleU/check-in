import { Card } from "react-bootstrap"
import DeleteModal from "components/modals/DeleteModal"
import { deleteCheckIn } from "api-wrapper"
import { type ICardProps } from "interfaces/ICardProps"
import { FULL_FORMAT, type CheckIn } from "api-wrapper/types/apiTypes"
import moment from "moment"

type CheckInCardProps = ICardProps<CheckIn>

function CheckInDeleteModal({ data, fetchData }: CheckInCardProps) {
  const handleDelete = () => {
    return deleteCheckIn(data.locationId, data.id)
  }

  return (
    <DeleteModal
      handler={handleDelete}
      fetchData={fetchData}
      typeName="check-In"
    />
  )
}

export default function CheckInCard({
  data,
  user,
  fetchData
}: CheckInCardProps) {
  return (
    <>
      <Card>
        <Card.Body>
          <div className="d-flex flex-row">
            <div>
              <Card.Title>{`${moment.utc(data.createdAt).format(FULL_FORMAT)}`}</Card.Title>
              <Card.Subtitle className="mb-2 text-muted">
                ID: {data.id}
              </Card.Subtitle>
              <Card.Subtitle className="mb-2 text-muted">
                School: {data.schoolName}
              </Card.Subtitle>
            </div>
            {(user?.role === "admin" || user?.role === "manager")
               && (
                <div className="ms-auto">
                  <CheckInDeleteModal
                    data={data}
                    user={user}
                    fetchData={fetchData}
                  />
                </div>
              )}
          </div>
        </Card.Body>
      </Card>
      <br />
    </>
  )
}
