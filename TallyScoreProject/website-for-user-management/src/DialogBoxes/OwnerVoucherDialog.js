import React from 'react';

function OwnerVoucherDialog({ onClose }) {
  const handleButtonClick = (action) => {
    console.log('Owner Voucher Button Clicked:', action);
    onClose();
  };

  return (
    <div className='voucher-dialog'>
      <button className='close-dialog-button' onClick={() => handleButtonClick('Close Dialog')}>Back</button>  
      <div className="voucher-dialog-buttons">
          <button className="dialog-buttons" onClick={() => handleButtonClick('New Voucher')}>New Voucher</button>
          <button className="dialog-buttons" onClick={() => handleButtonClick('Cancel Voucher')}>Cancel Voucher</button>
          <button className="dialog-buttons" onClick={() => handleButtonClick('Update Asset Voucher')}>Update Asset Voucher</button>
          <button className="dialog-buttons" onClick={() => handleButtonClick('List of all Vouchers as Owner')}>List of all Vouchers as Owner</button>
      </div>
    </div>
    
  );
}

export default OwnerVoucherDialog;
