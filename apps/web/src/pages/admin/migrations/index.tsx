import { applyMigrationsCreate, applyMigrationsDown, applyMigrationsUp, applySeeder } from "my-api-wrapper"
import { Col, Row } from "react-bootstrap"
import CustomButton from "../../../components/CustomButton"
import AdminLayout from "../../../layouts/AdminLayout"

// eslint-disable-next-line @typescript-eslint/naming-convention
async function migrationCreate() {
  const filename = (await applyMigrationsCreate()).data

  document.getElementById("output-create-migrations")!.innerHTML = filename ?? "Error"
}

// eslint-disable-next-line @typescript-eslint/naming-convention
async function migrationUp() {
  const name = (await applyMigrationsUp()).data

  document.getElementById("output-migrations-up")!.innerHTML = name ?? "Error"
}

// eslint-disable-next-line @typescript-eslint/naming-convention
async function migrationDown() {
  const name = (await applyMigrationsDown()).data

  document.getElementById("output-migrations-down")!.innerHTML = name ?? "Error"
}

// eslint-disable-next-line @typescript-eslint/naming-convention
async function seed() {
  const ok = (await applySeeder()).ok

  document.getElementById("output-seed")!.innerHTML = ok ? "ACK" : "NACK"
}

export default function Migrations() {
  return (
    <AdminLayout title="Migrations">
      <Row>
        <Col md={2}>
          <CustomButton onClick={migrationCreate}>Create migrations</CustomButton>
          <br />
          <br />
          <CustomButton onClick={migrationUp}>Migrations up</CustomButton>
          <br />
          <br />
          <CustomButton onClick={migrationDown}>Migrations down</CustomButton>
          <br />
          <br />
          <CustomButton onClick={seed}>Seed</CustomButton>
        </Col>
        <Col md={2}>
          <p id="output-create-migrations">...</p>
          <br />
          <p id="output-migrations-up">...</p>
          <br />
          <p id="output-migrations-down">...</p>
          <br />
          <p id="output-seed">...</p>
        </Col>
      </Row>
    </AdminLayout>
  )
}
