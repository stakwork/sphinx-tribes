import React from 'react';
import styled from 'styled-components';
import MaterialIcon from '@material/react-material-icon';
import { WidgetProps } from './interfaces';

const Dot = styled.div`
  color: #1d1e24;
  position: absolute;
  top: 6px;
  right: 6px;
  opacity: 0.8;
`;

const Title = styled.div`
  width: 100%;
  display: flex;
  align-content: center;
  justify-content: center;
  margin-top: 10px;
  color: #000;
`;

const Wrap = styled.div`
  background: #ffffff;
  width: 145px;
  margin: 5px;
  padding: 10px 5px;
  display: flex;
  border-radius: 5px;
  flex-direction: column;
  align-content: center;
  align-items: center;
  justify-content: center;
  box-shadow: 0 2px 5px 0 rgba(0, 0, 0, 0.2);
  position: relative;
  cursor: pointer;
`;

export interface IconProps {
  source: string;
}

const Icon = styled.div<IconProps>`
  background-image: ${(p: any) => `url(${p.source})`};
  width: 70px;
  height: 70px;
  margin-top: 10px;
  background-position: center; /* Center the image */
  background-repeat: no-repeat; /* Do not repeat the image */
  background-size: contain; /* Resize the background image to cover the entire container */
  border-radius: 5px;
  overflow: hidden;
`;
export default function Widget(props: WidgetProps) {
  // highlight if state has any of these

  const { values, name, parentName, setFieldValue } = props;
  const state = values.extras && values.extras[name];

  function objectOrArrayHasLength(mystate: any) {
    let v = 0;
    if (mystate) {
      if (Array.isArray(mystate)) {
        v = mystate.length;
      } else {
        v = Object.keys(mystate).length;
      }
    }
    return v > 0;
  }

  function deleteSingleWidget() {
    setFieldValue(`${parentName}.${name}`, undefined);
  }

  const highlight = objectOrArrayHasLength(state);

  return (
    <Wrap onClick={() => props.setSelected(props)}>
      <Icon source={`/static/${props.icon || 'sphinx'}.png`} />

      <Title>{props.label}</Title>

      {highlight && (
        <Dot>
          <MaterialIcon
            icon={props.single ? 'close' : 'settings'}
            style={{
              fontSize: 20,
              color: '#1d1e24',
              cursor: props.single && 'pointer'
            }}
            onClick={(e: any) => {
              if (props.single) {
                e.stopPropagation();
                //delete state
                deleteSingleWidget();
              }
            }}
          />
        </Dot>
      )}
    </Wrap>
  );
}
