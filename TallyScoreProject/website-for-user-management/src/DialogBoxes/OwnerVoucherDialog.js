import React, { useState } from 'react';
import NewVoucherDialog from './NewVoucherDialog';
import CancelVoucherDialog from './CancelVoucherDialog';
import UpdateVoucherDialog from './UpdateVoucherDialog';
import ListOwnerVouchersDialog from './ListOwnerVouchersDialog';

function OwnerVoucherDialog({ onClose }) {
  const [showNewVoucherDialog, setShowNewVoucherDialog] = useState(false);
  const [showCancelVoucherDialog, setShowCancelVoucherDialog] = useState(false);
  const [showUpdateVoucherDialog, setShowUpdateVoucherDialog] = useState(false);
  const [showListOwnerVouchersDialog, setShowListOwnerVouchersDialog] = useState(false);

  const handleButtonClick = (action) => {
    console.log('Owner Voucher Button Clicked:', action);
    if (action === 'New Voucher') {
      setShowNewVoucherDialog(true);
    } else if (action === 'Cancel Voucher') {
      setShowCancelVoucherDialog(true);
    } else if (action === 'Update Voucher') {
      setShowUpdateVoucherDialog(true);
    } else if (action === 'List of all Vouchers as Owner') {
      setShowListOwnerVouchersDialog(true);
    } else {
      onClose();
    }
  };

  return (
    <div className="voucher-dialog">
      {!showNewVoucherDialog && !showCancelVoucherDialog && !showUpdateVoucherDialog && !showListOwnerVouchersDialog && (
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
              onClick={() => handleButtonClick('Update Voucher')}
            >
              Update Voucher
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
      {showUpdateVoucherDialog && (
        <UpdateVoucherDialog onClose={onClose} />
      )}
      {showListOwnerVouchersDialog && (
        <ListOwnerVouchersDialog onClose={() => handleButtonClick('Back')} />
      )}
    </div>
  );
}

export default OwnerVoucherDialog;
