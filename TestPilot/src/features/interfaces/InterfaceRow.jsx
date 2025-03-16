import styled from "styled-components";
import PropTypes from "prop-types";
// import { formatCurrency } from "../../utils/helpers";

import CreateInterfaceForm from "./CreateInterfaceForm";
import { useDeleteInterface } from "./useDeleteInterface";
import { HiEye, HiPencil, HiSquare2Stack, HiTrash } from "react-icons/hi2";
import { useCreateInterface } from "./useCreateInterface";
import Modal from "../../ui/Modal";
import ConfirmDelete from "../../ui/ConfirmDelete";
import Table from "../../ui/Table";
import Menus from "../../ui/Menus";
import { useNavigate } from "react-router-dom";

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
const Type = styled.div`
  font-family: "Sono";
  font-weight: 500;
  color: var(--color-grey-600);
`;

function InterfaceRow({ interfaceItem }) {
  const navigate = useNavigate();
  const { isCreating, createInterface } = useCreateInterface();

  const {
    id: interfaceId,
    name,
    project,
    url,
    params,
    method,
    header,
    body,
    type,
    creator,
    updater,
    ctime,
    utime,
  } = interfaceItem;

  function handleDuplicate() {
    createInterface({
      name: `Copy of ${name}`,
      project,
      url,
      params,
      method,
      header,
      body,
      type,
      creator,
      updater,
      ctime,
      utime,
    });
  }

  const { isDeleting, deleteCabin } = useDeleteInterface();

  return (
    <Table.Row>
      {/* <Img src={image} /> */}
      <Interface>{name}</Interface>
      <Project>{project} </Project>
      <Type>{type} </Type>
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
                icon={<HiEye />}
                onClick={() => navigate(`/interfaces/${interfaceId}`)}
              >
                See details
              </Menus.Button>
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
              <CreateInterfaceForm interfaceToEdit={interfaceItem} />
            </Modal.Window>

            <Modal.Window name="delete">
              <ConfirmDelete
                resourceName="interfaces"
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

InterfaceRow.propTypes = {
  interfaceItem: PropTypes.shape({
    id: PropTypes.number.isRequired,
    name: PropTypes.string.isRequired,
    project: PropTypes.string.isRequired,
    url: PropTypes.string.isRequired,
    params: PropTypes.string.isRequired,
    method: PropTypes.string.isRequired,
    body: PropTypes.string.isRequired,
    header: PropTypes.string.isRequired,
    type: PropTypes.string.isRequired,
    creator: PropTypes.number.isRequired,
    updater: PropTypes.number.isRequired,
    ctime: PropTypes.string.isRequired,
    utime: PropTypes.string.isRequired,
  }).isRequired,
};

export default InterfaceRow;
