# fee-crawler
A crawler that systematically browses SmartMEI website plan fee section.

## Table of Contents

- [Technology](#technology)
- [How it works](#how-it-works)
	- [Crawler](#crawler)
	- [Exchange rates API client](#exchange-rates-api-client)
	- [API](#api)
- [Developing](#developing)
	- [First Install](#first-install)
	- [Running the tests](#running-the-tests)
	- [Reseting your environment](#running-the-tests)
    

### Crawler
The SmartMEI fee crawler is powered by [colly](http://go-colly.org/). It has reduced crawl frontier and searches primarily for  **fees** content container element and maps it as a matrix.


The goal of mapping the container as a matrix is to increase change resilience, since web crawlers are vulnerable to failures caused by layout changes. So, as long as it remains a "matrix-like" container, the crawler should works fine.

### Exchange Rates API client
The Exchange Rates API Client is a REST client for [Exchange Rates API](https://exchangeratesapi.io/). According to [Europa Central Bank](https://www.ecb.europa.eu/stats/policy_and_exchange_rates/euro_reference_exchange_rates/html/index.en.html) Rate values are updated around 16:00 CET on every work day, so a [cache](#cache) model was built to improve response time.

#### Cache

Simple [Redis](https://redislabs.com/)-based cache.

#### Cache invalidate conditions

1 - Everyday at 16:30 CET (invalidate by an asynchronous task);  
2 - On every server startup.

#### Rates data source priority
1 - Local cache  
2 - Exchange Rates API

**OBS**: Exchange Rates API Client was built implements the [Repository](https://martinfowler.com/eaaCatalog/repository.html) design pattern.

### API
Single endpoint GraphQL API powered by [graphql-go](https://github.com/graphql-go/graphql):

|  Resource                        | Port  | HTTP Method |  
|----------------------------------|-------|-------------|  
| `http://localhost:9000/graphql`  | `9000`| `POST`      |

#### Querying plan transfer fees

Structure

`?query={ttransfer(domain:**<SMARTMEI_DOMAIN>**, plan:**<SMARTMEI_PLAN>**){    
  description, rates_date,BRL,USD,EUR  
}`


#### Parameters
| Name    |  Type    | Mandatory?  |
|---------|----------|-------------|
|  domain | `string` |  yes        |
|  plan   | `string` |  yes        |


Example:

Request:

```bash
curl --request POST \
  --url 'http://localhost:9000/graphql?query=%7Btransfer(domain%3A%22https%3A%2F%2Fwww.smartmei.com.br%22%2C%20plan%3A%22B%C3%A1sico%22)%7Bdescription%2C%20rates_date%2CBRL%2CUSD%2CEUR%7D%7D'

```

Response: 
```json
{
  "data": {
    "transfer": {
      "BRL": "7.00",
      "EUR": "1.25",
      "USD": "1.41",
      "description": "*Limitado a contas PF do dono da MEI",
      "rates_date": "2020-06-08"
    }
  }
}
```

**OBS**: The resolver uses parallelism to command actions to [Crawler](#crawler) and [Exchange rates API client](#exchange-rates-api-client) simultaneously.
