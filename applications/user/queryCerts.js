'use strict';

// Bring key classes into scope, most importantly Fabric SDK network class
const fs = require('fs');
const yaml = require('js-yaml');
const { Wallets, Gateway } = require('fabric-network');

async function main() {
    const wallet = await Wallets.newFileSystemWallet('../identity/user/shane/wallet');

    const gateway = new Gateway();

    try {
        const userName = 'shane';

        let connectionProfile = yaml.safeLoad(fs.readFileSync('../gateway/connection-org1.yaml', 'utf8'));

        let connectionOptions = {
            identity: userName,
            wallet: wallet,
            discovery: { enabled: true, asLocalhost: true }

        };

        console.log('Connect to Fabric gateway.');

        await gateway.connect(connectionProfile, connectionOptions);

        console.log('Use network channel: test.');

        const network = await gateway.getNetwork('test');

        console.log('Use cert smart contract.');

        const contract = await network.getContract('papercontract', 'org.papernet.commercialpaper');

        console.log('-----------------------------------------------------------------------------------------');
        console.log('****** Submitting cert queries ****** \n\n ');

        console.log('1. Query Academic Cert by ACertID ....');
        console.log('-----------------------------------------------------------------------------------------\n');
        let queryResponse = await contract.evaluateTransaction('ReadAcaCert', 'aCert1');
        let json = JSON.parse(queryResponse.toString());
        console.log(json);
        console.log('\n\n');
        console.log('\n  ReadAcaCert complete.');
        console.log('-----------------------------------------------------------------------------------------\n\n');

        console.log('2. Query Academic Cert of SWE1904873 by rich queries ....');
        console.log('-----------------------------------------------------------------------------------------\n');
        let queryResponse2 = await contract.evaluateTransaction('QueryAcaCertByStudentID', 'SWE1904873');
        json = JSON.parse(queryResponse2.toString());
        console.log(json);
        console.log('\n\n');
        console.log('\n  QueryAcaCertByStudentID complete.');
        console.log('-----------------------------------------------------------------------------------------\n\n');

        console.log('3. Query Extracurricular Cert by CCertID ....');
        console.log('-----------------------------------------------------------------------------------------\n');
        let queryResponse3 = await contract.evaluateTransaction('ReadCurrCert', 'cCert2');
        json = JSON.parse(queryResponse3.toString());
        console.log(json);
        console.log('\n\n');
        console.log('\n  ReadCurrCert complete.');
        console.log('-----------------------------------------------------------------------------------------\n\n');

        console.log('4. Query Academic Cert of SWE1904873 by rich queries ....');
        console.log('-----------------------------------------------------------------------------------------\n');
        let queryResponse4 = await contract.evaluateTransaction('QueryCurrCertByStudentID', 'SWE1904873');
        json = JSON.parse(queryResponse4.toString());
        console.log(json);
        console.log('\n\n');
        console.log('\n  QueryAcaCertByStudentID complete.');
        console.log('-----------------------------------------------------------------------------------------\n\n');

    } catch (error) {

        console.log(`Error processing transaction. ${error}`);
        console.log(error.stack);
    } finally {

        console.log('Disconnect from Fabric gateway.');
        gateway.disconnect();
    }
}

main().then(() => {

    console.log('Queryapp program complete.');

}).catch((e) => {

    console.log('Queryapp program exception.');
    console.log(e);
    console.log(e.stack);
    process.exit(-1);

});