import React, { useState } from 'react';

function Form({ onNewRegistration }) {
  const [name, setName] = useState('');
  const [pan, setPan] = useState('');
  const [license, setLicense] = useState('');
  const [score, setScore] = useState('');
  const [status, setStatus] = useState('');

  const handleFormSubmit = (e) => {
    e.preventDefault();
    const newRegistration = {
      name: name,
      pan: pan,
      license: license,
      score: score,
      status: status
    };
    onNewRegistration(newRegistration);
    setName('');
    setPan('');
    setLicense('');
    setScore('');
    setStatus('');
  };

  return (
    <form onSubmit={handleFormSubmit}>
      <label>
        Name:
        <input type="text" value={name} onChange={(e) => setName(e.target.value)} />
      </label>
      <label>
        PAN:
        <input type="text" value={pan} onChange={(e) => setPan(e.target.value)} />
      </label>
      <label>
        License:
        <input type="text" value={license} onChange={(e) => setLicense(e.target.value)} />
      </label>
      <label>
        Score:
        <input type="text" value={score} onChange={(e) => setScore(e.target.value)} />
      </label>
      <label>
        Status:
        <input type="text" value={status} onChange={(e) => setStatus(e.target.value)} />
      </label>
      <button type="submit">Submit</button>
    </form>
  );
}

export default Form;
