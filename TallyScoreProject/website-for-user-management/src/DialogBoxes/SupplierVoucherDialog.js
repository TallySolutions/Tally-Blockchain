import React, { useState } from 'react';
import ApproveVoucherDialogBox from './ApproveVoucherDialogBox';
import RejectVoucherDialogBox from './RejectVoucherDialogBox';

function SupplierVoucherDialog({ onClose }) {
  const [showApproveVoucherDialog, setShowApproveVoucherDialog] = useState(false);
  const [showRejectVoucherDialog, setShowRejectVoucherDialog] = useState(false);

  const handleButtonClick = (action) => {
    console.log('Supplier Voucher Button Clicked:', action);
    if (action === 'Approve Voucher') {
      setShowApproveVoucherDialog(true);
    } else if (action === 'Reject Voucher') {
      setShowRejectVoucherDialog(true);
    } else {
      onClose();
    }
  };

  return (
    <div className="voucher-dialog">
      {!showApproveVoucherDialog && !showRejectVoucherDialog && (
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
    </div>
  );
}

export default SupplierVoucherDialog;
