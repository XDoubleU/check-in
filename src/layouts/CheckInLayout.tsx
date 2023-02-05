import Modal from "@/components/Modal"
import BaseLayout from "./BaseLayout"

function SchoolModal(){
  return (
    <Modal name="school">
      Bleep
    </Modal>
  )
}

function CountButton(){
  const count = 1

  if (count <= 0){
    return <button className="btn btn-custom btn-check-in bold">VOLZET</button>
  }

  return (
    <>
      <h2>Nog <span id="count" className="bold">{count}</span> plekken vrij</h2>
      <br/>
      <button id="check-in-btn" className="btn btn-custom btn-check-in bold" data-bs-toggle="modal" data-bs-target="#schoolModal">CHECK-IN</button>
    </>
  )
}

export default function CheckInLayout(){
  return (
    <BaseLayout>
      <div className="container content">
        <SchoolModal />

        <div className="d-flex align-items-center min-vh-80">
          <div className="container text-center">
            <input type="hidden" id="name" value="{{ user.location.get_normalized_name }}"/>

            <h1 className="bold" style={{"fontSize": "5rem"}}>Welkom bij INSERT-NAME!</h1>
            <br/>
            
            <CountButton />
          </div>
        </div>
      </div>

      <br/>
      <br/>

      <footer className="text-center">
        <br/>

        <p>Made with <i className="bi bi-heart-fill" style={{"color": "red"}}></i> by XDoubleU for Brugge Studentenstad</p>
      </footer>
    </BaseLayout>
  )   
}