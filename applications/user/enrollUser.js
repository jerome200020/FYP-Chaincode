'use strict';

const FabricCAServices = require('fabric-ca-client');
const { Wallets } = require('fabric-network');
const fs = require('fs');
const yaml = require('js-yaml');
const path = require('path');

async function main() {
    try {
        let connectionProfile = yaml.safeLoad(fs.readFileSync('/gateway/connection-org1.yaml','utf-8'));

        const caInfo = connectionProfile.certificateAuthorities['ca.org1.example.com'];
        const caTLSCACerts = caInfo.tlsCACerts.pem;
        const ca = new FabricCAServices(caInfo.url, { trustedRoots: caTLSCACerts, verify: false }, caInfo.caName);

        const walletPath = path.join(process.cwd(), '/identity/user/shane/wallet');
        const wallet = await Wallets.newFileSystemWallet(walletPath);
        console.log(`Wallet path: ${walletPath}`);

        const userExists = await wallet.get('shane');
        if (userExists) {
            console.log('An identity for the client user "shane" already exists in the wallet');
            return;
        }

        const enrollment = await ca.enroll({ enrollmentID: 'user1', enrollmentSecret: 'user1pw' });
        const x509Identity = {
            credentials: {
                certificate: enrollment.certificate,
                privateKey: enrollment.key.toBytes(),
            },
            mspId: 'Org1MSP',
            type: 'X.509',
        };
        await wallet.put('shane', x509Identity);
        console.log('Successfully enrolled client user "shane" and imported it into the wallet');

    } catch {
        console.error(`Failed to enroll client shane "shane": ${error}`);
        process.exit(1);
    }
}

main();