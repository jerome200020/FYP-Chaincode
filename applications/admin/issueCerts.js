'use strict';

const fs = require('fs');
const yaml = require('js-yaml');
const { Wallets, Gateway } = require('fabric-network');

async function main() {
    const wallet = await Wallets.newFileSystemWallet('/identity/user/isabella/wallet');

    const gateway = new Gateway();

    try {
        const userName = 'isabella';

        let connectionProfile = yaml.safeLoad(fs.readFileSync('/gateway/connection-org2.yaml','utf-8'));

        let connectionOptions = {
            identity: userName,
            wallet: wallet,
            discovery: { enabled:true, asLocalhost: true }
        };

        console.log('Connect to Fabric gateway.');

        await gateway.connect(connectionProfile, connectionOptions);

        console.log('Use network channel: test.');

        const network = await gateway.getNetwork('test');

        console.log('Use smart contract cert.');

        const contract = await network.getContract('cert');

        console.log('Submit certs issue transaction.');

        const issueResponse = await contract.submitTransaction('InitLedger');

        console.log('Process issue transaction response.'+issueResponse);

        console.log('Transaction complete.');
    } catch (error) {

        console.log(`Error processing transaction. ${error}`);
        console.log(error.stack);

    } finally {
        console.log('Disconnect from Fabric gateway.');
        gateway.disconnect();
    }

}
main().then(() => {

    console.log('Issue program complete.');

}).catch((e) => {

    console.log('Issue program exception.');
    console.log(e);
    console.log(e.stack);
    process.exit(-1);

});