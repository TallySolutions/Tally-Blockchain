import React, { useState } from 'react';
import ApproveVoucherDialogBox from './ApproveVoucherDialogBox';
import RejectVoucherDialogBox from './RejectVoucherDialogBox';
import SendBackVoucherDialogBox from './SendBackVoucherDialogBox';
import ListSupplierVouchersDialog from './ListSupplierVouchersDialog';

function SupplierVoucherDialog({ onClose }) {
  const [showApproveVoucherDialog, setShowApproveVoucherDialog] = useState(false);
  const [showRejectVoucherDialog, setShowRejectVoucherDialog] = useState(false);
  const [showSendBackVoucherDialog, setShowSendBackVoucherDialog] = useState(false);
  const [showListVouchersDialog, setShowListVouchersDialog] = useState(false);

  const handleButtonClick = (action) => {
    console.log('Supplier Voucher Button Clicked:', action);
    if (action === 'Approve Voucher') {
      setShowApproveVoucherDialog(true);
    } else if (action === 'Reject Voucher') {
      setShowRejectVoucherDialog(true);
    } else if (action === 'Send Back Voucher') {
      setShowSendBackVoucherDialog(true);
    } else if(action === 'List of Vouchers as Supplier'){
      setShowListVouchersDialog(true);
    } else {
      onClose();
    }
  };

  return (
    <div className="voucher-dialog">
      {!showApproveVoucherDialog && !showRejectVoucherDialog && !showSendBackVoucherDialog && !showListVouchersDialog && (
        <>
          <button className="close-dialog-button" onClick={() => handleButtonClick('Close Dialog')}>
            Back
          </button>
          <div className="voucher-dialog-buttons">
            <button className="dialog-buttons" onClick={() => handleButtonClick('Approve Voucher')}>
              Approve Voucher
            </button>
            <button className="dialog-buttons" onClick={() => handleButtonClick('Reject Voucher')}>
              Reject Voucher
            </button>
            <button className="dialog-buttons" onClick={() => handleButtonClick('Send Back Voucher')}>
              Send Back Voucher
            </button>
            <button className="dialog-buttons" onClick={() => handleButtonClick('List of Vouchers as Supplier')}>
              List of Vouchers as Supplier
            </button>
          </div>
        </>
      )}
      {showApproveVoucherDialog && (
        <ApproveVoucherDialogBox onClose={() => setShowApproveVoucherDialog(false)} />
      )}
      {showRejectVoucherDialog && (
        <RejectVoucherDialogBox onClose={() => setShowRejectVoucherDialog(false)} />
      )}
      {showSendBackVoucherDialog && (
        <SendBackVoucherDialogBox onClose={() => setShowSendBackVoucherDialog(false)} />
      )}
      {showListVouchersDialog && (
        <ListSupplierVouchersDialog onClose={() => setShowListVouchersDialog(false)} />
      )}
    </div>
  );
}

export default SupplierVoucherDialog;
