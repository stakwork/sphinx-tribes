import { EuiText } from '@elastic/eui';
import React, { useEffect, useState } from 'react';
import IconButton from '../../sphinxUI/icon_button';
import { useStores } from '../../store';

const BountyHeader = ({ selectedWidget, activeList, setShowFocusView }) => {
  const { main, ui } = useStores();
  const [peopleList, setPeopleList] = useState<Array<any> | null>(null);
  useEffect(() => {
    async function getPeopleList() {
      if (selectedWidget === 'wanted') {
        try {
          const response = await main.getPeople({ page: 1 });
          setPeopleList(response);
        } catch (error) {
          console.log(error);
        }
      } else {
        setPeopleList(null);
      }
    }
    getPeopleList();
  }, [main, selectedWidget]);

  return (
    <>
      <div
        style={{
          display: 'flex',
          flexDirection: 'row',
          justifyContent: 'space-between',
          padding: '10px 20px',
          alignItems: 'center'
        }}>
        <div
          style={{
            display: 'flex',
            flexDirection: 'row',
            justifyContent: 'space-evenly',
            alignItems: 'center'
          }}>
          <IconButton
            text={'Post a Bounty'}
            endingIcon={'add'}
            width={225}
            height={48}
            color={'success'}
            style={{
              color: '#fff',
              fontSize: '16px',
              fontWeight: '600'
            }}
            iconStyle={{
              fontSize: '16px',
              fontWeight: '600'
            }}
            onClick={() => {
              if (ui.meInfo && ui.meInfo?.owner_alias) {
                setShowFocusView(true);
              } else {
                ui.setShowSignIn(true);
              }
            }}
          />
          <IconButton
            text={`${activeList?.length} Bounties opened`}
            leadingIcon={'content_copy'}
            width={230}
            height={48}
            color={'transparent'}
            style={{
              color: '#909BAA',
              fontSize: '16px',
              fontWeight: '500',
              cursor: 'default',
              textDecoration: 'none'
            }}
            iconStyle={{
              fontSize: '18px',
              fontWeight: '500'
            }}
          />
          <IconButton
            text={'Filter'}
            color={'transparent'}
            leadingIcon={'tune'}
            width={80}
            height={48}
            style={{ color: '#909BAA', fontSize: '16px', fontWeight: '500' }}
            iconStyle={{
              fontSize: '18px',
              fontWeight: '500'
            }}
            onClick={() => {
              console.log('filter');
            }}
          />
        </div>
        <div
          style={{
            display: 'flex',
            flexDirection: 'row',
            alignItems: 'center',
            padding: '0 20px'
          }}>
          <EuiText
            color={'#909BAA'}
            style={{
              fontSize: '16px'
            }}>
            Developers
          </EuiText>
          <div
            style={{
              display: 'flex',
              flexDirection: 'row',
              alignItems: 'center',
              color: '#909BAA',
              padding: '0 10px'
            }}>
            {peopleList &&
              peopleList?.slice(0, 3).map((val, index) => {
                return (
                  <div
                    style={{
                      height: '28px',
                      width: '28px',
                      borderRadius: '50%',
                      background: '#fff',
                      overflow: 'hidden',
                      position: 'static',
                      display: 'flex',
                      justifyContent: 'center',
                      alignItems: 'center',
                      zIndex: 3 - index,
                      marginLeft: index > 0 ? '-14px' : ''
                    }}>
                    <img
                      height={'23px'}
                      width={'23px'}
                      src={val?.img || '/static/person_placeholder.png'}
                      alt={''}
                      style={{
                        borderRadius: '50%'
                      }}
                    />
                  </div>
                );
              })}
          </div>
          {peopleList && peopleList?.length}
        </div>
      </div>
    </>
  );
};

export default BountyHeader;
