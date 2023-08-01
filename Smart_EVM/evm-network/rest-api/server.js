const express = require('express');
const { spawn } = require('child_process');
const app = express();

app.get('/create', (req, res) => {
  channel_name = req.query.channel
  console.log(`Received channel name = "${channel_name}"`);
  if ( channel_name === undefined){
    res.status = 500;
    res.send('Channel name parameter not passed');
  }else if ( channel_name === ""){
    res.status = 500;
    res.send('Channel name parameter value is empty');
  }else{
    res.write('Creating channel "' + channel_name + '" ...\n');
        const cc = spawn('./network.sh', ['createChannel', '-c', channel_name], {cwd : '..' });

        cc.stdout.on('data', (data) => {
            res.write(`stdout: ${data}`);
        });

        cc.stderr.on('data', (data) => {
            res.write(`stderr: ${data}`);
        });

        cc.on('close', (code) => {
            res.write(`child process exited with code ${code}\n`);
            if ( code != 0 ){
                res.status(500)
            }
            res.end()
        });

        cc.on('error', (code) => {
            res.write(`child process exceution error ${code}\n`);
        });

  }   
});

// Listen to the App Engine-specified port, or 8080 otherwise
const PORT = process.env.PORT || 8080;
app.listen(PORT, () => {
  console.log(`Server listening on port ${PORT}...`);
});