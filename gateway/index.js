const { ApolloServer } = require('apollo-server')
const { ApolloGateway } = require('@apollo/gateway')

const gateway = new ApolloGateway({
    serviceList: [
        { name: 'accounts', url: 'https://accounts.summer-api.com/query' },
        { name: 'products', url: 'https://products.summer-api.com/query' },
        { name: 'reviews', url: 'https://reviews.summer-api.com/query' },
    ],
})

const server = new ApolloServer({
    gateway,
    subscriptions: false,
})

server
    .listen({
        port: process.env.PORT || 8080,
    })
    .then(({ url }) => {
        console.log(`🚀 Server ready at ${url}`)
    })
