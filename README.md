# gocryptoadmin
This module contains functions to read, write, calculate on exchange _CSV_ exports. It is also possible to blend manual edited _CSV_ with
exchanges.

## Example

### Group and Keep Accounting Globally and each Individual Exchange 

This example is taken from an unit test where it uses three _CSV_ files.

**NOTE: This sample uses the _coinbase pro_ format since it is completely manually edited and fake!** 

1. Bank - buy _EUR_ from _Svensk Krona_ and transfer
2. Kraken - receive _EUR_ to buy _LTC_, transfer, and sell _LTC_
3. Coinbase Pro - receive _LTC_ and sell it there

_lf.csv (Bank)_ 
```bash
portfolio,trade id,product,side,sideid,created at,size,size unit,price,fee,total,price/fee/total unit
default,1,EUR-SEK,BUY,lf,2017-12-06T08:00:00.000Z,50.0,EUR,10,50.000000000,-550,SEK
default,2,EUR-EUR,TRANSFER,kr,2017-12-06T09:00:00.000Z,48,EUR,1,2,-50,EUR
```
_kr.csv (Kraken)_
```bash
portfolio,trade id,product,side,sideid,created at,size,size unit,price,fee,total,price/fee/total unit
default,1,EUR-EUR,RECEIVE,lf,2017-12-06T10:00:00.000Z,48,EUR,1,3.000000000,45,EUR
default,2,LTC-EUR,BUY,kr,2017-12-06T11:00:00.000Z,2,LTC,10,4,-24,EUR
default,3,LTC-LTC,TRANSFER,cb,2017-12-06T12:00:00.000Z,0.9,LTC,1,0.1,-1,LTC
default,4,LTC-EUR,SELL,kr,2017-12-06T14:00:00.000Z,1,LTC,20,2,18,EUR
```

_cb.csv (Coinbase Pro)_
```bash
portfolio,trade id,product,side,sideid,created at,size,size unit,price,fee,total,price/fee/total unit
default,1,LTC-LTC,RECEIVE,kr,2017-12-06T13:00:00.000Z,0.9,LTC,1,0.1,0.8,LTC
default,2,LTC-EUR,SELL,cb,2017-12-06T15:00:00.000Z,0.8,LTC,30,2,22,EUR
```

Running read all _CSV_ , default _chronological_ sorter, transaction group _processor_ (20h window), 
and multi account _processor_, the following output can be emitted.  

```bash

Exchange: all

Exchange	Side	Side Identifier		Date			    Pair	Size		Price		Fee		    Total		EUR		    LTC		    SEK
-------------------------------------------------------------------------------------------------------------------------------------------------------
lf		BUY		    lf			        2017-12-06 08:00:00	EUR-SEK	50.000000	10.000000	50.000000	-550.000000	50.000000	0.000000	-550.000000
lf		TRANSFER	kr			        2017-12-06 09:00:00	EUR-EUR	48.000000	1.000000	2.000000	-50.000000	0.000000	0.000000	-550.000000
kr		RECEIVE		lf			        2017-12-06 10:00:00	EUR-EUR	48.000000	1.000000	3.000000	45.000000	45.000000	0.000000	-550.000000
kr		BUY		    kr			        2017-12-06 11:00:00	LTC-EUR	2.000000	10.000000	4.000000	-24.000000	21.000000	2.000000	-550.000000
kr		TRANSFER	cb			        2017-12-06 12:00:00	LTC-LTC	0.900000	1.000000	0.100000	-1.000000	21.000000	1.000000	-550.000000
cb		RECEIVE		kr			        2017-12-06 13:00:00	LTC-LTC	0.900000	1.000000	0.100000	0.800000	21.000000	1.800000	-550.000000
kr		SELL		kr			        2017-12-06 14:00:00	LTC-EUR	1.000000	20.000000	2.000000	18.000000	39.000000	0.800000	-550.000000
cb		SELL		cb			        2017-12-06 15:00:00	LTC-EUR	0.800000	30.000000	2.000000	22.000000	61.000000	0.000000	-550.000000

Exchange: lf

Exchange	Side	Side Identifier		Date			    Pair	Size		Price		Fee		    Total		EUR		    SEK
-------------------------------------------------------------------------------------------------------------------------------------------
lf		BUY		    lf			        2017-12-06 08:00:00	EUR-SEK	50.000000	10.000000	50.000000	-550.000000	50.000000	-550.000000
lf		TRANSFER	kr			        2017-12-06 09:00:00	EUR-EUR	48.000000	1.000000	2.000000	-50.000000	0.000000	-550.000000

Exchange: kr

Exchange	Side	Side Identifier		Date			    Pair	Size		Price		Fee		    Total		EUR		    LTC
----------------------------------------------------------------------------------------------------------------------------------------
kr		RECEIVE		lf			        2017-12-06 10:00:00	EUR-EUR	48.000000	1.000000	3.000000	45.000000	45.000000	0.000000
kr		BUY		    kr			        2017-12-06 11:00:00	LTC-EUR	2.000000	10.000000	4.000000	-24.000000	21.000000	2.000000
kr		TRANSFER	cb			        2017-12-06 12:00:00	LTC-LTC	0.900000	1.000000	0.100000	-1.000000	21.000000	1.000000
kr		SELL		kr			        2017-12-06 14:00:00	LTC-EUR	1.000000	20.000000	2.000000	18.000000	39.000000	0.000000

Exchange: cb

Exchange	Side	Side Identifier		Date			    Pair	Size		Price		Fee		    Total		EUR		    LTC
----------------------------------------------------------------------------------------------------------------------------------------
cb		RECEIVE		kr			        2017-12-06 13:00:00	LTC-LTC	0.900000	1.000000	0.100000	0.800000	0.000000	0.800000
cb		SELL		cb			        2017-12-06 15:00:00	LTC-EUR	0.800000	30.000000	2.000000	22.000000	22.000000	0.000000
```


## Development

* This project uses [golines](https://github.com/segmentio/golines) - do a `go get -u github.com/segmentio/golines` (outside of gocryptoadmin) to install

* Follow the [setup instructions for golines for vs code](https://github.com/segmentio/golines#visual-studio-code)