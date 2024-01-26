import React, { FC, useState } from 'react';
import styled from 'styled-components';
import { observer } from 'mobx-react-lite';
import { colors } from '../../../config/colors';
import { useIsMobile } from '../../../hooks';
import IconButton from '../../../components/common/IconButton2';
import { useStores } from '../../../store';
import StartUpModal from '../../utils/StartUpModal';
import { PostModal, PostModalProps } from './PostModal';

const color = colors['light'];
const StyledIconButton = styled(IconButton)`
  color: ${color.pureWhite};
  font-size: 16px;
  font-weight: 600;
  text-decoration: none;
`;

const iconStyle = {
  fontSize: '16px',
  fontWeight: 400
};

interface Props extends Omit<PostModalProps, 'onClose' | 'isOpen'> {
  title?: string;
  buttonProps?: {
    endingIcon?: string;
    leadingIcon?: string;
    color?: 'primary' | 'secondary';
  };
}

const mapBtnColorProps = {
  primary: {
    color: 'success',
    hovercolor: color.button_primary.hover,
    activecolor: color.button_primary.active,
    shadowcolor: color.button_primary.shadow
  },
  secondary: {
    color: 'primary',
    hovercolor: color.button_secondary.hover,
    activecolor: color.button_secondary.active,
    shadowcolor: color.button_secondary.shadow
  }
};

export const PostBounty: FC<Props> = observer(
  ({
    title = 'Post a Bounty',
    buttonProps = {
      color: 'primary'
    },
    ...modalProps
  }: any) => {
    const { ui } = useStores();
    const [isOpenPostModal, setIsOpenPostModal] = useState(false);
    const [isOpenStartUpModel, setIsOpenStartupModal] = useState(false);

    const isMobile = useIsMobile();
    const showSignIn = () => {
      if (isMobile) {
        ui.setShowSignIn(true);
        return;
      }
      setIsOpenStartupModal(true);
    };

    const clickHandler = () => {
      if (ui.meInfo && ui.meInfo?.owner_alias) {
        setIsOpenPostModal(true);
      } else {
        showSignIn();
      }
    };

    const icon = (() => {
      if (buttonProps.endingIcon && buttonProps.leadingIcon) {
        return { leadingIcon: buttonProps.leadingIcon };
      }
      if (buttonProps.leadingIcon) {
        return { leadingIcon: buttonProps.leadingIcon };
      }
      if (buttonProps.endingIcon) {
        return { endingIcon: buttonProps.endingIcon };
      }
      return { endingIcon: 'add' };
    })();

    return (
      <>
        <StyledIconButton
          {...icon}
          {...mapBtnColorProps[buttonProps.color || 'primary']}
          text={title}
          width={204}
          height={isMobile ? 36 : 48}
          iconStyle={iconStyle}
          onClick={clickHandler}
        />
        <PostModal
          isOpen={isOpenPostModal}
          onClose={() => setIsOpenPostModal(false)}
          {...modalProps}
        />
        {isOpenStartUpModel && (
          <StartUpModal
            closeModal={() => setIsOpenStartupModal(false)}
            dataObject={'createWork'}
            buttonColor={'success'}
          />
        )}
      </>
    );
  }
);
