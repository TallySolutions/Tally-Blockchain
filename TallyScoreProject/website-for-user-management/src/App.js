import React, { useState } from 'react';
import RegistrationForm from './RegistrationForm';
import RegistrationTable from './RegistrationTable';

function App() {
  const [registrations, setRegistrations] = useState([]);

  const handleNewRegistration = (registration) => {
    const newRegistration = {
      ...registration,
      userid: generateUniqueId(),
    };
    setRegistrations([...registrations, newRegistration]);
  };

  const generateUniqueId = () => {
    return '_' + Math.random().toString(36).substr(2, 9);
  };

  return (
    <div className="App">
            
      <RegistrationTable registrations={registrations} />
    </div>
  );
}

export default App;
