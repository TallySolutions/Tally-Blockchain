import React from 'react';

function SupplierVoucherDialog({ onClose }) {
  const handleButtonClick = (action) => {
    console.log('Supplier Voucher Button Clicked:', action);
    onClose();
  };

  return (

    <div className='voucher-dialog'>
        <button className='close-dialog-button' onClick={() => handleButtonClick('Close Dialog')}>Back</button>  
        <div className="voucher-dialog-buttons">
            <button className="dialog-buttons" onClick={() => handleButtonClick('Approve Voucher')}>Approve Voucher</button>
            <button className="dialog-buttons" onClick={() => handleButtonClick('Reject Voucher')}>Reject Voucher</button>
            <button className="dialog-buttons" onClick={() => handleButtonClick('Send Back Voucher')}>Send Back Voucher</button>
            <button className="dialog-buttons" onClick={() => handleButtonClick('List of Vouchers as Supplier')}>List of Vouchers as Supplier</button>
        </div>
    </div>

  );
}

export default SupplierVoucherDialog;
