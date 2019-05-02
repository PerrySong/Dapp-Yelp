# Dapp-Yelp

A decentralized review system built on blockchain.

# Why

In current day, restaurants in Yelp may generate fake reviews in order to have a good
rating in Yelp to attract more customer. In order to prevent this to happen, Yelp
developed algorithms to filter the “fake” reviews which cause a heated controversy.
Whereas, all the ratings are stored in Yelp’s database, and Yelp can manipulate the data
very easily for their interest.

Therefore, I propose to build a decentralized rating system based on blockchain that not
be control by any entity.

# Funtionalities
Bussiness informations and reviews are stored in the block's MPT. When a client want to post a bussiness or a review, the client need to pay the transaction fee to the miner.

1. A Client can post a bussiness or a review with a transaction fee.
2. A Minner can take a transaction, publish a block and take the transaction fee.
```
  // 1. Post /business
  Req:
  {
    business_name: “KO Ramen”,
    business_location: “San Francisco”,
    business_tag: “Restaurant”,
  }
  Res:
  {
    business_id: 1
  }
  // 2. Get /business?id=1
  Res:
  {
    average_rating: 3.4,
    reviews: [
    {
      rating: 4,
      comment: “Good restaurant”
    },
  {
    rating: 1,
    comment: “Bad restaurant”
  }]
  }
  3. Post /rate?id=1
  Req:
  {
    rating: 3,
    comment: “not bad”
  }
```

# How does the reviews get stored in blockchain?
Both reviews and business information are stored in the MPT of the block

# Who pay the transaction fee?
Users will have unlimited starting gas for now. And they will attach transaction fee when they submit a transaction.

# How it works
When a user create a transaction, the user will send the transaction to all its peers and peers will forward the transaction to their peers, so on and so forth.

### Local transaction pool: 
When a miner recieve a transaction, it will store the transaction in the local transaction pool.
When a miner tries to create a block, it polls a pending transaction (Check if the transaction is already in the chain) from local transaction pool, and tries to solve PoW. 
Once the miner succeed, it forward the new block to its peers.

### Client
Client is a node int 

# Milestones
1. Build data structure for business info  - [x]
2. Build data structure for transaction  - [x]
3. Implement integrity using public / private key  - [x]
4. Implement Post business
5. Implement Post business's review
6. Implement Get business's review
7. Build interactive front end 

