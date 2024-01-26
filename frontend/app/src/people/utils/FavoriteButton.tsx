import { EuiLoadingSpinner } from '@elastic/eui';
import MaterialIcon from '@material/react-material-icon';
import React, { useState } from 'react';

export default function Heart() {
  const [loading, setLoading] = useState(false);
  const selected = false;
  function clickIt() {
    setLoading(true);
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
