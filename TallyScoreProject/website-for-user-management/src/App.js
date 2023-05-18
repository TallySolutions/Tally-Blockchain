import Form from './Form';
import React, { useState } from 'react';

function App() {
  const [registrations, setRegistrations] = useState([]);

  const handleNewRegistration = (newRegistration) => {
    setRegistrations([...registrations, newRegistration]);
  };

  return (
    <div className="App">
      <Form onNewRegistration={handleNewRegistration} />
      <table>
        <thead>
          <tr>
            <th>Name</th>
            <th>PAN</th>
            <th>License</th>
            <th>Score</th>
            <th>Status</th>
          </tr>
        </thead>
        <tbody>
          {registrations.map((registration, index) => (
            <tr key={index}>
              <td>{registration.name}</td>
              <td>{registration.pan}</td>
              <td>{registration.license}</td>
              <td>{registration.score}</td>
              <td>{registration.status}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}

export default App;
