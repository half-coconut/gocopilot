// import { useState } from "react";
import Button from "../../ui/Button";
import Modal from "../../ui/Modal";

import CreateCabinForm from "./CreateCabinForm";

function AddCabin() {
  return (
    <Modal>
      <Modal.Open opens="cabin-form">
        <Button>Add new interface</Button>
      </Modal.Open>
      <Modal.Window name="cabin-form">
        <CreateCabinForm />
      </Modal.Window>
    </Modal>
  );
}

// function AddCabin() {
//   const [isOpenModel, setIsOpenModel] = useState(false);
//   return (
//     <div>
//       <Button onClick={() => setIsOpenModel((show) => !show)}>
//         Add new cabin
//       </Button>
//       {isOpenModel && (
//         <Modal onClose={() => setIsOpenModel(false)}>
//           <CreateCabinForm onCloseModel={() => setIsOpenModel(false)} />
//         </Modal>
//       )}
//     </div>
//   );
// }

export default AddCabin;
