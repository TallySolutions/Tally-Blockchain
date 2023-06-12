import React from 'react';

function RegistrationTable({ registrations, setRegistrations }) {
  const handleIncrementClick = (registration) => {
    const updateStruct = {
      PAN: registration.PAN,
      IncVal: "10"
    };

    const forRequest = {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(updateStruct),
    };

    fetch('http://43.204.226.103:8080/TallyScoreProject/increaseTallyScore', forRequest)
      .then(response => {
        if (response.ok) {
          return response.json();
        } else {
          throw new Error('Error in registration.');
        }
      })
      .then(data => {
        const updatedRegistration = { ...registration, Score: data.score };
        const updatedRegistrations = registrations.map(reg =>
          reg.PAN === registration.PAN ? updatedRegistration : reg
        );
        setRegistrations(updatedRegistrations);
        console.log("businessCertDetails:", JSON.stringify(data));
      })
      .catch(error => {
        console.error('Error:', error);
      });
  };

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
          <th>Action</th>
        </tr>
      </thead>
      <tbody>
        {registrations.map((registration) => (
          <tr key={registration.PAN}>
            <td>{registration.PAN}</td>
            <td>{registration.Name}</td>
            <td>{registration.PhoneNo}</td>
            <td>{registration.Address}</td>
            <td>{registration.LicenseType}</td>
            <td><div id ="scorerow">
                    <button id="scorebutton">-</button>
                    {registration.Score}
                    <button id="scorebutton" onClick={(e) => handleIncrementClick(registration, e)}>+</button>
              </div>
            </td>
            <td>{registration.status}</td>
          </tr>
        ))}
      </tbody>
    </table>
  );
}

export default RegistrationTable;
