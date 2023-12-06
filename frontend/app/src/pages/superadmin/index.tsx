import React from "react";
import { MyTable } from "./TableComponent";
import { bounties } from "./TableComponent/mockBountyData";
import { Header } from "./Header";
import { Statistics } from "./Staistics";
export const SuperAdmin = () => {
    
  console.log("super admin");
  return (
    <>
    <Header/>
    <Statistics/>
    <MyTable bounties={bounties} selectedButtonIndex={2} />
   
    </>
  );
};
