import styled from "styled-components";
import PropTypes from "prop-types";
// import { format, isToday } from "date-fns";
import {
  HiOutlineChatBubbleBottomCenterText,
  HiOutlineCheckCircle,
  //   HiOutlineCurrencyDollar,
  HiOutlineHomeModern,
} from "react-icons/hi2";

import DataItem from "../../ui/DataItem.jsx";
// import { Flag } from "../../ui/Flag.jsx";

// import { formatDistanceFromNow, formatCurrency } from "../../utils/helpers.js";

const StyledBookingDataBox = styled.section`
  /* Box */
  background-color: var(--color-grey-0);
  border: 1px solid var(--color-grey-100);
  border-radius: var(--border-radius-md);

  overflow: hidden;
`;

const Header = styled.header`
  background-color: var(--color-brand-500);
  padding: 2rem 4rem;
  color: #e0e7ff;
  font-size: 1.8rem;
  font-weight: 500;
  display: flex;
  align-items: center;
  justify-content: space-between;

  svg {
    height: 3.2rem;
    width: 3.2rem;
  }

  & div:first-child {
    display: flex;
    align-items: center;
    gap: 1.6rem;
    font-weight: 600;
    font-size: 1.8rem;
  }

  & span {
    font-family: "Sono";
    font-size: 2rem;
    margin-left: 4px;
  }
`;

const Section = styled.section`
  padding: 3.2rem 4rem 1.2rem;
`;

const Guest = styled.div`
  display: flex;
  align-items: center;
  gap: 1.2rem;
  margin-bottom: 1.6rem;
  color: var(--color-grey-500);

  & p:first-of-type {
    font-weight: 500;
    color: var(--color-grey-700);
  }
`;

// const Price = styled.div`
//   display: flex;
//   align-items: center;
//   justify-content: space-between;
//   padding: 1.6rem 3.2rem;
//   border-radius: var(--border-radius-sm);
//   margin-top: 2.4rem;

//   background-color: ${(props) =>
//     props.isPaid ? "var(--color-green-100)" : "var(--color-yellow-100)"};
//   color: ${(props) =>
//     props.isPaid ? "var(--color-green-700)" : "var(--color-yellow-700)"};

//   & p:last-child {
//     text-transform: uppercase;
//     font-size: 1.4rem;
//     font-weight: 600;
//   }

//   svg {
//     height: 2.4rem;
//     width: 2.4rem;
//     color: currentColor !important;
//   }
// `;

const Footer = styled.footer`
  padding: 1.6rem 4rem;
  font-size: 1.2rem;
  color: var(--color-grey-500);
  text-align: right;
`;

// A purely presentational component
function DebugResultBox({ data }) {
  const observations = true;
  return (
    <StyledBookingDataBox>
      <Header>
        <div>
          <HiOutlineHomeModern />
          <p>
            response
            {/* {numNights} nights in Cabin <span>{cabinName}</span> */}
          </p>
        </div>

        <p></p>
      </Header>

      <Section>
        <Guest>
          response
          {/* {countryFlag && <Flag src={countryFlag} alt={`Flag of ${country}`} />} */}
          <p>
            response
            {/* {guestName} {numGuests > 1 ? `+ ${numGuests - 1} guests` : ""} */}
          </p>
          <span>&bull;</span>
          <p>response</p>
          <span>&bull;</span>
          {/* <p>National ID {nationalID}</p> */}
        </Guest>

        {observations && (
          <DataItem
            icon={<HiOutlineChatBubbleBottomCenterText />}
            label="Observations"
          >
            {/* {observations} */}
            response9999
          </DataItem>
        )}

        <DataItem icon={<HiOutlineCheckCircle />} label="Breakfast included?">
          {/* {hasBreakfast ? "Yes" : "No"} */}
          <pre>{data}</pre>
        </DataItem>

        {/* <Price isPaid={isPaid}>
          <DataItem icon={<HiOutlineCurrencyDollar />} label={`Total price`}>
            {formatCurrency(totalPrice)}

            {hasBreakfast &&
              ` (${formatCurrency(cabinPrice)} cabin + ${formatCurrency(
                extrasPrice
              )} breakfast)`}
          </DataItem>

          <p>{isPaid ? "Paid" : "Will pay at property"}</p>
        </Price> */}
      </Section>

      <Footer>
        {/* <p>Booked {format(new Date(created_at), "EEE, MMM dd yyyy, p")}</p> */}
      </Footer>
    </StyledBookingDataBox>
  );
}

DebugResultBox.propTypes = {
  data: PropTypes.string,
};

export default DebugResultBox;
