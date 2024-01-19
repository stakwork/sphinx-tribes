import styled from 'styled-components';
import { Button } from 'components/common';

interface SmallBtnProps {
  selected: boolean;
}

interface UserProps {
  inactive: boolean;
}

export const ModalTitle = styled.h3`
  font-size: 1.2rem;
`;

export const CheckUl = styled.ul`
  list-style: none;
  padding: 0;
  margin-bottom: 20px;
`;

export const CheckLi = styled.li`
  display: flex;
  flex-direction: row;
  align-items: center;
  padding: 0px;
  margin-bottom: 10px;
`;

export const Check = styled.input`
  width: 20px;
  height: 20px;
  border-radius: 5px;
  padding: 0px;
  margin-right: 10px;
`;

export const CheckLabel = styled.label`
  padding: 0px;
  margin: 0px;
`;

export const ViewBounty = styled.p`
  padding: 0px;
  margin: 0px;
  cursor: pointer;
  font-size: 0.9rem;
  color: green;
  font-size: bold;
`;

export const Container = styled.div`
  width: 100%;
  min-height: 100%;
  background: white;
  z-index: 100;
`;

export const HeadWrap = styled.div`
  display: flex;
  align-items: center;
  padding: 25px 10px;
  padding-right: 40px;
  border-bottom: 1px solid #ebedef;
  @media only screen and (max-width: 800px) {
    padding: 15px 0px;
  }
  @media only screen and (max-width: 700px) {
    padding: 12px 0px;
  }
  @media only screen and (max-width: 500px) {
    padding: 0px;
    padding-bottom: 15px;
    flex-direction: column;
    align-items: start;
    padding: 20px 30px;
  }
`;

export const HeadNameWrap = styled.div`
  display: flex;
  align-items: center;
  @media only screen and (max-width: 500px) {
    margin-bottom: 20px;
  }
`;

export const OrgImg = styled.img`
  width: 48px;
  height: 48px;
  border-radius: 50%;
  margin-left: 20px;
  @media only screen and (max-width: 700px) {
    width: 42px;
    height: 42px;
  }
  @media only screen and (max-width: 500px) {
    width: 38px;
    height: 38px;
  }
  @media only screen and (max-width: 470px) {
    width: 35px;
    height: 35px;
  }
`;

export const OrgName = styled.h3`
  padding: 0px;
  margin: 0px;
  font-size: 1.5rem;
  color: #3c3f41;
  margin-left: 25px;
  font-weight: 700;
  margin-left: 20px;
  @media only screen and (max-width: 800px) {
    font-size: 1.05rem;
  }
  @media only screen and (max-width: 700px) {
    font-size: 1rem;
  }
  @media only screen and (max-width: 600px) {
    font-size: 0.9rem;
  }
  @media only screen and (max-width: 470px) {
    font-size: 0.8rem;
  }
`;

export const HeadButtonWrap = styled.div<{ forSmallScreen: boolean }>`
  margin-left: auto;
  display: flex;
  flex-direction: row;
  gap: 15px;
  @media only screen and (max-width: 700px) {
    gap: 10px;
    margin-left: auto;
  }
  @media only screen and (max-width: 500px) {
    gap: 8px;
    margin-left: 0px;
    width: 100vw;
    margin-left: ${(p: any) => (p.forSmallScreen ? '50px' : '0px')};
    flex-wrap: wrap;
    display: flex;
  }
  @media only screen and (max-width: 470px) {
    gap: 6px;
  }
`;

export const DetailsWrap = styled.div`
  width: 100%;
  min-height: 100%;
  padding: 0px 20px;
`;

export const ActionWrap = styled.div`
  display: flex;
  align-items: center;
  padding-right: 40px;
  border-bottom: 1px solid #ebedef;
  box-shadow: 0px 1px 2px 0px rgba(0, 0, 0, 0.15);
  @media only screen and (max-width: 700px) {
    padding: 25px 0px;
  }
  @media only screen and (max-width: 500px) {
    flex-direction: column;
    width: 100%;
    padding: 20px 30px;
  }
`;

export const BudgetWrap = styled.div`
  padding: 25px 60px;
  width: 55%;
  display: flex;
  flex-direction: column;
  @media only screen and (max-width: 700px) {
    width: 100%;
    padding: 22px 0px;
  }
  @media only screen and (max-width: 500px) {
    width: 100%;
    padding: 20px 0px;
    padding-top: 0;
  }
`;

export const NoBudgetWrap = styled.div`
  display: flex;
  flex-direction: row;
  align-items: center;
  width: 100%;
  border: 1px solid #ebedef;
`;

export const ViewBudgetWrap = styled.div`
  display: flex;
  flex-direction: column;
  width: 100%;
`;

export const ViewBudgetTextWrap = styled.div`
  display: flex;
  width: 100%;
  align-items: center;
  margin-top: 12px;
`;

export const BudgetSmall = styled.h6`
  padding: 0px;
  font-size: 0.8rem;
  color: #8e969c;
  @media only screen and (max-width: 500px) {
    font-size: 0.75rem;
  }
`;

export const BudgetSmallHead = styled.h6`
  padding: 0px;
  font-size: 0.625rem;
  color: #8e969c;
  margin: 0;
`;

export const Budget = styled.h4`
  color: #3c3f41;
  font-size: 1.0625rem;
  font-weight: 600;

  &.budget-small {
    border-left: 1px solid #ebedef;
    padding-left: 22px;
    margin-left: 22px;
  }

  @media only screen and (max-width: 500px) {
    font-size: 1rem;
  }
`;

export const Grey = styled.span`
  color: #8e969c;
  font-weight: 400;
`;

export const NoBudgetText = styled.p`
  font-size: 0.85rem;
  padding: 0px;
  margin: 0px;
  color: #8e969c;
  width: 90%;
  margin-left: auto;
`;

export const UserWrap = styled.div`
  display: flex;
  flex-direction: column;
  background-color: rgb(240, 241, 243);
  @media only screen and (max-width: 700px) {
    width: 100%;
    padding: 20px 0px;
  }
  @media only screen and (max-width: 500px) {
    padding: 20px 0px;
  }
`;

export const UsersHeadWrap = styled.div`
  display: flex;
  align-items: center;
  width: 100%;
  padding: 20px 60px;
  padding-right: 40px;
  border-bottom: 1px solid #dde1e5;
  @media only screen and (max-width: 500px) {
    width: 100%;
    padding: 0 30px;
    padding-bottom: 20px;
  }
`;

export const UsersHeader = styled.h4`
  font-size: 0.8125rem;
  font-weight: 700;
  padding: 0;
  margin: 0;
  color: #3c3f41;
  @media only screen and (max-width: 500px) {
    font-size: 0.8rem;
    margin-right: 55%;
  }
`;

export const UsersList = styled.div`
  padding: 0 60px;
  padding-right: 40px;
  border-bottom: 1px solid #dde1e5;

  @media only screen and (max-width: 500px) {
    width: 100%;
    padding: 0 30px;
  }
`;

export const UserImage = styled.img`
  width: 40px;
  height: 40px;
  border-radius: 50%;
`;

export const User = styled.div`
  padding: 15px 0px;
  border-bottom: 1px solid #ebedef;
  display: flex;
  align-items: center;
  @media only screen and (max-width: 500px) {
    padding: 10px 0px;
    width: 100%;
  }
`;

export const UserDetails = styled.div`
  display: flex;
  flex-direction: column;
  margin-left: 2%;
  width: 30%;
  @media only screen and (max-width: 500px) {
    width: 60%;
    margin-left: 5%;
  }
`;

export const UserName = styled.p`
  padding: 0px;
  margin: 0px;
  font-size: 0.9375rem;
  text-transform: capitalize;
  font-weight: 600;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  width: 100%;
  color: #3c3f41;
`;

export const UserPubkey = styled.p`
  padding: 0px;
  margin: 0px;
  font-size: 0.75rem;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  width: 100%;
  color: #5f6368;
`;

export const UserAction = styled.div`
  display: flex;
  align-items: center;
  margin-left: auto;
`;

export const ActionBtn = styled.button`
  border: 0px;
  padding: 0px;
`;

export const IconWrap = styled.div`
  :first-child {
    margin-right: 40px;
    @media only screen and (max-width: 700px) {
      margin-right: 20px;
    }
    @media only screen and (max-width: 500px) {
      margin-right: 8px;
    }
  }
`;

export const HeadButton = styled(Button)`
  border-radius: 5px;
`;

export const ImgText = styled.h3`
  color: #b0b7bc;
  text-align: center;
  font-family: 'Barlow';
  font-size: 1.875rem;
  font-style: normal;
  font-weight: 800;
  line-height: 1.0625rem;
  letter-spacing: 0.01875rem;
  text-transform: uppercase;
  opacity: 0.5;
  margin-bottom: 0;
`;

export const ImgDashContainer = styled.div`
  width: 8.875rem;
  height: 8.875rem;
  border-radius: 50%;
  display: flex;
  justify-content: center;
  align-items: center;
  border: 1px dashed #d0d5d8;
  padding: 0.5rem;
  position: relative;
`;

export const UploadImageContainer = styled.div`
  display: flex;
  align-items: center;
  justify-content: center;
  width: 2.37756rem;
  height: 2.37756rem;
  position: absolute;
  bottom: 0;
  right: 0;
  cursor: pointer;
`;

export const ImgContainer = styled.div`
  width: 7.875rem;
  height: 7.875rem;
  border-radius: 50%;
  display: flex;
  justify-content: center;
  align-items: center;
  background-color: #ebedf1;
`;

export const SelectedImg = styled.img`
  width: 7.875rem;
  height: 7.875rem;
  border-radius: 50%;
  object-fit: cover;
`;

export const ImgTextContainer = styled.div`
  display: flex;
  justify-content: center;
  align-items: center;
  flex-direction: column;
  margin-top: 1rem;
`;

export const InputFile = styled.input`
  display: none;
`;

export const ImgInstructionText = styled.p`
  color: #5f6368;
  text-align: center;
  font-family: 'Roboto';
  font-size: 0.8125rem;
  font-style: normal;
  font-weight: 400;
  line-height: 1.0625rem;
  letter-spacing: 0.00813rem;
  margin-bottom: 0;
`;

export const ImgInstructionSpan = styled.span`
  color: #618aff;
  cursor: pointer;
`;

export const ImgDetailInfo = styled.p`
  color: #b0b7bc;
  text-align: center;
  font-family: 'Roboto';
  font-size: 0.625rem;
  font-style: normal;
  font-weight: 400;
  line-height: 1.125rem;
  margin-bottom: 0;
  margin-top: 1rem;
`;

export const OrgInputContainer = styled.div`
  width: 16rem;
  height: 223px;
  display: flex;
  flex-direction: column;
  @media only screen and (max-width: 500px) {
    width: 100%;
    margin-top: 1rem;
  }
`;

export const OrgLabel = styled.label`
  color: #5f6368;
  font-family: 'Barlow';
  font-size: 13px;
  font-style: normal;
  font-weight: 500;
  margin-bottom: 0.75rem;
  height: 0.5625rem;
`;

export const TextInput = styled.input`
  padding: 8px 14px;
  border-radius: 6px;
  border: 2px solid #dde1e5;
  outline: none;
  caret-color: #618aff;
  color: #3c3f41;
  font-family: 'Barlow';
  font-size: 14px;
  font-style: normal;
  font-weight: 500;
  line-height: 20px;
  width: 16rem;
  height: 2.4rem;

  ::placeholder {
    color: #b0b7bc;
    font-family: 'Barlow';
    font-size: 14px;
    font-style: normal;
    font-weight: 400;
    line-height: 20px;
  }

  :focus {
    border: 2px solid #82b4ff;
  }
`;

export const TextAreaInput = styled.textarea`
  padding: 0.5rem 1rem;
  border-radius: 0.375rem;
  border: 2px solid #dde1e5;
  outline: none;
  caret-color: #618aff;
  color: #3c3f41;
  font-family: 'Barlow';
  font-size: 13px;
  font-style: normal;
  font-weight: 500;
  line-height: 20px;
  width: 16rem;
  resize: none;
  height: 13.9375rem;

  ::placeholder {
    color: #b0b7bc;
    font-family: 'Barlow';
    font-size: 13px;
    font-style: normal;
    font-weight: 400;
    line-height: 20px;
  }
  :focus {
    border: 2px solid #82b4ff;
  }
`;
export const SecondaryText = styled.p`
  color: #b0b7bc;
  font-family: 'Barlow';
  font-size: 0.813rem;
  font-style: normal;
  font-weight: 400;
  margin-bottom: 18px;
  height: 0.5625rem;
`;
export const RouteHintText = styled.p`
  font-size: 0.9rem;
  text-align: center;
  color: #9157f6;
`;

export const AddUserContainer = styled.div`
  display: flex;
  flex-direction: column;
  width: 100%;
`;

export const AddUserHeaderContainer = styled.div`
  display: flex;
  padding: 1.875rem;
  flex-direction: column;
`;

export const AddUserHeader = styled.h2`
  color: #3c3f41;
  font-family: 'Barlow';
  font-size: 1.625rem;
  font-style: normal;
  font-weight: 800;
  line-height: normal;
  margin-bottom: 1.25rem;
`;

export const SearchUserInput = styled.input`
  padding: 0.9375rem 0.875rem;
  border-radius: 0.375rem;
  border: 1px solid #dde1e5;
  background: #fff;
  width: 100%;
  color: #292c33;
  font-family: 'Barlow';
  font-size: 0.8125rem;
  font-style: normal;
  font-weight: 400;

  ::placeholder {
    color: #8e969c;
  }
`;

export const UsersListContainer = styled.div`
  display: flex;
  flex-direction: column;
  padding: 1rem 1.875rem;
  background-color: #f2f3f5;
  height: 16rem;
  overflow-y: auto;
`;

export const UserContianer = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 1rem;
`;

export const UserInfo = styled.div<UserProps>`
  display: flex;
  align-items: center;
  opacity: ${(p: any) => (p.inactive ? 0.3 : 1)};
`;

export const UserImg = styled.img`
  width: 2rem;
  height: 2rem;
  border-radius: 50%;
  margin-right: 0.63rem;
  object-fit: cover;
`;

export const Username = styled.p`
  color: #3c3f41;
  font-family: 'Barlow';
  font-size: 0.8125rem;
  font-style: normal;
  font-weight: 500;
  line-height: 1rem;
  margin-bottom: 0;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 150px;
`;

export const SmallBtn = styled.button<SmallBtnProps>`
  width: 5.375rem;
  height: 2rem;
  padding: 0.625rem;
  border-radius: 0.375rem;
  background: ${(p: any) => (p.selected ? '#618AFF' : '#dde1e5')};
  color: ${(p: any) => (p.selected ? '#FFF' : '#5f6368')};
  font-family: 'Barlow';
  font-size: 0.8125rem;
  font-style: normal;
  font-weight: 600;
  line-height: 0rem; /* 0% */
  letter-spacing: 0.00813rem;
  border: none;
`;

export const FooterContainer = styled.div`
  display: flex;
  padding: 1.125rem 1.875rem;
  flex-direction: column;
  justify-content: center;
  align-items: center;
`;

export const AddUserBtn = styled.button`
  height: 3rem;
  padding: 0.5rem 1rem;
  width: 100%;
  border-radius: 0.375rem;
  font-family: 'Barlow';
  font-size: 0.9375rem;
  font-style: normal;
  font-weight: 500;
  line-height: 0rem;
  letter-spacing: 0.00938rem;
  background: #618aff;
  box-shadow: 0px 2px 10px 0px rgba(97, 138, 255, 0.5);
  border: none;
  color: #fff;
  &:disabled {
    border: 1px solid rgba(0, 0, 0, 0.07);
    background: rgba(0, 0, 0, 0.04);
    color: rgba(142, 150, 156, 0.85);
    cursor: not-allowed;
    box-shadow: none;
  }
`;

export const AssignUserContainer = styled.div`
  display: flex;
  flex-direction: column;
  width: 100%;
  position: relative;
`;

export const UserInfoContainer = styled.div`
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  position: absolute;
  top: -2.5rem;
  left: 0;
  right: 0;
`;

export const AssignRoleUserImage = styled.img`
  width: 5rem;
  height: 5rem;
  border-radius: 50%;
  background: #dde1e5;
  border: 4px solid #fff;
  object-fit: cover;
`;

export const AssignRoleUsername = styled.p`
  color: #3c3f41;
  text-align: center;
  font-family: 'Barlow';
  font-size: 1.25rem;
  font-style: normal;
  font-weight: 500;
  line-height: 1.625rem;
  margin-top: 0.69rem;
  margin-bottom: 0;
  text-transform: capitalize;
`;

export const UserRolesContainer = styled.div`
  padding: 3.25rem 3rem 3rem 3rem;
  margin-top: 3.25rem;
`;

export const UserRolesTitle = styled.h2`
  color: #3c3f41;
  font-family: 'Barlow';
  font-size: 1.625rem;
  font-style: normal;
  font-weight: 800;
  line-height: 1.625rem;
  margin-bottom: 2.81rem;
`;

export const RolesContainer = styled.div`
  display: flex;
  flex-direction: column;
`;

export const RoleContainer = styled.div`
  display: flex;
  align-items: center;
  margin-bottom: 1rem;
`;

export const Checkbox = styled.input`
  margin-right: 1rem;
  width: 1rem;
  height: 1rem;
`;

export const Label = styled.label`
  margin-bottom: 0;
  color: #1e1f25;
  font-family: 'Barlow';
  font-size: 0.9375rem;
  font-style: normal;
  font-weight: 500;
  line-height: 1.125rem;
`;

export const AssingUserBtn = styled.button`
  height: 3rem;
  padding: 0.5rem 1rem;
  width: 100%;
  border-radius: 0.375rem;
  font-family: 'Barlow';
  font-size: 0.9375rem;
  font-style: normal;
  font-weight: 500;
  line-height: 0rem;
  letter-spacing: 0.00938rem;
  background: #618aff;
  box-shadow: 0px 2px 10px 0px rgba(97, 138, 255, 0.5);
  border: none;
  margin-top: 3rem;
  color: #fff;
  &:disabled {
    border: 1px solid rgba(0, 0, 0, 0.07);
    background: rgba(0, 0, 0, 0.04);
    color: rgba(142, 150, 156, 0.85);
    cursor: not-allowed;
    box-shadow: none;
  }
`;

export const BudgetButton = styled.button`
  width: 100%;
  padding: 1rem;
  border-radius: 0.375rem;
  margin-top: 1.25rem;
  font-family: 'Barlow';
  font-size: 0.9375rem;
  font-style: normal;
  font-weight: 500;
  letter-spacing: 0.00938rem;
  background: #49c998;
  box-shadow: 0px 2px 10px 0px rgba(73, 201, 152, 0.5);
  border: none;
  color: #fff;
  &:disabled {
    border: 1px solid rgba(0, 0, 0, 0.07);
    background: rgba(0, 0, 0, 0.04);
    color: rgba(142, 150, 156, 0.85);
    cursor: not-allowed;
    box-shadow: none;
  }
`;
