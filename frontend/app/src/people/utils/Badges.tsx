import React, { useCallback, useEffect, useState } from 'react';
import styled from 'styled-components';
import MaterialIcon from '@material/react-material-icon';
import { observer } from 'mobx-react-lite';
import { BadgesProps } from 'people/interfaces';
import { useIsMobile } from '../../hooks';
import { useStores } from '../../store';
import { Modal, Button, Divider, TextInput } from '../../components/common';
import PageLoadSpinner from './PageLoadSpinner';

interface BProps {
  readonly isMobile?: boolean;
}
const Wrap = styled.div<BProps>`
  display: flex;
  flex-wrap: ${(p: any) => (p.isMobile ? '' : 'wrap')};
  width: 100%;
  overflow-x: hidden;
`;

const BWrap = styled.div<BProps>`
  display: flex;
  // cursor:pointer;
  flex-direction: ${(p: any) => (p.isMobile ? 'row' : 'column')};
  position: relative;
  width: ${(p: any) => (p.isMobile ? '100%' : '192px')};
  min-width: ${(p: any) => (p.isMobile ? '100%' : '192px')};
  height: ${(p: any) => (p.isMobile ? '' : '272px')};
  min-height: ${(p: any) => (p.isMobile ? '' : '272px')};
  max-width: ${(p: any) => (p.isMobile ? '100%' : '192px')};
  align-items: center;
  padding: ${(p: any) => (p.isMobile ? '10px' : '20px 10px 10px')};
  background: #fff;
  margin-bottom: 10px;
  border-radius: 4px;
  box-shadow: 0px 1px 2px rgb(0 0 0 / 15%);

  width: ${(p: any) => (p.isMobile ? '100%' : 'auto')};
  margin-right: ${(p: any) => (p.isMobile ? '0px' : '20px')};
`;
const T = styled.div<BProps>`
  font-size: 15px;
  width: 100%;
  text-align: ${(p: any) => (p.isMobile ? 'left' : 'center')};

  font-family: Roboto;
  font-style: normal;
  font-weight: 600;
  font-size: ${(p: any) => (p.isMobile ? '20px' : '15px')};
  line-height: 20px;
  /* or 133% */

  /* Primary Text 1 */

  color: #292c33;
`;
const S = styled.div<BProps>`
  font-size: 15px;
  margin-left: ${(p: any) => (p.isMobile ? '15px' : '10px')};
  width: 100%;
  text-align: ${(p: any) => (p.isMobile ? '' : 'center')};

  font-family: Roboto;
  font-style: normal;
  font-weight: 400;
  font-size: 15px;
  line-height: 15px;
  /* or 133% */

  text-align: center;

  /* Primary Text 1 */

  color: #5f6368;
`;
const D = styled.div`
  width: 100%;
  font-size: 12px;
  display: flex;
  justify-content: center;
  align-items: center;

  line-height: 30px;
  /* or 400% */
  min-height: 30px;
  display: flex;
  align-items: center;
  text-align: center;

  /* Text 2 */

  color: #3c3f41;
`;

const Status = styled.div`
  margin: 15px 20px 25px;
  // width:100%;
  font-size: 12px;
  display: flex;
  justify-content: center;
`;
const StatusText = styled.div`
  display: flex;
  justify-content: center;
  align-items: center;
  // height:26px;
  color: #b0b7bc;
  font-size: 10px;
  padding: 0 10px;
  cursor: pointer;

  &:hover {
    color: #618aff99;
  }
`;
const Counter = styled.div``;

interface ImageProps {
  readonly src?: string;
  readonly isMobile?: boolean;
}
const Img = styled.div<ImageProps>`
  position: relative;
  background-image: url('${(p: any) => p.src}');
  background-position: center;
  background-size: cover;
  position: relative;
  min-width: ${(p: any) => (p.isMobile ? '108px' : '132px')};
  width: ${(p: any) => (p.isMobile ? '108px' : '132px')};
  min-height: ${(p: any) => (p.isMobile ? '108px' : '132px')};
  height: ${(p: any) => (p.isMobile ? '108px' : '132px')};
  margin: ${(p: any) => (p.isMobile ? '24px' : '0px 20px 20px')};
`;
const SmallImg = styled.div<ImageProps>`
  background-image: url('${(p: any) => p.src}');
  background-position: center;
  background-size: cover;
  position: relative;
  min-width: 90px;
  width: 90px;
  min-height: 90px;
  height: 90px;
  margin: 10px 30px;
`;

function Flag(props: BadgesProps) {
  return (
    <svg width="30" height="32" viewBox="0 0 30 32" fill="none" xmlns="http://www.w3.org/2000/svg">
      <g filter="url(#filter0_d_3736_56289)">
        <path d="M4 27V3H26V27L15 24.3333L4 27Z" fill={props.color} />
      </g>
      <defs>
        <filter
          id="filter0_d_3736_56289"
          x="0"
          y="0"
          width="30"
          height="32"
          filterUnits="userSpaceOnUse"
          colorInterpolationFilters="sRGB"
        >
          <feFlood floodOpacity="0" result="BackgroundImageFix" />
          <feColorMatrix
            in="SourceAlpha"
            type="matrix"
            values="0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 127 0"
            result="hardAlpha"
          />
          <feOffset dy="1" />
          <feGaussianBlur stdDeviation="2" />
          <feComposite in2="hardAlpha" operator="out" />
          <feColorMatrix
            type="matrix"
            values="0 0 0 0 0.286275 0 0 0 0 0.788235 0 0 0 0 0.596078 0 0 0 0.5 0"
          />
          <feBlend mode="normal" in2="BackgroundImageFix" result="effect1_dropShadow_3736_56289" />
          <feBlend
            mode="normal"
            in="SourceGraphic"
            in2="effect1_dropShadow_3736_56289"
            result="shape"
          />
        </filter>
      </defs>
    </svg>
  );
}

function BadgeStatus(props: BadgesProps) {
  const { txid } = props;

  return (
    <div>
      <StatusText>
        {txid ? (
          <>
            <MaterialIcon icon="link" style={{ fontSize: 13 }} />
            <div style={{ marginLeft: 5, fontWeight: 500 }}>ON-CHAIN</div>
          </>
        ) : (
          <div
            style={{
              display: 'flex',
              fontSize: 11,
              alignItems: 'center',
              color: '#618AFF',
              letterSpacing: '0.3px'
            }}
          >
            OFF-CHAIN
          </div>
        )}
      </StatusText>
    </div>
  );
}

function Badges(props: BadgesProps) {
  const { main, ui } = useStores();
  const { badgeList, meInfo } = ui || {};

  const [balancesTxns, setBalancesTxns]: any = useState({});
  const [loading, setLoading] = useState(true);
  const [selectedBadge, setSelectedBadge]: any = useState(null);
  const [badgeToPush, setBadgeToPush]: any = useState(null);

  const [liquidAddress, setLiquidAddress]: any = useState('');
  const [memo, setMemo]: any = useState('');
  const [claiming, setClaiming]: any = useState(false);

  const isMobile = useIsMobile();
  const { person } = props;

  const thisIsMe = meInfo?.owner_pubkey === person?.owner_pubkey;

  useEffect(() => {
    main.getBadgeList();
  }, [main]);

  const getBadges = useCallback(
    async function () {
      setLoading(true);
      setSelectedBadge(null);
      setBadgeToPush(null);
      if (person?.owner_pubkey) {
        const b = await main.getBalances(person?.owner_pubkey);
        setBalancesTxns(b);
      }
      setLoading(false);
    },
    [main, person?.owner_pubkey]
  );

  useEffect(() => {
    getBadges();
  }, [getBadges]);

  async function claimBadge() {
    setClaiming(true);
    try {
      /*
      const body: ClaimOnLiquid = {
        amount: badgeToPush.balance,
        to: liquidAddress,
        asset: badgeToPush.id,
        memo: memo
      };*/

      //const token = await main.claimBadgeOnLiquid(body);
      // refresh badges
      getBadges();
    } catch (e) {
      console.log('e', e);
    }

    setClaiming(false);
  }

  function redirectToBlockstream(txId: string) {
    const el = document.createElement('a');
    el.target = '_blank';
    el.href = `https://blockstream.info/liquid/tx/${txId}`;
    el.click();
  }

  // metadata should be json to support badge details
  const topLevelBadges = balancesTxns?.balances?.map((b: any, i: number) => {
    const badgeDetails = badgeList?.find((f: any) => f.id === b.asset_id);
    // if early adopter badge
    let counter = '';
    const theseTxType = balancesTxns?.txs?.find((f: any) => f.asset_id === b.asset_id);
    const metadata = theseTxType?.metadata;
    const liquidTxId =
      balancesTxns?.txs?.find((f: any) => f.asset_id === b.asset_id && f.txid)?.txid || '';
    let flagColor = '#41c292';

    if (metadata && !isNaN(parseInt(metadata))) {
      counter = metadata;

      // flag colors
      // 1 - 100 #41c292
      // 100 - 500 #35c3cc
      // 500 - 1000 #628afd
      // > 1000 no show
      const intCount = parseInt(counter);
      if (intCount < 100) flagColor = '#41c292';
      else if (intCount < 500) flagColor = '#35c3cc';
      else if (intCount < 1001) flagColor = '#628afd';
    }

    const showFlag = counter && parseInt(counter) < 1001 ? true : false;

    const packedBadge = {
      ...b,
      ...badgeDetails,
      txid: liquidTxId,
      counter,
      metadata,
      deck: balancesTxns?.txs?.filter((f: any) => f.asset_id === b.asset_id) || []
    };

    // console.log('packedBadge', packedBadge)

    if (isMobile) {
      return (
        <BWrap
          key={`${i}badges`}
          isMobile={isMobile}
          onClick={() => {
            // setSelectedBadge(packedBadge)
          }}
        >
          <Img src={`${badgeDetails?.icon}`} isMobile={isMobile}>
            {showFlag && counter && (
              <div style={{ position: 'absolute', background: '#fff', bottom: -6, left: 12 }}>
                <Flag color={flagColor} />

                <div
                  style={{
                    fontSize: 10,
                    height: '90%',
                    position: 'absolute',
                    top: 0,
                    left: 0,
                    display: 'flex',
                    justifyContent: 'center',
                    width: '100%',
                    alignItems: 'center',
                    color: '#fff'
                  }}
                >
                  {counter}
                </div>
              </div>
            )}
          </Img>

          <div
            style={{
              display: 'flex',
              flexDirection: 'column',
              justifyContent: 'center',
              width: '100%',
              paddingRight: 20
            }}
          >
            <div style={{ width: 'auto' }}>
              <T>
                {badgeDetails?.name} {b.balance > 1 && `(${b.balance})`}
              </T>
              {badgeDetails?.description && <S isMobile={isMobile}>{badgeDetails?.description}</S>}
            </div>

            <Status
              style={{
                pointerEvents: !thisIsMe && !packedBadge.txid ? 'none' : 'auto',
                margin: 0,
                marginTop: 10
              }}
              onClick={() => {
                if (thisIsMe && !packedBadge.txid) {
                  //user on own badge, off-chain
                  setBadgeToPush(packedBadge);
                } else if (packedBadge.txid) {
                  //on-chain, click to see on blockstream
                  redirectToBlockstream(packedBadge.txid);
                }
              }}
            >
              <BadgeStatus {...packedBadge} />
            </Status>
          </div>
        </BWrap>
      );
    }

    // is desktop
    return (
      <BWrap
        key={`${i}badges`}
        isMobile={isMobile}
        onClick={() => {
          // setSelectedBadge(packedBadge)
        }}
      >
        <Img src={`${badgeDetails?.icon}`} isMobile={isMobile}>
          {showFlag && counter && (
            <div style={{ position: 'absolute', background: '#fff', bottom: -6, left: 12 }}>
              <Flag color={flagColor} />

              <div
                style={{
                  fontSize: 10,
                  height: '90%',
                  position: 'absolute',
                  top: 0,
                  left: 0,
                  display: 'flex',
                  justifyContent: 'center',
                  width: '100%',
                  alignItems: 'center',
                  color: '#fff'
                }}
              >
                {counter}
              </div>
            </div>
          )}
        </Img>

        <div style={{ width: '100%', minWidth: 160 }}>
          <T isMobile={isMobile}>
            {badgeDetails?.name} {b.balance > 1 && `(${b.balance})`}
          </T>
          {badgeDetails?.description && <S isMobile={isMobile}>{badgeDetails?.description}</S>}
        </div>

        <D>
          {counter && (
            <>
              {counter}
              {badgeDetails?.amount && (
                <div style={{ color: '#8E969C' }}>&nbsp;/&nbsp;{badgeDetails?.amount}</div>
              )}
            </>
          )}
        </D>

        <Status
          style={{ pointerEvents: !thisIsMe && !packedBadge.txid ? 'none' : 'auto' }}
          onClick={() => {
            if (thisIsMe && !packedBadge.txid) {
              //user on own badge, off-chain
              setBadgeToPush(packedBadge);
            } else if (packedBadge.txid) {
              //on-chain, click to see on blockstream
              redirectToBlockstream(packedBadge.txid);
            }
          }}
        >
          <BadgeStatus {...packedBadge} />
        </Status>
      </BWrap>
    );
  });

  return (
    <Wrap>
      <PageLoadSpinner show={loading} />
      {selectedBadge ? (
        <div style={{ width: '100%' }}>
          <Button
            color="noColor"
            leadingIcon="arrow_back"
            text="Back to all badges"
            onClick={() => setSelectedBadge(null)}
            style={{ marginBottom: 20 }}
          />

          {selectedBadge.deck?.map((badge: any, i: number) => (
            <BWrap
              key={`${i}badges`}
              isMobile={isMobile}
              style={{ height: 'auto', minHeight: 'auto', cursor: 'default' }}
            >
              <SmallImg src={`${selectedBadge?.icon}`} isMobile={isMobile} />
              <div
                style={{
                  width: '100%',
                  minWidth: 160,
                  display: 'flex',
                  flexDirection: 'column',
                  justifyContent: 'center'
                }}
              >
                <T isMobile={isMobile}>
                  {selectedBadge?.name}{' '}
                  {selectedBadge?.balance > 1 && `(${selectedBadge?.balance})`}
                </T>
                {selectedBadge?.counter && (
                  <D>
                    <Counter>
                      {selectedBadge?.counter} / {selectedBadge?.amount}
                    </Counter>
                  </D>
                )}

                <div style={{ marginTop: 20, width: '100%' }}>
                  {thisIsMe ? (
                    <div
                      style={{
                        display: 'flex',
                        flexDirection: 'column',
                        justifyContent: 'center',
                        width: '100%',
                        textAlign: 'center'
                      }}
                    >
                      <Divider />
                      <Button
                        style={{
                          margin: 0,
                          marginTop: 2,
                          padding: 0,
                          minHeight: 40,
                          border: 'none'
                        }}
                        color="link"
                        text="Claim on Liquid"
                        onClick={() => setBadgeToPush(badge)}
                      />
                    </div>
                  ) : (
                    <Status>
                      <StatusText>{'Off-chain'}</StatusText>
                    </Status>
                  )}
                </div>
              </div>
            </BWrap>
          ))}
        </div>
      ) : (
        topLevelBadges
      )}

      <Modal
        visible={badgeToPush ? true : false}
        close={() => {
          setBadgeToPush(null);
        }}
      >
        <div
          style={{
            padding: 20,
            height: 300,
            display: 'flex',
            flexDirection: 'column',
            justifyContent: 'center'
          }}
        >
          <>
            <TextInput
              style={{ width: 240 }}
              label={'Liquid Address'}
              value={liquidAddress}
              onChange={(e: any) => setLiquidAddress(e)}
            />

            <TextInput
              style={{ width: 240 }}
              label={'Memo (optional)'}
              value={memo}
              onChange={(e: any) => setMemo(e)}
            />

            <Button
              color="primary"
              text="Claim on Liquid"
              loading={claiming}
              disabled={!liquidAddress || claiming}
              onClick={() => claimBadge()}
            />
          </>
        </div>
      </Modal>
    </Wrap>
  );
}

export default observer(Badges);
