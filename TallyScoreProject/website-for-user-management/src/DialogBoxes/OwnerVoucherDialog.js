import React, { useState } from 'react';
import NewVoucherDialog from './NewVoucherDialog';
import CancelVoucherDialog from './CancelVoucherDialog';
import UpdateVoucherDialog from './UpdateVoucherDialog';
import ListOwnerVouchersDialog from './ListOwnerVouchersDialog';

function OwnerVoucherDialog({ onClose }) {
  const [activeDialog, setActiveDialog] = useState(null);

  const handleDialogOpen = (dialogName) => {
    setActiveDialog(dialogName);
  };

  const handleDialogClose = () => {
    setActiveDialog(null);
  };

  return (
    <div>
      {activeDialog === null ? (
        <>
          <button className="close-dialog-button" onClick={onClose}>
            Back
          </button>
          <div className="voucher-dialog-buttons">
            <button
              className="dialog-buttons"
              onClick={() => handleDialogOpen('NewVoucherDialog')}
            >
              New Voucher
            </button>
            <button
              className="dialog-buttons"
              onClick={() => handleDialogOpen('CancelVoucherDialog')}
            >
              Cancel Voucher
            </button>
            <button
              className="dialog-buttons"
              onClick={() => handleDialogOpen('UpdateVoucherDialog')}
            >
              Update Voucher
            </button>
            <button
              className="dialog-buttons"
              onClick={() => handleDialogOpen('ListOwnerVouchersDialog')}
            >
              List of all Vouchers as Owner
            </button>
          </div>
        </>
      ) : (
        <>
          {activeDialog === 'NewVoucherDialog' && (
            <NewVoucherDialog onClose={handleDialogClose} />
          )}
          {activeDialog === 'CancelVoucherDialog' && (
            <CancelVoucherDialog onClose={handleDialogClose} />
          )}
          {activeDialog === 'UpdateVoucherDialog' && (
            <UpdateVoucherDialog onClose={handleDialogClose} />
          )}
          {activeDialog === 'ListOwnerVouchersDialog' && (
            <ListOwnerVouchersDialog onClose={handleDialogClose} />
          )}
        </>
      )}
    </div>
  );
}

export default OwnerVoucherDialog;
