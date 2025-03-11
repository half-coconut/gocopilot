import styled from "styled-components";
import PropTypes from "prop-types";
// import { formatCurrency } from "../../utils/helpers";

import CreateCabinForm from "./CreateCabinForm";
import { useDeleteCabiin } from "./useDeleteCabin";
import { HiPencil, HiSquare2Stack, HiTrash } from "react-icons/hi2";
import { useCreateCabiin } from "./useCreateCabin";
import Modal from "../../ui/Modal";
import ConfirmDelete from "../../ui/ConfirmDelete";
import Table from "../../ui/Table";
import Menus from "../../ui/Menus";

// const Img = styled.img`
//   display: block;
//   width: 6.4rem;
//   aspect-ratio: 3 / 2;
//   object-fit: cover;
//   object-position: center;
//   transform: scale(1.5) translateX(-7px);
// `;

const Interface = styled.div`
  font-size: 1.6rem;
  font-weight: 600;
  color: var(--color-grey-600);
  font-family: "Sono";
`;

// const Price = styled.div`
//   font-family: "Sono";
//   font-weight: 600;
// `;

const Project = styled.div`
  font-family: "Sono";
  font-weight: 500;
  color: var(--color-green-700);
`;

function CabinRow({ cabin }) {
  const { isCreating, createCabin } = useCreateCabiin();

  const {
    id: interfaceId,
    name,
    project,
    url,
    type,
    creator,
    updater,
    ctime,
    utime,
  } = cabin;

  function handleDuplicate() {
    createCabin({
      name: `Copy of ${name}`,
      project,
      url,
      type,
      creator,
      updater,
      ctime,
      utime,
    });
  }

  const { isDeleting, deleteCabin } = useDeleteCabiin();

  return (
    <Table.Row>
      {/* <Img src={image} /> */}
      <Interface>{name}</Interface>
      <Project>{project} </Project>
      <div>{type} </div>
      <div>{creator} </div>
      <div>{updater} </div>
      <div>{ctime} </div>
      <div>{utime} </div>

      <div>
        <Modal>
          <Menus.Menu>
            <Menus.Toggle id={interfaceId} />

            <Menus.List id={interfaceId}>
              <Menus.Button
                icon={<HiSquare2Stack />}
                onClick={handleDuplicate}
                disabled={isCreating}
              >
                Duplicate
              </Menus.Button>

              <Modal.Open opens="edit">
                <Menus.Button icon={<HiPencil />}>Edit</Menus.Button>
              </Modal.Open>

              <Modal.Open opens="delete">
                <Menus.Button icon={<HiTrash />}>Delete</Menus.Button>
              </Modal.Open>
            </Menus.List>

            <Modal.Window name="edit">
              <CreateCabinForm cabinToEdit={cabin} />
            </Modal.Window>

            <Modal.Window name="delete">
              <ConfirmDelete
                resourceName="cabins"
                disabled={isDeleting}
                onConfirm={() => deleteCabin(interfaceId)}
              />
            </Modal.Window>
          </Menus.Menu>
        </Modal>
      </div>
    </Table.Row>
  );
}

CabinRow.propTypes = {
  cabin: PropTypes.shape({
    id: PropTypes.number.isRequired,
    name: PropTypes.string.isRequired,
    project: PropTypes.string.isRequired,
    url: PropTypes.string.isRequired,
    type: PropTypes.string.isRequired,
    creator: PropTypes.string.isRequired,
    updater: PropTypes.string.isRequired,
    ctime: PropTypes.string.isRequired,
    utime: PropTypes.string.isRequired,
  }).isRequired,
};

export default CabinRow;
