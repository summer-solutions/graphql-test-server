const { ApolloServer } = require('apollo-server');
const { ApolloGateway } = require("@apollo/gateway");

const gateway = new ApolloGateway({
    serviceList: [
        { name: 'accounts', url: 'https://accounts.supper-api.com/query' },
        { name: 'products', url: 'https://products.supper-api.com/query' },
        { name: 'reviews', url: 'https://reviews.supper-api.com/query' }
    ],
});

const server = new ApolloServer({
    gateway,
    subscriptions: false,
});

server.listen().then(({ url }) => {
    console.log(`🚀 Server ready at ${url}`);
});
