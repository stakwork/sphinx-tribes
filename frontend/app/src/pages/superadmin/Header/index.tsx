import React from "react";
import {
  NavWrapper,
  AlternateWrapper,
  ButtonWrapper,
  Title,
  Button,
  AlternateTitle,
  ExportButton,
  ExportText,
  Month,
  ArrowButton,
  DropDown,
  DropDownOption,
  Flex,
} from "./headerstyles"
import signout from "./icons/signout.svg";
import arrowback from "./icons/arrowback.svg";
import arrowforward from "./icons/arrowforward.svg";


export const Header = () => (
    <>
  
          <NavWrapper>
            <Title>
              Bounties
              <AlternateTitle>Super Admin</AlternateTitle>
            </Title>
            <Button>
              Sign out 
              <img src={signout} alt="Sign Out" />

            </Button>
          </NavWrapper>
          <AlternateWrapper>
            <Flex>
              <ButtonWrapper>
                <ArrowButton>
                  <img src={arrowback} alt=""  />
                </ArrowButton>
                <ArrowButton><img src={arrowforward} alt=""  /></ArrowButton>
              </ButtonWrapper>
              <Month>01 Oct-31 Dec 2023</Month>
            </Flex>
            <Flex>
              <ExportButton>
                <ExportText>
                    Export CSV
                </ExportText>
              </ExportButton>
              <DropDown>
                <DropDownOption>7 Days</DropDownOption>
                <DropDownOption>30 Days</DropDownOption>
                <DropDownOption>45 Days</DropDownOption>
              </DropDown>
            </Flex>
          </AlternateWrapper>
      
    </>
  )

