import styled from "styled-components";

import PropTypes from "prop-types";
import {
  // HiArrowDownOnSquare,
  // HiArrowUpOnSquare,
  HiEye,
  HiTrash,
} from "react-icons/hi2";

// import Tag from "../../ui/Tag.jsx";
import Table from "../../ui/Table.jsx";
import Modal from "../../ui/Modal.jsx";
import Menus from "../../ui/Menus.jsx";
import ConfirmDelete from "../../ui/ConfirmDelete.jsx";

import { useNavigate } from "react-router-dom";
// import { useCheckout } from "../check-in-out/useCheckout.js";
import { useDeleteBooking } from "./useDeleteBooking.js";
import { convertNanosecondsToHMS } from "../../utils/helpers.js";

const Cabin = styled.div`
  font-size: 1.6rem;
  font-weight: 600;
  color: var(--color-grey-600);
  font-family: "Sono";
`;

const Stacked = styled.div`
  display: flex;
  flex-direction: column;
  gap: 0.2rem;

  & span:first-child {
    font-weight: 500;
  }

  & span:last-child {
    color: var(--color-grey-500);
    font-size: 1.2rem;
  }
`;

// const Amount = styled.div`
//   font-family: "Sono";
//   font-weight: 500;
// `;

function TaskRow({ taskItem }) {
  const {
    id: taskId,
    name,
    apis,
    durations,
    workers,
    max_workers,
    rate,
    creator,
    // updater,
    ctime,
    // utime,
  } = taskItem || {};

  const navigate = useNavigate();
  // const { checkout, isCheckingOut } = useCheckout();
  const { isBookingDeleting, deleteBooking } = useDeleteBooking();
  // console.log(updater, utime);

  // 根据状态返回不同的颜色
  // const statusToTagName = {
  //   unconfirmed: "blue",
  //   "checked-in": "green",
  //   "checked-out": "silver",
  // };
  let result;

  if (apis.length > 2) {
    // 取前 3 个元素并连接
    result = apis.slice(0, 2).join(", ") + " ...";
  } else {
    // 如果元素少于等于 3，直接连接所有
    result = apis.join(", ");
  }

  return (
    <Table.Row>
      <Cabin>{name}</Cabin>
      <div>{result} </div>

      <Stacked>
        <span>{convertNanosecondsToHMS(durations)}</span>
      </Stacked>
      <Stacked>
        <div>
          <span>{workers}</span> ~<span>{max_workers}</span>
        </div>
      </Stacked>
      <Stacked>{rate}</Stacked>
      <Stacked>{creator}</Stacked>
      <Stacked>{ctime}</Stacked>
      <Modal>
        <Menus.Menu>
          <Menus.Toggle id={taskId} />
          <Menus.List id={taskId}>
            <Menus.Button
              icon={<HiEye />}
              onClick={() => navigate(`/tasks/${taskId}`)}
            >
              See details
            </Menus.Button>

            <Modal.Open opens="delete">
              <Menus.Button icon={<HiTrash />}>Delete task</Menus.Button>
            </Modal.Open>
          </Menus.List>
        </Menus.Menu>

        <Modal.Window name="delete">
          <ConfirmDelete
            resourceName="booking"
            onConfirm={() => deleteBooking(taskId)}
            disabled={isBookingDeleting}
          />
        </Modal.Window>
      </Modal>
    </Table.Row>
  );
}

TaskRow.propTypes = {
  taskItem: PropTypes.shape({
    id: PropTypes.number,
    name: PropTypes.string,
    a_ids: PropTypes.List,
    durations: PropTypes.string,
    workers: PropTypes.number,
    max_workers: PropTypes.number,
    rate: PropTypes.string,
    creator: PropTypes.string,
    updater: PropTypes.string,
    ctime: PropTypes.string,
    utime: PropTypes.string,
  }),
};

export default TaskRow;
