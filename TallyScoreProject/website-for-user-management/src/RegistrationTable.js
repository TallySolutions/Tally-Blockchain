import React from 'react';

function RegistrationTable({ registrations }) {
  return (
    <table>
      <thead>
        <tr>
          <th>PAN</th>
          <th>Name</th>
          <th>Phone No.</th>
          <th>Address</th>
          <th>License</th>
          <th>Score</th>
          <th>Status</th>
        </tr>
      </thead>
      <tbody>
        {registrations.map((registration) => (
          <tr key={registration.pan}>
            <td>{registration.pan}</td>
            <td>{registration.name}</td>
            <td>{registration.phonenumber}</td>
            <td>{registration.address}</td>
            <td>{registration.license}</td>
            <td>{registration.score}</td>
            <td>{registration.status}</td>
          </tr>
        ))}
      </tbody>
    </table>
  );
}

export default RegistrationTable;
