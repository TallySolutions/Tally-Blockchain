import React, { useState } from 'react';
import ApproveVoucherDialogBox from './ApproveVoucherDialogBox';
import RejectVoucherDialogBox from './RejectVoucherDialogBox';
import SendBackVoucherDialogBox from './SendBackVoucherDialogBox';
import ListSupplierVouchersDialogBox from './ListSupplierVouchersDialogBox';

function SupplierVoucherDialog({ onClose , pan}) {
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
              onClick={() => handleDialogOpen('ApproveVoucherDialogBox')}
            >
              Approve Voucher
            </button>
            <button
              className="dialog-buttons"
              onClick={() => handleDialogOpen('RejectVoucherDialogBox')}
            >
              Reject Voucher
            </button>
            <button
              className="dialog-buttons"
              onClick={() => handleDialogOpen('SendBackVoucherDialogBox')}
            >
              Send Back Voucher
            </button>
            <button
              className="dialog-buttons"
              onClick={() => handleDialogOpen('ListSupplierVouchersDialogBox')}
            >
              List of all Vouchers as Supplier
            </button>
          </div>
        </>
      ) : (
        <>
          {activeDialog === 'ApproveVoucherDialogBox' && (
            <ApproveVoucherDialogBox onClose={handleDialogClose} pan={pan} />
          )}
          {activeDialog === 'RejectVoucherDialogBox' && (
            <RejectVoucherDialogBox onClose={handleDialogClose} pan={pan} />
          )}
          {activeDialog === 'SendBackVoucherDialogBox' && (
            <SendBackVoucherDialogBox onClose={handleDialogClose} pan={pan} />
          )}
          {activeDialog === 'ListSupplierVouchersDialogBox' && (
            <ListSupplierVouchersDialogBox onClose={handleDialogClose} pan={pan}/>
          )}
        </>
      )}
    </div>
  );
}

export default SupplierVoucherDialog;
