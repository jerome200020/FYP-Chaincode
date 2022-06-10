'use strict';

const FabricCAServices = require('fabric-ca-client');
const { Wallets } = require('fabric-network');
const fs = require('fs');
const yaml = require('js-yaml');
const path = require('path');

async function main(){
    try {
        let connectionProfile = yaml.safeLoad(fs.readFileSync('/gateway/connection-org2.yaml','utf-8'));

        const caInfo = connectionProfile.certificateAuthorities['ca.org2.example.com'];
        const caTLSCACerts = caInfo.caTLSCACerts.pem;
        const ca = new FabricCAServices(caInfo.url, { trustedRoosts: caTLSCACerts, verify: false }, caInfo.caName);

        const walletPath = path.join(process.cwd(), '/identify/user/isabella/wallet');
        const wallet = await Wallets.newFileSystemWallet(walletPath);
        console.log(`Wallet path: ${walletPath}`);

        const userExists = await wallet.get('isabella');
        if (userExists) {
            console.log('An identity for the client "user1" already exists in the wallet');
            return;
        }

        const enrollment = await ca.enroll({ enrollmentID: 'user1', enrollmentSecret: 'user1pw' });
        const x509Identity = {
            credentials: {
                certificate: enrollment.certificate,
                privateKey: enrollment.key.toBytes(),
            },
            mspId: 'Org2MSP',
            type: 'X.509',
        };
        await wallet.put('isabella', x509Identity);
        console.log('Successfully enrolled client user "isabella" and imported it into the wallet');

    } catch (error) {
        console.error(`Failed to enroll client user "isabella": ${error}`);
        process.exit(1);
    }
}

main();