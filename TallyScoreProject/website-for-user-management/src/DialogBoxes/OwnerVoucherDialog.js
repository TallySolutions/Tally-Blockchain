import React, { useState } from 'react';
import NewVoucherDialog from './NewVoucherDialog';
import CancelVoucherDialog from './CancelVoucherDialog';

function OwnerVoucherDialog({ onClose }) {
  const [showNewVoucherDialog, setShowNewVoucherDialog] = useState(false);
  const [showCancelVoucherDialog, setShowCancelVoucherDialog] = useState(false);

  const handleButtonClick = (action) => {
    console.log('Owner Voucher Button Clicked:', action);
    if (action === 'New Voucher') {
      setShowNewVoucherDialog(true);
    } else if (action === 'Cancel Voucher') {
      setShowCancelVoucherDialog(true);
    } else {
      onClose();
    }
  };

  return (
    <div className="voucher-dialog">
      {!showNewVoucherDialog && !showCancelVoucherDialog && (
        <>
          <button
            className="close-dialog-button"
            onClick={() => handleButtonClick('Close Dialog')}
          >
            Back
          </button>
          <div className="voucher-dialog-buttons">
            <button
              className="dialog-buttons"
              onClick={() => handleButtonClick('New Voucher')}
            >
              New Voucher
            </button>
            <button
              className="dialog-buttons"
              onClick={() => handleButtonClick('Cancel Voucher')}
            >
              Cancel Voucher
            </button>
            <button
              className="dialog-buttons"
              onClick={() => handleButtonClick('Update Asset Voucher')}
            >
              Update Asset Voucher
            </button>
            <button
              className="dialog-buttons"
              onClick={() =>
                handleButtonClick('List of all Vouchers as Owner')
              }
            >
              List of all Vouchers as Owner
            </button>
          </div>
        </>
      )}
      {showNewVoucherDialog && (
        <NewVoucherDialog onClose={onClose} />
      )}
      {showCancelVoucherDialog && (
        <CancelVoucherDialog onClose={onClose} />
      )}
    </div>
  );
}

export default OwnerVoucherDialog;
