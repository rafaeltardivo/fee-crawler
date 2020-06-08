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
	- [Reseting your environment](#reseting-your-environment)
- [More API Request Examples](#more-api-request-examples)

    

## Technology
- [Golang](https://www.python.org/) 1.14
- [go-colly](http://go-colly.org/docs/) v1.2.0
- [graphql-go](https://github.com/graphql-go/graphql) v0.7.9
- [gocron](https://github.com/go-co-op/gocron) v0.2.0
- [goquery](github.com/PuerkitoBio/goquery) v1.5.1
- [decimal](https://github.com/shopspring/decimal) v1.2.0
- [logrus](github.com/sirupsen/logrus) v1.6.0
- [gomega](github.com/onsi/gomega) v1.10.1
- [Docker](https://www.docker.com/) 19.03.6
- [Docker Compose](https://docs.docker.com/compose/) 1.25.0
- [Redis](https://redislabs.com/) 6.0


## How it Works

### Crawler
The SmartMEI fee crawler is powered by [colly](http://go-colly.org/). It has reduced crawl frontier and searches primarily for  **fees** content container element and maps it as a matrix:


|      ...       | Plano x      | Plan  y |
|----------------|--------------|---------|
| Transferência  |    [1,1]     |   [1,2] |
|      ...       |    [2,1]     |   [2,2] |



The goal of mapping the container as a matrix is to increase change resilience, since web crawlers are vulnerable to failures caused by layout changes. So, as long as it remains a "matrix-like" container, the crawler will won't be affected.

### Exchange Rates API client
The Exchange Rates API Client is a REST client for [Exchange Rates API](https://exchangeratesapi.io/). According to [Europa Central Bank](https://www.ecb.europa.eu/stats/policy_and_exchange_rates/euro_reference_exchange_rates/html/index.en.html), Rate values are updated around 16:00 CET on every work day, so a [cache](#cache) model was built to improve response time.

#### Cache

Simple [Redis](https://redislabs.com/)-based cache.

#### Cache invalidate conditions

1 - Everyday at 16:30 CET (invalidate by an asynchronous task);  
2 - On every server startup.

#### Rates data source priority
1 - Local cache  
2 - Exchange Rates API

**OBS**: Exchange Rates API Client implements the [Repository](https://martinfowler.com/eaaCatalog/repository.html) design pattern.

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

**OBS**: 
 - The API resolver uses parallelism to command actions to [Crawler](#crawler) and [Exchange rates API client](#exchange-rates-api-client) simultaneously;
 - `rates_date` represents the rate calculation date.

## Developing

### First Install
1 - Clone the project
```
git clone https://github.com/rafaeltardivo/fee-crawler  
```
2 - Build the application:  
```
make build
```  
4 - Run the application:  
```  
make up
```  
### Running the tests
```
make test  
```
**OBS**:
 - fee-crawler will run over port `9000`;  
 - redis service will run over port `6379`.

### Reseting your environment
If eventually you want to reset your environment, execute:
```
make destroy
```
After that, in order to run the application you'll need repeat the [First Install](#first-install) proccess.

### More API request examples

Query for plan **Básico** transfer fee (all fields): 

```bash
curl --request POST \
  --url 'http://localhost:9000/graphql?query=%7Btransfer(domain%3A%22https%3A%2F%2Fwww.smartmei.com.br%22%2C%20plan%3A%22B%C3%A1sico%22)%7Bdescription%2C%20rates_date%2CBRL%2CUSD%2CEUR%7D%7D'
```

Query for plan **Profissional** transfer fee (all fields):

```bash
curl --request POST \
  --url 'http://localhost:9000/graphql?query=%7Btransfer(domain%3A%22https%3A%2F%2Fwww.smartmei.com.br%22%2C%20plan%3A%22Profissional%22)%7Bdescription%2C%20rates_date%2CBRL%2CUSD%2CEUR%7D%7D'
```

Query for plan **Básico** transfer fee description:
```bash
curl --request POST \
  --url 'http://localhost:9000/graphql?query=%7Btransfer(domain%3A%22https%3A%2F%2Fwww.smartmei.com.br%22%2C%20plan%3A%22B%C3%A1sico%22)%7Bdescription%2C%7D%7D'
```

Query for plan **Profissional** transfer fee (USD and EUR):

```bash
curl --request POST \
  --url 'http://localhost:9000/graphql?query=%7Btransfer(domain%3A%22https%3A%2F%2Fwww.smartmei.com.br%22%2C%20plan%3A%22Profissional%22)%7BUSD%2CEUR%7D%7D'
```