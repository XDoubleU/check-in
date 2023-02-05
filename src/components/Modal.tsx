import { ReactNode } from "react"

type ModalProps = {
  children: ReactNode,
  name: string
}

export default function Modal({children, name}: ModalProps){
  return (
    <div className="modal fade" id={`${name}Modal`} tabIndex={-1} data-bs-backdrop="static" aria-labelledby={`${name}ModalLabel`} aria-hidden="true">
        <div className="modal-dialog modal-dialog-scrollable modal-fullscreen">
            <div className="modal-content">
                <div className="modal-body">
                    {children}
                </div>
            </div>
        </div>
    </div>
  )
}