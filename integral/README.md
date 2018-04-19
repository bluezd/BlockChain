# init telecom with 100 integral
`curl -H "Content-type:application/json" -X POST http://localhost:5110/bcsgw/rest/v1/transaction/invocation -d '{"channel":"integral.supply.chain","chaincode":"integralTrace","method":"initIntegral","args":["xiaoou", "telecom", "100", ""],"chaincodeVer":"v1"}'`

# query xiaoou telcom integral
`curl -H "Content-type:application/json" -X POST http://localhost:5110/bcsgw/rest/v1/transaction/query -d '{"channel":"integral.supply.chain","chaincode":"integralTrace","method":"queryIntegral","args":["xiaoou", "telecom"],"chaincodeVer":"v1"}'`

# query history xiaoou telcom integral
`curl -H "Content-type:application/json" -X POST http://localhost:5110/bcsgw/rest/v1/transaction/query -d '{"channel":"integral.supply.chain","chaincode":"integralTrace","method":"queryHistoryIntegral","args":["xiaoou", "telecom"],"chaincodeVer":"v1"}'`

`curl -H "Content-type:application/json" -X POST http://localhost:5100/bcsgw/rest/v1/transaction/query -d '{"channel":"integral.supply.chain","chaincode":"integralTrace","method":"queryHistoryIntegral","args":["xiaoou", "telecom"],"chaincodeVer":"v1"}'`

`curl -H "Content-type:application/json" -X POST http://localhost:5120/bcsgw/rest/v1/transaction/query -d '{"channel":"integral.supply.chain","chaincode":"integralTrace","method":"queryHistoryIntegral","args":["xiaoou", "telecom"],"chaincodeVer":"v1"}'`

# add 100 integral to xiaoou telcom
`curl -H "Content-type:application/json" -X POST http://localhost:5110/bcsgw/rest/v1/transaction/invocation -d '{"channel":"integral.supply.chain","chaincode":"integralTrace","method":"addIntegral","args":["xiaoou", "telecom", "50"],"chaincodeVer":"v1"}'`

# convert 50 integral from telcom to bank
`curl -H "Content-type:application/json" -X POST http://localhost:5110/bcsgw/rest/v1/transaction/invocation -d '{"channel":"integral.supply.chain","chaincode":"integralTrace","method":"convertIntegral","args":["xiaoou", "telecom", "bank", "50"],"chaincodeVer":"v1"}'`

`curl -H "Content-type:application/json" -X POST http://localhost:5110/bcsgw/rest/v1/transaction/invocation -d '{"channel":"integral.supply.chain","chaincode":"integralTrace","method":"convertIntegral","args":["xiaoou", "shopping_mall", "bank", "50"],"chaincodeVer":"v1"}'`
