import React, { useState } from 'react';
import { useStores } from '../../store';
import MaterialIcon from '@material/react-material-icon';
import { EuiLoadingSpinner } from '@elastic/eui';

export default function Heart(props) {
  const { main } = useStores();
  const [loading, setLoading] = useState(false);
  const selected = false;
  function clickIt(e) {
    setLoading(true);
    // if (selected) {
    //     main.deleteFavorite()
    // } else {
    //     main.addFavorite()
    // }
    setLoading(false);
  }

  if (loading) {
    return <EuiLoadingSpinner />;
  } else {
    return (
      <MaterialIcon
        onClick={clickIt}
        style={{ color: '#B0B7BC', cursor: 'pointer', userSelect: 'none' }}
        icon={selected ? 'favorite' : 'favorite_outline'}
      />
    );
  }
}
