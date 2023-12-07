import React , {useState} from "react";
import styled from "styled-components";
import { MyTable } from "./TableComponent";
import { bounties } from "./TableComponent/mockBountyData";
import { Header } from "./Header";
import { Statistics } from "./Staistics";

const Container = styled.body`
  height: 100vh; /* Set a fixed height for the container */
  overflow-y: auto; /* Enable vertical scrolling */
`;


export const SuperAdmin = () => {

  console.log("super admin",bounties);
  return (
    <Container>
    <Header/>
    <Statistics/>
    <MyTable 
    bounties={bounties}
    />
   
    </Container>
  );
};
