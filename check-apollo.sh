set -e
cd services/accounts
./check-apollo.sh
cd ../products
./check-apollo.sh
cd ../reviews
./check-apollo.sh