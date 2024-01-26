import React, { useState, useEffect } from 'react';
import styled from 'styled-components';
import { EuiButton } from '@elastic/eui';
import MaterialIcon from '@material/react-material-icon';
import Input from '..';
import WidgetList from './WidgetList';
import { FocusedWidgetProps } from './interfaces';

const Wrap = styled.div`
  color: #fff;
  display: flex;
  flex-direction: column;
  align-content: center;
  justify-content: space-evenly;
`;

const Nav = styled.div`
  color: #fff;
  display: flex;
  align-content: center;
  justify-content: space-evenly;
  margin-bottom: 20px;
  margin-top: 2px;
  height: 42px;
  min-height: 42px;
`;

export interface IconProps {
  source: string;
}

const Icon = styled.div<IconProps>`
  background-image: ${(p: any) => `url(${p.source})`};
  width: 30px;
  height: 30px;
  background-position: center; /* Center the image */
  background-repeat: no-repeat; /* Do not repeat the image */
  background-size: contain; /* Resize the background image to cover the entire container */
  border-radius: 5px;
  overflow: hidden;
`;
export default function FocusedWidget(props: FocusedWidgetProps) {
  const {
    name,
    values,
    errors,
    initialValues,
    setFieldTouched,
    setFieldValue,
    item,
    setShowFocused,
    setDisableFormButtons
  } = props;
  const { single } = item;
  const [selectedIndex, setSelectedIndex] = useState(-1);

  let newState = values[name] && values[name][item.name];
  let newErrors = errors[name] && errors[name][item.name];

  if (selectedIndex > -1) {
    newState = newState && newState[selectedIndex];
    newErrors = newErrors && newErrors[selectedIndex];
  }

  const [prevState, setPrevState] = useState(newState);

  useEffect(() => {
    setDisableFormButtons(true);

    return function cleanup() {
      setDisableFormButtons(false);
    };
  }, []);

  function getFieldToUpdate(e: any) {
    let valueToUpdate = `${name}.${item.name}.${e.name}`;
    if (!single) {
      valueToUpdate = `${name}.${item.name}[${selectedIndex}].${e.name}`;
    }
    return valueToUpdate;
  }

  function getInitialValueByType(type: any) {
    let value: any = '';
    if (type === 'number') value = 0;
    if (type === 'gallery') value = [];
    return value;
  }

  function thereAreErrors() {
    let result = false;

    if (newErrors && Array.isArray(newErrors) && newErrors.length) {
      result = true;
    } else if (newErrors && Object.keys(newErrors).length) {
      result = true;
    }
    return result;
  }
  function getFormState() {
    return (values[name] && values[name][item.name]) || [];
  }
  useEffect(() => {
    // in order to
    if (single) {
      const firstInput = item.fields[0];
      const firstInputName = getFieldToUpdate(firstInput);
      let vl = newState && newState[firstInput.name];
      if (!vl) {
        // get initial value by type
        vl = getInitialValueByType(firstInput.type);
      }
      setFieldValue(firstInputName, vl);
    }
  }, [selectedIndex]);

  function cancel(dismount: any) {
    // new widget cancelled, revert form state
    let returnState = prevState;
    if (!single && selectedIndex < 0) {
      returnState = getFormState();
    }

    setFieldValue(`${name}.${item.name}`, returnState);

    if (single || dismount) setShowFocused(false);
    // else setSelectedIndex(-1)
  }

  function done() {
    // new widget added, keep form state and reset selected index
    // if errors, prevent!
    if (thereAreErrors()) {
      alert('validation errors!');
      return;
    }
    if (single) setShowFocused(false);
    else {
      setSelectedIndex(-1);
      setPrevState(getFormState());
    }
  }

  function startCreate() {
    const cloneformState = getFormState();
    setPrevState(cloneformState);
    const formState = [...cloneformState];

    if (single) setFieldTouched(item.name, true);

    const obj = {};
    item.fields.forEach((o: any) => {
      let val: any = '';
      if (o.type === 'number') val = 0;
      obj[o.name] = val;
    });

    formState.unshift(obj);

    setFieldValue(`${name}.${item.name}`, formState);
    setSelectedIndex(0);
  }

  function startEdit(obj: any, i: any) {
    const cloneformState = getFormState();
    setPrevState(cloneformState);
    const formState = [...cloneformState];

    let index = i;
    if (obj.id) {
      index = formState && formState.findIndex((f: any) => f.id === obj.id);
    }
    setSelectedIndex(index);
  }

  function deleteItem(obj: any, i: any) {
    const formState = (values[name] && [...values[name][item.name]]) || [];
    let index = i;
    if (obj.id) {
      index = formState && formState.findIndex((f: any) => f.id === obj.id);
    }

    formState.splice(index, 1);
    setFieldValue(`${name}.${item.name}`, formState);
  }

  const showingList = single ? false : selectedIndex > -1 ? false : true;

  const widgetHeader = (
    <div
      style={{
        display: 'flex',
        marginBottom: 20,
        justifyContent: 'space-between',
        width: '100%',
        alignItems: 'center',
        cursor: 'pointer',
        height: 45
      }}
    >
      <div style={{ display: 'flex' }} onClick={() => cancel(true)}>
        <MaterialIcon icon={'west'} />
        <Icon source={`/static/${item.icon || 'sphinx'}.png`} style={{ marginLeft: 5 }} />
        <div style={{ marginLeft: 10 }}>{item.label}</div>
      </div>
      {showingList ? (
        <EuiButton onClick={() => startCreate()} style={{ fontSize: 12, fontWeight: 600 }}>
          Add New
        </EuiButton>
      ) : (
        <div />
      )}
    </div>
  );

  return (
    <Wrap>
      {widgetHeader}

      {/* single widgets will only show these fields */}
      {!showingList && (
        <>
          {item.fields.map((e: any) => (
            <Input
              {...e}
              key={e.name}
              value={newState && newState[e.name]}
              error={newErrors && newErrors[e.name]}
              initialValues={initialValues}
              handleChange={(c: any) => {
                const f = getFieldToUpdate(e);
                setFieldValue(f, c);
              }}
              handleBlur={() => {
                const f = getFieldToUpdate(e);
                setFieldTouched(f, false);
              }}
              handleFocus={() => {
                const f = getFieldToUpdate(e);
                setFieldTouched(f, true);
              }}
            />
          ))}
          <Nav style={{ marginTop: 20 }}>
            <EuiButton onClick={() => cancel(false)} style={{ fontSize: 12, fontWeight: 600 }}>
              Cancel
            </EuiButton>
            <EuiButton
              onClick={() => done()}
              disabled={thereAreErrors()}
              style={{ fontSize: 12, fontWeight: 600 }}
            >
              Save{item.itemLabel && ` ${item.itemLabel}`}
            </EuiButton>
          </Nav>
        </>
      )}

      {/* only multi widgets will only show the list */}

      {showingList && (
        <WidgetList
          schema={item}
          values={values[name] && values[name][item.name]}
          setSelected={startEdit}
          deleteItem={deleteItem}
        />
      )}
    </Wrap>
  );
}
