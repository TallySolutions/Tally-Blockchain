import React from 'react';

                                                    // NOTE: CODE HAS TO BE EDITED TO INCORPORATE ENDPOINTS

function ListOwnerVouchersDialog({ onClose }) {
  const dummyVouchersList = [
    { voucherid: 'voucher 1', supplierid: 'DUMMY PANuserCctest2' },
    { voucherid: 'voucher 2', supplierid: 'DUMMY PANuserCctest3' },
    { voucherid: 'voucher 3', supplierid: 'DUMMY PANuserCctest2' },
  ];

  return (
    <div className="list-vouchers-dialog">
      <h3 className="dialog-title">List of Vouchers as Owner</h3>
      <table className="vouchers-table">
        <thead>
          <tr>
            <th>Voucher ID</th>
            <th>Supplier ID</th>
          </tr>
        </thead>
        <tbody>
          {dummyVouchersList.map((voucher) => (
            <tr key={voucher.voucherid}>
              <td>{voucher.voucherid}</td>
              <td>{voucher.supplierid}</td>
            </tr>
          ))}
        </tbody>
      </table>
      <button className="close-list-dialog-button" onClick={onClose}>
        Close
      </button>
    </div>
  );
}

export default ListOwnerVouchersDialog;
