import styled from 'styled-components';
import { Button } from 'components/common';

export const ModalTitle = styled.h3`
  font-size: 1.2rem;
`;

export const CheckUl = styled.ul`
  list-style: none;
  padding: 0;
  margin-top: 20px;
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
  padding: 20px 0px;
  padding-top: 0px;
  z-index: 100;
`;

export const HeadWrap = styled.div`
  display: flex;
  align-items: center;
  padding: 25px 40px;
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
  font-size: 1.3rem;
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
  }
  @media only screen and (max-width: 470px) {
    gap: 6px;
  }
`;

export const DetailsWrap = styled.div`
  width: 100%;
  min-height: 100%;
  margin-top: 17px;
  padding: 0px 20px;
`;

export const ActionWrap = styled.div`
  display: flex;
  align-items: center;
  padding: 25px 40px;
  border-bottom: 1px solid #ebedef;
  @media only screen and (max-width: 700px) {
    padding: 25px 0px;
  }
  @media only screen and (max-width: 500px) {
    flex-direction: column;
    width: 100%;
    padding: 25px 0px;
  }
`;

export const BudgetWrap = styled.div`
  padding: 25px 40px;
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
  font-size: 0.7rem;
  color: #8e969c;
`;

export const Budget = styled.h4`
  color: #3c3f41;
  font-size: 1.15rem;
  @media only screen and (max-width: 500px) {
    font-size: 1rem;
  }
`;

export const Grey = styled.span`
  color: #8e969c;
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
  padding: 25px 40px;
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
  border-bottom: 1px solid #ebedef;
  padding-top: 5px;
  padding-bottom: 20px;
  @media only screen and (max-width: 500px) {
    width: 100%;
  }
`;

export const UsersHeader = styled.h4`
  font-size: 0.9rem;
  font-weight: 600;
  padding: 0;
  margin: 0;
  @media only screen and (max-width: 500px) {
    font-size: 0.8rem;
    margin-right: 55%;
  }
`;

export const UsersList = styled.div`
  @media only screen and (max-width: 500px) {
    width: 100%;
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
  flex-gap: 12px;
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
  font-size: 0.9rem;
  text-transform: capitalize;
  font-weight: bold;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  width: 100%;
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
