const { ApolloServer } = require('apollo-server')
const { ApolloGateway } = require('@apollo/gateway')
const { InMemoryLRUCache } = require('apollo-server-caching')

const gateway = new ApolloGateway({
    serviceList: [
        { name: 'accounts', url: 'http://localhost:4001/query' },
        { name: 'products', url: 'http://localhost:4002/query' },
        { name: 'reviews', url: 'http://localhost:4003/query' },
    ],
})

const server = new ApolloServer({
    gateway,
    cache: new InMemoryLRUCache({maxSize: 1000}),
    subscriptions: false,
})

server
    .listen({
        port: process.env.PORT || 8080,
    })
    .then(({ url }) => {
        console.log(`🚀 Server ready at ${url}`)
    })
