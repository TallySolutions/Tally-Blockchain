import React from 'react';

function SupplierDialog() {
  const handleApproveVoucher = () => {
    
  };

  const handleRejectVoucher = () => {
   
  };

  const handleSendBackVoucher = () => {
    
  };

  const handleListOfVouchers = () => {
    
  };

  return (
    <div className="owner-supplier-dialog">
      <button onClick={handleApproveVoucher}>Approve Voucher</button>
      <button onClick={handleRejectVoucher}>Reject Voucher</button>
      <button onClick={handleSendBackVoucher}>Send Back Voucher</button>
      <button onClick={handleListOfVouchers}>List of Vouchers as Supplier</button>
    </div>
  );
}

export default SupplierDialog;
