import Heading from "../ui/Heading";
import Row from "../ui/Row";
import UpdateUserDataForm from "../features/authentication/UpdateUserDataForm";
import UpdatePasswordForm from "../features/authentication/UpdatePasswordForm";
import { useUser } from "../features/authentication/useUser";
import { styled } from "styled-components";

const Profile = styled.div`
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

function Account() {
  const { userData } = useUser();

  return (
    <>
      <Heading as="h1">Profile</Heading>

      <Row>
        <Heading as="h3">
          <Profile>
            <p>{userData.fullName}</p>
            <span>&bull;</span>
            <p>{userData.email}</p>
          </Profile>
        </Heading>

        <UpdateUserDataForm userToEidt={userData} />
      </Row>

      <Row>
        <Heading as="h3">Update password(Under construction)</Heading>
        <UpdatePasswordForm />
      </Row>
    </>
  );
}

export default Account;
