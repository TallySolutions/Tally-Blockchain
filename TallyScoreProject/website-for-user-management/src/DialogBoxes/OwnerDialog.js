import React from 'react';

function OwnerDialog() {
  const handleNewVoucher = () => {
   
  };

  const handleCancelVoucher = () => {
    
  };

  const handleUpdateAssetVoucher = () => {
   
  };

  const handleListOfVouchers = () => {
    
  };

  return (
    <div className="owner-supplier-dialog">
      <button onClick={handleNewVoucher}>New Voucher</button>
      <button onClick={handleCancelVoucher}>Cancel Voucher</button>
      <button onClick={handleUpdateAssetVoucher}>Update Asset Voucher</button>
      <button onClick={handleListOfVouchers}>List of all Vouchers as Owner</button>
    </div>
  );
}

export default OwnerDialog;
